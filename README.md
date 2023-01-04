# Go gRPC Streaming with Custom Reconnect Mechanism
<p>
    A simple restaurant ordering system that implements <strong>gRPC</strong> and <strong>websocket</strong>.
    In this project, the usage of all gRPC methods will be demonstrated. Nevertheless, bi-directional streaming with gRPC and websocket is the most important part among other gRPC methods in this project. A custom reconnect mechanism together with <strong>concurrent</strong> approach will be implemented in client side gRPC bi-directional streaming, to ensure gRPC client will reconnect to gRPC server whenever the server is down, and running again. 
    Besides, client side will be using <a href="https://github.com/cloudwego/hertz" target="_blank" rel="noopener noreferrer">Hertz</a> as HTTP framework, a open-source project created by ByteDance, for handling REST APIs and websocket.
</p>

<br/>

## Tech Stacks
- Go
- Hertz
- Websocket
- gRPC

<br/>

## gRPC Methods
- Unary
- Server streaming
- Client streaming
- Bi-directional streaming

<br/>

## Protocol Buffer Compiler
Currently, pb.go files already generated in client and server folder. If you have to compile again for the **.proto** files, please run the following:
```sh
protoc --go_out=./client --go-grpc_out=./client ./proto/order.proto
protoc --go_out=./server --go-grpc_out=./server ./proto/order.proto
```
for compiling proto files into pb.go files.

<br/>

## Reconnect Mechanism
<p>
    A reconnect mechanism with <strong>concurrent</strong> approach is implemented in gRPC client side <strong>bi-directional streaming</strong>, depending on websocket connection.
</p>

### Reestablish Stream Connection
A goroutine will be spawned that running an infinite for loop that listening on reconnect channel. Every task needs to reconnect to gRPC server side will trigger this channel.
```go
case <-reconnectCh:
```
In this reconnect logic, it will wait until the connection is ready, then spawn a new goroutine for listening gRPC server side.

### Listen gRPC Server Side
A goroutine will be spawned for listening streaming from gRPC server side as it is bi-directional streaming,. Once any error occured during the listening process, it will trigger the goroutine 'Reestablish Stream Connection'. After that, this goroutine will be killed.

```go
// whenever need to trigger reconnect stream connection goroutine
reconnectCh <- struct{}{}
```

### Wait Until Ready
A for loop is called until the client connection state is 'ready'. A short duration is reserved for the state to change from 'idle' or 'transient_failure' to 'ready. Any reconnect time interval and/or reconnect attempts can be definite here. This function is called whenever the client connection has error connecting to gRPC server side.

```go
for {

    // some checking...

    if client.Conn.GetState() != connectivity.Ready {
        client.Conn.Connect()
    }

    // reserve a short duration (customizable) for conn to change state from idle to ready if grpc server is up
    time.Sleep(500 * time.Millisecond)

    if client.Conn.GetState() == connectivity.Ready {
        return true
    }

    // define reconnect time interval (backoff) or/and reconnect attempts here
    time.Sleep(2 * time.Second)
}
```

### Websocket Connection
Once the websocket connection is established, a global boolean variable is also created to determine whether websocket connection still alive. After the websocket connection is demolished, it will be false.

```go
var isConnectedWebSocket *bool
t := true
isConnectedWebSocket = &t

defer func() {
    *isConnectedWebSocket = false
}()
```

This is to ensure some tasks will only run when connected to websocket.

```go
if *isConnectedWebSocket {
    // tasks
}
```

### Memory Leak Prevention
Although there are goroutines spawned in order to perform reconnect gRPC stream connection tasks, the goroutines also must be killed when the websocket connection is demolished. When the websocket connection is closed, make sure that the following logs are shown:
```
reestablishStreamConnection goroutine killed
listenProcessOrderServerSide goroutine killed
```

To ensure goroutines are killed when websocket connection is closed, a context has been together with the websocket connection.
```go
websocketCtx, websocketCancel := context.WithCancel(context.Background())
defer websocketCancel()
```

Then once the function 'websocketCancel()' is called, the goroutines will be killed.
```go
for {
    select {
    case <-ctx.Done():
        return
    
    // remaining code ...
    }
}
```

Besides, the global boolean variable can also determine some tasks will be halted when the websocket connection is closed.

```go
if *isConnectedWebSocket {
    // tasks
}
```

<br/>

## Docker
Kindly change the variable _grpcServerHost_ in file _client/grpc_client/orderServiceClient.go_ to "host.docker.internal" before docker build.
