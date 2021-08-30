package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"dongchamao/global"
	"dongchamao/services/dyimg"
	"dongchamao/services/hbaseService/hbase"
	"dongchamao/services/mutex"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"html/template"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

/**
 * string转换int
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func ParseInt(b string, defInt int) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		return defInt
	} else {
		return id
	}
}

/**
 * int转换string
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func ParseString(b int) string {
	id := strconv.Itoa(b)
	return id
}

/**
 * int64转换string
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func ParseStringInt64(b int64) string {
	id := strconv.FormatInt(b, 10)
	return id
}

func ParseInt64String(b string) int64 {
	id, _ := strconv.ParseInt(b, 10, 64)
	return id
}

/**
 * 转换浮点数为string
 * @method func
 * @param  {[type]} t *             Tools [description]
 * @return {[type]}   [description]
 */
func ParseFlostToString(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func ParseStringFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func ParseByteInt64(byteData []byte) int64 {
	if len(byteData) < 8 {
		return 0
	}
	buf := bytes.NewBuffer(byteData)
	var i2 int64
	binary.Read(buf, binary.BigEndian, &i2)
	return i2
}

func ParseByteInt32(byteData []byte) int32 {
	if len(byteData) < 4 {
		return 0
	}
	buf := bytes.NewBuffer(byteData)
	var i2 int32
	binary.Read(buf, binary.BigEndian, &i2)
	return i2
}

func ParseByteFloat64(byteData []byte) float64 {
	if len(byteData) < 8 {
		return 0
	}
	bits := binary.BigEndian.Uint64(byteData)
	result := math.Float64frombits(bits)
	if math.IsNaN(result) {
		return 0
	} else if math.IsInf(result, 0) {
		return 0
	}
	return result
}

func ParseByteFloat32(byteData []byte) float32 {
	if len(byteData) < 4 {
		return 0
	}
	bits := binary.BigEndian.Uint32(byteData)
	return math.Float32frombits(bits)
}

/**
 * 结构体转成json 字符串
 * @method StruckToString
 * @param  {[type]}       data interface{} [description]
 */
func StructToString(data interface{}) (str string, err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	} else {
		str = string(b)
	}
	return
}

/**
 * 结构体转换成map对象
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func StructToMap(obj interface{}) map[string]interface{} {
	k := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < k.NumField(); i++ {
		data[strings.ToLower(k.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

func MapToStruct(objMap interface{}, obj interface{}) {
	jsonByte, _ := json.Marshal(objMap)
	err := json.Unmarshal(jsonByte, &obj)
	if err != nil {
		logs.Error(err)
	}
}

func InterfaceToString(v interface{}) string {
	var ret string
	//fmt.Println(v.(type))
	switch v.(type) {
	case []byte:
		ret = string(v.([]byte))
	case string:
		ret = v.(string)
	case int:
		ret = ParseString(v.(int))
	case float64:
		//ret = ParseString((int(v.(float64))))
		ret = ParseStringInt64((int64(v.(float64))))
	case []interface{}:
		ret = ""
	case error:
		return ""
	default:
		ret = ""
	}
	return ret
}

func DoubleToBytes(n float64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, n)
	return bytesBuffer.Bytes()
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func HideNickname(nickname string) string {
	characters := []rune(nickname)
	length := len(characters)
	if length >= 2 {
		first := characters[0]
		last := characters[length-1]
		nickname = string(first) + "**" + string(last)
	} else {
		nickname = nickname + "**"
	}
	return nickname
}

func Int64ToBytes(n int64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, n)
	return bytesBuffer.Bytes()
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func BytesToInt32(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return x
}

func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)

	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return x
}

//生成随机字符串
func GetRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()+[]{}/<>;:=.,?"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func GetRandomStringByLetter(n int, letterBytes string) string {
	lettersArr := strings.Split(letterBytes, "")
	retStr := ""
	b := make([]byte, n)
	for range b {
		randInt := rand.Int63() % int64(len(lettersArr))
		retStr += lettersArr[randInt]
	}
	return retStr
}

//生成随机字符串2
func GetRandomStringNew(n int) string {
	const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

//生成随机字符串2
func GetRandomInt(n int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

//不含max
func GenerateRangeNum(min, max int) int {
	if min == max {
		return min
	}
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}

/**
 * 字符串截取
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func SubString(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	rune_str := []rune(str)
	len_str := len(rune_str)

	if start < 0 {
		start = len_str + start
	}
	if start > len_str {
		start = len_str
	}
	end := start + length
	if end > len_str {
		end = len_str
	}
	if length < 0 {
		end = len_str + length
	}
	if start > end {
		start, end = end, start
	}
	return string(rune_str[start:end])
}

/**
 * base64 解码
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func Base64Decode(str string) string {
	s, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(s)
}

/**
 * 格式化时间数据
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func TimeFormat(time *time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

/**
 * 格式化时间数据
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func StringFormatTime(timeDo string, timeLayout string) int64 {
	if timeLayout == "" {
		timeLayout = "2006-01-02 15:04:05" //转化所需模板
	}
	loc, _ := time.LoadLocation("Local")                        //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, timeDo, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()
	return sr
}

/**
 * 旧token解析
 * @method func
 * @param token(string)
 * @return string
 */

