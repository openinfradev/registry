package service

import (
	"builder/config"
	"builder/util/logger"
	"bytes"
	"net/http"
)

// WebhookService is relative registry hook services
type WebhookService struct{}

func init() {
}

// Toss is delivered from registry to app
func (w *WebhookService) Toss(body []byte) {
	buff := bytes.NewBuffer(body)

	conf := config.GetConfig()
	if conf.Webhook != nil && conf.Webhook.Listener != nil {
		paths := conf.Webhook.Listener
		for _, path := range paths {
			resp, err := http.Post(path, "application/json", buff)
			if err != nil {
				logger.ERROR("service/webhook.go", "toss", err.Error())
				return
			}
			defer resp.Body.Close()
		}
	}
}
