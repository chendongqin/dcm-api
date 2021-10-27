package command

import (
	"dongchamao/business"
	"fmt"
	"strconv"
	"time"
)

//监听除了小时榜的x:30榜单
func CheckRank() {
	currentHour := time.Now().Hour()
	currentHourString := strconv.Itoa(currentHour)
	//currentHourString := "15"
	currentHourString = fmt.Sprintf("%s:30", currentHourString)
	business.NewRankBusiness().LoopCheck(currentHourString)
	return
}

//监控整点榜单
func CheckGoodsRank() {
	currentHour := time.Now().Hour()
	currentHourString := strconv.Itoa(currentHour)
	//currentHourString := "15"
	business.NewRankBusiness().LoopCheck(currentHourString)
	return
}

//监听小时榜
func CheckRankHour() {
	currentHourString := "every"
	business.NewRankBusiness().LoopCheck(currentHourString)
	return
}
