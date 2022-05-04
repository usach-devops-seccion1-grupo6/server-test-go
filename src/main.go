package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
	"time"
    "runtime"
	"math/rand"
)

func cpuIntensive(p *int64) {
    for i := int64(1); i <= 10000000; i++ {
      *p = i
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
	rand.Seed(time.Now().UnixNano())
    n := rand.Intn(5)
    environment := os.Getenv("ENVIRONMENT")
    if (len(environment) >= 0 && environment == "SMOKE") {
        time.Sleep(time.Duration(n)*time.Second)
    }

	duration := time.Since(started)

    urlPath := r.URL.Path
    log.Printf("Web request received on url path %s (%v)", urlPath, duration.Seconds())

	msg := ""
    switch {
        case urlPath == "/health":
            if duration.Seconds() >= 4 {
                w.WriteHeader(500)
                msg += fmt.Sprintf("Error: %v\n", duration.Seconds())
            } else {
                msg += fmt.Sprintf("Ok: %v\n", duration.Seconds())
            }
        case urlPath == "/cpu":
            // for stress cpu
            runtime.GOMAXPROCS(2)
            for true {
                x := int64(0)
                go cpuIntensive(&x)
            }
        case urlPath == "/exit":
            os.Exit(1)
        default:
            msg += fmt.Sprintf("Hello world (%s) (%v)\n", urlPath, duration.Seconds())
    }
	
    _, err := w.Write([]byte(msg))
    if err != nil {
        fmt.Printf("Failed to write response, err: %s", err)
    }
}

func main() {
    StartWebServer()
}