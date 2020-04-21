package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/http2"
	//"github.com/ninetyfivejae/avs"
)

const url = "https://prod-ni-cic.clova.ai/v1/directives"

func  main()  {
	client := http.Client{
		// InsecureTLSDial is temporary and will likely be
		// replaced by a different API later.
		Transport: &http2.Transport{AllowHTTP: true,},
	}

	if request, err := ioutil.ReadFile("input.json"); err != nil {
		panic(err)
	} else {
		if req, err := http.NewRequest("POST", "https://prod-ni-cic.clova.ai/v1/events", bytes.NewBuffer(request)); err != nil {
			panic(err)
		} else {
			req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
			req.Header.Add("Authorization", "Bearer _CjtPUecQPKO6-_W_udQeQ")
			req.Header.Add("Content-Type", "multipart/form-data; boundary=Boundary-Text")
			req.Header.Add("Content-Disposition", "form-data; name=\"metadata\"")

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
//	// Record your request into request.wav.
//	//audio, _ := os.Open("./request.wav")
//	request := avs.NewRequest(ACCESS_TOKEN)
//	response, err := avs.DefaultClient.Do(request)
//	//response, err := avs.PostRecognize(ACCESS_TOKEN, "b120c3e0-e6b9-4a3d-96de-71539e5f6214", "bc834682-6d22-4bbb-8352-4a49df2ed3d7", audio)
//	if err != nil {
//		fmt.Printf("Failed to call AVS: %v\n", err)
//		return
//	}
//	// AVS might not return any directives in some cases.
//	if len(response.Directives) == 0 {
//		fmt.Println("Alexa had nothing to say.")
//		return
//	}
//	// A response can have multiple directives in the response.
//	for _, directive := range response.Directives {
//		switch d := directive.Typed().(type) {
//		case *avs.ExpectSpeech:
//			fmt.Printf("Alexa wants you to speak within %s!\n", d.Timeout())
//		case *avs.Play:
//			// The Play directive can point to attached audio or remote streams.
//			if cid := d.Payload.AudioItem.Stream.ContentId(); cid != "" {
//				//save(response, cid)
//			} else {
//				fmt.Println("Remote stream:", d.Payload.AudioItem.Stream.URL)
//			}
//		case *avs.Speak:
//			// The Speak directive always points to attached audio.
//			//save(response, d.ContentId())
//		default:
//			fmt.Println("No code to handle directive:", d)
//		}
//	}
//}