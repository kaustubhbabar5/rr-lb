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
	log.Println("waiting........")
	time.Sleep(20 * time.Second)
	for {
		res, err := http.Get("http://0.0.0.0:8080/proxy/")
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
