# websockets package
__Each struct that implements the IWebSocketServer interface must define 4 methods:__
- Start()
- Stop()
- AddChannel(p_channelName string)
- Broadcast(p_channelName string, p_data interface{})

__The websockets package contains a default implementation of the IWebSocketServer interface, the DefaultWebSocketServer struct. I provide guidelines on how to use it below.__

## backend setup for websockets package
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

## frontend setup for websockets package
__The DefaultWebSocketServer broadcasts data in the form of json strings, so you have to decode them.__
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