package wrapper

//type Bar struct {
//	Symbol    string  `json:"symbol"`
//	Bob       float64 `json:"bob"`       ///Bar的开始时间
//	Eob       float64 `json:"eob"`       ///Bar的结束时间
//	Open      float32 `json:"open"`      ///<开盘价
//	Close     float32 `json:"close"`     ///<收盘价
//	High      float32 `json:"high"`      ///<最高价
//	Low       float32 `json:"low"`       ///<最低价
//	Volume    float64 `json:"volume"`    ///<成交量
//	Amount    float64 `json:"amount"`    ///<成交金额
//	PreClose  float32 `json:"pre_close"` ///昨收盘价，只有日频数据赋值
//	Position  int64   `json:"position"`  ///<持仓量
//	Frequency string  `json:"frequency"` ///Bar频度
//}
//
//type Quote struct {
//	BidPrice  float32 `json:"bid_price"`  ///本档委买价
//	BidVolume int64   `json:"bid_volume"` ///本档委买量
//	AskPrice  float32 `json:"ask_price"`  ///本档委卖价
//	AskVolume int64   `json:"ask_volume"` ///本档委卖量
//}
//
//type Tick struct {
//	Symbol      string  `json:"symbol"`
//	CreatedAt   float64 `json:"created_at"`   ///<Utc时间，精确到毫秒
//	Price       float32 `json:"price"`        ///<最新价
//	Open        float32 `json:"open"`         ///<开盘价
//	High        float32 `json:"high"`         ///<最高价
//	Low         float32 `json:"low"`          ///<最低价
//	CumVolume   float64 `json:"cum_volume"`   ///<成交总量
//	CumAmount   float64 `json:"cum_amount"`   ///<成交总金额/最新成交额,累计值
//	CumPosition int64   `json:"cum_position"` ///<合约持仓量(期),累计值
//	LastAmount  float64 `json:"last_amount"`  ///<瞬时成交额
//	LastVolume  int32   `json:"last_volume"`  ///<瞬时成交量
//	TradeType   int32   `json:"trade_type"`   ///(保留)交易类型,对应多开,多平等类型
//	Quotes      []Quote `json:"quotes"`       ///报价, 下标从0开始，0-表示第一档，1-表示第二档，依次类推
//}
//
//type L2Transaction struct {
//	Symbol        string  `json:"symbol"`
//	CreatedAt     float64 `json:"created_at"`      ///成交时间，Utc时间
//	Price         float32 `json:"price"`           ///成交价
//	Volume        int64   `json:"volume"`          ///成交量
//	Side          byte    `json:"side"`            ///内外盘标记
//	ExecType      byte    `json:"exec_type"`       ///成交类型
//	ExecIndex     int64   `json:"exec_index"`      ///成交编号
//	AskOrderIndex int64   `json:"ask_order_index"` ///叫卖委托编号
//	BidOrderIndex int64   `json:"bid_order_index"` ///叫买委托编号
//}
//
//type L2Order struct {
//	Symbol     string  `json:"symbol"`
//	CreatedAt  float64 `json:"created_at"`  ///委托时间，Utc时间
//	Price      float32 `json:"price"`       ///委托价
//	Volume     int64   `json:"volume"`      ///委托量
//	Side       byte    `json:"side"`        ///买卖方向
//	OrderType  byte    `json:"order_type"`  ///委托类型
//	OrderIndex int64   `json:"order_index"` ///委托编号
//}
//
//type L2OrderQueue struct {
//	Symbol       string  `json:"symbol"`
//	CreatedAt    float64 `json:"created_at"`    ///行情时间，Utc时间
//	Price        float32 `json:"price"`         ///最优委托价
//	Volume       int64   `json:"volume"`        ///委托量
//	Side         byte    `json:"side"`          ///买卖方向
//	QueueOrders  int32   `json:"queue_orders"`  ///委托量队列中元素个数(最多50)
//	QueueVolumes []int32 `json:"queue_volumes"` ///委托量队列(最多50个，有可能小于50, 有效数据长度取决于queueOrders)
//}
//
//type Order struct {
//	StrategyId  string `json:"strategy_id"`  //策略ID
//	AccountId   string `json:"account_id"`   //账号ID
//	AccountName string `json:"account_name"` //账户登录名
//
//	ClOrdId       string `json:"cl_ord_id"`      //委托客户端ID
//	OrderId       string `json:"order_id"`       //委托柜台ID
//	ExOrdId       string `json:"ex_ord_id"`      //委托交易所ID
//	AlgoOrderId   string `json:"algo_order_id"`  //算法母单ID
//	OrderBusiness int32  `json:"order_business"` //业务类型
//
//	Symbol         string `json:"symbol"`          //Symbol
//	Side           int32  `json:"side"`            //买卖方向，取值参考enum OrderSide
//	PositionEffect int32  `json:"position_effect"` //开平标志，取值参考enum PositionEffect
//	PositionSide   int32  `json:"position_side"`   //持仓方向，取值参考enum PositionSide
//
//	OrderType      int32 `json:"order_type"`      //委托类型，取值参考enum OrderType
//	OrderDuration  int32 `json:"order_duration"`  //委托时间属性，取值参考enum OrderDuration
//	OrderQualifier int32 `json:"order_qualifier"` //委托成交属性，取值参考enum OrderQualifier
//	OrderSrc       int32 `json:"order_src"`       //委托来源，取值参考enum OrderSrc
//	PositionSrc    int32 `json:"position_src"`    //头寸来源（仅适用融资融券），取值参考 Enum PositionSrc
//
//	Status             int32  `json:"status"`                //委托状态，取值参考enum OrderStatus
//	OrdRejReason       int32  `json:"ord_rej_reason"`        //委托拒绝原因，取值参考enum OrderRejectReason
//	OrdRejReasonDetail string `json:"ord_rej_reason_detail"` //委托拒绝原因描述
//
//	Price     float64 `json:"price"`      //委托价格
//	StopPrice float64 `json:"stop_price"` //委托止损/止盈触发价格
//
//	OrderStyle    int32   `json:"order_style"`    //委托风格，取值参考 Enum OrderStyle
//	Volume        int64   `json:"volume"`         //委托量
//	Value         float64 `json:"value"`          //委托额
//	Percent       float64 `json:"percent "`       //委托百分比
//	TargetVolume  int64   `json:"target_volume"`  //委托目标量
//	TargetValue   float64 `json:"target_value"`   //委托目标额
//	TargetPercent float64 `json:"target_percent"` //委托目标百分比
//
//	FilledVolume     int64   `json:"filled_volume"`     //已成量
//	FilledVWap       float64 `json:"filled_vwap"`       //已成均价
//	FilledAmount     float64 `json:"filled_amount"`     //已成金额
//	FilledCommission float64 `json:"filled_commission"` //已成手续费
//
//	CreatedAt int64 `json:"created_at"` //委托创建时间
//	UpdatedAt int64 `json:"updated_at"` //委托更新时间
//}
//
//type AlgoOrder struct {
//	StrategyId  string `json:"strategy_id"`  //策略ID
//	AccountId   string `json:"account_id"`   //账号ID
//	AccountName string `json:"account_name"` //账户登录名
//
//	ClOrdId       string `json:"cl_ord_id"`      //委托客户端ID
//	OrderId       string `json:"order_id"`       //委托柜台ID
//	ExOrdId       string `json:"ex_ord_id"`      //委托交易所ID
//	OrderBusiness int32  `json:"order_business"` //业务类型
//
//	Symbol         string `json:"symbol"`          //Symbol
//	Side           int32  `json:"side"`            //买卖方向，取值参考enum OrderSide
//	PositionEffect int32  `json:"position_effect"` //开平标志，取值参考enum PositionEffect
//	PositionSide   int32  `json:"position_side"`   //持仓方向，取值参考enum PositionSide
//
//	OrderType      int32 `json:"order_type"`      //委托类型，取值参考enum OrderType
//	OrderDuration  int32 `json:"order_duration"`  //委托时间属性，取值参考enum OrderDuration
//	OrderQualifier int32 `json:"order_qualifier"` //委托成交属性，取值参考enum OrderQualifier
//	OrderSrc       int32 `json:"order_src"`       //委托来源，取值参考enum OrderSrc
//	PositionSrc    int32 `json:"position_src"`    //头寸来源（仅适用融资融券），取值参考 Enum PositionSrc
//
//	Status             int32  `json:"status"`                //委托状态，取值参考enum OrderStatus
//	OrdRejReason       int32  `json:"ord_rej_reason"`        //委托拒绝原因，取值参考enum OrderRejectReason
//	OrdRejReasonDetail string `json:"ord_rej_reason_detail"` //委托拒绝原因描述
//
//	Price     float64 `json:"price"`      //委托价格
//	StopPrice float64 `json:"stop_price"` //委托止损/止盈触发价格
//
//	OrderStyle    int32   `json:"order_style"`    //委托风格，取值参考 Enum OrderStyle
//	Volume        int64   `json:"volume"`         //委托量
//	Value         float64 `json:"value"`          //委托额
//	Percent       float64 `json:"percent"`        //委托百分比
//	TargetVolume  int64   `json:"target_volume"`  //委托目标量
//	TargetValue   float64 `json:"target_value"`   //委托目标额
//	TargetPercent float64 `json:"target_percent"` //委托目标百分比
//
//	FilledVolume     int64   `json:"filled_volume"`     //已成量
//	FilledVWap       float64 `json:"filled_vwap"`       //已成均价
//	FilledAmount     float64 `json:"filled_amount"`     //已成金额
//	FilledCommission float64 `json:"filled_commission"` //已成手续费
//
//	AlgoName    string `json:"algo_name"`    //算法策略名
//	AlgoParam   string `json:"algo_param"`   //算法策略参数
//	AlgoStatus  int32  `json:"algo_status"`  //算法策略状态,仅作为AlgoOrder Pause请求入参，取值参考 Enum AlgoOrderStatus
//	AlgoComment string `json:"algo_comment"` //算法单备注
//
//	CreatedAt int64 `json:"created_at"` //委托创建时间
//	UpdatedAt int64 `json:"updated_at"` //委托更新时间
//}
//
//type ExecRpt struct {
//	StrategyId  string `json:"strategy_id"`  //策略ID
//	AccountId   string `json:"account_id"`   //账号ID
//	AccountName string `json:"account_name"` //账户登录名
//
//	ClOrdId string `json:"cl_ord_id"` //委托客户端ID
//	OrderId string `json:"order_id"`  //委托柜台ID
//	ExecId  string `json:"exec_id"`   //委托回报ID
//
//	Symbol string `json:"symbol"` //Symbol
//
//	PositionEffect     int32  `json:"position_effect"`       //开平标志，取值参考enum PositionEffect
//	Side               int32  `json:"side"`                  //买卖方向，取值参考enum OrderSide
//	OrdRejReason       int32  `json:"ord_rej_reason"`        //委托拒绝原因，取值参考enum OrderRejectReason
//	OrdRejReasonDetail string `json:"ord_rej_reason_detail"` //委托拒绝原因描述
//	ExecType           int32  `json:"exec_type"`             //执行回报类型, 取值参考enum ExecType
//
//	Price      float64 `json:"price"`      //委托成交价格
//	Volume     int64   `json:"volume"`     //委托成交量
//	Amount     float64 `json:"amount"`     //委托成交金额
//	Commission float64 `json:"commission"` //委托成交手续费
//	Cost       float64 `json:"cost"`       //委托成交成本金额
//	CreatedAt  int64   `json:"created_at"` //回报创建时间
//}
//
//type Cash struct {
//	AccountId   string `json:"account_id"`   //账号ID
//	AccountName string `json:"account_name"` //账户登录名
//
//	Currency int32 `json:"currency"` //币种
//
//	Nav         float64 `json:"nav"`          //净值(CumInout + CumPnl + FPnl - CumCommission)
//	Pnl         float64 `json:"pnl"`          //净收益(Nav-CumInout)
//	FPnl        float64 `json:"fpnl"`         //浮动盈亏(Sum(Each Position FPnl))
//	Frozen      float64 `json:"frozen"`       //持仓占用资金
//	OrderFrozen float64 `json:"order_frozen"` //挂单冻结资金
//	Available   float64 `json:"available"`    //可用资金
//	//No  Leverage:  Available=(CumInout + CumPnl - CumCommission - Frozen - OrderFrozen)
//	//Has Leverage:  FPnl     =(FPnl>0 ? FPnl : (Frozen < |FPnl|) ? (Frozen-|FPnl|) : 0)
//	//               Available=(CumInout + CumPnl - CumCommission - Frozen - OrderFrozen + FPnl)
//	Balance       float64 `json:"balance"`        //资金余额
//	MarketValue   float64 `json:"market_value"`   //持仓市值
//	CumInout      float64 `json:"cum_inout"`      //累计出入金
//	CumTrade      float64 `json:"cum_trade"`      //累计交易额
//	CumPnl        float64 `json:"cum_pnl"`        //累计平仓收益(没扣除手续费)
//	CumCommission float64 `json:"cum_commission"` //累计手续费
//
//	LastTrade      float64 `json:"last_trade"`      //上一次交易额
//	LastPnl        float64 `json:"last_pnl"`        //上一次收益
//	LastCommission float64 `json:"last_commission"` //上一次手续费
//	LastInout      float64 `json:"last_inout"`      //上一次出入金
//	ChangeReason   int32   `json:"change_reason"`   //资金变更原因，取值参考enum CashPositionChangeReason
//	ChangeEventId  string  `json:"change_event_id"` //触发资金变更事件的ID
//
//	CreatedAt int64 `json:"created_at"` //资金初始时间
//	UpdatedAt int64 `json:"updated_at"` //资金变更时间
//}
//
//type Position struct {
//	AccountId   string `json:"account_id"`   //账号ID
//	AccountName string `json:"account_name"` //账户登录名
//
//	Symbol      string  `json:"symbol"`       //Symbol
//	Side        int32   `json:"side"`         //持仓方向，取值参考enum PositionSide
//	Volume      int64   `json:"volume"`       //总持仓量 昨持仓量(Volume-VolumeToday)
//	VolumeToday int64   `json:"volume_today"` //今日持仓量
//	VWap        float64 `json:"vwap"`         //持仓均价(股票为基于开仓价的持仓均价，期货为基于结算价的持仓均价)
//	VWapDiluted float64 `json:"vwapdiluted"`  //摊薄成本价
//	VWapOpen    float64 `json:"vwap_open"`    //基于开仓价的持仓均价(期货)
//	Amount      float64 `json:"amount"`       //持仓额(Volume*VWap*Multiplier)
//
//	Price            float64 `json:"price"`              //当前行情价格
//	FPnl             float64 `json:"fpnl"`               //持仓浮动盈亏((Price-VWap)*Volume*Multiplier)
//	FPnlOpen         float64 `json:"fpnl_open"`          //持仓浮动盈亏,基于开仓均价，适用于期货((Price-VWapOpen)*Volume*Multiplier)
//	Cost             float64 `json:"cost"`               //持仓成本(VWap*Volume*Multiplier*MarginRatio)
//	OrderFrozen      int64   `json:"order_frozen"`       //挂单冻结仓位
//	OrderFrozenToday int64   `json:"order_frozen_today"` //挂单冻结今仓仓位
//	Available        int64   `json:"available"`          //可用总仓位(Volume-OrderFrozen) 可用昨仓位(Available-AvailableToday)
//	AvailableToday   int64   `json:"available_today"`    //可用今仓位(VolumeToday-OrderFrozenToday)
//	AvailableNow     int64   `json:"available_now"`      //当前可平仓位
//	MarketValue      float64 `json:"market_value"`       //持仓市值
//
//	LastPrice     float64 `json:"last_price"`      //上一次成交价
//	LastVolume    int64   `json:"last_volume"`     //上一次成交量
//	LastInout     int64   `json:"last_inout"`      //上一次出入持仓量
//	ChangeReason  int32   `json:"change_reason"`   //仓位变更原因，取值参考enum CashPositionChangeReason
//	ChangeEventId string  `json:"change_event_id"` //触发资金变更事件的ID
//
//	HasDividend int32 `json:"has_dividend"` //持仓区间有分红配送
//	CreatedAt   int64 `json:"created_at"`   //建仓时间
//	UpdatedAt   int64 `json:"updated_at"`   //仓位变更时间
//}
//
//type Account struct {
//	AccountId   string `json:"account_id"`   //账号ID
//	AccountName string `json:"account_name"` //账户登录名
//	Title       string `json:"title"`        //账号名称
//	Intro       string `json:"intro"`        //账号描述
//	Comment     string `json:"comment"`      //账号备注
//}
//
//type AccountStatus struct {
//	AccountId   string `json:"account_id"`   //账号ID
//	AccountName string `json:"account_name"` //账户登录名
//	State       int32  `json:"state"`        //账户状态
//	ErrorCode   int32  `json:"error_code"`   //错误码
//	ErrorMsg    string `json:"error_msg"`    //错误信息
//}
//
//type Parameter struct {
//	Key      string  `json:"key"`      //参数键
//	Value    float64 `json:"value"`    //参数值
//	Min      float64 `json:"min"`      //可设置的最小值
//	Max      float64 `json:"max"`      //可设置的最大值
//	Name     string  `json:"name"`     //参数名
//	Intro    string  `json:"intro"`    //参数说明
//	Group    string  `json:"group"`    //组名
//	Readonly bool    `json:"readonly"` //是否只读
//}
//
//type Indicator struct {
//	AccountId      string  `json:"account_id"`       //账号ID
//	PnlRatio       float64 `json:"pnl_ratio"`        //累计收益率(Pnl/CumInout)
//	PnlRatioAnnual float64 `json:"pnl_ratio_annual"` //年化收益率
//	SharpRatio     float64 `json:"sharp_ratio"`      //夏普比率
//	MaxDrawDown    float64 `json:"max_drawdown"`     //最大回撤
//	RiskRatio      float64 `json:"risk_ratio"`       //风险比率
//	OpenCount      int32   `json:"open_count"`       //开仓次数
//	CloseCount     int32   `json:"close_count"`      //平仓次数
//	WinCount       int32   `json:"win_count"`        //盈利次数
//	LoseCount      int32   `json:"lose_count"`       //亏损次数
//	WinRatio       float64 `json:"win_ratio"`        //胜率
//
//	CreatedAt int64 //指标创建时间
//	UpdatedAt int64 //指标变更时间
//}
//
//type CollateralInstrument struct {
//	Symbol     string  `json:"symbol"`      //担保证券标的
//	Name       string  `json:"name"`        //名称
//	PledgeRate float64 `json:"pledge_rate"` //折算率
//}
//
//type BorrowableInstrument struct {
//	Symbol                string  `json:"symbol"`                   //可融证券标的
//	Name                  string  `json:"name"`                     //名称
//	MarginRateForCash     float64 `json:"margin_rate_for_cash"`     //融资保证金比率
//	MarginRateForSecurity float64 `json:"margin_rate_for_security"` //融券保证金比率
//}
//
//type BorrowableInstrumentPosition struct {
//	Symbol    string  `json:"symbol"`    //可融证券标的
//	Name      string  `json:"name"`      //名称
//	Balance   float64 `json:"balance"`   //证券余额
//	Available float64 `json:"available"` //证券可用
//}
//
//type CreditContract struct {
//	Symbol           string  `json:"symbol"`           //证券代码 Stkcode
//	Name             string  `json:"name"`             //名称
//	OrderDate        int32   `json:"orderdate"`        //委托日期
//	OrderSno         string  `json:"ordersno"`         //委 托 号
//	CreditDirect     byte    `json:"creditdirect"`     //融资融券方向
//	OrderQty         float64 `json:"orderqty"`         //委托数量
//	MatchQty         float64 `json:"matchqty"`         //成交数量
//	OrderAmt         float64 `json:"orderamt"`         //委托金额
//	OrderFrzAmt      float64 `json:"orderfrzamt"`      //委托冻结金额
//	MatChAmt         float64 `json:"matchamt"`         //成交金额
//	ClearAmt         float64 `json:"clearamt"`         //清算金额
//	LifeStatus       byte    `json:"lifestatus"`       //合约状态
//	EndDate          int32   `json:"enddate"`          //负债截止日期
//	OldEndDate       int32   `json:"oldenddate"`       //原始的负债截止日期
//	CreditRePay      float64 `json:"creditrepay"`      //T日之前归还金额
//	CreditRePayUnFrz float64 `json:"creditrepayunfrz"` //T日归还金额
//	FunReMain        float64 `json:"funremain"`        //应还金额
//	StkRePay         float64 `json:"stkrepay"`         //T日之前归还数量
//	StkRePayUnFrz    float64 `json:"stkrepayunfrz"`    //T日归还数量
//	StkReMain        float64 `json:"stkremain"`        //应还证券数量
//	StkReMainValue   float64 `json:"stkremainvalue"`   //应还证券市值
//	Fee              float64 `json:"fee"`              //融资融券息、费
//	OverDueFee       float64 `json:"overduefee"`       //逾期未偿还息、费
//	FeeRepay         float64 `json:"feerepay"`         //己偿还息、费
//	PUniFee          float64 `json:"punifee"`          //利息产生的罚息
//	PUniFeeRepay     float64 `json:"punifeerepay"`     //己偿还罚息
//	Rights           float64 `json:"rights"`           //未偿还权益金额
//	OverDueRights    float64 `json:"overduerights"`    //逾期未偿还权益
//	RightsRepay      float64 `json:"rightsrepay"`      //己偿还权益
//	LastPrice        float64 `json:"lastprice"`        //最新价
//	ProfitCost       float64 `json:"profitcost"`       //浮动盈亏
//	Sysdate          int32   `json:"sysdate"`          //系统日期
//	Sno              string  `json:"sno"`              //合约编号
//	LastDate         int32   `json:"lastdate"`         //最后一次计算息费日期
//	CloseDate        int32   `json:"closedate"`        //合约全部偿还日期
//	PUniDebts        float64 `json:"punidebts"`        //逾期本金罚息
//	PUniDebtsRepay   float64 `json:"punidebtsrepay"`   //本金罚息偿还
//	PUniDebtSunFrz   float64 `json:"punidebtsunfrz"`   //逾期本金罚息
//	PUniFeeUnFrz     float64 `json:"punifeeunfrz"`     //逾期息费罚息
//	PUniRights       float64 `json:"punirights"`       //逾期权益罚息
//	PUniRightsRepay  float64 `json:"punirightsrepay"`  //权益罚息偿还
//	PUniRightSunFrz  float64 `json:"punirightsunfrz"`  //逾期权益罚息
//	FeeUnFrz         float64 `json:"feeunfrz"`         //实时偿还利息
//	OverDueFeeUnFrz  float64 `json:"overduefeeunfrz"`  //实时偿还逾期利息
//	RightsQty        float64 `json:"rightsqty"`        //未偿还权益数量
//	OverDueRightsQty float64 `json:"overduerightsqty"` //逾期未偿还权益数量
//}
//
//type CreditCash struct {
//	FundIntrRate             float64 `json:"fundintrrate"`             //融资利率
//	StkIntrRate              float64 `json:"stkintrrate"`              //融券利率
//	PUniShIntrRate           float64 `json:"punishintrrate"`           //罚息利率
//	CreditStatus             byte    `json:"creditstatus"`             //信用状态
//	MarginRates              float64 `json:"marginrates"`              //维持担保比例
//	RealRate                 float64 `json:"realrate"`                 //实时担保比例
//	Asset                    float64 `json:"asset"`                    //总资产
//	Liability                float64 `json:"liability"`                //总负债
//	MargiNAvl                float64 `json:"marginavl"`                //保证金可用数
//	FundBal                  float64 `json:"fundbal"`                  //资金余额
//	FundAvl                  float64 `json:"fundavl"`                  //资金可用数
//	DsaleamtBal              float64 `json:"dsaleamtbal"`              //融券卖出所得资金
//	GuaranteeOut             float64 `json:"guaranteeout"`             //可转出担保资产
//	GageMktAvl               float64 `json:"gagemktavl"`               //担保证券市值
//	FDealAvl                 float64 `json:"fdealavl"`                 //融资本金
//	FFee                     float64 `json:"ffee"`                     //融资息费
//	FTotalDebts              float64 `json:"ftotaldebts"`              //融资负债合计
//	DealFMktAvl              float64 `json:"dealfmktavl"`              //应付融券市值
//	DFee                     float64 `json:"dfee"`                     //融券息费
//	DTotalDebts              float64 `json:"dtotaldebts"`              //融券负债合计
//	FCreditBal               float64 `json:"fcreditbal"`               //融资授信额度
//	FCreditAvl               float64 `json:"fcreditavl"`               //融资可用额度
//	FCreditFrz               float64 `json:"fcreditfrz"`               //融资额度冻结
//	DCreditBal               float64 `json:"dcreditbal"`               //融券授信额度
//	DCreditAvl               float64 `json:"dcreditavl"`               //融券可用额度
//	DCreditFrz               float64 `json:"dcreditfrz"`               //融券额度冻结
//	Rights                   float64 `json:"rights"`                   //红利权益
//	ServiceUnComErqRights    float64 `json:"serviceuncomerqrights"`    //红利权益(在途)
//	RightsQty                float64 `json:"rightsqty"`                //红股权益
//	ServiceUnComErqRightsQty float64 `json:"serviceuncomerqrightsqty"` //红股权益(在途)
//	ACreditBal               float64 `json:"acreditbal"`               //总额度
//	ACreditAvl               float64 `json:"acreditavl"`               //总可用额度
//	ACashCapital             float64 `json:"acashcapital"`             //所有现金资产（所有资产、包括融券卖出）
//	AStkMktValue             float64 `json:"astkmktvalue"`             //所有证券市值（包含融资买入、非担保品）
//	WithDrawable             float64 `json:"withdrawable"`             //可取资金
//	NetCapital               float64 `json:"netcapital"`               //净资产
//	FCreditPnl               float64 `json:"fcreditpnl"`               //融资盈亏
//	DCreditPnl               float64 `json:"dcreditpnl"`               //融券盈亏
//	FCreditMarginOccupied    float64 `json:"fcreditmarginoccupied"`    //融资占用保证金
//	DCreditMarginOccupied    float64 `json:"dcreditmarginoccupied"`    //融券占用保证金
//	CollateralBuyAbleAmt     float64 `json:"collateralbuyableamt"`     //可买担保品资金
//	RePayAbleAmt             float64 `json:"repayableamt"`             //可还款金额
//	DCreditCashAvl           float64 `json:"dcreditcashavl"`           //融券可用资金
//}
//
//// IPOQI 新股申购额度
//type IPOQI struct {
//	Exchange     string  `json:"exchange"`       //市场代码
//	Quota        float64 `json:"quota"`          //市场配额
//	SseStarQuota float64 `json:"sse_star_quota"` //上海科创板配额
//}
//
//type IPOInstruments struct {
//	Symbol string  `json:"symbol"`  //申购新股symbol
//	Price  float64 `json:"price"`   //申购价格
//	MinVol int32   `json:"min_vol"` //申购最小数量
//	MaxVol int32   `json:"max_vol"` //申购最大数量
//}
//
//type IPOMatchNumber struct {
//	OrderId     string `json:"order_id"`     //委托号
//	Symbol      string `json:"symbol"`       //新股symbol
//	Volume      int32  `json:"volume"`       //成交数量
//	MatchNumber string `json:"match_number"` //申购配号
//	OrderAt     int32  `json:"order_at"`     //委托日期
//	MatchAt     int32  `json:"match_at"`     //配号日期
//}
//
//type IPOLotInfo struct {
//	Symbol       string  `json:"symbol"`         //新股symbol
//	OrderAt      int32   `json:"order_at"`       //委托日期
//	LotAt        int32   `json:"lot_at"`         //中签日期
//	LotVolume    int32   `json:"lot_volume"`     //中签数量
//	GiveUpVolume int32   `json:"give_up_volume"` //放弃数量
//	Price        float64 `json:"price"`          //中签价格
//	Amount       float64 `json:"amount"`         //中签金额
//	PayVolume    float64 `json:"pay_volume"`     //已缴款数量
//	PayAmount    float64 `json:"pay_amount"`     //已缴款金额
//}

