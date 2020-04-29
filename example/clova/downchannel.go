package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"fmt"

	"golang.org/x/net/http2"
	//"github.com/ninetyfivejae/avs"
)

var ch chan int

func  main() {
	// 정수형 채널을 생성한다
	ch = make(chan int)

	go func() {
		ch <- 123   //채널에 123을 보낸다
	}()
	i := <- ch  // 채널로부터 123을 받는다
	println(i)

	done := make(chan bool)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
		done <- true
	}()
	// 위의 Go루틴이 끝날 때까지 대기
	<-done
	fmt.Println("###")

	const clovaAccessToken = "iuV3BMFjQfGNS9Aqh4pjRg"
	//const clovaAccessToken = "GxTj4dO8S4mYdxZJUtE3_Q"

	client := http.Client{
		Transport: &http2.Transport{AllowHTTP: true},
	}

	if req, err := http.NewRequest("GET", "https://prod-ni-cic.clova.ai/v1/directives", nil); err != nil {
		panic(err)
	} else {
		req.Header.Add("User-Agent", "Mozilla/5=.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
		req.Header.Add("Authorization", "Bearer "+clovaAccessToken)

		if _, err := client.Do(req); err != nil {
			panic(err)
		}
	}

	if input, err := ioutil.ReadFile("input.json"); err != nil {
		panic(err)
	} else {
		if req, err := http.NewRequest("POST", "https://prod-ni-cic.clova.ai/v1/events", bytes.NewBuffer(input)); err != nil {
			panic(err)
		} else {
			req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
			req.Header.Add("Authorization", "Bearer "+clovaAccessToken)
			req.Header.Add("Content-Type", "multipart/form-data; boundary=Boundary-Text")

			if resp, err := client.Do(req); err != nil {
				panic(err)
			} else {
				defer resp.Body.Close()
				result, _ := ioutil.ReadAll(resp.Body)
				str := string(result)
				fmt.Println(str)
			}
		}
	}
}

//const ACCESS_TOKEN = "_CjtPUecQPKO6-_W_udQeQ"
//
//func main() {
//	directives, err := avs.CreateDownchannel(ACCESS_TOKEN)
//	if err != nil {
//		fmt.Printf("Failed to open downchannel: %v\n", err)
//		return
//	}
//	// Wait for directives to come in on the downchannel.
//	for directive := range directives {
//		switch d := directive.Typed().(type) {
//		case *avs.DeleteAlert:
//			fmt.Println("Unset alert:", d.Payload.Token)
//		case *avs.SetAlert:
//			fmt.Printf("Set alert %s (%s) for %s\n", d.Payload.Token, d.Payload.Type, d.Payload.ScheduledTime)
//		default:
//			fmt.Println("No code to handle directive:", d)
//		}
//	}
//	fmt.Println("Downchannel closed. Bye!")
//}