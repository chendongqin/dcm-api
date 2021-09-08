package example

import "time"

type ExaFileUploadAndDownload struct {
	Id        int       `xorm:"not null pk INT(11)"`
	Name      string    `json:"name" gorm:"comment:文件名"`        // 文件名
	Url       string    `json:"url" gorm:"comment:文件地址"`        // 文件地址
	Tag       string    `json:"tag" gorm:"comment:文件标签"`        // 文件标签
	Key       string    `json:"key" gorm:"comment:编号"`          // 编号
	MediaId   string    `json:"media_id" gorm:"comment:微信媒体id"` // 微信媒体id
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}
