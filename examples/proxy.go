package main

import (
        "../"
        "fmt"
        "log"
        "net/http"
        "time"
)

func main() {
        // 8484ポートでdeltaがLISTEN
        server := delta.NewServer("0.0.0.0", 8484)
        //既存サーバー?
        server.AddMasterBackend("production", "127.0.0.1", 8080)
        //新サーバー?
        server.AddBackend("testing", "127.0.0.1", 8081)

        server.OnSelectBackend(func(req *http.Request) []string {
                // GETリクエストならproductionとtestingの両方にリクエストを送るが
                // それ以外はproductionのみに送る
                if req.Method == "GET" {
                        return []string{"production", "testing"}
                } else {
                        return []string{"production"}
                }
        })

        server.OnMungeHeader(func(backend string, header *http.Header) {
                if backend == "testing" {
                        header.Add("X-Delta-Sandbox", "1")
                }
        })

        server.OnBackendFinished(func(responses map[string]*delta.Response) {
                for backend, response := range responses {
                        log.Printf("%s [%d ms]: ", backend, (response.Elapsed / time.Millisecond))
                        for k, v := range response.HttpResponse.Header {
                                //log.Printf("%s:  %+V\n", k, v)
                                fmt.Printf("%s: %s\n", k, v[0])
                        }
                        log.Printf("\n%s", response.Data)
                }
        })

        server.Run()
}