// SubscribePara 订阅行情参数
type SubscribePara struct {
	Symbol           string `json:"symbol"`             // 订阅代码
	Frequency        string `json:"frequency"`          // 订阅频率
	IsUnSubscribePre bool   `json:"is_unsubscribe_pre"` // 是否退订以前的订阅,默认为false
}

const (
	FixedVolumeEntrust = 1 // 固定数量委托
	FixedPriceEntrust  = 2 // 固定价格委托
	PercentEntrust     = 3 // 百分比委托
)

// EntrustParam 委托订单参数
type EntrustParam struct {
	Symbol         string  `json:"symbol"`          // 指定委托代码符号
	Volume         int32   `json:"volume"`          // 指定委托数量
	Value          float64 `json:"value"`           // 指定委托价格
	Percent        float64 `json:"percent"`         // 指定委托比例
	Side           int32   `json:"side"`            // 指定交易方向, 0: 未知 1: 买入 2: 卖出
	OrderType      int32   `json:"order_type"`      // 指定交易类型, 0: 未知 1: 限价委托 2: 市价委托 3: 止盈止损委托
	PositionEffect int32   `json:"position_effect"` // 指定开仓平仓标志, 0: 未知 1: 开仓 2: 平仓 3: 平今日仓 4: 平昨日仓 取值参考gmdef.h enum PositionEffect
	Price          float64 `json:"price"`           // 指定委托价格
	Account        string  `json:"account"`         // 指定交易账户

	EntrustType int32 `json:"entrust_type"` // 指定委托订单类型,  1:定量委托 2:定价委托 3: 按总资产比例委托
}

