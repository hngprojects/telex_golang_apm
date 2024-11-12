package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hngprojects/go-apm-sdk/telex"
	"github.com/hngprojects/go-apm-sdk/telexmux"
)

func main() {
	apmClient, err := telex.Init(
		telex.APMOptions{
			WebhookURL: "https://ping.telex.im/v1/webhooks/cf82320045eb?username=collins",
		},
	)
	if err != nil {
		log.Printf("Failed to initialize APM: %v", err)
	}

	router := mux.NewRouter()
	router.Use(telexmux.NewMux(apmClient, telexmux.Options{
		Repanic:         false,
		WaitForDelivery: true,
		Timeout:         5 * time.Second,
	}))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	router.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("Something went wrong")
	})

	http.ListenAndServe(":8081", router)
}
