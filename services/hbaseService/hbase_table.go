package hbaseService

const (
	TableFamily = "r"

	HbaseDyLiveInfo                 = "dy_live_info"              //直播间信息
	HbaseDyLiveReputation           = "dy_live_reputation"        //直播间带货口碑信息
	HbaseDyAuthorLiveSalesData      = "dy_author_live_sales_data" //直播间带货数据
	HbaseDyLivePmt                  = "dy_live_pmt"               //直播pmt(直播商品列表)
	HbaseDyLiveCurProduct           = "dy_live_cur_product"
	HbaseDyLiveRoomUserInfo         = "dy_live_room_user_info"         //直播粉丝数据
	HbaseDyRoomProduct              = "dy_room_product"                //直播间商品全网销量
	HbaseDyLiveRankTrend            = "dy_live_rank_trend"             //直播间榜单排名数据
	HbaseDyReputation               = "dy_reputation"                  //带货口碑
	HbaseDyAuthorProductDateMapping = "dy_author_product_date_mapping" //达人带货

	HbaseDyAweme                        = "dy_aweme"                            //抖音视频
	HbaseDyAuthorAwemeAgg               = "dy_author_aweme_agg"                 //达人抖音视频
	HbaseDyAwemeDiggCommentForwardCount = "dy_aweme_digg_comment_forward_count" //达人抖音视频
	HbaseDyProductAwemeSalesTrend       = "dy_product_aweme_sales_trend"        //商品视频分销趋势

	HbaseDyAuthor                 = "dy_author"       //抖音用户信息
	HbaseDyAuthorBasic            = "dy_author_basic" //抖音用户基本数据信息
	HbaseDyAuthorFans             = "dy_author_fans"  //抖音粉丝信息
	HbaseDyAuthorLiveFansClubUser = "dy_live_fansclub_user"
	HbaseDyLiveFansClub           = "dy_live_fansclub"           //抖音粉丝团数据信息
	HbaseDyLiveChatMessage        = "dy_live_chat_message"       //直播间弹幕
	HbaseDyAuthorRoomMapping      = "dy_author_room_mapping"     //达人直播间
	HbaseDyAuthorProductAnalysis  = "dy_author_product_analysis" //达人电商分析
	HbaseXtAuthorDetail           = "xt_author_detail"           //星图达人详情
	HbaseDyAuthorLiveTags         = "dy_live_hour_rank_sell"     //达人带货行业

	HbaseDyProduct                     = "dy_product"
	HbaseDyProductBrand                = "dy_product_brand"
	HbaseDyProductDaily                = "dy_product_daily"
	HbaseDyProductAwemeDailyDistribute = "dy_product_aweme_daily_distribute"
	HbaseDyProductLiveSalesTrend       = "dy_product_live_sales_trend"
	HbaseDyLivePromotionMonth          = "dy_live_promotion_month"
	HbaseDyProductAuthorAnalysis       = "dy_product_author_analysis"
	HbaseDyProductAwemeAuthorAnalysis  = "dy_product_aweme_author_analysis"
	HbaseAdsDyProductGpmDi             = "ads_dy_product_gpm_di" //gpm每日数据

	HbaseXtHotAwemeAuthorRank     = "xt_hot_aweme_author"
	HbaseXtHotLiveAuthorRank      = "xt_hot_live_author"
	HbaseDyLiveHourRank           = "dy_live_hour_rank"
	HbaseDyLiveHourRankSell       = "dy_live_hour_rank_sell_top"
	HbaseDyLiveHourRankPopularity = "dy_live_hour_rank_popularity"
	HbaseDyLiveTopRank            = "dy_live_top"
	HbaseDyLiveShareWeekRank      = "dy_live_share_top"
	HbaseDyAwemeShareRank         = "dy_aweme_share_top"

	HbaseDyAwemeTopComment = "dws_dy_aweme_ext_topcomment_di"
	//todo 修改表名
	HbaseDyProductTopComment = "dws_dy_aweme_ext_topcomment_di"

	HbaseDyShop       = "dy_shop"        //小店
	HbaseDyShopDetail = "dy_shop_detail" //小店
)
