package coinbasepro

import (
	"time"
)

var (
	lastRequest       time.Time
	RequestsPerSecond float64                                             = 10
	BeforeRequest     func(client *Client, method, endpoint string) error = nil
	AfterRequest      func()                                              = nil
)

func init() {
	BeforeRequest = func(client *Client, method, endpoint string) error {
		elapsed := time.Since(lastRequest)
		if elapsed.Seconds() < (float64(1) / RequestsPerSecond) {
			time.Sleep(time.Duration((float64(time.Second) / RequestsPerSecond)) - elapsed)
		}
		return nil
	}
	AfterRequest = func() {
		lastRequest = time.Now()
	}
}
