package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/go-apm-sdk/telex"
	"github.com/hngprojects/go-apm-sdk/telexgin"
)

func main() {
	//Initialize the APM
	apmClient, err := telex.Init(telex.APMOptions{
		WebhookURL:        "https://ping.telex.im/v1/webhooks/cf82320045eb?username=collins",
		Async:             false,
		EnableTracing:     true,
		TracingSampleRate: 1.0,
	})

	if err != nil {
		log.Fatalf("Failed to initialize APM: %v", err)
	}

	router := gin.Default()
	router.Use(telexgin.NewGin(apmClient, telexgin.Options{
		Repanic:         false, // Set to true only in development and debugging environments
		WaitForDelivery: true,
		Timeout:         5,
	}))

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Hello, World!"})
	})

	router.GET("/panic", func(ctx *gin.Context) {
		panic("Something went wrong")
	})

	router.Run(":8081")

}
