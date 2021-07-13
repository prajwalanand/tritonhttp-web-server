package tritonhttp

import (
	"net"
	"time"
	"fmt"
	"path/filepath"
	"strings"
	"io"
	"os"
)


/* 
For a connection, keep handling requests until 
	1. a timeout occurs or
	2. client closes connection or
	3. client sends a bad request //-Anything other than GET,HTTP/1.1, url doesn't start with /, invalid header (string:string), host header not found
*/
func (hs *HttpServer) handleConnection(conn net.Conn) {
	//panic("todo - handleConnection")

	DELIM := "\r\n"
	EMPTY_BUF := string(make([]byte,50))

	// Start a loop for reading requests continuously
	remaining := ""
	isFirstLine := true 
	request := &HttpRequestHeader{
						URL: "",
						FieldMap: make(map[string]string),
					}
	url := ""
	isConnectionClose := false
	closeConn := false

	//ch := make(chan bool, 1)
	//ch <- true
	//timeout :=  make(chan bool, 1)
	//defer close(ch)
	//defer conn.Close()
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	timeout <- true
	// 	ch <- false
	// }()
	
	// Set a timeout for read operation
	//timer := time.NewTimer(5 * time.Second)
	//defer timer.Stop()

	buf := make([]byte, 50)


	


	for{
		//fmt.Println("Inside for")
		//select {
		//	case <-ch:
		err := conn.SetReadDeadline(time.Now().Add(5*time.Second))
		//if err != nil {
		
			fmt.Println("Read from ch")
		    // Read from the connection socket into a buffer
			size, err := conn.Read(buf)
			print("No of bytes read = ",size)
			if size == 0 {
				fmt.Println("Size zero error! Partial request read and timeout happened")
				hs.handleBadRequest(conn)
				conn.Close()
				return
			}
			if err != nil {
        		if err != io.EOF {
            		fmt.Println("Read error:", err) //TO-DO: WHAT SHOULD WE DO?	
        		}
        		if os.IsTimeout(err){
					 fmt.Println("Timed out")
					 if string(buf) != EMPTY_BUF {
					  	fmt.Println("ERROR: Partial request read and timeout happened")
					  	hs.handleBadRequest(conn)
					  }
					  fmt.Println("ALERT: Server is closing the connection")
					  conn.Close()
					  //conn = nil
					  return
			//return nil, err
				}
        		fmt.Println("Read error: ", err)
        		hs.handleBadRequest(conn)
        		fmt.Println("ALERT: Server is closing the connection")
				conn.Close()
				//conn = nil
				return
    		}
    		data := buf[:size]
    		remaining = remaining + string(data)
    		for strings.Contains(remaining, DELIM) {
				idx := strings.Index(remaining, DELIM)
				requestLine := remaining[:idx]
				remaining = remaining[idx+2:]

				fmt.Println("Processing requestline: ",requestLine)
				requestLine = strings.TrimSpace(requestLine)
				// Validate the request lines that were read
				if isFirstLine {
					//elements := strings.Split(requestLine," ")
					elements := strings.Fields(requestLine)
					if (len(elements)!=3 || len(elements[1])==0 || elements[0] != "GET" || elements[2] != "HTTP/1.1" || elements[1][0] != '/') {
						fmt.Println("ERROR: Request line length = ",len(requestLine))
						fmt.Println("ERROR: Bad request type ",elements)
						//hs.handleBadRequest(conn)
						//fmt.Println("ALERT: Server is closing the connection")
						closeConn = true
						//conn.Close()
						//conn = nil
						//return
					} 
					url = elements[1]
					isFirstLine = false
				} else {
					// Handle any complete requests
					if requestLine == ""{
						if closeConn{
							fmt.Println("ALERT: Server is closing the connection")
							hs.handleBadRequest(conn)
							conn.Close()
							return
						} 
						if _,ok := request.FieldMap["Host"];ok{
							if url[len(url)-1:]=="/" {
								url = url + "index.html"
							}
							url = hs.DocRoot + url
							absoluteURL,err1 := filepath.Abs(url)
							if err1 != nil {
								fmt.Println("ERROR: Error in resolving absolute filepath of url ",url)
							}
							absoluteDocRoot,err2 := filepath.Abs(hs.DocRoot)
							if err2 != nil {
								fmt.Println("ERROR: Error in resolving absolute filepath of docroot ",hs.DocRoot)
							}
							request.URL = absoluteURL
							if !strings.HasPrefix(absoluteURL,absoluteDocRoot) || !fileExists(absoluteURL){
								fmt.Println("ERROR: Invalid FileURL - ",absoluteURL)
								hs.handleFileNotFoundRequest(request,conn)
								isFirstLine = true
								request = &HttpRequestHeader{
												URL: "",
												FieldMap: make(map[string]string),
											}
								url = ""
								//buf = make([]byte, 50)
								 //ch <- true
								 //func()()
								//timer.Reset(5 * time.Second)
							} else {
								hs.handleResponse(request,conn)	
								// If reusing read buffer, truncate it before next read
								isFirstLine = true
								request = &HttpRequestHeader{
												URL: "",
												FieldMap: make(map[string]string),
											}
								url = ""
								//buf = make([]byte, 50)
								 //ch <- true
								//func()()
								//timer.Reset(5 * time.Second)
							}
							if isConnectionClose{
								fmt.Println("ALERT: Closing the connection because header says so!")
								conn.Close()
								return
							}
						} else {
							fmt.Println(request.FieldMap)
							fmt.Println("ERROR: Bad request type. Missing host header")
							hs.handleBadRequest(conn)
							fmt.Println("ALERT: Server is closing the connection")
							conn.Close()
							return
						}
					} else {
						elements := strings.Split(requestLine,":")
						if len(elements)<2 {
							fmt.Println("ERROR: Bad request type. Missing header(s) of the format string:string")
							//hs.handleBadRequest(conn)
							//fmt.Println("ALERT: Server is closing the connection")
							closeConn = true
							//conn.Close()
							//return
						} else {
							// Update any ongoing requests
							key := strings.TrimSpace(elements[0])
							value := strings.TrimSpace(elements[1])
							if key=="Connection" && value=="close" {
								isConnectionClose = true
							}
							request.FieldMap[key] = value
						}
					}
				}
			}
		}
	}