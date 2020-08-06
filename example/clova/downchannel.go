package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fstab/h2c/http2client"

	"github.com/fstab/h2c/cli/util"
	"github.com/fstab/h2c/http2client/frames"
	"golang.org/x/net/http2"
)

//var ch chan int
var Stream *http.Response

func test1(client *http.Client, clovaAccessToken string, done chan bool) {
	if req, err := http.NewRequest("GET", "https://beta-cic.clova.ai/ping", nil); err != nil {
		panic(err)
	} else {
		req.Header.Add("Authorization", "Bearer "+clovaAccessToken)

		if resp, err := client.Do(req); err != nil {
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Body)
			panic(err)
		} else {
			fmt.Println(resp.StatusCode)
			defer resp.Body.Close()
		}
	}

	if req, err := http.NewRequest("GET", "https://beta-cic.clova.ai/v1/directives", nil); err != nil {
		panic(err)
	} else {
		req.Header.Add("Authorization", "Bearer "+clovaAccessToken)

		if resp, err := client.Do(req); err != nil {
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Body)
			panic(err)
		} else {
			fmt.Println(resp.StatusCode)
			Stream = resp
			defer resp.Body.Close()
		}
	}

	if input, err := ioutil.ReadFile("input2ApplicationRequest.json"); err != nil {
		panic(err)
	} else {
		if req, err := http.NewRequest("POST", "https://beta-cic.clova.ai/v1/events", bytes.NewBuffer(input)); err != nil {
			panic(err)
		} else {
			req.Header.Add("Authorization", "Bearer "+clovaAccessToken)
			req.Header.Add("Content-Type", "multipart/form-data; boundary=Boundary-Text")

			if resp, err := client.Do(req); err != nil {
				panic(err)
			} else {
				defer resp.Body.Close()
				result, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(result))
				fmt.Println(resp.StatusCode)
				//if resp.ContentLength == 0 {
				//	fmt.Println("###ASFG")
				//}
			}
		}
	}
	done <- true
}

func test2(client *http.Client, clovaAccessToken string, done chan bool) {
	if input, err := ioutil.ReadFile("input3.json"); err != nil {
		panic(err)
	} else {
		if req, err := http.NewRequest("POST", "https://beta-cic.clova.ai/v1/events", bytes.NewBuffer(input)); err != nil {
			panic(err)
		} else {
			req.Header.Add("Authorization", "Bearer "+clovaAccessToken)
			req.Header.Add("Content-Type", "multipart/form-data; boundary=Boundary-Text")

			if resp, err := client.Do(req); err != nil {
				panic(err)
			} else {
				defer resp.Body.Close()
				result, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(result))
			}
		}
	}
	done <- true
}

func ping(client *http.Client, clovaAccessToken string) {
	if req, err := http.NewRequest("GET", "https://beta-cic.clova.ai/ping", nil); err != nil {
		panic(err)
	} else {
		req.Header.Add("Authorization", "Bearer "+clovaAccessToken)

		if resp, err := http.DefaultClient.Do(req); err != nil {
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Body)
			panic(err)
		} else {
			fmt.Println(resp.StatusCode)
			defer resp.Body.Close()
		}
	}
}

func schedule(client *http.Client, clovaAccessToken string) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			ping(client, clovaAccessToken)
			select {
			case <-time.After(time.Second):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func makeFrameFilter(dumpFunction func(frames.Frame), frameTypesToBeDumped []frames.Type) func(frames.Frame) frames.Frame {
	return func(frame frames.Frame) frames.Frame {
		if util.SliceContainsFrameType(frameTypesToBeDumped, frame.Type()) {
			dumpFunction(frame)
		}
		return frame
	}
}

func testhandler(w http.ResponseWriter, r *http.Request) {
	const clovaAccessToken = "PthFbj2QQmWFbKpX1bB-Wg"

	client := http.Client{
		//Transport: &http2.Transport{
		//	AllowHTTP:       true,
		//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//},
		Transport: &http2.Transport{
			DialTLS:                    nil,
			TLSClientConfig:            &tls.Config{InsecureSkipVerify: true},
			ConnPool:                   nil,
			DisableCompression:         false,
			AllowHTTP:                  true,
			MaxHeaderListSize:          0,
			StrictMaxConcurrentStreams: false,
			ReadIdleTimeout:            0,
			PingTimeout:                60,
		},
		Timeout: 10 * time.Second,
	}

	ctx := r.Context()

	if req, err := http.NewRequest("GET", "https://beta-cic.clova.ai/ping", nil); err != nil {
		panic(err)
	} else {
		req.Header.Add("Authorization", "Bearer "+clovaAccessToken)
		req = req.WithContext(ctx)

		if resp, err := client.Do(req); err != nil {
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Body)
			panic(err)
		} else {
			fmt.Println(resp.StatusCode)
			defer resp.Body.Close()
		}
	}

	var h2c = http2client.New()
	result, err := h2c.Connect("https", "beta-cic.clova.ai", 443)
	if err != nil {
		panic(err)
	}
	fmt.Print(result)

	result, err = h2c.PingOnce()
	if err != nil {
		panic(err)
	}
	fmt.Print(result)

	result, err = h2c.PingRepeatedly(1 * time.Second)
	if err != nil {
		fmt.Println("Ping error")
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Print(result)

	done := make(chan bool)

	go func() {
		time.Sleep(30 * time.Second)
		h2c.StopPingRepeatedly()
		done <- true
	}()

	<-done
	fmt.Println(*h2c)

	if req, err := http.NewRequest("GET", "https://beta-cic.clova.ai/v1/directives", nil); err != nil {
		panic(err)
	} else {
		req.Header.Add("Authorization", "Bearer "+clovaAccessToken)
		req = req.WithContext(ctx)

		if resp, err := client.Do(req); err != nil {
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Body)
			panic(err)
		} else {
			fmt.Println(resp.StatusCode)
			Stream = resp
			//defer resp.Body.Close()
		}

		if transport, ok := client.Transport.(*http2.Transport); ok {
			conn, err := tls.Dial("tcp", "beta-cic.clova.ai:443", &tls.Config{InsecureSkipVerify: true})
			if err != nil {
				panic(err)
			}
			transport.PingTimeout = 1 * time.Second
			fmt.Println(transport.ConnPool)
			fmt.Println(conn)

			//cc, err := transport.ConnPool.GetClientConn(req, "beta-cic.clova.ai:443")
			cc, err := transport.NewClientConn(conn)
			if err != nil {
				panic(err)
			}
			ok := cc.CanTakeNewRequest()
			fmt.Println(ok)
			response, err := cc.RoundTrip(req)
			fmt.Println(response)
			//ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
			//defer cancel()
			done := make(chan bool)
			go func() {
				if err := cc.Ping(ctx); err != nil {
					panic(err)
				}
				done <- true
			}()
			<-done
		}
	}
}

