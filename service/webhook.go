package service

import (
	"bytes"
	"encoding/json"
	"github.com/openinfradev/registry-builder/config"
	"github.com/openinfradev/registry-builder/util/logger"
	"net/http"
)

// WebhookService is relative registry hook services
type WebhookService struct{}

func init() {
}

// Toss is delivered from registry to app
func (w *WebhookService) Toss(body *map[string]interface{}) {
	b, _ := json.Marshal(body)
	buff := bytes.NewBuffer(b)

	// logger.DEBUG("service/webhook.go", "toss", "start toss")

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

			// r, err := ioutil.ReadAll(resp.Body)
			// if err != nil {
			// 	logger.ERROR("service/webhook.go", "Toss", err.Error())
			// 	return
			// }
			// logger.DEBUG("service/webhook.go", "Toss", string(r))
		}
	}
}
