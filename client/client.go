package client

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	host     string
	port     string
	logger   *log.Logger
	shutdown chan int
}

func NewClient(host, port string, logger *log.Logger) (*Client, error) {
	return &Client{
		host:   host,
		port:   port,
		logger: logger,
	}, nil
}

func (this *Client) Run(mch chan int) {
	this.shutdown = make(chan int, 1)
	// serv
	var shutdown int
	var url string = "http://" + this.host + ":" + this.port + "/"
	this.logger.Println("client begin notice.")

forBreak:
	for {
		select {
		case shutdown = <-this.shutdown:
			break forBreak
		default:
			this.handle(url)
		}
		time.Sleep(time.Second * 10)
	}

	// shutdown
	this.logger.Println("receive close & close listen.")
	mch <- shutdown
}

func (this *Client) Shutdown() {
	this.shutdown <- 1
}

func (this *Client) handle(uri string) {
	_, err := http.PostForm(uri, url.Values{"up": {"true"}})
	if err != nil {
		this.logger.Println(err)
	} else {
		// ok
	}
}
