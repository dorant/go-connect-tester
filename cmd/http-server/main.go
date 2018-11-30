package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	port = ":50000"
)

type request struct {
	Hostname string `json: hostname`
}

var hostname string
var respondWithClose bool

func hello(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}

	req := request{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	name := req.Hostname
	log.Printf("Call from client '%s' received\n", name)

	rsp := request{Hostname: hostname}
	rspData, err := json.Marshal(rsp)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")

	if respondWithClose {
		// Making sure the client disconnects, and not re-use connection
		w.Header().Set("Connection", "close")
	}
	w.Write(rspData)
}

func main() {
	respondWithClose = false
	if strings.ToLower(os.Getenv("FORCE_CLOSE")) == "yes" {
		log.Println("Force close of client request enabled")
		respondWithClose = true
	}

	hostname, _ = os.Hostname()

	http.HandleFunc("/", hello)

	log.Printf("Start listen to: %s\n", port)
	http.ListenAndServe(port, nil)
}
