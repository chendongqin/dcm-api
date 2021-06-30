package hbaseService
//
//import (
//	douyinmodelsV2 "douyin-api/models/douyinv2"
//	"fmt"
//)
//
//type HbaseThriftService struct {
//}
//
//
//
//func NewHbaseThriftService() *HbaseThriftService{
//	return new(HbaseThriftService)
//}
//
//func (this *HbaseThriftService)Get(){

	//defaultCtx :=
	//protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//trans, err := thrift.NewTHttpClient(HOST)
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "error resolving address:", err)
	//	os.Exit(1)
	//}
	//// 设置用户名密码
	//httClient := trans.(*thrift.THttpClient)
	//httClient.SetHeader("ACCESSKEYID", USER)
	//httClient.SetHeader("ACCESSSIGNATURE", PASSWORD)
	//client := hbase.NewTHBaseServiceClientFactory(trans, protocolFactory)
	//
	//if err := trans.Open(); err != nil {
	//	fmt.Fprintln(os.Stderr, "Error opening "+HOST, err)
	//	os.Exit(1)
	//}()

	//client := global.HbasePools.Get()
	//defer client.Close()
	//
	//
	//timenow := time.Now().UnixNano()
	//fmt.Println(timenow)
	//tableInbytes := []byte("dy:aweme")
	//result, _ := client.Get(context.Background(), tableInbytes, &hbase.TGet{Row: []byte("8a97da082497e1562c23e44c792741363058e2a9")})
	//
	//fmt.Println((time.Now().UnixNano()-timenow)/1000000)
	////fmt.Println("Get result:")
	////fmt.Println(result)
	////#fmt.Println(err)
	//fmt.Println(string(result.Row))
	//fmt.Println(result.ColumnValues)
	//fmt.Println(result.Stale)
	//fmt.Println(result.Partial)
	//for _,v := range result.ColumnValues{
	//	//fmt.Print(string(v.Family))
	//	fmt.Print("field:"+string(v.Qualifier))
	//	//fmt.Print(v.Value)
	//	fmt.Print("value:"+string(v.Value))
	//	//fmt.Print(v.Timestamp)
	//	//fmt.Print(string(v.Tags))
	//	//fmt.Print(v.Type)
	//	fmt.Println("")
	//}
	//
	//return
//	awemeId := "6684182690855472391"
//	dh := douyinmodelsV2.NewDyHbaseModel()
//	awemeInfo ,_ := dh.GetAwemeById(awemeId)
//	fmt.Println(awemeInfo)
//
//	authorId := "62782088590"
//	authorInfo ,_ := dh.GetAuthorById(authorId)
//	fmt.Println(authorInfo)
//	//
//	//authorId := "88120071564"
//	//
//	//authorInfo ,_ := dh.GetAuthorById(authorId)
//	//fmt.Println(authorInfo)
//
//	//awemeId := "6733030004751371531"
//	//dh := douyinmodelsV2.NewDyCommentModel()
//	//_ = dh.GetCommentByAwemeId(awemeId,false)
//	//fmt.Println(dh)
//	//
//	//fmt.Println(dh)
//	//fmt.Println(dh.Author.UpdateTime)
//	//dh := douyinmodelsV2.NewDyHbaseModel()
//	//ret,_ := dh.GetCommentByAwemeId(awemeId)
//	//fmt.Println(ret)
//
//
//}