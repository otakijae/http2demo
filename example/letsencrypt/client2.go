package main

import (
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
)

func main() {
    log.SetFlags(log.Lshortfile)

    transport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: transport}
    response, err := client.Get("localhost:5050")
    if err != nil {
        fmt.Println(err)
    }
    defer response.Body.Close()

    content, _ := ioutil.ReadAll(response.Body)
    s := strings.TrimSpace(string(content))

    fmt.Println(s)
}
