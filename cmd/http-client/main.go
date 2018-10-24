package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type request struct {
	Hostname string `json: hostname`
}

func main() {

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "50000"
	}

	url := fmt.Sprintf("http://%s:%s", host, port)
	log.Printf("Client connects to server at: %s\n", url)

	hostname, _ := os.Hostname()
	b, err := json.Marshal(request{Hostname: hostname})
	if err != nil {
		log.Printf("json failure: %v\n", err)
	}
	for {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{} //{Timeout: time.Second * 10}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failure while Get(): %v\n", err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error reading the body", err)
			}

			bodyRsp := request{}
			err = json.Unmarshal(body, &bodyRsp)
			if err != nil {
				log.Fatal("Decoding error: ", err)
			}

			name := bodyRsp.Hostname
			log.Printf("Response from server '%s' received\n", name)
		}
		req.Close = true
		client = nil

		time.Sleep(10 * time.Second)
	}
}
