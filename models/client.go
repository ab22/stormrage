package models

type Client struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Target         string `json:"target"`
	MaxLimit       string `json:"max_limit"`
	BurstLimit     string `json:"burst_limit"`
	BurstThreshold string `json:"burst_threshold"`
	BurstTime      string `json:"burst_time"`
}
