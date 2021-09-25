# http2demo
>  Source: https://posener.github.io/http2/
>
>  Watch out. Here are examples for using the **h2conn package**. This project here is just for test.

## certificate

HTTP/2 enforces TLS. In order to achieve this we first need a private key and a certificate.

```
openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
```

reference mac os letsencrypt certbot으로 인증서 발급

```
sudo certbot certonly --manual --email jh33836782@gmail.com -d 6ea4-183-101-7-107.ngrok.io:8443
//do the challenge specified
```

## conn.go

```go
package http2demo

import (
	"context"
	"io"
	"sync"
)

// Conn is client/server symmetric connection.
// It implements the io.Reader/io.Writer/io.Closer to read/write or close the connection to the other side.
// It also has a Send/Recv function to use channels to communicate with the other side.
type Conn struct {
	r  io.Reader
	wc io.WriteCloser

	cancel context.CancelFunc

	wLock sync.Mutex
	rLock sync.Mutex
}

func newConn(ctx context.Context, r io.Reader, wc io.WriteCloser) (*Conn, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Conn{
		r:      r,
		wc:     wc,
		cancel: cancel,
	}, ctx
}

// Write writes data to the connection
func (c *Conn) Write(data []byte) (int, error) {
	c.wLock.Lock()
	defer c.wLock.Unlock()
	return c.wc.Write(data)
}

// Read reads data from the connection
func (c *Conn) Read(data []byte) (int, error) {
	c.rLock.Lock()
	defer c.rLock.Unlock()
	return c.r.Read(data)
}

// Close closes the connection
func (c *Conn) Close() error {
	c.cancel()
	return c.wc.Close()
}
```

## server.go

```go
package http2demo

import (
	"fmt"
	"io"
	"net/http"
)

// ErrHTTP2NotSupported is returned by Accept if the client connection does not
// support HTTP2 connection.
// The server than can response to the client with an HTTP1.1 as he wishes.
var ErrHTTP2NotSupported = fmt.Errorf("HTTP2 not supported")

// Server can "accept" an http2 connection to obtain a read/write object
// for full duplex communication with a client.
type Server struct {
	StatusCode int
}

// Accept is used on a server http.Handler to extract a full-duplex communication object with the client.
// See h2conn.Accept documentation for more info.
func (u *Server) Accept(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	if !r.ProtoAtLeast(2, 0) {
		return nil, ErrHTTP2NotSupported
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, ErrHTTP2NotSupported
	}

	c, ctx := newConn(r.Context(), r.Body, &flushWrite{w: w, f: flusher})

	// Update the request context with the connection context.
	// If the connection is closed by the server, it will also notify everything that waits on the request context.
	*r = *r.WithContext(ctx)

	w.WriteHeader(u.StatusCode)
	flusher.Flush()

	return c, nil
}

var defaultUpgrader = Server{
	StatusCode: http.StatusOK,
}

// Accept is used on a server http.Handler to extract a full-duplex communication object with the client.
// The server connection will be closed when the http handler function will return.
// If the client does not support HTTP2, an ErrHTTP2NotSupported is returned.
//
// Usage:
//
//      func (w http.ResponseWriter, r *http.Request) {
//          conn, err := h2conn.Accept(w, r)
//          if err != nil {
//		        log.Printf("Failed creating http2 connection: %s", err)
//		        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		        return
//	        }
//          // use conn
//      }
func Accept(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	return defaultUpgrader.Accept(w, r)
}

type flushWrite struct {
	w io.Writer
	f http.Flusher
}

func (w *flushWrite) Write(data []byte) (int, error) {
	n, err := w.w.Write(data)
	w.f.Flush()
	return n, err
}

func (w *flushWrite) Close() error {
	// Currently server side close of connection is not supported in Go.
	// The server closes the connection when the http.Handler function returns.
	// We use connection context and cancel function as a work-around.
	return nil
}
```

## client.go

```go
package http2demo

import (
	"context"
	"io"
	"net/http"

	"golang.org/x/net/http2"
)

// Client provides HTTP2 client side connection with special arguments
type Client struct {
	// Method sets the HTTP method for the dial
	// The default method, if not set, is HTTP POST.
	Method string
	// Header enables sending custom headers to the server
	Header http.Header
	// Client is a custom HTTP client to be used for the connection.
	// The client must have an http2.Transport as it's transport.
	Client *http.Client
}

// Connect establishes a full duplex communication with an HTTP2 server with custom client.
// See h2conn.Connect documentation for more info.
func (c *Client) Connect(ctx context.Context, urlStr string) (*Conn, *http.Response, error) {
	reader, writer := io.Pipe()

	// Create a request object to send to the server
	req, err := http.NewRequest(c.Method, urlStr, reader)
	if err != nil {
		return nil, nil, err
	}

	// Apply custom headers
	if c.Header != nil {
		req.Header = c.Header
	}

	// Apply given context to the sent request
	req = req.WithContext(ctx)

	// If an http client was not defined, use the default http client
	httpClient := c.Client
	if httpClient == nil {
		httpClient = defaultClient.Client
	}

	// Perform the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	// Create a connection
	conn, ctx := newConn(req.Context(), resp.Body, writer)

	// Apply the connection context on the request context
	resp.Request = req.WithContext(ctx)

	return conn, resp, nil
}

var defaultClient = Client{
	Method: http.MethodPost,
	Client: &http.Client{Transport: &http2.Transport{}},
}

// Connect establishes a full duplex communication with an HTTP2 server.
//
// Usage:
//
//      conn, resp, err := h2conn.Connect(ctx, url)
//      if err != nil {
//          log.Fatalf("Initiate client: %s", err)
//      }
//      if resp.StatusCode != http.StatusOK {
//          log.Fatalf("Bad status code: %d", resp.StatusCode)
//      }
//      defer conn.Close()
//
//      // use conn
//
func Connect(ctx context.Context, urlStr string) (*Conn, *http.Response, error) {
	return defaultClient.Connect(ctx, urlStr)
}
```