package models

type CommonMap map[string]interface{}

func ParsePageAndPageSize(params CommonMap) (rePage, rePageSize, limitOffSet int) {
	rePage = params.GetInt("page")
	rePageSize = params.GetInt("pageSize")
	if rePage <= 0 {
		rePage = 1
	}
	if rePageSize <= 0 {
		rePageSize = 50
	}
	limitOffSet = (rePage - 1) * rePageSize
	return
}

func (cm *CommonMap) GetInt(key string) int {
	if v, ok := (*cm)[key]; ok {
		return v.(int)
	}
	return 0
}

func (cm *CommonMap) GetInt64(key string) int64 {
	if v, ok := (*cm)[key]; ok {
		return v.(int64)
	}
	return 0
}

func (cm *CommonMap) GetString(key string) string {
	if v, ok := (*cm)[key]; ok {
		return v.(string)
	}
	return ""
}
