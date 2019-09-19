package service

import (
	"builder/model"
	"builder/util/logger"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// RegisterService is builder registered in Redis
type RegisterService struct {
	Register *redis.Client
}

// Health returns health check
func (g *RegisterService) Health() bool {
	if g.Register == nil {
		g.InitClient()
	}
	pong, err := g.Register.Ping().Result()
	if err != nil || pong != "PONG" {
		logger.ERROR("service/builder-register.go", "Health", err.Error())
		return false
	}
	return true
}

// Regist is builder registered in Redis
func (g *RegisterService) Regist() {
	if !g.Health() {
		logger.ERROR("service/builder-register.go", "Regist", "redis is unhealthy")
		return
	}
	builderListJSON, err := g.Register.Get("builderList").Result()
	if err != nil {
		logger.ERROR("service/builder-register.go", "Regist", err.Error())
	}
	builderList := &model.BuilderList{}
	json.Unmarshal([]byte(builderListJSON), &builderList)
	port, _ := strconv.Atoi(basicinfo.ServicePort)
	// ???
	// self := &model.Builder{
	// 	Host: util.GetOutboundIP(),
	// 	Port: port,
	// }
	self := &model.Builder{
		Host: basicinfo.ServiceDomain,
		Port: port,
	}
	exist := false
	for _, b := range builderList.Builders {
		if b.Host == self.Host && b.Port == self.Port {
			exist = true
			break
		}
	}
	if !exist {
		builderList.Builders = append(builderList.Builders, *self)
	}

	bytes, _ := json.Marshal(builderList)
	result, err := g.Register.Set("builderList", string(bytes), 0).Result()
	if err != nil {
		logger.ERROR("service/builder-register.go", "Regist", err.Error())
		return
	}
	logger.INFO("service/builder-register.go", "Regist", result+"::"+string(bytes))
}

// Sync is builder host sync to redis
func (g *RegisterService) Sync() {
	ticker := time.NewTicker(time.Second * 15)
	go func() {
		for t := range ticker.C {
			logger.DEBUG("service/builder-register.go", "Regist", t.String())
			g.Regist()
		}
	}()
}

// InitClient is register create
func (g *RegisterService) InitClient() {
	g.Register = redis.NewClient(&redis.Options{
		Addr:     basicinfo.RedisEndpoint,
		Password: "",
		DB:       0,
	})
}
