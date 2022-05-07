package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

func cpuIntensive(p *int64) {
	for i := int64(1); i <= 10000000; i++ {
		*p = i
	}
}

func Fib2(n int) uint64 {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return Fib2(n-1) + Fib2(n-2)
	}
}

const homepageEndPoint = "/"

// StartWebServer the webserver
// code from https://ashishb.net/tech/docker-101-a-basic-web-server-displaying-hello-world/
func StartWebServer() {
	http.HandleFunc(homepageEndPoint, handleHomepage)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		panic("Environment variable PORT is not set")
	}

	log.Printf("Starting web server to listen on endpoints [%s] and port %s",
		homepageEndPoint, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func handleHomepage(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	environment := os.Getenv("ENVIRONMENT")
	rand.Seed(time.Now().UnixNano())
	myrand := rand.Intn(5)
	if len(environment) >= 0 && environment == "SMOKE" {
		time.Sleep(time.Duration(myrand) * time.Second)
	}

	duration := time.Since(started)
	urlPath := r.URL.Path

	log.Printf("Web request received on url path %s", urlPath)

	msg := ""
	tmp := ""
	switch {
	case urlPath == "/health":
		if duration.Seconds() >= 4 {
			w.WriteHeader(500)
			msg += fmt.Sprintf("Error: %v\n", duration.Seconds())
		} else {
			msg += fmt.Sprintf("Ok: %v\n", duration.Seconds())
		}
	case urlPath == "/cpu":
		tmp = r.URL.Query().Get("count")
		count := 10
		if len(tmp) > 0 {
			count, _ = strconv.Atoi(r.URL.Query().Get("count"))
		}

		tmp = r.URL.Query().Get("procs")
		procs := 8
		if len(tmp) > 0 {
			procs, _ = strconv.Atoi(r.URL.Query().Get("procs"))
		}

		runtime.GOMAXPROCS(procs)
		//for i := int64(1); i <= int64(count); i++ {
		//x := int64(0)
		//go cpuIntensive(&x)
		r := uint64(0)
		for i := myrand; i <= count; i++ {
			r = Fib2(i)
		}

		msg += fmt.Sprintf("Cpu: %d (%v)\n", r, duration.Seconds())
	case urlPath == "/mem":
		tmp = r.URL.Query().Get("iter")
		iter := 1000000
		if len(tmp) > 0 {
			iter, _ = strconv.Atoi(r.URL.Query().Get("iter"))
		}
		var s []int64
		for j := int64(myrand); j <= int64(iter); j++ {
			s = append(s, j)
		}

		msg += fmt.Sprintf("Mem: (%v)\n", duration.Seconds())
	case urlPath == "/exit":
		msg += fmt.Sprintf("Exit (%s) (%v)\n", urlPath, duration.Seconds())
		w.Write([]byte(msg))
		os.Exit(1)
	default:
		msg += fmt.Sprintf("Hi (%s) (%v)\n", urlPath, duration.Seconds())
	}

	_, err := w.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Failed to write response, err: %s", err)
	}
}

func main() {
	StartWebServer()
}