const (
	AdjustPositionToVolume  = 1 // 按数量调仓
	AdjustPositionToAmount  = 2 // 按额度调仓
	AdjustPositionToPercent = 3 // 按比例调仓
)

// AdjustPositionParam 调整持仓参数
type AdjustPositionParam struct {
	Symbol         string  `json:"symbol"`          // 指定代码符号
	Volume         int32   `json:"volume"`          // 指定调整仓位目标数量
	Value          float64 `json:"value"`           // 指定调整仓位目标持仓额度
	Percent        float64 `json:"percent"`         // 指定调整仓位目标持仓比例(总持资产额度计算)
	PositionSide   int32   `json:"position_side"`   // 指定交易方向, 0: 未知 1: 买入 2: 卖出
	OrderType      int32   `json:"order_type"`      // 指定交易类型, 0: 未知 1: 限价委托 2: 市价委托 3: 止盈止损委托
	PositionEffect int32   `json:"position_effect"` // 指定开仓平仓标志, 0: 未知 1: 开仓 2: 平仓 3: 平今日仓 4: 平昨日仓 取值参考gmdef.h enum PositionEffect
	Price          float64 `json:"price"`           // 指定调仓价格
	Account        string  `json:"account"`         // 指定交易账户

	AdjustType int32 `json:"adjust_type"` // 调仓类型,  1:调仓位到目标持仓数量 2:调仓位到目标持仓额 3: 调仓位到目标持仓比例(总资产的比例)
}

