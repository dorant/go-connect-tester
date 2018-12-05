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

const (
	timeBetweenConnects = 10 * time.Second
)

type request struct {
	Hostname string `json: hostname`
}

func sendRequest(url string, dataToSend []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataToSend))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	// Make sure the connection is timed out before next connect
	client.Timeout = timeBetweenConnects - 2*time.Second

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failure while Get(): %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status code %d - %s",
			resp.StatusCode, http.StatusText(resp.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading the body: %s", err)
		return
	}

	bodyRsp := request{}
	err = json.Unmarshal(body, &bodyRsp)
	if err != nil {
		log.Printf("Decoding error: %s", err)
		return
	}

	name := bodyRsp.Hostname
	log.Printf("Response from server '%s' received\n", name)
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

	// Prepare data to send
	hostname, _ := os.Hostname()
	b, err := json.Marshal(request{Hostname: hostname})
	if err != nil {
		log.Printf("json failure: %v\n", err)
	}

	for {
		sendRequest(url, b)
		time.Sleep(timeBetweenConnects)
	}
}
