from socket import socket
# Create connection to the server
s = socket()
s.connect(("localhost",8080))
# Compose the message/HTTP request we want to send to the server
#msgPart1 = b"GET /index.html HTTP/1.1\r\nHost: Ha\r\n\r\n"

msgPart1 = b"GET /index.html  HTTP/1.1\r\nHost: Ha1\r\nRequestNo: 1\r\n\r\nGET / HTTP/1.1\r\nHost: Ha1\r\nRequestNo: 2\r\nConnection\r\n\r\nGET /index.html HTTP/1.1\r\nHost: Ha1\r\nRequestNo: 3\r\n\r\n"
# Send out the request
s.sendall(msgPart1)
# Listen for response and print it out
resp = s.recv(100)
while resp!=b'': 
	try:
		#print("I am printing my response!")
		#print (resp)
		resp = s.recv(100)
		print(str(resp))
	except Exception as e:
		continue
		#print("Testing read error: ",e)
#print("Client is closing the connection...")
s.close()
#print("Client closed the connection")