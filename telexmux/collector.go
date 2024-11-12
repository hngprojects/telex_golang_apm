package telexmux


import (
	"runtime"
	"time"

	"github.com/hngprojects/go-apm-sdk/telex"
)

func (h *Handler) RequestMetricsCollector(startTime time.Time, path string, method string, statusCode int) telex.Metrics {
	latency := time.Since(startTime)
	//request count metrics

	requestMetrics := telex.Metrics{
		RequestMetrics: telex.RequestMetrics{
			Path:       path,
			Method:     method,
			Latency:    latency.String(),
			StatusCode: statusCode,
		},
	}

	return requestMetrics
}

func (h *Handler) PerformanceMetricsCollector() telex.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	h.totalRequests++
	h.totalLatency += float64(memStats.TotalAlloc)
	avgLatency := h.totalLatency / float64(h.totalRequests)

	performanceMetrics := telex.Metrics{
		PerformanceMetrics: telex.PerformanceMetrics{
			MemoryAllocation:   memStats.Alloc,
			CPUUsage:           runtime.NumCPU(),
			MemoryUsage:        memStats.Sys,
			NetworkTraffic:     0,
			GCCycles:           memStats.NumGC,
			Goroutines:         runtime.NumGoroutine(),
			AverageRequestTime: avgLatency,
			MaxMemoryUsage:     max(memStats.TotalAlloc, h.maxMemoryUsage),
		},
	}
	return performanceMetrics
}

func (h *Handler) ErrorMetricsCollector() telex.Metrics {
	h.errorCount++

	if len(h.recentErrors) >= 5 {
		h.recentErrors = h.recentErrors[1:]
	}

	errorMetrics := telex.Metrics{
		ErrorMetrics: telex.ErrorMetrics{
			ErrorCount:   h.errorCount,
			ErrorMessage: h.errorMsg,
			LastErrors:   h.recentErrors,
		},
	}
	return errorMetrics

}
