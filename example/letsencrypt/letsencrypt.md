### certificates

```
openssl req -newkey rsa:2048 -nodes -keyout cert.key -x509 -days 365 -out cert.crt
```

### server

```go
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// generate a `Certificate` struct
	cert, _ := tls.LoadX509KeyPair("cert.crt", "cert.key")

	// create a custom server with `TLSConfig`
	s := &http.Server{
		Addr: ":5050",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello World!")
	})

	log.Fatal(s.ListenAndServeTLS("cert.crt", "cert.key"))
}
```

### request with public key

```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	caCert, err := ioutil.ReadFile("cert.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	response, err := client.Get("https://localhost:5050")
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	s := strings.TrimSpace(string(content))

	fmt.Println(s)
}
```

### request without ssl verification for test

```go
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
```

---

```
httpstat -k https://..navercorp.com:5050
```

```
Connected to ....:5050

HTTP/2.0 200 OK
Content-Length: 12
Content-Type: text/plain; charset=utf-8
Date: Sun, 11 Jul 2021 11:45:14 GMT

Body discarded

  DNS Lookup   TCP Connection   TLS Handshake   Server Processing   Content Transfer
[      1ms  |          13ms  |         10ms  |              8ms  |             0ms  ]
            |                |               |                   |                  |
   namelookup:1ms            |               |                   |                  |
                       connect:14ms          |                   |                  |
                                   pretransfer:25ms              |                  |
                                                     starttransfer:34ms             |
                                                                                total:34ms
```

