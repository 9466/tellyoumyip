tellyoumyip 0.1 beta

this is a test util for ipaddress send/receive tool, it is useful for get intranet gateway ipaddress. 

```
Usage: tellyoumyip -m {client|server} [OPTIONS]
  -m <client|server> 	 run as a server or client mode. 
  -h <host> 		 server mode <host> is listen ipaddress, default 0.0.0.0 
            		 client mode <host> is server ipaddress. 
  -p <port> 		 server mode <port> is listen port, default 9404 
            		 client mode <port> is server port. 
  -L <file> 		 logfile, default none log. 
  -P <file> 		 pidfile, default none pidfile, client mode not need. 
  --help    		 Output this help and exit. 
  --version 		 Output version and and exit. 

Examples:
  tellyoumyip -m client -h 192.168.2.3 -L /var/log/tellyoumyip.client.log
  tellyoumyip -m server -P /var/run/tellyoumyip.pid -L /var/log/tellyoumyip.server.log
```
