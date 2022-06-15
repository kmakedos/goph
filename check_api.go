package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type APIResponse struct {
	Thread string `json:"thread_executed"`
}

func (a APIResponse) Get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return "", err
	}
	return a.Thread, nil
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{ API Tester }"))
}

func main() {
	url, ok := os.LookupEnv("API_URL")
	if !ok {
		log.Println("INFO: API_URL is not set")
		url = "http://localhost:8080"
	}
	log.Println("INFO: API URL is: ", url)
	port := flag.String("p", "8000", "Port of this server to listen to")
	flag.Parse()
	api := APIResponse{}
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		thread, err := api.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"Thread": thread,
				"took":   time.Since(begin).String(),
			})
		}
	})
	http.HandleFunc("/health", health)
	log.Println("Starting server and listening to port 8000...")
	l, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalln("Could not listen: ", err)
	}
	srv := http.Server{}
	as := sync.WaitGroup{}
	go func() {
		if err := srv.Serve(l); err != nil {
			log.Fatalln("Could not start server: ", err)
		}
		as.Done()
	}()
	as.Add(1)
	as.Wait()
}