// PlaceOrderParam 委托下单参数
type PlaceOrderParam struct {
	Symbol         string  `json:"symbol"`          // 指定代码符号
	Volume         int32   `json:"volume"`          // 指定委托下单数量
	Side           int32   `json:"side"`            // 指定委托下单方向, 0: 未知 1: 买入 2: 卖出
	OrderType      int32   `json:"order_type"`      // 指定委托下单类型, 0: 未知 1: 限价委托 2: 市价委托 3: 止盈止损委托
	PositionEffect int32   `json:"position_effect"` // 指定委托下单开仓平仓标志, 0: 未知 1: 开仓 2: 平仓 3: 平今日仓 4: 平昨日仓 取值参考gmdef.h enum PositionEffect
	Price          float64 `json:"price"`           // 指定委托下单价格
	OrderDuration  int32   `json:"order_duration"`  // 指定委托下单有效期, 0: 未知 1: 即时成交剩余撤销 2: 即时全额成交或撤销 3: 当日有效 4: 本节有效 5: 指定日期前有效 6: 撤销前有效 7: 集合竞前有效
	OrderQualifier int32   `json:"order_qualifier"` // 指定委托下单限定条件, 0: 未知 1: 对方最优价格 2: 己方最优价格 3: 最优五档剩余撤销 4:最优五档剩余转限价
	StopPrice      int32   `json:"stop_price"`      // 指定委托下单后止损止盈触发价格
	OrderBusiness  int32   `json:"order_business"`  // 指定委托订单业务类型
	Account        string  `json:"account"`         // 指定交易账户
}

