package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	for i := 1; i < 10; i++ {
		address := fmt.Sprintf("0.0.0.0:808%v", i)

		go StartHttpServer(address)

		body, err := json.Marshal(map[string]string{"EndPoint": "http://" + address + "/"})
		if err != nil {
			log.Fatalln(err)
		}

		res, err := http.Post("http://0.0.0.0:8080/url/register", "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Fatalln(err)
		}

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			log.Fatalln("unknown status code for:", address, string(resBody))
		}
	}
	log.Println("waiting 20 secs")
	time.Sleep(20 * time.Second)
	for {
		req, err := http.NewRequest("GET", "http://0.0.0.0:8080/proxy/", nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Add("Content-Type", " application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			log.Fatalln("unknown status code from proxy:", res.StatusCode, string(resBody))
		}

		log.Println("response: ", string(resBody))

		time.Sleep(200 * time.Millisecond)

	}

}

func StartHttpServer(address string) {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello from " + address))
	})
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatalln(err)
	}
}
