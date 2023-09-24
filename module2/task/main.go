package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/health", HealthCheck())
	serverChan := make(chan error, 1)
	serverDone := make(chan struct{}, 1)
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case err := <-serverChan:
				if err != nil {
					panic("server error: " + err.Error())
				}
				serverDone <- struct{}{}
			case <-ticker.C:
				fmt.Println("server is running")
			}
		}
	}()
	err := http.ListenAndServe("0.0.0.0:8081", nil)
	serverChan <- err
	<-serverDone
}

func HealthCheck() func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		// 读取 req header 信息写入 resp header
		for k, v := range req.Header {
			if len(v) > 0 {
				resp.Header().Add(k, v[0])
			}
		}
		fmt.Printf("req remote-addr: %s\n", req.RemoteAddr)
		fmt.Printf("server version: %s\n", os.Getenv("TEST_VERSION"))
		resp.Write([]byte("success"))
	}
}
