package business

import (
	"crypto/tls"
	"dongchamao/global/utils"
	dy2 "dongchamao/models/repost/dy"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	AddLiveTopDefault        = 0
	AddLiveTopStar           = 2
	AddLiveTopMonitored      = 4
	AddLiveTopConcerned      = 8
	AddLiveTopAfanti         = 16
	AddLiveTopAfantiPolling  = 32
	AddLiveTopHighLevelStar  = 64
	AddLiveTopSuperLevelStar = 128
	BaseSpiderUrl            = "http://api.spider.dongchamao.cn/"
	LiveSpiderUrl            = "http://dy-live.spider.dongchamao.cn/"
	ZHIMASpiderUrl           = "http://zhima-proxy.spider.dongchamao.cn/"
	AuthorInfoUrl            = "https://webcast-hl.amemv.com/webcast/room/reflow/info/?room_id=70&user_id=%s&live_id=1&app_id=1128"
)

var SpiderNames = map[string]int{
	"author":   1,
	"post":     1,
	"follower": 1,
	"product":  1,
	"month":    1,
	"brand":    1,
	"aweme":    1,
	"comment":  1,
}

var H5UserAgents = []string{
	"Mozilla/5.0 (Linux; Android 6.0.1; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 8.0; Pixel 2 Build/OPD3.170816.012) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15A5341f Safari/604.1",
	"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 10 Build/MOB31T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36",
	"Mozilla/5.0 (Linux; U; Android 4.3; en-us; SM-N900T Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
}

func GetH5UserAgent() string {
	rand.Seed(time.Now().Unix())
	return H5UserAgents[rand.Intn(len(H5UserAgents)-1)]
}

type SpiderBusiness struct {
}

func NewSpiderBusiness() *SpiderBusiness {
	return new(SpiderBusiness)
}

// 爬虫api
// 爬虫抓取加速
// spidername 例如：爬虫名称
// id  123,124  url id
func (s *SpiderBusiness) SpiderSpeedUp(spiderName, id string) string {
	if _, ok := SpiderNames[spiderName]; !ok {
		logs.Error("[爬虫加速] [%s] id:[%s] 推送错误: %s", spiderName, id)
		return ""
	}
	if len(spiderName) == 0 || len(id) == 0 {
		logs.Error("[爬虫加速] [%s] id:[%s] 推送错误: %s", spiderName, id)
		return ""
	}
	pushUrl := BaseSpiderUrl + "crawl?spider=" + spiderName + "&id=" + id
	res := utils.SimpleCurl(pushUrl, "GET", "", "")
	return res
}

// 将达人添加到直播库
// 传入达人id，粉丝数，top 过期时间
func (s *SpiderBusiness) AddLive(authorId string, followerCount int64, top int, expireTime ...int64) {
	if followerCount < 1000 {
		followerCount = 2000 //默认2000 不然爬虫不会抓取
	}
	pushUrl := LiveSpiderUrl + "add_live?uid=" + authorId + "&top=" + utils.ToString(top) + "&follower_count=" + fmt.Sprintf("%d", followerCount)
	if len(expireTime) > 0 {
		pushUrl += "&expire_time=" + strconv.Itoa(int(expireTime[0]))
	}
	data, err := utils.Curl(pushUrl, "GET", "", "")
	if err != nil {
		logs.Error("[直播加速] 达人 [%s] 推送错误: %s", authorId, err)
		return
	}
	code, _ := data.Get("code").Int()
	if code == -1 {
		message, _ := data.Get("msg").String()
		logs.Error("[直播加速] top:[%d] author:[%s] 推送错误: %s", top, authorId, message)
		return
	}
	logs.Info("[直播加速] 达人 [%s] 推送成功 top [%s]", authorId, top)
}

