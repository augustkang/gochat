# gochat

Simple TCP chat application

## How it works(flow)

1. Run server `go run server.go`
2. Run client1, 2 `go run client.go`
3. (client1,2) : type username
4. (client1,2) : join room
5. (client1,2) : send message

## Limitations
- Better printing
- Duplex but not realtime in client terminal(have to type something to receive message from other users)
- Need more command(features)
