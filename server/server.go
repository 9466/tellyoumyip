package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	host     string
	port     string
	pidfile  string
	logger   *log.Logger
	shutdown chan int
}

var ips []string

func NewServer(host, port, pidfile string, logger *log.Logger) (*Server, error) {
	if pidfile != "" {
		pid := os.Getpid()
		err := ioutil.WriteFile(pidfile, []byte(strconv.Itoa(pid)), 0666)
		if err != nil {
			return &Server{}, errors.New("pid " + pidfile + err.Error())
		}
	}

	return &Server{
		host:    host,
		port:    port,
		pidfile: pidfile,
		logger:  logger,
	}, nil
}

func (this *Server) Run(mch chan int) {
	// listen
	ln, err := net.Listen("tcp", this.host+":"+this.port)
	if err != nil {
		this.logger.Fatalln("listen error: " + err.Error())
	}

	this.shutdown = make(chan int, 1)
	go func(ch, mch chan int) {
		shutdown := <-ch
		this.logger.Println("receive close & close listen.")
		ln.Close()
		mch <- shutdown
	}(this.shutdown, mch)

	// serv
	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)
	srv := &http.Server{
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ips = make([]string, 1)
	ips[0] = "127.0.0.1|" + time.Now().Format("2006-01-02 15:04:05")

	this.logger.Println("server begin receive.")
	if err = srv.Serve(ln.(*net.TCPListener)); err != nil {
		this.logger.Println(err)
	}
}

func (this *Server) Shutdown() {
	if this.pidfile != "" {
		syscall.Unlink(this.pidfile)
	}
	this.shutdown <- 1
}

func handle(w http.ResponseWriter, req *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	// use post update ip
	if req.Method == "POST" {
		ip := strings.Split(req.RemoteAddr, ":")[0]
		time := time.Now().Format("2006-01-02 15:04:05")
		ipnum := len(ips)
		if strings.Split(ips[ipnum-1], "|")[0] != ip {
			ips = append(ips, ip+"|"+time)
		} else {
			// now, ip not update
		}
	}

	nowip := strings.Split(ips[len(ips)-1], "|")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Welcome to tellyouip server page! </h1><br/>")
	fmt.Fprintf(w, "you current gateway ip: <strong>%s</strong> <br/><br/>", nowip[0])
	fmt.Fprintf(w, "ip update list:<br/>")
	for i := len(ips) - 1; i >= 0; i-- {
		nowip = strings.Split(ips[i], "|")
		fmt.Fprintf(w, "%s \t %s <br/>", nowip[0], nowip[1])
	}
	fmt.Fprintf(w, "%s", req.Method)
}