//抖音号搜索
//搜索抖音号 即时接口 先根据抖音号在es中查找，查找不到在调用该方法，调用爬虫实时接口
func (s *SpiderBusiness) GetAuthorByKeyword(keyword string) *dy2.DyAuthorIncome {
	retData := ""
	keyword = url.QueryEscape(keyword)
	pushUrl := BaseSpiderUrl + "searchAuthor?" + "keyword=" + keyword
	for i := 0; i < 5; i++ {
		retData = utils.SimpleCurl(pushUrl, "GET", "", "")
		jd := gjson.Parse(retData)
		if jd.Get("data.nickname").Exists() == false {
			logs.Error("[搜索达人] keyword:[%s] 失败", keyword)
		} else {
			uniqueId := jd.Get("data.unique_id").String()
			if uniqueId == "0" || uniqueId == "" {
				uniqueId = jd.Get("data.short_id").String()
			}
			authorIncome := &dy2.DyAuthorIncome{
				AuthorId:     jd.Get("data.author_id").String(),
				Avatar:       jd.Get("data.avatar").String(),
				Nickname:     jd.Get("data.nickname").String(),
				UniqueId:     uniqueId,
				IsCollection: 0,
			}
			return authorIncome
		}
	}
	return nil
}

// 获取正在直播的达人商品列表
func (s *SpiderBusiness) GetRoomPmt(authorId string) string {
	if len(authorId) == 0 {
		logs.Error("[获取直播间商品列表] author:[%s] 达人id为空", authorId)
		return ""
	}
	pushUrl := LiveSpiderUrl + "room_pmt?uid=" + authorId
	res := utils.SimpleCurl(pushUrl, "GET", "", "")
	return res
}

// AddRecordLive 将达人添加到录制直播、定制大屏
func (s *SpiderBusiness) AddRecordLive(authorId string, level ...int) {
	defLevel := 0
	if len(level) > 0 {
		defLevel = level[0]
	}
	top := utils.ToString(AddLiveTopStar)
	if defLevel == 1 {
		top = utils.ToString(AddLiveTopHighLevelStar) //高级版推64
	}
	pushUrl := LiveSpiderUrl + "add_live?uid=" + authorId + "&top=" + top + "&record=1"
	pushUrl += "&expire_time=" + strconv.Itoa(int(time.Now().AddDate(1, 0, 0).Unix()))

	_, err := utils.Curl(pushUrl, "GET", "", "")
	if err != nil {
		logs.Error("[定制大屏达人加速] 达人 [%s] 推送错误: %s", authorId, err)
		return
	}
	logs.Info("[定制大屏达人加速] 达人 [%s] 推送成功 level [%d] top [%s]", authorId, defLevel, top)
}

//将达人移除录制直播  定制大屏
func (s *SpiderBusiness) DelRecordLive(authorId string, top int) {
	topParam := "&top=-" + utils.ToString(AddLiveTopStar)
	if top == 0 {
		topParam = ""
	}
	pushUrl := LiveSpiderUrl + "add_live?uid=" + authorId + topParam + "&record=0"
	pushUrl += "&expire_time=" + strconv.Itoa(int(time.Now().AddDate(1, 0, 0).Unix()))

	_, err := utils.Curl(pushUrl, "GET", "", "")
	if err != nil {
		logs.Error("[定制大屏达人移除] 达人 [%s] 推送错误: %s", authorId, err)
		return
	}
}

//抓取达人基本信息
func (s *SpiderBusiness) GetAuthorBaseInfo(authorId string) *dy2.DyAuthorIncome {
	pushUrl := fmt.Sprintf(AuthorInfoUrl, authorId)
	headers := map[string]string{fasthttp.HeaderUserAgent: GetH5UserAgent()}
	body := utils.TryDoReq(pushUrl, "", false, nil, headers)
	jd := gjson.ParseBytes(body)
	if jd.Get("status_code").Int() != 0 {
		logs.Error("[搜索达人] author_id:[%s] 失败", authorId)
	} else {
		uniqueId := jd.Get("data.user.display_id").String()
		if uniqueId == "0" || uniqueId == "" {
			uniqueId = jd.Get("data.user.short_id").String()
		}
		authorIncome := &dy2.DyAuthorIncome{
			AuthorId:     jd.Get("data.user.id_str").String(),
			Avatar:       jd.Get("data.user.avatar_thumb.url_list.0").String(),
			Nickname:     jd.Get("data.user.nickname").String(),
			UniqueId:     uniqueId,
			IsCollection: 0,
		}
		return authorIncome
	}
	return nil
}

