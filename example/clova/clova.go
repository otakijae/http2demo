package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"golang.org/x/net/http2"
)

const url = "https://prod-ni-cic.clova.ai/v1/directives"

var httpVersion = flag.Int("version", 2, "HTTP version")

func main() {
	flag.Parse()
	client := &http.Client{}

	// Create a pool with the server certificate since it is not signed by a known CA
	caCert, err := ioutil.ReadFile("server.crt")
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	//Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	// Use the proper transport in the client
	switch *httpVersion {
	case 1:
		client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	case 2:
		client.Transport = &http2.Transport{TLSClientConfig: tlsConfig}
	}

	// Request 객체 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	//필요시 헤더 추가 가능
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	req.Header.Add("Authorization", "Bearer AAAAu1xf2KwhK3A1JVbk+sXdp857iriGgf2/8xqlcU68XtEnjnPadfgjRKIjCTaA0gFt/+EW7PCf4irJ7E3SS0akPjE6bkNiE3gX1veldEJBh7OBP0kQ02HFNcvm5/+8WXjBW8Fn1Yw3tJL8+/We6TUZewATPNzght8Z7m2NEnS26Mc14l3nA93FmDJTy8H4Hxd84pTaK/5yghahHADLNw/qCPvsNB6+WTwFNUpuYQYAwAtbCn8x4YjXLankk4Qj/cVE+A==")

	// Client객체에서 Request 실행
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	fmt.Println(str)
}

var publicKey *rsa.PublicKey

const publicKeyDownloadURL = "https://clova.ai/.well-known/signature-public-key.pem"

func downloadPublicKey() bool {
	tokens := strings.Split(publicKeyDownloadURL, "://")
	if tokens[0] != "https" {
		return false
	}

	response, err := http.Get(publicKeyDownloadURL)
	if err != nil {
		log.Println("Error during downloading", publicKeyDownloadURL, "-", err)
		return false
	}
	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error during reading")
		return false
	}

	block, _ := pem.Decode([]byte(read))
	downloadedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey = downloadedKey.(*rsa.PublicKey)
	log.Println("Download public key complete")
	return true
}

func CheckSignature(r *http.Request, body []byte) bool {
	signatureStr := r.Header.Get("SignatureCEK")

	if publicKey == nil && !downloadPublicKey() {
		return false
	}

	hash := crypto.SHA256.New()
	hash.Write(body)
	hashData := hash.Sum(nil)
	signature, _ := base64.StdEncoding.DecodeString(signatureStr)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashData, signature)
	if err != nil {
		return false
	}
	return true
}

func performRequest(r http.Handler, method, path string, reqBody io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, reqBody)
	responseRecorder := httptest.NewRecorder()
	r.ServeHTTP(responseRecorder, req)
	return responseRecorder
}
