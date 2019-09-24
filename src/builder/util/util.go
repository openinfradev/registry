package util

import (
	"builder/util/logger"
	"encoding/json"
	"net"
	"strconv"
	"time"
)

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
func MapToStruct(raw interface{}, target interface{}) error {
	jsonbody, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonbody, target)
	if err != nil {
		return err
	}
	return nil
}
