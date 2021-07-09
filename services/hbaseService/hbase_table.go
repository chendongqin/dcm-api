package hbaseService

const (
	TableFamily = "r"

	HbaseDyLiveInfo   = "dy_live_info"  //直播间信息
	HbaseDyLivePmt    = "dy_live_pmt"   //直播pmt(直播商品列表)
	HbaseDyReputation = "dy_reputation" //带货口碑

	HbaseDyAweme                        = "dy_aweme"                            //抖音视频
	HbaseDyAuthorAwemeAgg               = "dy_author_aweme_agg"                 //达人抖音视频
	HbaseDyAwemeDiggCommentForwardCount = "dy_aweme_digg_comment_forward_count" //达人抖音视频

	HbaseDyAuthor       = "dy_author"        //抖音用户信息
	HbaseDyAuthorBasic  = "dy_author_basic"  //抖音用户基本数据信息
	HbaseDyAuthorFans   = "dy_author_fans"   //抖音粉丝信息
	HbaseDyLiveFansClub = "dy_live_fansclub" //抖音粉丝团数据信息
	HbaseXtAuthorDetail = "xt_author_detail" //星图达人详情

)
