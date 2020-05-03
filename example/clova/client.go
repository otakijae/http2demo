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
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const clovaAccessToken = "4PRu4SqcRBmHleOBK3x2AQ"

func sendEventMessage(conn *Conn, ctx context.Context, urlStr string) {
	d := defaultClient

	input, err := ioutil.ReadFile("input.json")
	if err != nil {
		panic(err)
	}

	conn.Write(input)
	conn.r = bytes.NewBuffer(input)

	//req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(input))
	req, err := http.NewRequest("POST", urlStr, conn.r)
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
		result, _ := ioutil.ReadAll(resp.Body)
		//str := string(result)
		//fmt.Println(str)

		conn.r = resp.Body
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

	//in  = json.NewDecoder(conn)
	//out = json.NewEncoder(conn)
	)

	defer log.Println("Exited")

	d.Method = http.MethodPost
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second * 3)
			sendEventMessage(conn, ctx, "https://prod-ni-cic.clova.ai/v1/events")
		}
	}()

	fmt.Println("Echo session starts, press ctrl-C to terminate.")
	for ctx.Err() == nil {

		//err = out.Encode(msg)
		//if err != nil {
		//	log.Fatalf("Failed sending message: %v", err)
		//}

		if _, err := io.Copy(os.Stdout, conn.r); err != nil {
			log.Fatal(err)
		}
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


type Response struct {
	RequestId string
	Directives []*Message
	Content map[string][]byte
}

type Message struct {
	Header  map[string]string `json:"header"`
	Payload json.RawMessage   `json:"payload,omitempty"`
}

func newMultipartReaderFromResponse(resp *http.Response) (*multipart.Reader, error) {
	contentType := strings.Replace(resp.Header.Get("Content-Type"), "type=application/json", `type="application/json"`, 1)
	mediatype, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(mediatype, "multipart/") {
		return nil, fmt.Errorf("unexpected content type %s", mediatype)
	}
	return multipart.NewReader(resp.Body, params["boundary"]), nil
}

type responsePart struct {
	Directive *Message
}