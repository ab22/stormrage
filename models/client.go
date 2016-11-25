package models

type Client struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Target         string `json:"target"`
	MaxLimit       string `json:"maxLimit"`
	BurstLimit     string `json:"burstLimit"`
	BurstThreshold string `json:"burstThreshold"`
	BurstTime      string `json:"burstTime"`
}
