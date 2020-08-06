package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

const clovaAccessToken = "peSGYmiDQ4SzuyXryr0uag"

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

		conn.wc.Write(result)
		//conn, ctx := newConn(req.Context(), resp.Body, writer)
		resp.Request = req.WithContext(ctx)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go catchSignal(cancel)

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
		stdin = bufio.NewReader(os.Stdin)

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
		fmt.Print("Input: ")
		msg, err := stdin.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed reading stdin: %v", err)
		}
		msg = strings.TrimRight(msg, "\n")

		//err = out.Encode(msg)
		//if err != nil {
		//	log.Fatalf("Failed sending message: %v", err)
		//}

		go func() {
			//for i := 0; i < 5; i++ {
			//	time.Sleep(time.Second * 5)
			//}
			sendEventMessage(conn, ctx, "https://prod-ni-cic.clova.ai/v1/events")
		}()

		fmt.Println()

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

		n, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("done Reading the server response's body n=(%v) err=%v\nres.headers: %v\n", n, err, resp.Header)

		//a, _ := ioutil.ReadAll(conn.r)
		//str := string(a)
		//fmt.Println(str)

		//var result map[string]interface{}
		//in.Decode(&result)
		//log.Println(result)
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

func catchSignal(cancel context.CancelFunc) {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("Cancelling due to interrupt")
	cancel()
}