func KeyDecode(token string) string {
	tokens := strings.Replace(token, "-", "+", -1)
	tokens = strings.Replace(tokens, "_", "/", -1)
	tokens = strings.Replace(tokens, " ", "", -1)
	mod := len(tokens) % 4
	if mod != 0 {
		buchong := string("===="[mod:])
		tokens = tokens + buchong
	}
	str := Base64Decode(tokens)
	return str
}

/**
 * curl
 * @method func
 * @param token(string)
 * @return string
 */

func CurlData(url string, method string, headers map[string]string, postData map[string]interface{}) string {
	postJson, _ := jsoniter.Marshal(postData)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(postJson))) //建立一个请求
	if err != nil {
		return ""
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req) //提交
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(body)
	}
	return ""
}

func SimpleCurl(url string, method string, postData string, contentType string) string {
	var resp *http.Response
	var err error
	if method == "POST" {
		resp, err = http.Post(url, contentType, strings.NewReader(postData))
	} else {
		resp, err = http.Get(url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func Curl(url string, method string, postData string, contentType string) (*simplejson.Json, error) {
	var resp *http.Response
	var err error
	if method == "POST" {
		resp, err = http.Post(url, contentType, strings.NewReader(postData))
	} else {
		resp, err = http.Get(url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res, _ := simplejson.NewJson(body)

	return res, nil
}

/**
 * @method func
 * @param token(string)
 * @return string
 */

func SortbyKey(sortArrays map[string]string) string {
	var keys []string
	var keyMap = make(map[string]string)
	for k := range sortArrays {
		//orginKey = append(keys, k)
		formatkey := strings.ToLower(k)
		keyMap[formatkey] = k
		keys = append(keys, formatkey)
	}
	sort.Strings(keys)
	buildStr := ""
	//To perform the opertion you want
	for _, k := range keys {
		realKey := keyMap[k]
		if sortArrays[realKey] != "" {
			buildStr = buildStr + realKey + "=" + sortArrays[realKey] + "&"
		}
	}
	return buildStr
}

func Md5_encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/**
 * curl
 * @method func
 * @param token(string)
 * @return string
 */

func Wechat_makeSign(params map[string]string, key string) string {
	buildStr := SortbyKey(params)
	signStr := buildStr + "key=" + key
	sign := strings.ToUpper(Md5_encode(signStr))
	return sign
}

func MapToXmlSimple(postData map[string]string) string {
	xmlStr := "<xml>\n"
	for k, v := range postData {
		xmlNode := "<" + k + ">" + v + "</" + k + ">\n"
		xmlStr = xmlStr + xmlNode
	}
	xmlStr += "</xml>"
	return xmlStr
}

func RandomOrderId() string {
	randStr := time.Now().Format("200601021504")
	randInt := GetRandomInt(8)
	return randStr + randInt
}

func GetClientIP(ctx *context.Context) string {
	return ctx.Input.IP()
}

func InArray(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func InArrayInt64(need int64, needArr []int64) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func InArrayInt(need int, needArr []int) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func InArrayString(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func Sha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

func Sha512(data string) string {
	sha := sha512.New()
	sha.Write([]byte(data))
	return hex.EncodeToString(sha.Sum([]byte("")))
}

func Sha256(data string) string {
	sha := sha256.New()
	sha.Write([]byte(data))
	return hex.EncodeToString(sha.Sum([]byte("")))
}

func CreateToken(claims jwt.MapClaims, SecretKey string) (tokenString string, err error) {
	//claims := jwt.MapClaims{
	//	data
	//}
	//claims := &jwt.StandardClaims{
	//	NotBefore: int64(time.Now().Unix()),
	//	ExpiresAt: int64(time.Now().Unix() + 1000),
	//	Issuer:    "Bitch",
	//}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(SecretKey))
	return
}

func Time() int64 {
	return time.Now().Unix()
}

func GetZeroHour(date time.Time) time.Time {
	zeroHour, _ := time.ParseInLocation("20060102", date.Format("20060102"), time.Local)
	return zeroHour
}

//Date(1524799394,"02/01/2006 15:04:05 PM")
func Date(format string, timestamp int64) string {
	if timestamp == 0 {
		timestamp = Time()
	}
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return time.Unix(timestamp, 0).Format(format)
}

//获取每月的第一天和最后一天日期
func GetFirstDateAndLastDate(date time.Time) (firstOfMonth time.Time, lastOfMonth time.Time) {
	firstOfMonth = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	lastOfMonth = firstOfMonth.AddDate(0, 1, -1)
	return
}

func GetFirstDateOfWeek() (weekStartDate time.Time) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return
}

func GetCurrentMonthBeginAndEndDate() (firstOfMonth time.Time, lastOfMonth time.Time) {
	firstOfMonth, _ = time.ParseInLocation("2006-01-02", time.Now().Format("2006-01")+"-01", time.Local)
	lastOfMonth = firstOfMonth.AddDate(0, 1, -1)
	return
}

func GetLastWeekBeginAndEndDate() (firstOfWeek time.Time, lastOfWeek time.Time) {
	offset := int(time.Sunday - time.Now().Weekday())
	endTimestamp := time.Now().Unix() + 86400*int64(offset)
	lastOfWeek = time.Unix(endTimestamp, 0)
	firstOfWeek = time.Unix(endTimestamp-6*86400, 0)
	return
}

func GetCurrentWeekBeginAndEndDate() (firstOfWeek time.Time, lastOfWeek time.Time) {
	offset := int(time.Saturday - time.Now().Weekday())
	endTimestamp := time.Now().Unix() + 86400*int64(offset+1)
	lastOfWeek = time.Unix(endTimestamp, 0)
	firstOfWeek = time.Unix(endTimestamp-6*86400, 0)
	return
}

// Strtotime strtotime()
// Strtotime("02/01/2016 15:04:05","02/01/2006 15:04:05") == 1451747045
// Strtotime("3 04 PM", "8 41 PM") == -62167144740
func Strtotime(strtime, format string) (int64, error) {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	t, err := time.ParseInLocation(format, strtime, time.Local)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func GetFormatDay(strtime, format string, partday int) string {
	ts, _ := Strtotime(strtime, format)
	tsd := ts + 86400*int64(partday)
	return Date(format, tsd)
}

func ArraryMergeInt64(arr1 []int64, arr2 []int64) []int64 {
	if len(arr1) < 1 {
		return arr2
	}
	if len(arr2) < 1 {
		return arr1
	}
	ret := []int64{}
	for _, v1 := range arr1 {
		ret = append(ret, v1)
	}
	for _, v2 := range arr2 {
		inArr := false
		for _, v3 := range ret {
			if v2 == v3 {
				inArr = true
			}
		}
		if inArr == false {
			ret = append(ret, v2)
		}

	}
	return ret
}

/**
 * 数组去重 去空
 */
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func SafeMobile(mobile string) string {
	if len(mobile) < 11 {
		return mobile
	}
	return mobile[0:3] + "****" + mobile[len(mobile)-4:]
}

//通用输入类型比较
func CheckType(str string, ct string) bool {
	if str == "" || ct == "" {
		return false
	}
	var pattern string
	switch strings.ToLower(ct) {
	case "email":
		pattern = `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	case "phone":
		pattern = `^1\d{10}$`
	case "url":
		pattern = `^(http|https):\/\/.*$`
	case "number":
		pattern = `^\d+$`
	case "md5":
		pattern = `([a-f\d]{32}|[A-F\d]{32})`
	default:
		return false
	}
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(str)
}

/**
修正City参数
*/
func FixCity(city *string) {
	*city = FormatProvince(*city)
	cityArr := []string{"北京", "上海", "天冿", "重庆"}
	if InArray(*city, cityArr) == true {
		*city = ""
	}
}

func ParseDyVideoShare(str string) string {
	return ParseDyShareUrl(str)
}

func ParseDyShareUrl(str string) string {
	//#在抖音，记录美好生活#真的特别好用，用了一个月皮肤变得白皙光滑，姐妹们块安排上#好物推荐 https://v.douyin.com/xg2xVr/ 复制此链接，打开【抖音短视频】，直接观看视频！
	pattern := `https:\/\/v.douyin.com\/.*?\/`
	reg := regexp.MustCompile(pattern)
	dyurl := reg.FindString(str)
	return dyurl
}

func RedShareURL2ID(str string) (ret string) {
	pattern := `\/user\/profile\/([a-z0-9]+)\?*\S*`
	reg := regexp.MustCompile(pattern)
	da := reg.FindAllStringSubmatch(str, -1)
	if len(da) > 0 {
		if len(da[0]) > 1 {
			ret = da[0][1]
		}
		return
	}
	return
}

func ParseTaobaoShare(str string) string {
	//#在抖音，记录美好生活#真的特别好用，用了一个月皮肤变得白皙光滑，姐妹们块安排上#好物推荐 https://v.douyin.com/xg2xVr/ 复制此链接，打开【抖音短视频】，直接观看视频！
	pattern := `https:\/\/m.tb.cn\/\S+\s*`
	reg := regexp.MustCompile(pattern)
	dyurl := reg.FindString(str)
	return dyurl
}

//抖音url解析
func ParseDyVideoUrl(url string) (ret string) {
	//"https://www.iesdouyin.com/share/video/6726733413115579661/?region=CN&mid=6696792007840484099"
	pattern := `\/video\/(\d*?)\/`
	reg := regexp.MustCompile(pattern)
	da := reg.FindAllStringSubmatch(url, -1)
	if len(da) > 0 {
		if len(da[0]) > 1 {
			ret = da[0][1]
		}
		return
	}
	return
}

//抖音作者url解析
func ParseDyAuthorUrl(url string) (ret string) {
	//https://www.iesdouyin.com/share/user/96407975163
	pattern := `\/user\/(\d*?)\?`
	reg := regexp.MustCompile(pattern)
	da := reg.FindAllStringSubmatch(url, -1)
	if len(da) > 0 {
		if len(da[0]) > 1 {
			ret = da[0][1]
		}
		return
	}
	return
}

// 还原京东短链
func ReversedJDShortUrl(url string) string {
	var str string
	err := SimpleCacheGet("reversed:jd:short:url:"+Md5_encode(url), func() interface{} {
		jdaUrl := ExtraJDShortUrl(url)
		client := &http.Client{
			Timeout: time.Second * 5,
		}
		request, _ := http.NewRequest("GET", jdaUrl, nil)
		res, err := client.Do(request)
		if err != nil {
			return err
		}

		finalUrl := res.Request.URL.Scheme + "://" + res.Request.URL.Host + res.Request.URL.RequestURI()
		return finalUrl
	}, &str)
	if err != nil {
		return ""
	}
	return str
}

//短url还原解析
func GetLocation(url string) string {
	url = strings.TrimSpace(url)
	var finalURL string
	err := SimpleCacheGet("str:any:short:url:"+Md5_encode(url), func() interface{} {
		client := &http.Client{}
		request, _ := http.NewRequest("GET", url, nil)
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		retURL := response.Request.URL.String()
		return retURL
	}, &finalURL, 600)
	if err != nil {
		logs.Error("[短链转换] err: %s", err)
		return ""
	}
	return finalURL
}

// 提取京东短链JDA地址
func ExtraJDShortUrl(url string) string {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	request, _ := http.NewRequest("GET", url, nil)
	response, err := client.Do(request)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	pattern := `hrl=\'(https://u.jd.com/jda.*?)\'`
	reg := regexp.MustCompile(pattern)
	matches := reg.FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func FormatStartEnd(data string) (start int, end int) {
	dataArr := strings.Split(data, "-")
	if len(dataArr) >= 1 {
		start = ParseInt(dataArr[0], 0)
		if len(dataArr) == 2 {
			end = ParseInt(dataArr[1], 0)
		}
	}
	return
}

func FormatProvince(province string) string {
	province = strings.Replace(province, "市", "", 1)
	province = strings.Replace(province, "省", "", 1)
	province = strings.Replace(province, "壮族自治区", "", 1)
	province = strings.Replace(province, "回族自治区", "", 1)
	province = strings.Replace(province, "维吾尔自治区", "", 1)
	province = strings.Replace(province, "自治区", "", 1)
	return province
}

func FormatCity(province string) string {
	province = strings.Replace(province, "市", "", 1)
	return province
}

func ExtraFields(items []map[string]interface{}, fields []string) map[string][]interface{} {
	data := make(map[string][]interface{})
	for _, v := range items {
		for _, v1 := range fields {
			data[v1] = append(data[v1], v[v1])
		}
	}
	return data
}

func TaobaoPriceFix(price float64, h5price float64, marketPrice float64) (float64, float64) {
	if marketPrice == 0 {
		marketPrice = price
	}
	if h5price == 0 {
		h5price = price
	}
	minPrice := math.Min(price, math.Min(h5price, marketPrice))
	maxPrice := math.Max(price, math.Max(h5price, marketPrice))
	return minPrice, maxPrice
}

func ExtraStringFields(items []map[string]interface{}, fields []string) map[string][]string {
	data := make(map[string][]string)
	for _, v := range items {
		for _, v1 := range fields {
			if val, ok := v[v1].(string); ok {
				data[v1] = append(data[v1], val)
			}
		}
	}
	return data
}

func BuildMapByField(items []map[string]interface{}, field string) map[interface{}]interface{} {
	data := make(map[interface{}]interface{})
	for _, v := range items {
		data[v[field]] = v
	}
	return data
}

// deprecated
// see CacheHelper
func GetFromCache(name string, source func() interface{}, expire ...int) (interface{}, error) {
	var finalData interface{}
	cacheData := global.Cache.Get(name)
	if cacheData == "" {
		token := GetRandomString(16)
		lock, ok, err := mutex.TryLockWithTimeout(global.Cache.GetInstance().(redis.Conn), name+":lock", token, 10)
		if err != nil {
			return nil, err
		}
		if !ok {
			tryTimes := 30
			for {
				if tryTimes <= 0 {
					break
				}
				//TODO: optimize
				data, err := GetFromCache(name, source, expire...)
				if err == nil {
					return data, err
				}
				tryTimes--
				time.Sleep(time.Microsecond * 100)
			}
		}
		defer lock.Unlock()
		finalData = source()
		if finalData != nil {
			dataSize := 1
			if reflect.TypeOf(finalData).Kind() == reflect.Slice {
				dataSize = Size(finalData)
			}
			if dataSize > 0 {
				timeout := 86400
				if len(expire) > 0 {
					timeout = expire[0]
				}
				err := global.Cache.SetMap(name, finalData, time.Duration(timeout))
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		err := json.Unmarshal([]byte(cacheData), &finalData)
		if err != nil {
			return nil, err
		}
	}
	return finalData, nil
}

func Size(a interface{}) int {
	if reflect.TypeOf(a).Kind() != reflect.Slice {
		return -1
	}
	ins := reflect.ValueOf(a)
	return ins.Len()
}

func ToUint8(val interface{}) uint8 {
	if v, ok := val.(uint8); ok {
		return v
	}
	return 0
}

func ToInt(val interface{}) int {
	switch val.(type) {
	case float64:
		return int(val.(float64))
	case int32:
		return int(val.(int32))
	case int:
		return val.(int)
	case int64:
		return int(val.(int64))
	case uint8:
		return int(val.(uint8))
	case string:
		return ParseInt(val.(string), 0)
	}
	return 0
}

func ToInt64(val interface{}) int64 {
	switch val.(type) {
	case float64:
		return int64(val.(float64))
	case int:
		return int64(val.(int))
	case int32:
		return int64(val.(int32))
	case int64:
		return val.(int64)
	case string:
		return ParseInt64String(val.(string))
	}
	return 0
}

func ToBool(v interface{}) bool {
	if res, ok := v.(bool); ok {
		return res
	}
	return false
}

func BoolToInt(val bool) int {
	if val {
		return 1
	}
	return 0
}

func ToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		return fmt.Sprint(result)
	}
}

func ToStringSlice(v interface{}) []string {
	list := make([]string, 0)
	slice, ok := v.([]interface{})
	if !ok {
		return list
	}
	for _, v := range slice {
		if str, ok := v.(string); ok {
			list = append(list, str)
		}
	}
	return list
}

func ToFloat64(value interface{}) float64 {
	num, err := strconv.ParseFloat(ToString(value), 64)
	if err != nil {
		return 0
	}
	return num
}

//保留两位小数
func FriendlyFloat64(fv float64) float64 {
	val, err := strconv.ParseFloat(fmt.Sprintf("%.2f", fv), 64)
	if err != nil {
		return fv
	}
	return val
}

func FriendlyFloat64One(fv float64) float64 {
	val, err := strconv.ParseFloat(fmt.Sprintf("%.1f", fv), 64)
	if err != nil {
		return fv
	}
	return val
}

func FriendlyFloat64String(fv float64) string {
	return fmt.Sprintf("%.2f", fv)
}

func FriendlyStringFloat64(fv string, delta ...int) string {
	fl, _ := strconv.ParseFloat(fv, 64)
	if len(delta) > 0 {
		fl = fl * float64(delta[0])
	}
	val := fmt.Sprintf("%.2f", fl)
	return val
}

//四舍五入保留一位小数
func RoundAverage(fv float64) float64 {
	val, err := strconv.ParseFloat(fmt.Sprintf("%.1f", fv), 64)
	if err != nil {
		return fv
	}
	return val
}

func StrPaddingLeft(originStr string, length int, paddingStr string) string {
	strLength := len([]rune(originStr))
	if strLength >= length {
		return originStr
	}
	differenceLength := length - strLength
	for i := 0; i < differenceLength; i++ {
		originStr = paddingStr + originStr
	}
	return originStr
}

func CalcPages(totalCount int64, size int) int {
	return int(math.Ceil(float64(totalCount) / float64(size)))
}

func WhichWeek(t time.Time) int {
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}
	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		week = (yearDay-firstWeekDays)/7 + 2
	}
	return week
}

func ToMapStringInterface(v interface{}) (data map[string]interface{}, ok bool) {
	data = make(map[string]interface{})
	val := make(map[string]interface{})
	if val, ok = v.(map[string]interface{}); ok {
		data = val
		return
	}
	return
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func GetProductPrice(val interface{}) string {
	var price float64
	if v, ok := val.(float64); ok {
		price = v / 100
	} else {
		price = float64(ToInt64(val)) / 100
	}
	return fmt.Sprintf("%.2f", price)
}

func JsonDecode(str string, result interface{}) {
	if str == "" {
		return
	}
	err := json.Unmarshal([]byte(str), result)
	if err != nil {
		logs.Error("json decode err:", err)
	}
}

func StringLen(str string) int {
	return len([]rune(str))
}

func UniqueStringSlice(slice []string) []string {
	exists := make(map[interface{}]bool)
	temp := make([]string, 0)
	for _, v := range slice {
		if _, ok := exists[v]; ok {
			continue
		}
		temp = append(temp, v)
		exists[v] = true
	}
	return temp
}

func ExplainNumberInterval(str string, delta int) (firstNum float64, secondNum float64, err error) {
	firstNum = 0
	secondNum = 0
	num := strings.Split(str, "-")
	firstNum = ParseStringFloat64(num[0])

	if len(num) > 1 && num[1] != "" {
		secondNum = ParseStringFloat64(num[1])
	}

	if delta != 1 {
		firstNum *= float64(delta)
		secondNum *= float64(delta)
	}

	return
}

func CalcPageStartEnd(page, size, totalCount int) (pageStart, pageEnd int) {
	pageStart = (page - 1) * size
	pageEnd = pageStart + size
	if pageStart > totalCount {
		pageStart = 0
		pageEnd = 0
		return
	}
	if pageEnd > totalCount {
		pageEnd = totalCount
	}
	return
}

func BuildIntervalCondition(val string, delta ...int) (condition map[string]float64, err error) {
	condition = make(map[string]float64)
	realDelta := 1
	if len(delta) > 0 {
		realDelta = delta[0]
	}
	min, max, err := ExplainNumberInterval(val, realDelta)
	if err != nil {
		return
	}
	if max < min {
		condition["gte"] = min
	} else if min == 0 {
		condition["lt"] = max
	} else {
		condition["gte"] = min
		condition["lt"] = max
	}
	return
}

func ToInterfaceSlice(val interface{}) []interface{} {
	if v, ok := val.([]interface{}); ok {
		return v
	}
	return make([]interface{}, 0)
}

func ProductImageFix(image string, imgType ...string) string {
	return dyimg.Product(image, imgType...)
	//image = strings.Replace(image, "http://", "https://", 1)
	//image = strings.Replace(cache.GetString(image), "tiktokcdn.com", "pstatp.com", 1)
	//if len(imgType) > 0 {
	//	image = strings.Replace(cache.GetString(image), "/obj/", "/"+imgType[0]+"/", 1)
	//}
	//return image
}

func AvatarImageFix(image string, px ...string) string {
	image = strings.Replace(image, "http://", "https://", 1)
	if len(px) > 0 {
		image = strings.Replace(cache.GetString(image), "720x720", px[0], 1)
	}
	return image
}

func FixAuthorNickname(nickname string, markDelete int) string {
	if nickname == "已重置" {
		return nickname
	}
	if markDelete == 1 {
		return nickname + "（已封号）"
	}
	return nickname
}

func DevRun(callback func()) {
	if beego.BConfig.RunMode == "dev" || beego.BConfig.RunMode == "test" {
		callback()
	}
}

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func HasChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

func FixHttps(url string) string {
	return strings.ReplaceAll(url, "http:", "https:")
}

func HBaseColumns(family string, val interface{}) []*hbase.TColumn {
	var fieldArr []string
	switch fields := val.(type) {
	case string:
		fieldArr = strings.Split(fields, ",")
	case []string:
		fieldArr = fields
	}
	tColumns := make([]*hbase.TColumn, 0)
	for _, v := range fieldArr {
		tColumns = append(tColumns, &hbase.TColumn{Family: []byte(family), Qualifier: []byte(strings.Trim(v, " "))})
	}
	return tColumns
}

/*
 * source = 0， 没有获取到销量
 * source = 1， 人工统计销量
 * source = 2， 抖音月销量，按日分区
 * source = 3， 淘宝h5销量，对减
 * source = 4， 抖音近30天销量，对减
 * source = 5， 小店GMV销量来源
 * source = 6， 实时销量来源，实时分摊
 */
func DiagnosisSaleSource(source int) (hasVolume bool, isPredicted bool) {
	switch source {
	case 0:
		// 没获取到销量
		hasVolume = false
		isPredicted = false
	case 1, 2, 5:
		// 有准确销量
		hasVolume = true
		isPredicted = false
	case 3, 4, 6:
		// 有预估销量
		hasVolume = true
		isPredicted = true
	}
	return
}

func CheckDuration(startTime, endTime time.Time, maxDays int64) bool {
	if endTime.Before(startTime) {
		return false
	}
	if maxDays > 0 {
		if endTime.Sub(startTime).Hours() > float64(24*(maxDays+1)) {
			return false
		}
	}
	return true
}

func CalcNatureWeeks(startTime time.Time, endTime time.Time) int {
	offset := int(time.Monday - startTime.Weekday())
	if offset > 0 {
		offset = -6
	}
	startWeekTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	howManyWeek := int(math.Ceil(endTime.Sub(startWeekTime).Hours() / 24 / 7))
	return howManyWeek
}

func CalcNatureMonths(startTime time.Time, endTime time.Time) int {
	startYear := startTime.Year()
	endYear := endTime.Year()
	startMonth, _ := strconv.Atoi(startTime.Format("1"))
	endMonth, _ := strconv.Atoi(endTime.Format("1"))

	result := (endYear-startYear)*12 + (endMonth - startMonth) + 1
	return result
}

func CalcUserValue(roomTicketCount int64, amount float64, totalUser int64) float64 {
	if totalUser <= 0 {
		return 0
	}
	return FriendlyFloat64((float64(roomTicketCount)*0.1 + amount) / float64(totalUser))
}

func GetPromotionIdWithHBaseFormat(promotionId string) []string {
	more := []string{}
	temp := strings.Split(promotionId, ",")
	for _, pid := range temp {
		if pid == "" {
			continue
		}
		more = append(more, pid)
	}
	return more
}

func HidePhoneNum(nickname string) string {
	if len(nickname) == 11 {
		nickname = nickname[:3] + "****" + nickname[7:]
	}
	return nickname
}

func PickTime(timeType int) (startTime, endTime time.Time, err error) {
	today := GetZeroHour(time.Now())
	endTime = today.AddDate(0, 0, 0)
	switch timeType {
	case 1:
		startTime = today.AddDate(0, 0, -1)
	case 7:
		startTime = today.AddDate(0, 0, -7)
	case 15:
		startTime = today.AddDate(0, 0, -15)
	case 30:
		startTime = today.AddDate(0, 0, -30)
	case 90:
		startTime = today.AddDate(0, 0, -90)
	default:
		err = errors.New("time type was wrong")
	}
	return
}

func Percent(val float64, delta ...float64) float64 {
	defaultDelta := float64(100)
	if len(delta) > 0 {
		defaultDelta = delta[0]
	}
	return FriendlyFloat64(val * defaultDelta)
}

func AgeParamFix(age string) string {
	age = strings.ReplaceAll(age, "-18", "<18")
	age = strings.ReplaceAll(age, "50-", ">50")
	return age
}

// 数组去重
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// 校验 min- max 形式的字符串转换成数字是否合法
func CheckStrNumLimit(val string, lowerLimit, upperLimit float64) (min, max float64, err error) {
	if val == "" {
		err = errors.New("val was empty")
		return
	}
	min, max, err = ExplainNumberInterval(val, 1)
	if err != nil {
		return
	}
	// 上限下限校验,范围校验
	if min < lowerLimit || max > upperLimit {
		err = errors.New("min or max out of limit")
		return
	}
	if min > max {
		err = errors.New(" min > max")
	}
	return
}

func Timestamp2Str(ts int) (str string) {
	hours := ts / 3600
	mins := ts % 3600 / 60
	seconds := ts % 60
	if hours != 0 {
		str += fmt.Sprintf("%d时", hours)
	}
	if mins != 0 {
		str += fmt.Sprintf("%d分", mins)
	}
	if seconds != 0 {
		str += fmt.Sprintf("%d秒", seconds)
	}
	return
}

func MatchDouyinNewText(text string) string {
	reg := regexp.MustCompile(`^抖音搜索 (.*?) 复制全文就能搜`)
	matches := reg.FindStringSubmatch(text)
	if len(matches) > 0 {
		return matches[1]
	}
	return text
}

func FormatAppVersion(appVersion, GitCommit, BuildDate string) (string, error) {
	content := `
   Version: {{.Version}}
Go Version: {{.GoVersion}}
Git Commit: {{.GitCommit}}
     Built: {{.BuildDate}}
   OS/ARCH: {{.GOOS}}/{{.GOARCH}}
`
	tpl, err := template.New("version").Parse(content)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, map[string]string{
		"Version":   appVersion,
		"GoVersion": runtime.Version(),
		"GitCommit": GitCommit,
		"BuildDate": BuildDate,
		"GOOS":      runtime.GOOS,
		"GOARCH":    runtime.GOARCH,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), err
}

func GetNowTimeStamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetTimeStamp(t time.Time, layout string) string {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return t.Format(layout)
}

//验证手机号
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func DeserializeData(data string) string {
	if len(data) == 0 {
		return "nil"
	}
	buf, _ := base64.StdEncoding.DecodeString(data)
	zr, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return ""
	}
	dataJson, _ := ioutil.ReadAll(zr)
	return string(dataJson)
}

func SerializeData(data interface{}) string {
	var dataJson []byte
	val, ok := data.(string)
	if ok {
		dataJson = []byte(val)
	} else {
		dataJson, _ = jsoniter.Marshal(data)
	}
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(dataJson)
	if err != nil {
		return ""
	}
	if err := zw.Close(); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// NewRequest 请求包装
func NewRequest(method, url string, data []byte) (body []byte, err error) {

	if method == "GET" {
		url = fmt.Sprint(url, "?", string(data))
		data = nil
	}

	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return body, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return body, err
	}

	return body, err
}

// Struct2Map struct to map，依赖 json tab
func Struct2Map(r interface{}) (s map[string]string, err error) {
	var temp map[string]interface{}
	var result = make(map[string]string)

	bin, err := json.Marshal(r)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(bin, &temp); err != nil {
		return nil, err
	}
	for k, v := range temp {
		result[k], err = ToStringE(v)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// ToStringE interface to string
func ToStringE(i interface{}) (string, error) {
	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatInt(int64(s), 10), nil
	case uint64:
		return strconv.FormatInt(int64(s), 10), nil
	case uint32:
		return strconv.FormatInt(int64(s), 10), nil
	case uint16:
		return strconv.FormatInt(int64(s), 10), nil
	case uint8:
		return strconv.FormatInt(int64(s), 10), nil
	case []byte:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}
