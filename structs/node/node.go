package node

import "time"

type Node struct {
	Address  string    `json:"address"`
	LastPing time.Time `json:"lastPing"`
	Ping     int       `json:"ping"`
}