func main() {
	//// 정수형 채널을 생성한다
	//ch = make(chan int)
	//
	//go func() {
	//	ch <- 123   //채널에 123을 보낸다
	//}()
	//i := <- ch  // 채널로부터 123을 받는다
	//println(i)
	//
	//done := make(chan bool)
	//go func() {
	//	for i := 0; i < 10; i++ {
	//		fmt.Println(i)
	//	}
	//	done <- true
	//}()
	//// 위의 Go루틴이 끝날 때까지 대기
	//<-done
	//fmt.Println("###")

	//{
	//    "access_token": "PthFbj2QQmWFbKpX1bB-Wg",
	//    "expires_in": 12960000,
	//    "refresh_token": "16u-EyXxQd2HFFJIiqTYnA",
	//    "token_type": "Bearer"
	//}

	const clovaAccessToken = "PthFbj2QQmWFbKpX1bB-Wg"

	client := http.Client{
		//Transport: &http2.Transport{
		//	AllowHTTP:       true,
		//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//},
		Transport: &http2.Transport{
			DialTLS:                    nil,
			TLSClientConfig:            &tls.Config{InsecureSkipVerify: true},
			ConnPool:                   nil,
			DisableCompression:         false,
			AllowHTTP:                  true,
			MaxHeaderListSize:          0,
			StrictMaxConcurrentStreams: false,
			ReadIdleTimeout:            0,
			PingTimeout:                60,
		},
		Timeout: 10 * time.Second,
	}

	http.HandleFunc("/test", testhandler)
	http.ListenAndServe(":8080", nil)

	//ping
	//stop := schedule(&client, clovaAccessToken)

	//event message
	done1 := make(chan bool)
	done2 := make(chan bool)
	go test1(&client, clovaAccessToken, done1)
	<-done1
	//time.Sleep(10 * time.Second)
	go test2(&client, clovaAccessToken, done2)
	<-done2

	//stop <- true
	//fmt.Println("Ping Done")
}

//if req, err := http.NewRequest("GET", "https://beta-cic.clova.ai/ping", nil); err != nil {
//	panic(err)
//} else {
//	req.Header.Add("Authorization", "Bearer "+clovaAccessToken)
//
//	if resp, err := client.Do(req); err != nil {
//		fmt.Println(resp.StatusCode)
//		fmt.Println(resp.Body)
//		panic(err)
//	} else {
//		fmt.Println(resp.StatusCode)
//		defer resp.Body.Close()
//	}
//}

//a := http2.ClientConn{}
//ctx := context.Background()
//if err := a.Ping(ctx); err != nil {
//	fmt.Println(err)
//}
//	Dial("tcp", "golang.org:http")

//conn, err := net.Dial("tcp", "beta-cic.clova.ai:443")
//if err != nil {
//	panic(err)
//}
//transport, ok := client.Transport.(*http2.Transport)
//if ok {
//	cc, err := transport.NewClientConn(conn)
//	if err != nil {
//		panic(err)
//	}
//	ctx := context.Background()
//
//	if err := cc.Ping(ctx); err != nil {
//		panic(err)
//	}
//}

//pinger, err := ping.NewPinger("beta-cic.clova.ai")
//if err != nil {
//	fmt.Printf("ERROR: %s\n", err.Error())
//	return
//}
//
//pinger.OnRecv = func(pkt *ping.Packet) {
//	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
//}
//pinger.OnFinish = func(stats *ping.Statistics) {
//	fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
//	fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
//	fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n", stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
//}
//fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
//pinger.Run()
