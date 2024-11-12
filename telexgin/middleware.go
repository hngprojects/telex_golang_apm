package telexgin

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/go-apm-sdk/telex"
)

func (h *Handler) handlegin(c *gin.Context) {
	startTime := time.Now()

	username := c.GetString("username")
	if username == "" {
		username = "APM"
	}

	defer func() {
		if err := recover(); err != nil {

			if h.Options.Repanic {
				panic(err)
			}

			payload := h.apmClient.ReportError(err, username)

			// Synchronously or asynchronously send error metrics
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

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
	}()

	c.Next()

	// Metrics collection

	reqMetrics := h.RequestMetricsCollector(
		startTime,
		c.Request.URL.Path,
		c.Request.Method,
		c.Writer.Status(),
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
			Message:   msg,
			Status:    status,
			Username:  username,
		},
		h.Options.Timeout,
	)

	log.Printf("Request Metrics: %s", msg)
}
