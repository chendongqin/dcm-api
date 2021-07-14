package entity

var DyAuthorRoomMappingMap = HbaseEntity{
	"data": {AJson, "data"},
}

type DyAuthorRoomMapping struct {
	Data []DyAuthorLiveRoom `json:"data"`
}

type DyAuthorLiveRoom struct {
	RoomID     string `json:"room_id"`
	CreateTime int64  `json:"create_time"`
	RoomStatus int    `json:"room_status"`
}