type spiderRet struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type SpiderDouyinDevice struct {
	AC             string `json:"ac"`
	Carrier        string `json:"carrier"`
	CDID           string `json:"cdid"`
	DeviceModel    string `json:"device_model"`
	DeviceType     string `json:"device_type"`
	DyBuildNumber  string `json:"dy_build_number"`
	DyJSSDKVersion string `json:"dy_js_sdk_version"`
	DyVersion      string `json:"dy_version"`
	IDFA           string `json:"idfa"`
	OpenUDID       string `json:"openudid"`
	OSVersion      string `json:"os_version"`
	Screen         string `json:"screen"`
	UserAgent      string `json:"user_agent"`
	UUID           string `json:"uuid"`
	IID            string `json:"iid"`
	DeviceID       string `json:"device_id"`
	TTReq          string `json:"tt_req"`
}

type SpiderBaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (s *SpiderBusiness) GetProxyTransport() (tr *http.Transport, err error) {
	newIP, _, err := s.NewProxyIP()
	if err != nil {
		return
	}
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("//" + newIP)
	}
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	tr.Proxy = proxy
	return
}

func (s *SpiderBusiness) GetMusicURL(musicID string, proxy bool) (musicURL string) {
	api := "https://www.iesdouyin.com/web/api/v2/music/info/?music_id=" + musicID
	var client *http.Client
	if proxy {
		tr, err := s.GetProxyTransport()
		if err != nil {
			return s.GetMusicURL(musicID, false)
		}
		client = &http.Client{
			Transport: tr,
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "text/html,application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		if utils.IsProxyError(err) {
			return s.GetMusicURL(musicID, false)
		}
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	content, _ := ioutil.ReadAll(resp.Body)
	var urlList []string
	jsoniter.Get(content, "music_info").Get("play_url").Get("url_list").ToVal(&urlList)
	if len(urlList) > 0 {
		musicURL = urlList[0]
	}
	return
}

// NewProxyIP 返回一个新的代理IP
// returns ip, expire timestamp, error
func (s *SpiderBusiness) NewProxyIP() (string, int64, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	api := "http://zhima-acquire-ip.cmm-crawler-intranet.k8s.ajin.me/get_ip?topic=cmm-api&proxy_type=2&size=1"
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return "", 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	type result struct {
		Code int `json:"code"`
		Data []struct {
			Host   string `json:"host"`
			Expire string `json:"expire"`
		} `json:"data"`
	}
	item := &result{}
	err = jsoniter.Unmarshal(body, item)
	if err != nil {
		return "", 0, err
	}
	if len(item.Data) <= 0 {
		return "", 0, errors.New("data size wrong")
	}
	if item.Data[0].Host == "" {
		return "", 0, errors.New("found ip failed")
	}
	expiredTime, _ := time.ParseInLocation(time.RFC3339, item.Data[0].Expire, time.Local)
	return item.Data[0].Host, expiredTime.Unix(), nil
}

const ProxyJson = `
{
  "code": "1",
  "msg": "ok",
  "ret": 0,
  "ret_data": [
    {
      "date": [
        {
          "id": 511100,
          "name": "四川省乐山市"
        },
        {
          "id": 510300,
          "name": "四川省自贡市"
        },
        {
          "id": 511600,
          "name": "四川省广安市"
        },
        {
          "id": 510501,
          "name": "四川省泸州市市辖区"
        },
        {
          "id": 510600,
          "name": "四川省德阳市"
        },
        {
          "id": 511101,
          "name": "四川省乐山市市辖区"
        },
        {
          "id": 510700,
          "name": "四川省绵阳市"
        },
        {
          "id": 511800,
          "name": "四川省雅安市"
        },
        {
          "id": 512000,
          "name": "四川省资阳市"
        },
        {
          "id": 511900,
          "name": "四川省巴中市"
        }
      ],
      "id": 510000,
      "name": "四川省"
    },
    {
      "date": [
        {
          "id": 371200,
          "name": "山东省莱芜市"
        },
        {
          "id": 370600,
          "name": "山东省烟台市"
        },
        {
          "id": 371000,
          "name": "山东省威海市"
        },
        {
          "id": 371700,
          "name": "山东省荷泽市"
        },
        {
          "id": 370100,
          "name": "山东省济南市"
        },
        {
          "id": 370300,
          "name": "山东省淄博市"
        },
        {
          "id": 371300,
          "name": "山东省临沂市"
        },
        {
          "id": 370900,
          "name": "山东省泰安市"
        },
        {
          "id": 370400,
          "name": "山东省枣庄市"
        }
      ],
      "id": 370000,
      "name": "山东省"
    },
    {
      "date": [
        {
          "id": 341800,
          "name": "安徽省宣城市"
        },
        {
          "id": 340700,
          "name": "安徽省铜陵市"
        },
        {
          "id": 340800,
          "name": "安徽省安庆市"
        },
        {
          "id": 341000,
          "name": "安徽省黄山市"
        },
        {
          "id": 341700,
          "name": "安徽省池州市"
        },
        {
          "id": 340600,
          "name": "安徽省淮北市"
        },
        {
          "id": 340100,
          "name": "安徽省合肥市"
        },
        {
          "id": 340400,
          "name": "安徽省淮南市"
        },
        {
          "id": 341802,
          "name": "安徽省宣城市宣州区"
        },
        {
          "id": 340200,
          "name": "安徽省芜湖市"
        },
        {
          "id": 341100,
          "name": "安徽省滁州市"
        },
        {
          "id": 340300,
          "name": "安徽省蚌埠市"
        },
        {
          "id": 341500,
          "name": "安徽省六安市"
        }
      ],
      "id": 340000,
      "name": "安徽省"
    },
    {
      "date": [
        {
          "id": 350600,
          "name": "福建省漳州市"
        },
        {
          "id": 350400,
          "name": "福建省三明市"
        },
        {
          "id": 350900,
          "name": "福建省宁德市"
        },
        {
          "id": 350700,
          "name": "福建省南平市"
        },
        {
          "id": 350500,
          "name": "福建省泉州市"
        },
        {
          "id": 350100,
          "name": "福建省福州市"
        },
        {
          "id": 350300,
          "name": "福建省莆田市"
        }
      ],
      "id": 350000,
      "name": "福建省"
    },
    {
      "date": [
        {
          "id": 321200,
          "name": "江苏省泰州市"
        },
        {
          "id": 321300,
          "name": "江苏省宿迁市"
        },
        {
          "id": 320800,
          "name": "江苏省淮安市"
        },
        {
          "id": 320700,
          "name": "江苏省连云港市"
        },
        {
          "id": 320500,
          "name": "江苏省苏州市"
        },
        {
          "id": 321000,
          "name": "江苏省扬州市"
        },
        {
          "id": 320300,
          "name": "江苏省徐州市"
        },
        {
          "id": 320400,
          "name": "江苏省常州市"
        },
        {
          "id": 320900,
          "name": "江苏省盐城市"
        },
        {
          "id": 321100,
          "name": "江苏省镇江市"
        },
        {
          "id": 320100,
          "name": "江苏省南京市"
        },
        {
          "id": 320600,
          "name": "江苏省南通市"
        }
      ],
      "id": 320000,
      "name": "江苏省"
    },
    {
      "date": [
        {
          "id": 330600,
          "name": "浙江省绍兴市"
        },
        {
          "id": 330200,
          "name": "浙江省宁波市"
        },
        {
          "id": 330700,
          "name": "浙江省金华市"
        },
        {
          "id": 330800,
          "name": "浙江省衢州市"
        },
        {
          "id": 330500,
          "name": "浙江省湖州市"
        },
        {
          "id": 330100,
          "name": "浙江省杭州市"
        },
        {
          "id": 331000,
          "name": "浙江省台州市"
        },
        {
          "id": 330400,
          "name": "浙江省嘉兴市"
        },
        {
          "id": 330300,
          "name": "浙江省温州市"
        },
        {
          "id": 331100,
          "name": "浙江省丽水市"
        },
        {
          "id": 330900,
          "name": "浙江省舟山市"
        }
      ],
      "id": 330000,
      "name": "浙江省"
    },
    {
      "date": [
        {
          "id": 210200,
          "name": "辽宁省大连市"
        },
        {
          "id": 210300,
          "name": "辽宁省鞍山市"
        },
        {
          "id": 210800,
          "name": "辽宁省营口市"
        },
        {
          "id": 210701,
          "name": "辽宁省锦州市市辖区"
        },
        {
          "id": 211000,
          "name": "辽宁省辽阳市"
        },
        {
          "id": 211400,
          "name": "辽宁省葫芦岛市"
        },
        {
          "id": 210600,
          "name": "辽宁省丹东市"
        }
      ],
      "id": 210000,
      "name": "辽宁省"
    },
    {
      "date": [
        {
          "id": 420800,
          "name": "湖北省荆门市"
        },
        {
          "id": 420100,
          "name": "湖北省武汉市"
        },
        {
          "id": 420801,
          "name": "湖北省荆门市市辖区"
        },
        {
          "id": 421300,
          "name": "湖北省随州市"
        }
      ],
      "id": 420000,
      "name": "湖北省"
    },
    {
      "date": [
        {
          "id": 500300,
          "name": "重庆市市"
        }
      ],
      "id": 500000,
      "name": "重庆市"
    },
    {
      "date": [
        {
          "id": 411200,
          "name": "河南省三门峡市"
        }
      ],
      "id": 410000,
      "name": "河南省"
    },
    {
      "date": [
        {
          "id": 220500,
          "name": "吉林省通化市"
        },
        {
          "id": 220600,
          "name": "吉林省白山市"
        },
        {
          "id": 222400,
          "name": "吉林省延边朝鲜族自治州"
        },
        {
          "id": 220400,
          "name": "吉林省辽源市"
        },
        {
          "id": 220300,
          "name": "吉林省四平市"
        },
        {
          "id": 220800,
          "name": "吉林省白城市"
        },
        {
          "id": 220700,
          "name": "吉林省松原市"
        }
      ],
      "id": 220000,
      "name": "吉林省"
    },
    {
      "date": [
        {
          "id": 360322,
          "name": "江西省上栗县"
        },
        {
          "id": 361101,
          "name": "江西省上饶市市辖区"
        },
        {
          "id": 360500,
          "name": "江西省新余市"
        },
        {
          "id": 360900,
          "name": "江西省宜春市"
        },
        {
          "id": 360300,
          "name": "江西省萍乡市"
        },
        {
          "id": 360800,
          "name": "江西省吉安市"
        },
        {
          "id": 360100,
          "name": "江西省南昌市"
        }
      ],
      "id": 360000,
      "name": "江西省"
    },
    {
      "date": [
        {
          "id": 130200,
          "name": "河北省唐山市"
        },
        {
          "id": 130500,
          "name": "河北省邢台市"
        }
      ],
      "id": 130000,
      "name": "河北省"
    },
    {
      "date": [
        {
          "id": 530400,
          "name": "云南省玉溪市"
        },
        {
          "id": 532301,
          "name": "云南省楚雄市"
        },
        {
          "id": 532600,
          "name": "云南省文山壮族苗族自治州"
        },
        {
          "id": 532900,
          "name": "云南省大理白族自治州"
        },
        {
          "id": 532500,
          "name": "云南省红河哈尼族彝族自治州"
        },
        {
          "id": 533300,
          "name": "云南省怒江傈僳族自治州"
        },
        {
          "id": 530300,
          "name": "云南省曲靖市"
        },
        {
          "id": 530821,
          "name": "云南省普洱哈尼族彝族自治县"
        },
        {
          "id": 530901,
          "name": "云南省临沧市市辖区"
        },
        {
          "id": 530700,
          "name": "云南省丽江市"
        },
        {
          "id": 533421,
          "name": "云南省香格里拉县"
        },
        {
          "id": 530500,
          "name": "云南省保山市"
        }
      ],
      "id": 530000,
      "name": "云南省"
    },
    {
      "date": [
        {
          "id": 140100,
          "name": "山西省太原市"
        },
        {
          "id": 140500,
          "name": "山西省晋城市"
        }
      ],
      "id": 140000,
      "name": "山西省"
    },
    {
      "date": [
        {
          "id": 232723,
          "name": "黑龙江省漠河县"
        },
        {
          "id": 231000,
          "name": "黑龙江省牡丹江市"
        }
      ],
      "id": 230000,
      "name": "黑龙江省"
    },
    {
      "date": [
        {
          "id": 150600,
          "name": "内蒙古鄂尔多斯市"
        },
        {
          "id": 150200,
          "name": "内蒙古包头市"
        }
      ],
      "id": 150000,
      "name": "内蒙古自治区"
    },
    {
      "date": [
        {
          "id": 441900,
          "name": "广东省东莞市"
        },
        {
          "id": 445200,
          "name": "广东省揭阳市"
        },
        {
          "id": 442000,
          "name": "广东省中山市"
        },
        {
          "id": 440100,
          "name": "广东省广州市"
        },
        {
          "id": 440700,
          "name": "广东省江门市"
        },
        {
          "id": 440200,
          "name": "广东省韶关市"
        },
        {
          "id": 440500,
          "name": "广东省汕头市"
        },
        {
          "id": 440900,
          "name": "广东省茂名市"
        }
      ],
      "id": 440000,
      "name": "广东省"
    },
    {
      "date": [
        {
          "id": 110105,
          "name": "北京市朝阳区"
        }
      ],
      "id": 110000,
      "name": "北京市"
    },
    {
      "date": [
        {
          "id": 610100,
          "name": "陕西省西安市"
        },
        {
          "id": 610700,
          "name": "陕西省汉中市"
        },
        {
          "id": 610800,
          "name": "陕西省榆林市"
        }
      ],
      "id": 610000,
      "name": "陕西省"
    },
    {
      "date": [
        {
          "id": 430300,
          "name": "湖南省湘潭市"
        },
        {
          "id": 431000,
          "name": "湖南省郴州市"
        },
        {
          "id": 430600,
          "name": "湖南省岳阳市"
        }
      ],
      "id": 430000,
      "name": "湖南省"
    },
    {
      "date": [
        {
          "id": 310112,
          "name": "上海市闵行区"
        }
      ],
      "id": 310000,
      "name": "上海市"
    },
    {
      "date": [
        {
          "id": 640400,
          "name": "宁夏固原市"
        }
      ],
      "id": 640000,
      "name": "宁夏回族自治区"
    }
  ],
  "timestamp": 1611197609
}
`

type ProxyData struct {
	ID   int         `json:"id"`
	Name string      `json:"name"`
	Date []ProxyData `json:"date"`
}

func (s *SpiderBusiness) GetProxyList() (data []ProxyData) {
	data = make([]ProxyData, 0)

	jsoniter.Get([]byte(ProxyJson), "ret_data").ToVal(&data)
	return data
}

//NewProxyIPWithRequestIP 通过请求的ip获取对应的代理ip
func (s *SpiderBusiness) NewProxyIPWithRequestIP(requestIP string) (
	proxyIP string, expTime int64, provinceCode int, err error,
) {
	proxyDataList := s.GetProxyList()

	province := "福建"
	provinceCode = 350000

	data := utils.SimpleCurl("http://api.shike.ddashi.com/ipLocation?ip="+requestIP, "GET", "", "")
	province = jsoniter.Get([]byte(data), "data", "province").ToString()

	for _, p := range proxyDataList {
		if strings.Contains(p.Name, province) == true {
			province = p.Name
			provinceCode = p.ID
			break
		}
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	api := ZHIMASpiderUrl + "get_city_ip?proxy_type=0&city=" + utils.ToString(provinceCode) + "&size=2"
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	type result struct {
		Code int `json:"code"`
		Data []struct {
			Host   string `json:"host"`
			Expire string `json:"expire"`
		} `json:"data"`
	}
	item := &result{}
	err = jsoniter.Unmarshal(body, item)
	if err != nil {
		return
	}
	if len(item.Data) <= 0 {
		err = errors.New("data size wrong")
		return
	}
	if item.Data[0].Host == "" {
		err = errors.New("found ip failed")
		return
	}
	expiredTime, _ := time.ParseInLocation(time.RFC3339, item.Data[0].Expire, time.Local)
	proxyIP = item.Data[0].Host
	expTime = expiredTime.Unix()
	return
}
