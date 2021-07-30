package es

import (
	"dongchamao/models/es"
	"time"
)

type EsAuthorBusiness struct {
}

func NewEsAuthorBusiness() *EsAuthorBusiness {
	return new(EsAuthorBusiness)
}

func (receiver *EsAuthorBusiness) AuthorProductAnalysis(authorId string, startTime, stopTime time.Time) (list []es.EsDyAuthorProductAnalysis) {
	return
}
