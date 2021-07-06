package business

import (
	"strings"
	"time"
)

const (
	VoUserUniqueTokenPlatformWeb           = 1
	VoUserUniqueTokenPlatformH5            = 2
	VoUserUniqueTokenPlatformWxMiniProgram = 3
	VoUserUniqueTokenPlatformApp           = 4
	VoUserUniqueTokenPlatformWap           = 5
)

func GetAppPlatFormIdWithAppId(appId int) int {
	switch appId {
	case 10000:
		return VoUserUniqueTokenPlatformWeb
	case 10001:
		return VoUserUniqueTokenPlatformH5
	case 10002:
		return VoUserUniqueTokenPlatformWxMiniProgram
	case 10003, 10004:
		return VoUserUniqueTokenPlatformApp
	case 10005:
		return VoUserUniqueTokenPlatformWap
	}
	return 0
}

func GetAge(startTime string) int {
	var Age int64
	var pslTime string
	if strings.Index(startTime, ".") != -1 {
		pslTime = "2006.01.02"
	} else if strings.Index(startTime, "-") != -1 {
		pslTime = "2006-01-02"
	} else {
		pslTime = "2006/01/02"
	}
	t1, err := time.ParseInLocation(pslTime, startTime, time.Local)
	t2 := time.Now()
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		Age = diff / (3600 * 365 * 24)
		return int(Age)
	} else {
		return int(Age)
	}
}
