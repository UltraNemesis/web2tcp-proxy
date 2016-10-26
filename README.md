# web2tcp-proxy
web2tcp-proxy is a simple bridging proxy application to allow Web applications running the browser to connect to a TCP server through a SockJS/Websocket interface and communicate with it.

**Warning:** Note that while the basic functionality works, several features are yet to be implemented.

## Build

web2tcp-proxy is written in Go langauge. In order to use it, make sure you have Go tooling in the PATH and run the following command. This will get the code and its dependencies, build it and will put the executable for your current platform in the bin directory.   

    go get github.com/UltraNemesis/web2tcp-proxy
    

##Running

Update the provided sample configuration file as per requirement. Frontend part refers to the HTTP server part that will handle the SockJS/WebSocket connections. Backend refers the tcp server to which the connections are bridged. 

Note that the configuration file can be in yaml or toml format as well as long as the same overall structure is maintained.

Once the configuration is ready, ensure that the file is in the same directory as the executable and run the executable. 

    {
      "frontend" : {
        "endpoint": ":8000",
        "route": "gateway",
        "tls": {
          "enabled": false,
          "certFile": "cert.pem",
          "certKeyFile": "certkey.pem"
        }
      },
      "backend": {
        "endpoint": "localhost:3000",
        "proxyProtocol": false,
        "tls": {
          "enabled": true,
          "skipVerify": true,
          "certAuthorityFile": "cacert.pem"
        }     
      }
    }

##Client side

A Web app can connect to the proxy using either SockJS or pure Websocket. When SockJS is used, fallback's are available when the browser doesn't have support for websocket. The binary data send by the server arrives at the client side in the form of Base64 encoded strings and must be decoded to get the byte array and simplely, while sending, the byte array must be encoded into Base64 

**Usage with SockJS**

    // Usage with SockJS
     var sock = new SockJS('http://localhost:8000/gateway');
     sock.onopen = function() {
         console.log('open');
     };
     sock.onmessage = function(e) {
         console.log('message as base64 string : ', e.data);
         console.log('message as byte array :', base64.decode(e.data)
     };
     sock.onclose = function() {
         console.log('close');
     };
    sock.send(base64.encode(byteArray));


**Usage with WebSocket**

    // Usage with WebSocket
     var sock = new WebSocket('ws://localhost:8000/gateway/websocket');
     sock.onopen = function() {
         console.log('open');
     };
     sock.onmessage = function(e) {
         console.log('message as base64 string : ', e.data);
         console.log('message as byte array :', base64.decode(e.data)
     };
     sock.onclose = function() {
         console.log('close');
     };
    sock.send(base64.encode(byteArray));

