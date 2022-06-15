package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

type APIResponse struct {
	Thread string `json:"thread_executed"`
}

func (a APIResponse) Get(url string) (string, error) {
	resp, err := http.Get("http://" + url)
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
	url := flag.String("url", "http://localhost:8000/api", "A url to ask")
	flag.Parse()
	api := APIResponse{}
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		thread, err := api.Get(*url)
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
	http.ListenAndServe(":8000", nil)
}
