//package main
//
//import (
//	"bufio"
//	"context"
//	"crypto/tls"
//	"encoding/json"
//	"log"
//	"net/http"
//	"os"
//	"strings"
//	"fmt"
//
//	"github.com/ninetyfivejae/http2demo"
//	"golang.org/x/net/http2"
//)
//
//const url = "https://prod-ni-cic.clova.ai/v1/directives"

//func main() {
//	client := http.Client{
//		Transport: &http2.Transport{
//			AllowHTTP: true,
//			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
//				return net.Dial(network, addr)
//			},
//		},
//	}
//
//	resp, err := client.Get(url)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Client Proto: %d\n", resp.ProtoMajor)
//
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		panic(err)
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
//	defer cancel()
//
//	tr := &http2.Transport{
//		AllowHTTP: true,
//		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
//			return net.Dial(network, addr)
//		},
//	}
//
//	req.WithContext(ctx)
//	resp2, err := tr.RoundTrip(req)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Printf("RoundTrip Proto: %d\n", resp2.ProtoMajor)
//}

//func main() {
//	c := http.Client{}
//
//	// According to the CloudFront documentation for a request behavior, if the
//	// request is GET and includes a body, it returns a 403 Forbidden. See the
//	// documentation here:
//	// https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/RequestAndResponseBehaviorCustomOrigin.html#RequestCustom-get-body
//
//	// var body bytes.Buffer
//	r, err := http.NewRequest("GET", url, http.NoBody)
//	if err != nil {
//		log.Fatalf("error creating the request: %s", err)
//	}
//
//	//r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
//	//r.Header.Add("Authorization", "Bearer AAAAu1xf2KwhK3A1JVbk+sXdp857iriGgf2/8xqlcU68XtEnjnPadfgjRKIjCTaA0gFt/+EW7PCf4irJ7E3SS0akPjE6bkNiE3gX1veldEJBh7OBP")
//
//	res, err := c.Do(r)
//	if err != nil {
//		log.Fatalf("error doing the request: %s", err)
//	}
//	io.Copy(ioutil.Discard, res.Body)
//	res.Body.Close()
//
//	log.Printf("response status for an HTTP/2 request: %s", res.Status)
//
//	// doing the same request without HTTP/2 does work
//	c.Transport = &http.Transport{
//		TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
//	}
//	r, err = http.NewRequest("GET", url, http.NoBody)
//	if err != nil {
//		log.Fatalf("error creating the request: %s", err)
//	}
//
//	res, err = c.Do(r)
//	if err != nil {
//		log.Fatalf("error doing the request: %s", err)
//	}
//	io.Copy(ioutil.Discard, res.Body)
//	res.Body.Close()
//
//	log.Printf("response status for an HTTP/1 request: %s", res.Status)
//}