// AlgoOrderParam 委托算法下单参数
type AlgoOrderParam struct {
	Symbol         string  `json:"symbol"`          // 指定代码符号
	Volume         int32   `json:"volume"`          // 指定委托下单数量
	PositionEffect int32   `json:"position_effect"` // 指定委托下单开仓平仓标志, 0: 未知 1: 开仓 2: 平仓 3: 平今日仓 4: 平昨日仓 取值参考gmdef.h enum PositionEffect
	Side           int32   `json:"side"`            // 指定委托下单方向, 0: 未知 1: 买入 2: 卖出
	OrderType      int32   `json:"order_type"`      // 指定委托下单类型, 0: 未知 1: 限价委托 2: 市价委托 3: 止盈止损委托
	Price          float64 `json:"price"`           // 指定委托下单价格
	AlgoName       string  `json:"algo_name"`       // 指定算法名称
	AlgoParam      string  `json:"algo_param"`      // 指定算法参数
	Account        string  `json:"account"`         // 指定交易账户
}

const (
	CreditBuying          = 1 // 融资买入
	CreditSelling         = 2 // 融券卖出
	BuyingShareRepayShare = 3 // 买券还券
	SellingShareRepay     = 4 // 卖券还款
	CollateralBuying      = 5 // 担保品买入
	CollateralSelling     = 6 // 担保品卖出
)

