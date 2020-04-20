package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"fmt"

	"golang.org/x/net/http2"
	//"github.com/ninetyfivejae/avs"
)

const url = "https://prod-ni-cic.clova.ai/v1/directives"

func  main() {
	client := http.Client{
		// InsecureTLSDial is temporary and will likely be
		// replaced by a different API later.
		Transport: &http2.Transport{AllowHTTP: true,},
	}

	// Request 객체 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	//필요시 헤더 추가 가능
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	req.Header.Add("Authorization", "Bearer _CjtPUecQPKO6-_W_udQeQ")

	// Client객체에서 Request 실행
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	result, _ := ioutil.ReadAll(resp.Body)
	str := string(result) //바이트를 문자열로
	fmt.Println(str)

	if request, err := ioutil.ReadFile("input.json"); err != nil {
		panic(err)
	} else {
		if req, err := http.NewRequest("POST", "https://prod-ni-cic.clova.ai/v1/events", bytes.NewBuffer(request)); err != nil {
			panic(err)
		} else {
			req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
			req.Header.Add("Authorization", "Bearer _CjtPUecQPKO6-_W_udQeQ")
			req.Header.Add("Content-Type", "application/json; charset=UTF-8")
			req.Header.Add("Content-Disposition", "form-data; name=metadata")

			// Client객체에서 Request 실행
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			// 결과 출력
			result, _ := ioutil.ReadAll(resp.Body)
			str := string(result) //바이트를 문자열로
			fmt.Println(str)
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