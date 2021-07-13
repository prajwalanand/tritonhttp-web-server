package tritonhttp

import (
	"net"
	"fmt"
	"os"
	"io"
	"strconv"
	"strings"
	"bufio"
)

//400
func (hs *HttpServer) handleBadRequest(conn net.Conn) {
	fmt.Println("We are handling bad request (400) now ...")
	response := HttpResponseHeader{
						InitialLine: "HTTP/1.1 400 Bad Request",
						FieldMap: make(map[string]string),
				}
	response.FieldMap["Server"] = "Go-Triton-Server-OMS-224"
	//response.FieldMap["Content-Type"] = "text/html"
	//response.URL = "src/error_400.html"
	response.FieldMap["Content-Length"] = "0"//getContentLength(response.URL)
	hs.sendResponse(response,conn)
}

//404
func (hs *HttpServer) handleFileNotFoundRequest(requestHeader *HttpRequestHeader, conn net.Conn) {
	fmt.Println("We are handling FileNotFoundRequest (404) now ...")
	response := HttpResponseHeader{
						InitialLine: "HTTP/1.1 404 Not Found",
						FieldMap: make(map[string]string),
				}
	response.FieldMap["Server"] = "Go-Triton-Server-OMS-224"
	//response.FieldMap["Content-Type"] = "text/html"
	//response.URL = "src/error_404.html"
	response.FieldMap["Content-Length"] = "0"//getContentLength(response.URL)
	if(requestHeader.FieldMap["Connection"]=="close"){
		response.FieldMap["Connection"] = "close"
	}
	hs.sendResponse(response,conn)
}

//200 
func (hs *HttpServer) handleResponse(requestHeader *HttpRequestHeader, conn net.Conn) {
	fmt.Println("We are handling Valid Request (200) now ...")
	response := HttpResponseHeader{
						InitialLine: "HTTP/1.1 200 OK",
						FieldMap: make(map[string]string),
				}
	response.FieldMap["Server"] = "Go-Triton-Server-OMS-224"
	response.FieldMap["Last-Modified"] = getLastModifiedTime(requestHeader.URL)
	response.FieldMap["Content-Length"] = getContentLength(requestHeader.URL)
	response.FieldMap["Content-Type"] = getContentType(requestHeader.URL,hs.MIMEMap)
	if(requestHeader.FieldMap["Connection"]=="close"){
		response.FieldMap["Connection"] = "close"
	}
	response.URL = requestHeader.URL
	hs.sendResponse(response,conn)
}

func (hs *HttpServer) sendResponse(responseHeader HttpResponseHeader, conn net.Conn) {
	fmt.Println("Sending the response now...")
	DELIM := "\r\n"

	responseMsg := responseHeader.InitialLine + DELIM 

	// Send headers
	// valbytes := []byte(responseHeader.InitialLine + DELIM)
	// _, err := conn.Write(valbytes)
	// fmt.Println("Initial line is written...of size ",len(valbytes))
	// if err != nil {
	// 	return
	// }
	for k, v := range responseHeader.FieldMap { 
    	line := k+":"+v
    	//valbytes = []byte(line + DELIM)
    	responseMsg = responseMsg + line + DELIM
    	fmt.Println("Writing header key:"+k+" value:"+v+"...")
  //   	_, err = conn.Write(valbytes)
		// if err != nil {
		// 	return
		// }
	}	

	responseMsg = responseMsg + DELIM

	//valbytes = []byte(DELIM)
	fmt.Println("Writing end delimiter...")

	fmt.Println("Writing final response headers...")

	// valbytes := []byte(responseMsg)
	// _, err := conn.Write(valbytes)
	// if err != nil {
	// 	fmt.Println("Error in writing response headers to conn: ",err)
	// 	return
	// }

	
	// Hint - Use the bufio package to write response
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(responseMsg)
	if err != nil {
		fmt.Println("Error in writing response headers to conn: ",err)
		return
	}else{
		writer.Flush()
	}

	// Send file if required
	if strings.Contains(responseHeader.InitialLine,"OK"){
		fmt.Println("Writing response body ...")
		//readAndWriteFileInChunks(responseHeader,conn,writer)
		var BufferSize int64 = 100 
		filesize,err := strconv.ParseInt(responseHeader.FieldMap["Content-Length"],10,64)
		file, err := os.Open(responseHeader.URL)
		if err != nil {
		  fmt.Println(err)
	  	  return
		}
		defer file.Close()

		buffer := make([]byte, BufferSize)

		var i int64 = 0

		for i=0; i<filesize/BufferSize; i++ {
	  		_, err := file.Read(buffer)

	  		if err != nil {
	    		if err != io.EOF {
	      			fmt.Println(err)
	    		}
	    		break
	  		}
	  		_, err = writer.Write(buffer)//conn.Write(buffer)
			if err != nil {
				return
			}else{
				writer.Flush()
			}
		}
		buffer = make([]byte, filesize % BufferSize)
		_, err = file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
		}
		_, err = writer.Write(buffer)//conn.Write(buffer)
		if err != nil {
			return
		}else{
			writer.Flush()
		}

	}

	fmt.Println("Writing response completed!")
	writer.Flush()
}

// func readAndWriteFileInChunks(responseHeader HttpResponseHeader, conn net.Conn, writer *Writer){
// 	var BufferSize int64 = 100 
// 	filesize,err := strconv.ParseInt(responseHeader.FieldMap["Content-Length"],10,64)
// 	file, err := os.Open(responseHeader.URL)
// 	if err != nil {
// 	  fmt.Println(err)
//   	  return
// 	}
// 	defer file.Close()

// 	buffer := make([]byte, BufferSize)

// 	var i int64 = 0

// 	for i=0; i<filesize/BufferSize; i++ {
//   		_, err := file.Read(buffer)

//   		if err != nil {
//     		if err != io.EOF {
//       			fmt.Println(err)
//     		}
//     		break
//   		}
//   		_, err = writer.Write(buffer)//conn.Write(buffer)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	buffer = make([]byte, filesize % BufferSize)
// 	_, err = file.Read(buffer)
// 	if err != nil {
// 		if err != io.EOF {
// 			fmt.Println(err)
// 		}
// 	}
// 	_, err = writer.Writer(buffer)//conn.Write(buffer)
// 	if err != nil {
// 		return
// 	}else{
// 		writer.Flush()
// 	}
// }
