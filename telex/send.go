package telex

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"
)

func (c *Client) Message(metrics Metrics) (string, error) {
	const metricsTemplate = 
	`
		Request Metrics:
		Path: {{.RequestMetrics.Path}}
		Method: {{.RequestMetrics.Method}}
		StatusCode: {{.RequestMetrics.StatusCode}}
		Latency: {{.RequestMetrics.Latency}}
	Performance Metrics:
		CPU Usage: {{.PerformanceMetrics.CPUUsage}} %
		Memory Usage: {{.PerformanceMetrics.MemoryUsage}} B
		Network Traffic: {{.PerformanceMetrics.NetworkTraffic}} B/s
		GC Cycles: {{.PerformanceMetrics.GCCycles}}
		Goroutines: {{.PerformanceMetrics.Goroutines}}
		Average Request Time: {{.PerformanceMetrics.AverageRequestTime}} ms
		Max Memory Usage: {{.PerformanceMetrics.MaxMemoryUsage}} B
		Memory Allocation: {{.PerformanceMetrics.MemoryAllocation}} B
	Error Metrics:
		Error Message: {{.ErrorMetrics.ErrorMessage}}
		Error Count: {{.ErrorMetrics.ErrorCount}}
		Last Errors: {{.ErrorMetrics.LastErrors}}
	`
	tmpl, err := template.New("metrics").Parse(metricsTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, metrics)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *Client) SendMetricsToWebHook(webhookURL string, metrics APMMetrics, timeout time.Duration) {
	payload, jsonErr := json.Marshal(metrics)
	if jsonErr != nil {
		log.Printf("Error marshalling metrics payload: %v", jsonErr)
		return
	}

	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	req, err := http.NewRequest(
		"POST",
		webhookURL,
		bytes.NewBuffer(payload),
	)

	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending metrics to webhook: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Metrics sent to webhook: %s", webhookURL)
}

// func (c *Client)SendMetricsToKafka(metrics Metrics, timeout time.Duration) {}

// func (c *Client)SendMetricsToS3(metrics Metrics, timeout time.Duration) {}

// func (c *Client)SendMetricsToSNS(metrics Metrics, timeout time.Duration) {}

// func (c *Client)SendMetricsToSQS(metrics Metrics, timeout time.Duration) {}
