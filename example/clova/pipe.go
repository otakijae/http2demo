package main

import (
	"bytes"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	//r, w := io.Pipe()
	//
	//go func() {
	//	defer w.Close()
	//	json.NewEncoder(w).Encode(map[string]string{
	//		"Name": "Mohanson",
	//	})
	//}()
	//
	//resp, err := http.Post("https://httpbin.org/post", "application/json", r)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
	//	log.Fatal(err)
	//}

	reader, _ := io.Pipe()

	const clovaAccessToken = "F3g97g-IT6-UXpW-RIaHJQ"
	//const clovaAccessToken = "GxTj4dO8S4mYdxZJUtE3_Q"

	client := http.Client{
		Transport: &http2.Transport{AllowHTTP: true},
	}

	if req, err := http.NewRequest("GET", "https://prod-ni-cic.clova.ai/v1/directives", reader); err != nil {
		panic(err)
	} else {
		req.Header.Add("User-Agent", "Mozilla/5=.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
		req.Header.Add("Authorization", "Bearer "+clovaAccessToken)

		if _, err := client.Do(req); err != nil {
			panic(err)
		}
	}

	input, _ := ioutil.ReadFile("input.json")

	req, err := http.NewRequest("POST", "https://prod-ni-cic.clova.ai/v1/events", bytes.NewBuffer(input))
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	req.Header.Add("Authorization", "Bearer "+clovaAccessToken)
	req.Header.Add("Content-Type", "multipart/form-data; boundary=Boundary-Text")

	if resp, err := client.Do(req); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		//result, _ := ioutil.ReadAll(resp.Body)
		//str := string(result)
		//fmt.Println(str)
		if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
			log.Fatal(err)
		}
	}
}