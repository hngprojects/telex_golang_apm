package telex

type APMMetrics struct {
	Path      string `json:"path"`
	EventName string `json:"event_name"`
	Message   string `json:"message"`
	Status    string `json:"status"`
	Username  string `json:"username"`
}

type Metrics struct {
	APMMetrics         APMMetrics         `json:"apm_metrics"`
	RequestMetrics     RequestMetrics     `json:"request_metrics"`
	ErrorMetrics       ErrorMetrics       `json:"error_metrics"`
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
}

type RequestMetrics struct {
	Path       string `json:"path"`
	Method     string `json:"method"`
	Latency    string `json:"latency"`
	StatusCode int    `json:"status_code"`
}

type ErrorMetrics struct {
	ErrorMessage string   `json:"error_message"`
	ErrorCount   int      `json:"error_count"`
	LastErrors   []string `json:"last_error"`
}

type PerformanceMetrics struct {
	MemoryAllocation   uint64  `json:"memory_allocation"`
	MemoryUsage        uint64  `json:"memory_usage"`
	MaxMemoryUsage     uint64  `json:"max_memory_usage"`
	AverageRequestTime float64 `json:"average_request_time"`
	CPUUsage           int     `json:"cpu_usage"`
	NetworkTraffic     int     `json:"network_traffic"`
	GCCycles           uint32  `json:"gc_cycles"`
	Goroutines         int     `json:"goroutines"`
}
