package telexmux

import (
	"log"
	"net/http"
	"time"

	"github.com/hngprojects/go-apm-sdk/telex"
)

func (h *Handler) handlemux(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		username := r.Header.Get("username")
		if username == "" {
			username = "APM"
		}

		defer func() {
			if err := recover(); err != nil {
				if h.Options.Repanic {
					panic(err)
				}

				payload := h.apmClient.ReportError(err, "APM")

				if h.Options.WaitForDelivery {
					h.apmClient.SendMetricsToWebHook(
						h.apmClient.Options.WebhookURL,
						payload,
						h.Options.Timeout,
					)
				} else {
					go h.apmClient.SendMetricsToWebHook(
						h.apmClient.Options.WebhookURL,
						payload,
						h.Options.Timeout,
					)
				}

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		}()

		next.ServeHTTP(w, r)
		rw := NewResponseWriter(w)

		reqMetrics := h.RequestMetricsCollector(
			startTime,
			r.URL.Path,
			r.Method,
			rw.status,
		)
		perfMetrics := h.PerformanceMetricsCollector()

		metrics := telex.Metrics{
			RequestMetrics:     reqMetrics.RequestMetrics,
			PerformanceMetrics: perfMetrics.PerformanceMetrics,
		}

		m := reqMetrics.RequestMetrics
		status := "success"
		event := "request_completed"

		if m.StatusCode >= 300 {
			status = "error"
			event = "request_failed"
		}

		msg, err := h.apmClient.Message(metrics)
		if err != nil {
			log.Printf("Error creating message: %v", err)
			return
		}

		h.apmClient.SendMetricsToWebHook(
			h.apmClient.Options.WebhookURL,
			telex.APMMetrics{
				EventName: event,
				Status:    status,
				Message:   msg,
				Username:  username,
			},
			h.Options.Timeout,
		)

		log.Printf("Request Metrics : %v", msg)

	})
}
