package hbaseService

//import (
//	"context"
//	"fmt"
//	"github.com/tsuna/gohbase"
//	"github.com/tsuna/gohbase/hrpc"
//)
//
//type HbaseService struct {
//}
//
//func NewHbaseService() *HbaseService{
//
//
//	return new(HbaseService)
//}
//
//func (this *HbaseService)Get(){
//	option := gohbase.EffectiveUser("root")
//
//	client := gohbase.NewClient("ld-uf6w1y03mb950v52e-proxy-hbaseue-pub.hbaseue.rds.aliyuncs.com:30020",option)
//	fmt.Println(1)
//	getRequest, err := hrpc.NewGetStr(context.Background(), "table", "row")
//	fmt.Println(2)
//	getRsp, err := client.Get(getRequest)
//	fmt.Println(getRsp,err)
//}