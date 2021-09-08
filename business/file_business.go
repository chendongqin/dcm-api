package business

import (
	"dongchamao/global"
	"dongchamao/models/dcm"
	"dongchamao/models/example"
	"time"
)

type FileBusiness struct {
}

func NewFileBusiness() *FileBusiness {
	return new(FileBusiness)
}

func (f *FileBusiness) InsertFile(fileName, url, tag, mediaId string) (comErr global.CommonError) {
	if _, err := dcm.GetDbSession().Table("exa_file_upload_and_downloads").Insert(example.ExaFileUploadAndDownload{
		Name:      fileName,
		Url:       url,
		Tag:       tag,
		MediaId:   mediaId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return nil
	}
	return
}
