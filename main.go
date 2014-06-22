package main

import (
	"fmt"
	"github.com/9466/tellyoumyip/client"
	"github.com/9466/tellyoumyip/server"
	"github.com/9466/tellyoumyip/util"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

const (
	VERSION  = "0.1 beta"
	DEF_PORT = "9404"
)

type modeServer interface {
	Run(chan int)
	Shutdown()
}

var logger *log.Logger

func main() {
	// check args
	argc := len(os.Args)
	if argc < 2 {
		help()
		os.Exit(1)
	}

	var mode, host, port, logfile, pidfile string
	for key, val := range os.Args {
		switch val {
		case "-m":
			if argc > key+1 {
				mode = os.Args[key+1]
			}
		case "-h":
			if argc > key+1 {
				host = os.Args[key+1]
			}
		case "-p":
			if argc > key+1 {
				port = os.Args[key+1]
			}
		case "-L":
			if argc > key+1 {
				logfile = os.Args[key+1]
			}
		case "-P":
			if argc > key+1 {
				pidfile = os.Args[key+1]
			}
		case "--version":
			version()
			os.Exit(0)
		case "--help":
			help()
			os.Exit(0)
		}
	}

	if mode == "" || (mode == "client" && host == "") {
		help()
		os.Exit(0)
	}

	if port == "" {
		port = DEF_PORT
	}

	// daemon
	var err error
	_, err = util.Daemon(1, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// init log
	var file *os.File
	if logfile != "" {
		file, err = os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	} else {
		file, err = os.OpenFile("/dev/null", 0, 0)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	// log options

	logger.Println("mode \t", mode)
	logger.Println("host \t", host)
	logger.Println("port \t", port)
	logger.Println("logfile \t", logfile)
	logger.Println("pidfile \t", pidfile)

	// init mode
	mChan := make(chan int, 1)
	var srv modeServer
	if mode == "client" {
		srv, err = client.NewClient(host, port, logger)
		if err != nil {
			logger.Fatalln("client error: ", err)
		}
		go srv.Run(mChan)
	} else if mode == "server" {
		srv, err = server.NewServer(host, port, pidfile, logger)
		if err != nil {
			logger.Fatalln("server error: ", err)
		}
		go srv.Run(mChan)
	} else {
		logger.Println("mode[", mode, "] not support")
		return
	}

	// trap signal
	signalChan := make(chan os.Signal, 10)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT,
		syscall.SIGHUP, syscall.SIGSTOP, syscall.SIGQUIT)

	go func(ch <-chan os.Signal) {
		sig := <-ch
		logger.Println("signal recieved " + sig.String() + ", at: " + time.Now().String())
		srv.Shutdown()
		if sig == syscall.SIGHUP {
			logger.Println("restart now ...")
			procAttr := new(os.ProcAttr)
			procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
			procAttr.Dir = os.Getenv("PWD")
			procAttr.Env = os.Environ()
			process, err := os.StartProcess(os.Args[0], os.Args, procAttr)
			if err != nil {
				logger.Println("restart process failed: " + err.Error())
				return
			}
			waitMsg, err := process.Wait()
			if err != nil {
				logger.Println("restart wait error: " + err.Error())
			}
			logger.Println(waitMsg)
		} else {
			logger.Println("shutdown now...")
		}
	}(signalChan)

	// wait for shutdown
	<-mChan

	logger.Println("shutdown ok.")
}

func help() {
	prog := path.Base(os.Args[0])
	fmt.Println(prog + " " + VERSION)
	fmt.Println("")
	fmt.Println("this is a test util for ipaddress send/receive tool")
	fmt.Println("it is useful for get intranet gateway ipaddress. ")
	fmt.Println("")
	fmt.Println("Usage: " + prog + " -m {client|server} [OPTIONS]")
	fmt.Println("  -m <client|server> \t run as a server or client mode. ")
	fmt.Println("  -h <host> \t\t server mode <host> is listen ipaddress, default 0.0.0.0 ")
	fmt.Println("            \t\t client mode <host> is server ipaddress. ")
	fmt.Println("  -p <port> \t\t server mode <port> is listen port, default 9404 ")
	fmt.Println("            \t\t client mode <port> is server port. ")
	fmt.Println("  -L <file> \t\t logfile, default none log. ")
	fmt.Println("  -P <file> \t\t pidfile, default none pidfile, client mode not need. ")
	fmt.Println("  --help    \t\t Output this help and exit. ")
	fmt.Println("  --version \t\t Output version and and exit. ")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  " + prog + " -m client -h 192.168.2.3 -L /var/log/" + prog + ".client.log")
	fmt.Println("  " + prog + " -m server -P /var/run/" + prog + ".pid -L /var/log/" + prog + ".server.log")
	fmt.Println("")
}

func version() {
	prog := path.Base(os.Args[0])
	fmt.Println(prog + " " + VERSION)
	fmt.Println("")
}
