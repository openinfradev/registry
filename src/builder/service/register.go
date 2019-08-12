package service

import (
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
	if err != nil && pong != "PONG" {
		return false
	}
	// problem...
	return true
}

// Regist is builder registered in Redis
func (g *RegisterService) Regist() {
	// self := util.GetLocalIP() + ":" + basicinfo.ServicePort
	// builders, err := g.Register.Get("builders").Result()
	// if err != nil {
	// 	logger.ERROR("register.go", err.Error())
	// }
	// builderArray := json.Unmarshal([]byte(builders), []string{})
	// builderArray = append(builderArray, self)

}

// InitClient is register create
func (g *RegisterService) InitClient() {
	g.Register = redis.NewClient(&redis.Options{
		Addr:     basicinfo.RedisEndpoint,
		Password: "",
		DB:       0,
	})
}
