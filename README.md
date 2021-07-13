# TritonHTTP Web Server
## Usage

Step 1:
You need to add directory into ```.bash_profile``` (for OS X environment) or ```.bashrc``` to let the compiler knows where to find the dependencies
```
	export PATH=$PATH:/usr/local/go/bin     # making sure go is on path
	export GOPATH=<path-to-repo>
	export PATH=$PATH:$GOPATH/bin           
 ```
example:
```
	export GOPATH=/[The directory you put this folder]/Project-1/
	export PATH=$PATH:$GOPATH/bin
```

^ you need to add these lines in your .bashrc file and 
```source .bashrc``` to activate them.


Step 2:
When you run your code, we are currently using the run-server.sh script. It basically runs a immediate compilation and looks for dependency. You will encounter errors like:
```
	cannot find package "github.com/go-ini/ini" in any of:
		/usr/local/go/src/github.com/go-ini/ini (from $GOROOT)
		/Desktop/Project-1/src/github.com/go-ini/ini (from $GOPATH)
```
That's because we usually don't include the dependency directory into our repository. Thus, we need to download the dependency.
example:
```
	go get github.com/go-ini/ini
```


Step 3:
Now, you should be able to run the script. 
```
	run-server.sh [config_file]
```
In the config file, under the [httpd] section, set 
<li> `use_default_server=true` to run Go's internal http web server 
<li> `use_default_server=false` to run your TritonHTTP web server

## Building a simple web server that implements TritonHTTP

### main
The server program starts in
```
/src/main/main.go
```

### tritonhttp
This package contains the code for developing the tritonhttp web server. <br><br>

1. `httpd_structs.go` has the structures that can be used to initialize the new server
2. `httpd_server.go` has the methods for initializing and starting the new server
3. `http_request_handler.go` has the methods for handling an incoming connection and subsequent requests
4. `http_response_handler.go` has the methods for handling various types of responses based on the request
5. `server_utils.go` has utility functions



## Testing

For Simple testing : <br>

1. Set `use_default_server=false` in the config file and start the server using `run-server.sh [config_file]`.

2. Once the server is launched, test it's functionality using a browser, curl, netcat or a custom client.

3. A directory `sample_htdocs` (populated with some sample files) has been provided in the code to help with testing. Verify the TritonHTTP server's functionality by accessing the files in this directory

Starter code provided by George Porter, UCSD.