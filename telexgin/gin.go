package telexgin

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/go-apm-sdk/telex"
)

type Options struct {
	Repanic         bool
	WaitForDelivery bool
	Timeout         time.Duration
}

type Handler struct {
	apmClient *telex.Client
	//more clients can be added

	Options
	
	totalLatency   float64
	totalRequests  int
	maxMemoryUsage uint64
	errorCount     int
	recentErrors   []string
	errorMsg       string

}

func NewGin(apmClient *telex.Client, options Options) gin.HandlerFunc {
	if options.Timeout == 0 {
		options.Timeout = 5 * time.Second
	}

	handler := &Handler{
		apmClient: apmClient,
		Options:   options,
	}

	return handler.handlegin
}