// CreditParam 融资融券参数
type CreditParam struct {
	PositionSrc    int32   `json:"position_src"`    // 指定头寸来源, 0: 未知 1: 普通沲 2: 专项沲
	Symbol         string  `json:"symbol"`          // 指定代码符号
	Volume         int32   `json:"volume"`          // 指定买入数量
	Price          float64 `json:"price"`           // 指定买入价格
	OrderType      int32   `json:"order_type"`      // 指定委托下单类型, 0: 未知 1: 限价委托 2: 市价委托 3: 止盈止损委托
	OrderDuration  int32   `json:"order_duration"`  // 指定委托下单有效期, 0: 未知 1: 即时成交剩余撤销 2: 即时全额成交或撤销 3: 当日有效 4: 本节有效 5: 指定日期前有效 6: 撤销前有效 7: 集合竞前有效
	OrderQualifier int32   `json:"order_qualifier"` // 指定委托下单限定条件, 0: 未知 1: 对方最优价格 2: 己方最优价格 3: 最优五档剩余撤销 4:最优五档剩余转限价
	Account        string  `json:"account"`         // 指定交易账户

	CreditType int32 `json:"credit_type"` // 两融业务类型, 0: 未知 1: 融资买入 2: 融券卖出 3: 买券还券 4: 卖券还款 5: 担保品买入 6: 担保品卖出
}

