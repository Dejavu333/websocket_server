# Explanation of the code
This Go code defines a WebSocket server that can handle multiple channels, broadcast data to all clients in a given channel, and be started and stopped. It uses the gorilla/websocket package to manage WebSocket connections.
## IWebSocketServer interface
This interface defines the methods that any implementation of the WebSocket server must implement. It includes methods for starting and stopping the server, broadcasting data to clients in a channel, and adding a channel to the server.
## DefaultWebSocketServer struct
This struct is the default implementation of the IWebSocketServer interface. It has the following fields:
- upgrader: a websocket.Upgrader instance that is used to upgrade HTTP connections to WebSocket connections.
- clients: a map of WebSocket connections to a boolean value indicating whether they are still active.
- channels: a map of channel names to maps of WebSocket connections to a boolean value indicating whether they are subscribed to the channel.
The struct also includes the following methods:
- Start(): starts the WebSocket server by creating an HTTP server and listening on the specified host and port. It also sets up a route for handling HTTP requests and upgrading them to WebSocket connections.
- Stop(): stops the WebSocket server by closing all WebSocket connections and deleting them from the clients and channels maps.
- handleHTTPRequest(): handles incoming HTTP requests and upgrades them to WebSocket connections if the requested URL matches an existing channel. If the URL does not match an existing channel, it returns a 404 Not Found error.
- upgradeWebSocket(): upgrades an HTTP connection to a WebSocket connection and adds the connection to the clients and channels maps.
- Broadcast(): sends a message to all WebSocket connections in a given channel.
- AddChannel(): adds a new channel to the server.
## Utility functions
getEnvOrDefault(): retrieves the value of an environment variable or returns a default value if the variable is not set.
# Conclusion
This Go code defines a simple WebSocket server that can handle multiple channels and broadcast data to clients in a given channel. It can be easily customized to fit specific use cases and integrated into existing Go applications.
# Usage
## Backend setup
```go
    // SETUP THE WEBSOCKET SERVER
	websocketServer := websockets.NewDefaultWebSocketServer()
    // ADD CHANNELS YOU WANT TO USE
	websocketServer.AddChannel("/ws/test1")
	websocketServer.AddChannel("/ws/test2")

	go func() {
		websocketServer.Start()
	}()

	// SEND DATA TO THE WEBSOCKET SERVER IN A LOOP OR WHENEVER YOU WANT AND WHATEVER YOU WANT
	for {
		time.Sleep(5 * time.Second)
		websocketServer.Broadcast("/ws/test1", "test message one")
		websocketServer.Broadcast("/ws/test2", "test message two")
	}

    // YOU CAN ALSO SEND STRUCTS
    type TestStruct struct {
        Name string
        Age int
    }
    websocketServer.Broadcast("/ws/test1", TestStruct{Name: "Naruto", Age: 30})
    websocketServer.Broadcast("/ws/test2", TestStruct{Name: "Yagami", Age: 25})
```
## Frontend setup
```js
    // I RECOMMEND USING THIS HELPER FUNCTION TO ESTABLISH A CONNECTION WITH THE WEBSOCKET SERVER
    // IN A CALLBACK FUNCTION YOU CAN DO WHATEVER YOU WANT WITH THE DATA YOU RECEIVE FROM THE WEBSOCKET SERVER
    function connectToWebSocketServer(p_url, p_channelName, p_functionInvokedOnMessageFromServer) {

        const socket = new WebSocket(p_url);
        socket.onopen = (event) => {
            console.log("Connected to WebSocket server: " + p_url + " on channel: " + p_channelName);
            socket.send(p_channelName);
        };

        socket.onmessage = (event) => {
            const parsedData = JSON.parse(event.data);
            p_functionInvokedOnMessageFromServer(parsedData);
        };

        socket.onclose = (event) => {
            console.log("Disconnected from WebSocket server");
        };
    }

    // THIS IS THE CALLBACK FUNCTION THAT WILL BE INVOKED WHEN THE WEBSOCKETSERVER SENDS DATA, 
    // YOU CAN DO WHATEVER YOU WANT WITH THE DATA
    function outputData(p_data) {
        console.log("Received message:", data);
        // DO YOUR STUFF HERE
    }

    // INVOKING THE HELPER FUNCTION
    connectToWebSocketServer("ws://localhost:8080", "/ws/test1", outputData);
    connectToWebSocketServer("ws://localhost:8080", "/ws/test2", outputData);


    // OR YOU CAN USE IT EXPLICITLY LIKE THIS
    // SOCKET1
    const socket = new WebSocket('ws://localhost:8080/ws/test1');
    socket.onopen = (event) => {
    console.log("Connected to WebSocket server");
    };
    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log("Received message:", data);
        // DO YOUR STUFF HERE
    };
    socket.onclose = (event) => {
    console.log("Disconnected from WebSocket server");
    };

    // SOCKET2
    const socket2 = new WebSocket('ws://localhost:8080/ws/test2');
    socket2.onopen = (event) => {
    console.log("Connected to WebSocket server");
    };
    socket2.onmessage = (event) => {
        const data = JSON.parse(vent.data);
        console.log("Received message:", data);
        // DO YOUR STUFF HERE
    };
    socket2.onclose = (event) => {
    console.log("Disconnected from WebSocket server");
    };
```