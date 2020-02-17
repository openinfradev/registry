package service

import (
	"builder/config"
	"builder/util/logger"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// WebhookService is relative registry hook services
type WebhookService struct{}

func init() {
}

// Toss is delivered from registry to app
func (w *WebhookService) Toss(body []byte) {
	m := make(map[string]interface{})
	err := json.Unmarshal(body, &m)
	if err != nil {
		logger.ERROR("service/webhook.go", "toss", err.Error())
		return
	}

	b, _ := json.Marshal(m)
	buff := bytes.NewBuffer(b)

	logger.DEBUG("service/webhook.go", "toss", "start toss")

	conf := config.GetConfig()
	if conf.Webhook != nil && conf.Webhook.Listener != nil {
		paths := conf.Webhook.Listener
		for _, path := range paths {
			logger.DEBUG("service/webhook.go", "toss", path)
			resp, err := http.Post(path, "application/json", buff)
			if err != nil {
				logger.ERROR("service/webhook.go", "toss", err.Error())
				return
			}
			defer resp.Body.Close()

			r, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.ERROR("service/webhook.go", "Toss", err.Error())
				return
			}
			logger.DEBUG("service/webhook.go", "Toss", string(r))
		}
	}
}
