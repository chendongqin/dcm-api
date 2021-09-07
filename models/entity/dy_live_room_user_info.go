package entity

var DyLiveRoomUserInfoMap = HbaseEntity{
	"roomId":           {String, "room_id"},
	"gender":           {Json, "gender"},
	"province":         {Json, "province"},
	"city":             {Json, "city"},
	"word":             {Json, "word"},
	"ageDistrinbution": {Json, "age_distrinbution"},
}

type DyLiveRoomUserInfo struct {
	RoomId           string           `json:"room_id"`
	Gender           map[string]int64 `json:"gender"`
	Province         map[string]int64 `json:"province"`
	City             map[string]int64 `json:"city"`
	AgeDistrinbution map[string]int64 `json:"age_distrinbution"`
	Word             map[string]int64 `json:"word"`
}
