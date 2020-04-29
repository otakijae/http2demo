package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const clovaAccessToken = "iuV3BMFjQfGNS9Aqh4pjRg"

func sendEventMessage(conn *Conn, ctx context.Context, urlStr string) {
	d := defaultClient

	input, err := ioutil.ReadFile("input.json")
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(input))
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	req.Header.Add("Authorization", "Bearer "+clovaAccessToken)
	req.Header.Add("Content-Type", "multipart/form-data; boundary=Boundary-Text")

	req = req.WithContext(ctx)

	if resp, err := d.Client.Do(req); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		result, _ := ioutil.ReadAll(resp.Body)
		//str := string(result)
		//fmt.Println(str)

		conn.Write(result)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := defaultClient
	d.Client.Transport = &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, AllowHTTP: true}

	d.Method = http.MethodGet
	d.Header = http.Header{}
	d.Header.Add("User-Agent", "Mozilla/5=.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	d.Header.Add("Authorization", "Bearer "+clovaAccessToken)

	conn, resp, err := d.Connect(ctx, "https://prod-ni-cic.clova.ai/v1/directives")
	if err != nil {
		log.Fatalf("Initiate conn: %s", err)
	}
	defer conn.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad status code: %d", resp.StatusCode)
	}

	var (
		//stdin = bufio.NewReader(os.Stdin)

		in  = json.NewDecoder(conn)
		//out = json.NewEncoder(conn)
	)

	defer log.Println("Exited")

	fmt.Println("Echo session starts, press ctrl-C to terminate.")
	for ctx.Err() == nil {
		//fmt.Print("Send: ")
		//msg, err := stdin.ReadString('\n')
		//if err != nil {
		//	log.Fatalf("Failed reading stdin: %v", err)
		//}
		//msg = strings.TrimRight(msg, "\n")

		//err = out.Encode(msg)
		//if err != nil {
		//	log.Fatalf("Failed sending message: %v", err)
		//}

		go func() {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second * 5)
				sendEventMessage(conn, ctx, "https://prod-ni-cic.clova.ai/v1/events")
			}
		}()

		fmt.Println("###")
		fmt.Println(ctx.Err())

		//body := &bytes.Buffer{}
		//if _, err := body.ReadFrom(resp.Body); err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Println(resp.StatusCode)
		//fmt.Println(resp.Header)
		//fmt.Println(body)

		//var response map[string]string
		//err = in.Decode(&response)
		//if err != nil {
		//	log.Fatalf("Failed receiving message: %v", err)
		//}
		//fmt.Printf("Got response %q\n", response)

		//if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		//	log.Fatal(err)
		//}

		//a, _ := ioutil.ReadAll(conn.r)
		//str := string(a)
		//fmt.Println(str)

		var result map[string]interface{}
		in.Decode(&result)
		log.Println(result)
	}
}

type Client struct {
	Method string
	Header http.Header
	Client *http.Client
}

func (c *Client) Connect(ctx context.Context, urlStr string) (*Conn, *http.Response, error) {
	reader, writer := io.Pipe()

	req, err := http.NewRequest(c.Method, urlStr, reader)
	if err != nil {
		return nil, nil, err
	}

	if c.Header != nil {
		req.Header = c.Header
	}

	req = req.WithContext(ctx)

	httpClient := c.Client
	if httpClient == nil {
		httpClient = defaultClient.Client
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	conn, ctx := newConn(req.Context(), resp.Body, writer)
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

func (c *Conn) Write(data []byte) (int, error) {
	c.wLock.Lock()
	defer c.wLock.Unlock()
	return c.wc.Write(data)
}

func (c *Conn) Read(data []byte) (int, error) {
	c.rLock.Lock()
	defer c.rLock.Unlock()
	return c.r.Read(data)
}

func (c *Conn) Close() error {
	c.cancel()
	return c.wc.Close()
}
