package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
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
			return utils.DeserializeData(config)
		}
	}
	configM := dcm.DcConfigJson{}
	exist, _ := dcm.GetBy("key_name", keyName, &configM)
	if exist {
		jsonData := utils.SerializeData(configM.Value)
		_ = global.Cache.Set(memberKey, jsonData, 3600)
	}
	return configM.Value
}
