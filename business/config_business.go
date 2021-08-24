package business

import (
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

func (c *ConfigBusiness) GetConfigJson(keyName string, enableCache bool) string {
	cacheKey := cache.GetCacheKey(cache.LongTimeConfigKeyCache)
	redisService := services.NewRedisService()
	if enableCache == true {
		config := redisService.Hget(cacheKey, "author_cate")
		if config != "" {
			return utils.DeserializeData(config)
		}
	}
	configM := dcm.DcConfigJson{}
	exist, _ := dcm.GetBy("key_name", keyName, &configM)
	if exist {
		jsonData := utils.SerializeData(configM.Value)
		_ = redisService.Hset(cacheKey, "author_cate", jsonData)
	}
	return configM.Value
}
