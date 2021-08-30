package dyimg

//根据一定规则替换抖音CDN图片

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"github.com/astaxie/beego/cache"
	"math"
	"net/url"
	"regexp"
	"strings"
)

const (
	AvatarLittle = "200x200"
	AvatarTiny   = "100x100"

	ProductLarge  = "large"
	ProductMedium = "medium"
	ProductThumb  = "thumb"
)

var cdn = []string{
	"https://cdn-images.dongchamao.com",
}

func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func pickCDN(md5 string) string {
	length := len(cdn)
	bytesStr, _ := hex.DecodeString(md5)
	hashValue := BytesToInt64(bytesStr)
	count := hashValue % int64(length)
	count = int64(math.Abs(float64(count))) //md5 128bits > 64bits, 取绝对值
	baseURL := cdn[count]
	return baseURL
}

func buildURL(prefix, source string) string {
	infoURL, err := url.Parse(source)
	if err != nil {
		return source
	}
	//已经转换过的直接返回
	if strings.Contains(infoURL.Host, "dongchamao") {
		return source
	}
	source = strings.Replace(source, "https://", "", 1)
	source = strings.Replace(source, "http://", "", 1)
	return cdn[0] + "/" + source
	md5Str := Md5_encode(source)
	source = url.QueryEscape(source)
	fileName := md5Str + ".jpeg"
	return pickCDN(md5Str) + "/douyin/" + prefix + "/" + fileName + "?source=" + source
}

func Md5_encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Fix(image string) string {
	image = strings.Replace(image, "http://", "https://", 1)
	image = strings.Replace(image, "-ipv6-test", "", 1)
	image = strings.Replace(image, "p5-dy-ipv6", "p3-dy", 1)
	//image = strings.Replace(image, "-ipv6", "", 1)
	image = strings.Replace(image, ".heic", ".jpeg", 1)
	image = WebpToJpg(image)
	return Convert("", Fix(image))
}

func Aweme(image string) string {
	return Convert("aweme", Fix(image))
}

func Live(image string) string {
	return Convert("live", Fix(image))
}

func Convert(prefix string, source string) string {
	return buildURL(prefix, source)
}

func ConvertAvatar(image string, px ...string) string {
	if strings.Contains(image, "dongchamao") {
		return image
	}
	return buildURL("avatar", avatar(image, px...))
}

func avatar(image string, px ...string) string {
	image = Fix(image)
	if len(px) > 0 {
		if strings.Contains(image, "/obj/") {
			if px[0] == AvatarLittle {
				image = strings.Replace(cache.GetString(image), "/obj/", "/medium/", 1)
			} else if px[0] == AvatarTiny {
				image = strings.Replace(cache.GetString(image), "/obj/", "/thumb/", 1)
			}
		} else {
			image = strings.Replace(image, "720x720", px[0], 1)
		}
	}
	return pstatpFix(image)
}

func pstatpFix(image string) string {
	return strings.ReplaceAll(image, "p3.pstatp.com", "p3.douyinpic.com")
}

func Avatar(image string, px ...string) string {
	return ConvertAvatar(image, px...)
}

func Custom(image string, px string) string {
	//已经转换过的直接返回
	if strings.Contains(image, "dongchamao") {
		return image
	}
	image = Fix(image)
	if strings.Contains(image, "~tplv-obj.image") {
		image = strings.Replace(cache.GetString(image), "~tplv-obj.image", "~c5_"+px+".jpeg", 1)
	} else {
		image = strings.Replace(cache.GetString(image), "~tplv-dy-360p.jpeg", "~c5_"+px+".jpeg", 1)
	}
	return Convert("custom", image)
	//return image
}

func Product(image string, imgType ...string) string {
	if strings.Contains(image, "dongchamao") {
		return image
	}
	image = Fix(image)
	image = strings.Replace(cache.GetString(image), "tiktokcdn.com", "pstatp.com", 1)
	image = strings.Replace(image, "sf6-ttcdn-tos.pstatp.com", "p1.pstatp.com", 1)
	re := regexp.MustCompile(`(sf[\d+]-ttcdn-tos\.pstatp\.com)`)
	image = re.ReplaceAllLiteralString(image, "p1.pstatp.com")
	if len(imgType) > 0 {
		if strings.Contains(image, "/obj/") {
			image = strings.Replace(cache.GetString(image), "/obj/", "/"+imgType[0]+"/", 1)
		} else if strings.Contains(image, "/large/") {
			image = strings.Replace(cache.GetString(image), "/large/", "/"+imgType[0]+"/", 1)
		}
	}
	return Convert("product", image)
	//return image
}

func WebpToJpg(image string) string {
	imageLength := len(image)
	if imageLength > 5 && image[imageLength-5:] == ".webp" {
		image = image[:imageLength-5] + ".jpg"
	}
	return image
}
