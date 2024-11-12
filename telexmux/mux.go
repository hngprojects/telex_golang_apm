package telexmux

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/hngprojects/go-apm-sdk/telex"
)

type Options struct {
	Repanic         bool
	WaitForDelivery bool
	Timeout         time.Duration
}

type Handler struct {
	apmClient *telex.Client
	Options
	totalLatency   float64
	totalRequests  int
	maxMemoryUsage uint64
	errorCount     int
	recentErrors   []string
	errorMsg       string
}

func NewMux(apmClient *telex.Client, options Options) mux.MiddlewareFunc {
	if options.Timeout == 0 {
		options.Timeout = 5 * time.Second
	}

	handler := &Handler{
		apmClient: apmClient,
		Options:   options,
	}

	return handler.handlemux
}
