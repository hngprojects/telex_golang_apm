package telex

import (
	"time"

)

//global collector interface for all apm clients
type Collector interface {
	RequestMetricsCollector(startTime time.Time, path string, method string, statusCode int) Metrics
	ErrorMetricsCollector(errorMsg string) Metrics
	PerformanceMetricsCollector() Metrics
}
