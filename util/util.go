package util

import (
	"encoding/json"
	"github.com/openinfradev/registry-builder/util/logger"
	"net"
	"strconv"
	"strings"
	"time"
)

// GitRepositoryURL is protocol and url
type GitRepositoryURL struct {
	Protocol string
	URL      string
}

// StringToMap returns map interface to convert json string
func StringToMap(str string) *map[string]interface{} {
	var raw map[string]interface{}
	json.Unmarshal([]byte(str), &raw)
	return &raw
}

// GetTimeMillisecond returns current time to millisecond string
func GetTimeMillisecond() string {
	now := time.Now()
	nano := now.UnixNano()
	ms := nano / 1000000
	return strconv.FormatInt(ms, 10)
}

// GetLocalIP returns local ips
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.ERROR("util/util.go", "GetLocalIP", err.Error())
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// logger.DEBUG("util.go", ipnet.IP.String())
				// for test ... first ip only
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// GetOutboundIP returns outbound ip
func GetOutboundIP() string {
	// china not work
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.ERROR("util/util.go", "GetOutboundIP", err.Error())
		return ""
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	// logger.DEBUG("util.go", localAddr.IP.String())

	return localAddr.IP.String()
}

// MapToStruct returns raw map to target struct type
func MapToStruct(raw interface{}, dist interface{}) error {
	jsonbody, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonbody, dist)
	if err != nil {
		return err
	}
	return nil
}

// ExtractGitRepositoryURL returns extracted protocol and url
func ExtractGitRepositoryURL(raw string) *GitRepositoryURL {
	gitRepo := &GitRepositoryURL{}

	// 1. protocol
	if strings.HasPrefix(raw, "http://") {
		raw = strings.Replace(raw, "http://", "", 1)
		gitRepo.Protocol = "http"
	} else if strings.HasPrefix(raw, "https://") {
		raw = strings.Replace(raw, "https://", "", 1)
		gitRepo.Protocol = "https"
	} else {
		gitRepo.Protocol = "https"
	}

	// 2. @(at)
	if strings.Contains(raw, "@") {
		tmp := strings.Split(raw, "@")
		if len(tmp) > 1 {
			raw = tmp[1]
		}
	}

	gitRepo.URL = raw

	return gitRepo
}
