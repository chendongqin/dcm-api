package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/services"
)

type ConfigBusiness struct {
}

func NewConfigBusiness() *ConfigBusiness {
	return new(ConfigBusiness)
}

func (c *ConfigBusiness) GetHashConfigJson(keyName string, enableCache bool) string {
	cacheKey := cache.GetCacheKey(cache.LongTimeConfigKeyCache)
	redisService := services.NewRedisService()
	if enableCache == true {
		config := redisService.Hget(cacheKey, keyName)
		if config != "" {
			return utils.DeserializeData(config)
		}
	}
	configM := dcm.DcConfigJson{}
	exist, _ := dcm.GetBy("key_name", keyName, &configM)
	if exist {
		jsonData := utils.SerializeData(configM.Value)
		if enableCache == true {
			_ = redisService.Hset(cacheKey, keyName, jsonData)
		}
	}
	return configM.Value
}

func (c *ConfigBusiness) GetConfigJson(keyName string, enableCache bool) string {
	cacheKey := cache.GetCacheKey(cache.ConfigKeyCache, keyName)
	if enableCache == true {
		config := global.Cache.Get(cacheKey)
		if config != "" {
			return utils.DeserializeData(config)
		}
	}
	configM := dcm.DcConfigJson{}
	exist, _ := dcm.GetBy("key_name", keyName, &configM)
	if exist {
		jsonData := utils.SerializeData(configM.Value)
		if enableCache == true {
			_ = global.Cache.Set(cacheKey, jsonData, 300)
		}
	}
	return configM.Value
}
