// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"bytes"
	"crypto/md5"
	"dongchamao/global/cache"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// FileCacheItem is basic unit of file cache adapter.
// it contains data and expire time.
type FileCacheItem struct {
	Data       interface{}
	Lastaccess time.Time
	Expired    time.Time
}

// FileCache Config
var (
	FileCachePath           = "cache"     // cache directory
	FileCacheFileSuffix     = ".bin"      // cache file suffix
	FileCacheDirectoryLevel = 2           // cache file deep level if auto generated cache files.
	FileCacheEmbedExpiry    time.Duration // cache expire time, default is no expire forever.
)

// FileCache is cache adapter for file storage.
type FileCache struct {
	CachePath      string
	FileSuffix     string
	DirectoryLevel int
	EmbedExpiry    int
}

// NewFileCache Create new file cache with no config.
// the level and expiry need set in method StartAndGC as config string.
func NewFileCache() cache.CacheInterface {
	//    return &FileCache{CachePath:FileCachePath, FileSuffix:FileCacheFileSuffix}
	return &FileCache{}
}

// StartAndGC will start and begin gc for file cache.
// the config need to be like {CachePath:"/cache","FileSuffix":".bin","DirectoryLevel":2,"EmbedExpiry":0}
func (fc *FileCache) StartAndGC(config string) error {

	cfg := make(map[string]string)
	json.Unmarshal([]byte(config), &cfg)
	if _, ok := cfg["CachePath"]; !ok {
		cfg["CachePath"] = FileCachePath
	}
	if _, ok := cfg["FileSuffix"]; !ok {
		cfg["FileSuffix"] = FileCacheFileSuffix
	}
	if _, ok := cfg["DirectoryLevel"]; !ok {
		cfg["DirectoryLevel"] = strconv.Itoa(FileCacheDirectoryLevel)
	}
	if _, ok := cfg["EmbedExpiry"]; !ok {
		cfg["EmbedExpiry"] = strconv.FormatInt(int64(FileCacheEmbedExpiry.Seconds()), 10)
	}
	fc.CachePath = cfg["CachePath"]
	fc.FileSuffix = cfg["FileSuffix"]
	fc.DirectoryLevel, _ = strconv.Atoi(cfg["DirectoryLevel"])
	fc.EmbedExpiry, _ = strconv.Atoi(cfg["EmbedExpiry"])

	fc.Init()
	return nil
}

// Init will make new dir for file cache if not exist.
func (fc *FileCache) Init() {
	if ok, _ := exists(fc.CachePath); !ok { // todo : error handle
		_ = os.MkdirAll(fc.CachePath, os.ModePerm) // todo : error handle
	}
}

// get cached file name. it's md5 encoded.
func (fc *FileCache) getCacheFileName(key string) string {
	m := md5.New()
	io.WriteString(m, key)
	keyMd5 := hex.EncodeToString(m.Sum(nil))
	cachePath := fc.CachePath
	switch fc.DirectoryLevel {
	case 2:
		cachePath = filepath.Join(cachePath, keyMd5[0:2], keyMd5[2:4])
	case 1:
		cachePath = filepath.Join(cachePath, keyMd5[0:2])
	}

	if ok, _ := exists(cachePath); !ok { // todo : error handle
		_ = os.MkdirAll(cachePath, os.ModePerm) // todo : error handle
	}

	return filepath.Join(cachePath, fmt.Sprintf("%s%s", keyMd5, fc.FileSuffix))
}

// Get value from file cache.
// if non-exist or expired, return empty string.
func (fc *FileCache) Get(key string) string {
	fileData, err := FileGetContents(fc.getCacheFileName(key))
	if err != nil {
		return ""
	}
	var to FileCacheItem
	GobDecode(fileData, &to)
	if to.Expired.Before(time.Now()) {
		return ""
	}
	return InterfaceToString(to.Data)
}

// Put value into file cache.
// timeout means how long to keep this file, unit of ms.
// if timeout equals FileCacheEmbedExpiry(default is 0), cache this item forever.
func (fc *FileCache) Set(key string, val string, timeout time.Duration) error {
	timeout = timeout*time.Second
	gob.Register(val)
	item := FileCacheItem{Data: val}
	if timeout == FileCacheEmbedExpiry {
		item.Expired = time.Now().Add((86400 * 365 * 10) * time.Second) // ten years
	} else {
		item.Expired = time.Now().Add(timeout)
	}
	item.Lastaccess = time.Now()
	data, err := GobEncode(item)
	if err != nil {
		return err
	}
	return FilePutContents(fc.getCacheFileName(key), data)
}

// Delete file cache value.
func (fc *FileCache) Delete(key string) error {
	filename := fc.getCacheFileName(key)
	if ok, _ := exists(filename); ok {
		return os.Remove(filename)
	}
	return nil
}


func (fc *FileCache) GetMap(key string) interface{}{
	retData := fc.Get(key)
	var retMap interface{}
	if(retData != ""){
		err := json.Unmarshal([]byte(retData), &retMap)
		if(err == nil) {
			return retMap
		}else{
			return nil
		}
	}else {
		return nil
	}
}

func (fc *FileCache) SetMap(key string,val interface{},timeout time.Duration) error{
	jsonStr,err := json.Marshal(val)
	if(err != nil){
		return err
	}else{
		return fc.Set(key,string(jsonStr),timeout)
	}
}

func (fc *FileCache) GetInstance() interface{}{
	return nil
}

// check file exist.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileGetContents Get bytes to file.
// if non-exist, create this file.
func FileGetContents(filename string) (data []byte, e error) {
	return ioutil.ReadFile(filename)
}

// FilePutContents Put bytes to file.
// if non-exist, create this file.
func FilePutContents(filename string, content []byte) error {
	return ioutil.WriteFile(filename, content, os.ModePerm)
}

// GobEncode Gob encodes file cache item.
func GobEncode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// GobDecode Gob decodes file cache item.
func GobDecode(data []byte, to *FileCacheItem) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&to)
}


func InterfaceToString(v interface{})string{
	var ret string
	//fmt.Println(v.(type))
	switch v.(type){
	case []byte:
		ret = string(v.([]byte))
	case string:
		ret = v.(string)
	case []interface{}:
		ret = ""
	case error:
		return ""
	default:
		ret = ""
	}
	return ret
}

func init() {
	cache.Register("file", NewFileCache)
}

