# broadcaster

[![Build Status](https://travis-ci.org/rubenv/broadcaster.svg?branch=master)](https://travis-ci.org/rubenv/broadcaster) [![GoDoc](https://godoc.org/github.com/rubenv/broadcaster?status.png)](https://godoc.org/github.com/rubenv/broadcaster)

Package broadcaster implements a websocket server for broadcasting Redis pub/sub
messages to web clients.

A JavaScript client can be found here:
https://github.com/rubenv/broadcaster-client

Originally based on https://github.com/rubenv/node-broadcast-hub but
significantly improved while moving to Go.

## Installation
```
go get github.com/rubenv/broadcaster
```

Import into your application with:

```go
import "github.com/rubenv/broadcaster"
```

## Usage

```go
const (
	// Client: start authentication
	AuthMessage = "auth"

	// Server: Authentication succeeded
	AuthOKMessage = "authOk"

	// Server: Authentication failed
	AuthFailedMessage = "authError"

	// Client: Subscribe to channel
	SubscribeMessage = "subscribe"

	// Server: Subscribe succeeded
	SubscribeOKMessage = "subscribeOk"

	// Server: Subscribe failed
	SubscribeErrorMessage = "subscribeError"

	// Server: Broadcast message
	MessageMessage = "message"

	// Client: Unsubscribe from channel
	UnsubscribeMessage = "unsubscribe"

	// Server: Unsubscribe succeeded
	UnsubscribeOKMessage = "unsubscribeOk"

	// Server: Unsubscribe failed
	UnsubscribeErrorMessage = "unsubscribeError"

	// Client: Send me more messages
	PollMessage = "poll"

	// Client: I'm still alive
	PingMessage = "ping"

	// Server: Unknown message
	UnknownMessage = "unknown"

	// Server: Server error
	ServerErrorMessage = "serverError"
)
```
Message types used between server and client.

#### type Client

```go
type Client struct {
	Mode ClientMode

	// Data passed when authenticating
	AuthData map[string]interface{}

	// Set when disconnecting
	Error error

	// Incoming messages
	Messages chan clientMessage

	// Receives true when disconnected
	Disconnected chan bool

	// Timeout
	Timeout time.Duration

	// Reconnection attempts
	MaxAttempts int
}
```


#### func  NewClient

```go
func NewClient(urlStr string) (*Client, error)
```

#### func (*Client) Connect

```go
func (c *Client) Connect() error
```

#### func (*Client) Disconnect

```go
func (c *Client) Disconnect() error
```

#### func (*Client) Subscribe

```go
func (c *Client) Subscribe(channel string) error
```

#### func (*Client) Unsubscribe

```go
func (c *Client) Unsubscribe(channel string) error
```

#### type ClientMode

```go
type ClientMode int
```

Client connection mode.

```go
const (
	ClientModeAuto      ClientMode = 0
	ClientModeWebsocket ClientMode = 1
	ClientModeLongPoll  ClientMode = 2
)
```
Connection modes, can be used to force a specific connection type.

#### type Server

```go
type Server struct {
	// Invoked upon initial connection, can be used to enforce access control.
	CanConnect func(data map[string]interface{}) bool

	// Invoked upon channel subscription, can be used to enforce access control
	// for channels.
	CanSubscribe func(data map[string]interface{}, channel string) bool

	// Can be set to allow CORS requests.
	CheckOrigin func(r *http.Request) bool

	// Can be used to configure buffer sizes etc.
	// See http://godoc.org/github.com/gorilla/websocket#Upgrader
	Upgrader websocket.Upgrader

	// Redis host, used for data, defaults to localhost:6379
	RedisHost string

	// Redis pubsub channel, used for internal coordination messages
	// Defaults to "broadcaster"
	ControlChannel string

	// Namespace for storing session data.
	// Defaults to "bc:"
	ControlNamespace string

	// PubSub host, used for pubsub, defaults to RedisHost
	PubSubHost string

	// Timeout for long-polling connections
	Timeout time.Duration

	// Combine long poll message for given duration (more latency, less load)
	PollTime time.Duration
}
```

A Server is the main class of this package, pass it to http.Handle on a chosen
path to start a broadcast server.

#### func (*Server) Prepare

```go
func (s *Server) Prepare() error
```

#### func (*Server) ServeHTTP

```go
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
Main HTTP server.

#### func (*Server) Stats

```go
func (s *Server) Stats() (Stats, error)
```

#### type Stats

```go
type Stats struct {
	// Number of active connections
	Connections int

	// For debugging purposes only
	LocalSubscriptions map[string]int
}
```

## License

    (The MIT License)

    Copyright (C) 2013-2016 by Ruben Vermeersch <ruben@rocketeer.be>

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
    THE SOFTWARE.