// NationalDebtParam 国债参数
type NationalDebtParam struct {
	Symbol         string  `json:"symbol"`          // 指定代码符号
	Volume         int32   `json:"volume"`          // 指定委托数量
	Price          float64 `json:"price"`           // 指定委托价格
	OrderType      int32   `json:"order_type"`      // 指定委托类型, 0: 未知 1: 限价委托 2: 市价委托 3: 止盈止损委托
	OrderDuration  int32   `json:"order_duration"`  // 指定委托有效期, 0: 未知 1: 即时成交剩余撤销 2: 即时全额成交或撤销 3: 当日有效 4: 本节有效 5: 指定日期前有效 6: 撤销前有效 7: 集合竞前有效
	OrderQualifier int32   `json:"order_qualifier"` // 指定委托限定条件, 0: 未知 1: 对方最优价格 2: 己方最优价格 3: 最优五档剩余撤销 4:最优五档剩余转限价
	Account        string  `json:"account"`         // 指定交易账户
}

type HistoryTickParam struct {
	Symbols       string `json:"symbols"`    // 指定代码符号
	StartTime     string `json:"start_time"` // 指定开始时间
	EndTime       string `json:"end_time"`   // 指定结束时间
	AdjustEndTime string `json:"adjust_end_time"`
	FillMissing   string `json:"fill_missing"` // 指定是否填充缺失数据
	Adjust        int32  `json:"adjust"`
	SkipSuspended bool   `json:"skip_suspended"` // 指定是否跳过停牌日
	N             int32  `json:"n"`              // 指定查询数量
}

