package tritonhttp
import (
	"flag"
	"log"
	"net"
)
/** 
	Initialize the tritonhttp server by populating HttpServer structure
**/
func NewHttpdServer(port, docRoot, mimePath string) (*HttpServer, error) {
	//panic("todo - NewHttpdServer")

	// Initialize mimeMap for server to refer
	mimeMap,err := ParseMIME(mimePath)
	if err != nil {
		log.Panicln(err)
	}

	s := &HttpServer{
			ServerPort:	port,
			DocRoot:	docRoot,
			MIMEPath:	mimePath,
			MIMEMap:	mimeMap,
		}

	// Return pointer to HttpServer
	return s,err
}

/** 
	Start the tritonhttp server
**/
func (hs *HttpServer) Start() (err error) {
	//panic("todo - StartServer")

	// Start listening to the server port
	port := flag.String("port", hs.ServerPort, "Port to accept connections on.")
	host := flag.String("host", "127.0.0.1", "Host or IP to bind to")
	flag.Parse()

	l, err := net.Listen("tcp", *host+*port)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Listening to connections at '"+*host+"' on port", *port)
	defer l.Close()


	// Accept connection from client
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		// Spawn a go routine to handle request
		go hs.handleConnection(conn)
	}

}

