package global

import (
	"strconv"
)

/**
输入API类型
*/
type InputMap map[string]interface{}

func (this InputMap) _checkKey(key string) interface{} {
	if this == nil || key == "" {
		return nil
	}
	if _, ok := this[key]; ok {
		return this[key]
	} else {
		return nil
	}
}

func (this InputMap) GetInt(key string, defaultVal int) int {
	keyData := this._checkKey(key)
	if keyData != nil {
		//fmt.Println(keyData.(int))
		//if _,ok:=keyData.(int);ok{
		//	return keyData.(int)
		//}
		if _, ok := keyData.(float64); ok {
			return int(keyData.(float64))
		} else if val, ok := keyData.(string); ok {
			i, err := strconv.Atoi(val)
			if err != nil {
				return defaultVal
			}
			return i
		}
	}
	return defaultVal
}

func (this InputMap) Get(key string) interface{} {
	return this._checkKey(key)
}

func (this InputMap) GetArrArrString(key string) [][]string {
	keyData := this.GetArr(key)
	var ret [][]string
	if keyData != nil {
		for _, v := range keyData {
			var subRet []string
			if stt, ok := v.([]interface{}); ok {
				for _, v2 := range stt {
					if str, ok := v2.(string); ok {
						subRet = append(subRet, str)
					}
				}
				ret = append(ret, subRet)
			}
		}
	}
	return ret
}

//布尔型 0 1运算
func (this InputMap) GetIntBool(key string, defaultVal int) int {
	keyData := this._checkKey(key)
	if keyData != nil {
		if _, ok := keyData.(float64); ok {
			retbool := int(keyData.(float64))
			if retbool != 0 && retbool != 1 {
				retbool = defaultVal
			}
			return retbool
		}
	}
	return defaultVal
}

func (this InputMap) GetInt64(key string, defaultVal int64) int64 {
	keyData := this._checkKey(key)
	if keyData != nil {
		if _, ok := keyData.(float64); ok {
			return int64(keyData.(float64))
		}
	}
	return defaultVal
}

func (this InputMap) GetFloat64(key string, defaultVal float64) float64 {
	keyData := this._checkKey(key)
	if keyData != nil {
		if _, ok := keyData.(float64); ok {
			return keyData.(float64)
		}
	}
	return defaultVal
}

func (this InputMap) GetString(key string, defaultVal string) string {
	//args ...interface{}
	keyData := this._checkKey(key)
	if keyData != nil {
		if kd, ok := keyData.(string); ok {
			if kd == "" {
				return defaultVal
			} else {
				return kd
			}
		}
	}
	return defaultVal
}

func (this InputMap) GetStringWithEmpty(key string) string {
	return this.GetString(key, "")
}

func (this InputMap) GetStringWithQualified(key string, def string, qualified ...string) string {
	str := this.GetString(key, def)
	if str != def {
		allowed := false
		for _, v := range qualified {
			if v == str {
				allowed = true
			}
		}
		if allowed {
			return str
		} else {
			return def
		}
	}
	return str
}

func (this InputMap) GetPage(key ...string) int {
	keyName := "page"
	if len(key) >= 1 {
		keyName = key[0]
	}
	val := this.GetInt(keyName, 1)
	return val
}

func (this InputMap) GetPageSize(key string, defSize int, maxSize int) (size int) {
	//var err error
	size = defSize
	size = this.GetInt(key, defSize)
	if size > maxSize {
		size = maxSize
	}
	return
}

func (this InputMap) GetBool(key string, defaultVal bool) bool {
	keyData := this._checkKey(key)
	if keyData != nil {
		if _, ok := keyData.(bool); ok {
			return keyData.(bool)
		}
	}
	return defaultVal
}

func (this InputMap) Exists(key string) bool {
	keyData := this._checkKey(key)
	if keyData != nil {
		return true
	}
	return false
}

func (this InputMap) GetArr(key string) []interface{} {
	keyData := this._checkKey(key)
	if keyData != nil {
		if _, ok := keyData.([]interface{}); ok {
			return keyData.([]interface{})
		}
	}
	return nil
}

func (this InputMap) GetArrString(key string) []string {
	keyData := this.GetArr(key)
	if keyData != nil {
		ret := []string{}
		for _, v := range keyData {
			if _, ok := v.(string); ok {
				ret = append(ret, v.(string))
			}
		}
		return ret
	}
	return nil
}

func (this InputMap) GetArrInt64(key string) []int64 {
	keyData := this.GetArr(key)
	if keyData != nil {
		ret := []int64{}
		for _, v := range keyData {
			if fdata, ok := v.(float64); ok {
				ret = append(ret, int64(fdata))
			}
		}
		return ret
	}
	return nil
}

func (this InputMap) GetArrInt(key string) []int {
	keyData := this.GetArr(key)
	if keyData != nil {
		ret := []int{}
		for _, v := range keyData {
			if fdata, ok := v.(float64); ok {
				ret = append(ret, int(fdata))
			}
		}
		return ret
	}
	return nil
}
