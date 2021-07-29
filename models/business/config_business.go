package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/models/dcm"
)

type ConfigBusiness struct {
}

func NewConfigBusiness() *ConfigBusiness {
	return new(ConfigBusiness)
}

func (c *ConfigBusiness) GetConfigJson(keyName string, enableCache bool) string {
	memberKey := cache.GetCacheKey(cache.ConfigKeyCache, keyName)
	if enableCache == true {
		config := global.Cache.Get(memberKey)
		if config != "" {
			return config
		}
	}
	configM := dcm.DcConfigJson{}
	exist, _ := dcm.GetBy("key_name", keyName, &configM)
	if exist {
		_ = global.Cache.Set(memberKey, configM.Value, 3600)
	}
	return configM.Value
}