type HistoryBarParam struct {
	Symbols       string `json:"symbols"`         // 指定代码符号
	StartTime     string `json:"start_time"`      // 指定开始时间
	Frequency     string `json:"frequency"`       // 指定周期频率
	N             int32  `json:"n"`               // 指定查询数量
	EndTime       string `json:"end_time"`        // 指定结束时间
	AdjustEndTime string `json:"adjust_end_time"` // 复权基点时间, 默认当前时间
	FillMissing   string `json:"fill_missing"`    //  None  不填充,   'NaN' - 用空值填充, 'Last' - 用上一个值填充, 默认None
	Adjust        int32  `json:"adjust"`          //  0: 不复权  1: 前复权,   2: 后复权
	SkipSuspended bool   `json:"skip_suspended"`  // 指定是否跳过停牌日
}

type FundamentalParam struct {
	Table     string `json:"table"` // 指定代码符号
	Symbols   string `json:"symbols"`
	StartDate string `json:"start_date"` // 指定开始时间
	EndDate   string `json:"end_date"`   // 指定结束时间
	Fields    string `json:"fields"`     // 指定查询字段
	Filter    string `json:"filter"`     // 指定过滤条件
	OrderBy   string `json:"order_by"`   // 指定排序条件
	Limit     int32  `json:"limit"`      // 指定查询数量
	N         int32  `json:"n"`          // 指定查询数量
}
