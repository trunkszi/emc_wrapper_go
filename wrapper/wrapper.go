package wrapper

/*
#cgo CFLAGS: -I. -Iemc/src -Iemc/src/third_party/gmsdk/include
#cgo LDFLAGS: -L. -Lemc/build/windows/x64/release -lemc -lgmsdk
#include "emc/src/third_party/gmsdk/include/gmdef.h"
#include "emc/src/interface.h"
#include <stdlib.h>
#include <string.h>

typedef struct Tick Tick;
typedef struct AccountStatus AccountStatus;
typedef struct AlgoOrder AlgoOrder;
typedef struct Bar Bar;
typedef struct Cash Cash;
typedef struct ExecRpt ExecRpt;
typedef struct Indicator Indicator;
typedef struct L2Order L2Order;
typedef struct L2OrderQueue L2OrderQueue;
typedef struct L2Transaction L2Transaction;
typedef struct Order Order;
typedef struct Parameter Parameter;
typedef struct Position Position;
typedef struct Account Account;

void OnInit();
void OnTick(struct Tick* tick);
void OnBar(struct Bar* bar);
void Onl2Transaction(struct L2Transaction* l2transaction);
void Onl2Order(struct L2Order* l2order);
void OnL2OrderQueue(struct L2OrderQueue* l2queue);
void OnOrderStatus(struct Order* order);
void OnExecutionReport(struct ExecRpt* rpt);
void OnAlgoOrderStatus(struct AlgoOrder* order);
void OnCash(struct Cash* cash);
void OnPosition(struct Position* position);
void OnParameter(struct Parameter* param);
void OnSchedule( char *data_rule,  char *time_rule);
void OnBackTestFinished(struct Indicator *indicator);
void OnAccountStatus(struct AccountStatus *account_status);
void OnError(int error_code,  char *error_msg);
void OnStop();
void OnMarketDataConnected();
void OnTradeDataConnected();
void OnMarketDataDisconnected();
void OnTradeDataDisconnected();
*/
import "C"
import (
	"fmt"
	"github.com/bytedance/sonic"
	"unsafe"
)

var global GoStrategy

func init() {
	CreateInstance()
	C.register_callbacks(
		C.on_init_callback(C.OnInit),
		C.on_tick_callback(C.OnTick),
		C.on_bar_callback(C.OnBar),
		C.on_l2transaction_callback(C.Onl2Transaction),
		C.on_l2order_callback(C.Onl2Order),
		C.on_l2order_queue_callback(C.OnL2OrderQueue),
		C.on_order_status_callback(C.OnOrderStatus),
		C.on_execution_report_callback(C.OnExecutionReport),
		C.on_algo_order_status_callback(C.OnAlgoOrderStatus),
		C.on_cash_callback(C.OnCash),
		C.on_position_callback(C.OnPosition),
		C.on_parameter_callback(C.OnParameter),
		C.on_schedule_callback(C.OnSchedule),
		C.on_backtest_finished_callback(C.OnBackTestFinished),
		C.on_account_status_callback(C.OnAccountStatus),
		C.on_error_callback(C.OnError),
		C.on_stop_callback(C.OnStop),
		C.on_market_data_connected_callback(C.OnMarketDataConnected),
		C.on_trade_data_connected_callback(C.OnTradeDataConnected),
		C.on_market_data_disconnected_callback(C.OnMarketDataDisconnected),
		C.on_trade_data_disconnected_callback(C.OnTradeDataDisconnected))
}

//export OnInit
func OnInit() {

	global.GoInit()
}

//export OnTick
func OnTick(tick *C.Tick) {
	if tick == nil {
		return
	}
	var t Tick
	t.Symbol = C.GoString(&tick.symbol[0])
	t.CreatedAt = float64(tick.created_at)
	t.Price = float32(tick.price)
	t.Open = float32(tick.open)
	t.High = float32(tick.high)
	t.Low = float32(tick.low)
	t.CumVolume = float64(tick.cum_volume)
	t.CumAmount = float64(tick.cum_amount)
	t.CumPosition = int64(tick.cum_position)
	t.LastAmount = float64(tick.last_amount)
	t.LastVolume = int32(tick.last_volume)
	t.TradeType = int32(tick.trade_type)
	for i := 0; i < 10; i++ {
		var q Quote
		q.BidPrice = float32(tick.quotes[i].bid_price)
		q.BidVolume = int64(tick.quotes[i].bid_volume)
		q.AskPrice = float32(tick.quotes[i].ask_price)
		q.AskVolume = int64(tick.quotes[i].ask_volume)
		t.Quotes = append(t.Quotes, &q)
	}
	global.GoTick(&t)
}

//export OnBar
func OnBar(bar *C.Bar) {
	fmt.Println("OnBar")
	var b Bar
	b.Symbol = C.GoString(&bar.symbol[0])
	b.Bob = float64(bar.bob)
	b.Eob = float64(bar.eob)
	b.Open = float32(bar.open)
	b.Close = float32(bar.close)
	b.High = float32(bar.high)
	b.Low = float32(bar.low)
	b.Volume = float64(bar.volume)
	b.Amount = float64(bar.amount)
	b.PreClose = float32(bar.pre_close)
	b.Position = int64(bar.position)
	b.Frequency = C.GoString(&bar.frequency[0])
	global.GoBar(&b)
}

//export Onl2Transaction
func Onl2Transaction(l2transaction *C.L2Transaction) {
	var l2 L2Transaction
	l2.Symbol = C.GoString(&l2transaction.symbol[0])
	l2.CreatedAt = float64(l2transaction.created_at)
	l2.Price = float32(l2transaction.price)
	l2.Volume = int64(l2transaction.volume)
	l2.Side = string(l2transaction.side)
	l2.ExecType = string(l2transaction.exec_type)
	l2.ExecIndex = int64(l2transaction.exec_index)
	l2.AskOrderIndex = int64(l2transaction.ask_order_index)
	l2.BidOrderIndex = int64(l2transaction.bid_order_index)
	global.Gol2Transaction(&l2)
}

//export Onl2Order
func Onl2Order(l2order *C.L2Order) {
	var l2 L2Order
	l2.Symbol = C.GoString(&l2order.symbol[0])
	l2.CreatedAt = float64(l2order.created_at)
	l2.Price = float32(l2order.price)
	l2.Volume = int64(l2order.volume)
	l2.Side = string(l2order.side)
	l2.OrderType = string(l2order.order_type)
	l2.OrderIndex = int64(l2order.order_index)
	global.Gol2Order(&l2)
}

//export OnL2OrderQueue
func OnL2OrderQueue(l2queue *C.L2OrderQueue) {
	var l2 L2OrderQueue
	l2.Symbol = C.GoString(&l2queue.symbol[0])
	l2.CreatedAt = float64(l2queue.created_at)
	l2.Price = float32(l2queue.price)
	l2.Volume = int64(l2queue.volume)
	l2.Side = string(l2queue.side)
	l2.QueueOrders = int32(l2queue.queue_orders)
	for i := 0; i < int(l2.QueueOrders); i++ {
		l2.QueueVolumes = append(l2.QueueVolumes, int32(l2queue.queue_volumes[i]))
	}
	global.GoL2OrderQueue(&l2)
}

//export OnOrderStatus
func OnOrderStatus(order *C.Order) {
	var o Order
	o.StrategyId = C.GoString(&order.strategy_id[0])
	o.AccountId = C.GoString(&order.account_id[0])
	o.AccountName = C.GoString(&order.account_name[0])
	o.ClOrdId = C.GoString(&order.cl_ord_id[0])
	o.OrderId = C.GoString(&order.order_id[0])
	o.ExOrdId = C.GoString(&order.ex_ord_id[0])
	o.AlgoOrderId = C.GoString(&order.algo_order_id[0])
	o.OrderBusiness = int32(order.order_business)

	o.Symbol = C.GoString(&order.symbol[0])
	o.Side = int32(order.side)
	o.PositionEffect = int32(order.position_effect)
	o.PositionSide = int32(order.position_side)

	o.OrderType = int32(order.order_type)
	o.OrderDuration = int32(order.order_duration)
	o.OrderQualifier = int32(order.order_qualifier)
	o.OrderSrc = int32(order.order_src)
	o.PositionSrc = int32(order.position_src)

	o.Status = int32(order.status)
	o.OrdRejReason = int32(order.ord_rej_reason)
	o.OrdRejReasonDetail = C.GoString(&order.ord_rej_reason_detail[0])

	o.Price = float64(order.price)
	o.StopPrice = float64(order.stop_price)

	o.OrderStyle = int32(order.order_style)
	o.Volume = int64(order.volume)
	o.Value = float64(order.value)
	o.Percent = float64(order.percent)
	o.TargetVolume = int64(order.target_volume)
	o.TargetValue = float64(order.target_value)
	o.TargetPercent = float64(order.target_percent)

	o.FilledVolume = int64(order.filled_volume)
	o.FilledVwap = float64(order.filled_vwap)
	o.FilledAmount = float64(order.filled_amount)
	o.FilledCommission = float64(order.filled_commission)

	o.CreatedAt = int64(order.created_at)
	o.UpdatedAt = int64(order.updated_at)
	global.GoOrderStatus(&o)
}

//export OnExecutionReport
func OnExecutionReport(rpt *C.ExecRpt) {
	var r ExecRpt
	r.StrategyId = C.GoString(&rpt.strategy_id[0])
	r.AccountId = C.GoString(&rpt.account_id[0])
	r.AccountName = C.GoString(&rpt.account_name[0])
	r.ClOrdId = C.GoString(&rpt.cl_ord_id[0])
	r.OrderId = C.GoString(&rpt.order_id[0])
	r.ExecId = C.GoString(&rpt.exec_id[0])

	r.Symbol = C.GoString(&rpt.symbol[0])

	r.PositionEffect = int32(rpt.position_effect)
	r.Side = int32(rpt.side)
	r.OrdRejReason = int32(rpt.ord_rej_reason)
	r.OrdRejReasonDetail = C.GoString(&rpt.ord_rej_reason_detail[0])
	r.ExecType = int32(rpt.exec_type)

	r.Price = float64(rpt.price)
	r.Volume = int64(rpt.volume)
	r.Amount = float64(rpt.amount)
	r.Commission = float64(rpt.commission)
	r.Cost = float64(rpt.cost)
	r.CreatedAt = int64(rpt.created_at)

	global.GoExecutionReport(&r)
}

//export OnAlgoOrderStatus
func OnAlgoOrderStatus(order *C.AlgoOrder) {
	var o AlgoOrder
	o.StrategyId = C.GoString(&order.strategy_id[0])
	o.AccountId = C.GoString(&order.account_id[0])
	o.AccountName = C.GoString(&order.account_name[0])
	o.ClOrdId = C.GoString(&order.cl_ord_id[0])
	o.OrderId = C.GoString(&order.order_id[0])
	o.ExOrdId = C.GoString(&order.ex_ord_id[0])
	o.OrderBusiness = int32(order.order_business)

	o.Symbol = C.GoString(&order.symbol[0])
	o.Side = int32(order.side)
	o.PositionEffect = int32(order.position_effect)
	o.PositionSide = int32(order.position_side)

	o.OrderType = int32(order.order_type)
	o.OrderDuration = int32(order.order_duration)
	o.OrderQualifier = int32(order.order_qualifier)
	o.OrderSrc = int32(order.order_src)
	o.PositionSrc = int32(order.position_src)

	o.Status = int32(order.status)
	o.OrdRejReason = int32(order.ord_rej_reason)
	o.OrdRejReasonDetail = C.GoString(&order.ord_rej_reason_detail[0])

	o.Price = float64(order.price)
	o.StopPrice = float64(order.stop_price)

	o.OrderStyle = int32(order.order_style)
	o.Volume = int64(order.volume)
	o.Value = float64(order.value)
	o.Percent = float64(order.percent)
	o.TargetVolume = int64(order.target_volume)
	o.TargetValue = float64(order.target_value)
	o.TargetPercent = float64(order.target_percent)

	o.FilledVolume = int64(order.filled_volume)
	o.FilledVwap = float64(order.filled_vwap)
	o.FilledAmount = float64(order.filled_amount)
	o.FilledCommission = float64(order.filled_commission)

	o.AlgoName = C.GoString(&order.algo_name[0])
	o.AlgoParam = C.GoString(&order.algo_param[0])
	o.AlgoStatus = int32(order.algo_status)
	o.AlgoComment = C.GoString(&order.algo_comment[0])

	o.CreatedAt = int64(order.created_at)
	o.UpdatedAt = int64(order.updated_at)

	global.GoAlgoOrderStatus(&o)
}

//export OnCash
func OnCash(cash *C.Cash) {
	var c Cash
	c.AccountId = C.GoString(&cash.account_id[0])
	c.AccountName = C.GoString(&cash.account_name[0])

	c.Currency = int32(cash.currency)

	c.Nav = float64(cash.nav)
	c.Pnl = float64(cash.pnl)
	c.Fpnl = float64(cash.fpnl)
	c.Frozen = float64(cash.frozen)
	c.OrderFrozen = float64(cash.order_frozen)
	c.Available = float64(cash.available)
	c.Balance = float64(cash.balance)
	c.MarketValue = float64(cash.market_value)
	c.CumInout = float64(cash.cum_inout)
	c.CumTrade = float64(cash.cum_trade)
	c.CumPnl = float64(cash.cum_pnl)
	c.CumCommission = float64(cash.cum_commission)

	c.LastTrade = float64(cash.last_trade)
	c.LastPnl = float64(cash.last_pnl)
	c.LastCommission = float64(cash.last_commission)
	c.LastInout = float64(cash.last_inout)
	c.ChangeReason = int32(cash.change_reason)
	c.ChangeEventId = C.GoString(&cash.change_event_id[0])

	c.CreatedAt = int64(cash.created_at)
	c.UpdatedAt = int64(cash.updated_at)

	global.GoCash(&c)
}

//export OnPosition
func OnPosition(position *C.Position) {
	var p Position
	p.AccountId = C.GoString(&position.account_id[0])
	p.AccountName = C.GoString(&position.account_name[0])
	p.Symbol = C.GoString(&position.symbol[0])
	p.Side = int32(position.side)
	p.Volume = int64(position.volume)
	p.VolumeToday = int64(position.volume_today)
	p.Vwap = float64(position.vwap)
	p.VwapDiluted = float64(position.vwap_diluted)
	p.VwapOpen = float64(position.vwap_open)
	p.Amount = float64(position.amount)

	p.Price = float64(position.price)
	p.Fpnl = float64(position.fpnl)
	p.FpnlOpen = float64(position.fpnl_open)
	p.Cost = float64(position.cost)

	p.OrderFrozen = int64(position.order_frozen)
	p.OrderFrozenToday = int64(position.order_frozen_today)
	p.Available = int64(position.available)
	p.AvailableToday = int64(position.available_today)
	p.AvailableNow = int64(position.available_now)
	p.MarketValue = float64(position.market_value)

	p.LastPrice = float64(position.last_price)
	p.LastVolume = int64(position.last_volume)
	p.LastInout = int64(position.last_inout)
	p.ChangeReason = int32(position.change_reason)
	p.ChangeEventId = C.GoString(&position.change_event_id[0])

	p.HasDividend = int32(position.has_dividend)
	p.CreatedAt = int64(position.created_at)
	p.UpdatedAt = int64(position.updated_at)

	global.GoPosition(&p)
}

//export OnParameter
func OnParameter(param *C.Parameter) {
	var p Parameter
	p.Key = C.GoString(&param.key[0])
	p.Value = float64(param.value)
	p.Min = float64(param.min)
	p.Max = float64(param.max)
	p.Name = C.GoString(&param.name[0])
	p.Intro = C.GoString(&param.intro[0])
	p.Group = C.GoString(&param.group[0])
	p.Readonly = bool(param.readonly)

	global.GoParameter(&p)
}

//export OnSchedule
func OnSchedule(dataRule *C.char, timeRule *C.char) {
	global.GoSchedule(C.GoString(dataRule), C.GoString(timeRule))
}

//export OnBackTestFinished
func OnBackTestFinished(indicator *C.Indicator) {
	var i Indicator
	i.AccountId = C.GoString(&indicator.account_id[0])
	i.PnlRatio = float64(indicator.pnl_ratio)
	i.PnlRatioAnnual = float64(indicator.pnl_ratio_annual)
	i.SharpRatio = float64(indicator.sharp_ratio)
	i.MaxDrawdown = float64(indicator.max_drawdown)
	i.RiskRatio = float64(indicator.risk_ratio)
	i.OpenCount = int32(indicator.open_count)
	i.CloseCount = int32(indicator.close_count)
	i.WinCount = int32(indicator.win_count)
	i.LoseCount = int32(indicator.lose_count)
	i.WinRatio = float64(indicator.win_ratio)

	i.CreatedAt = int64(indicator.created_at)
	i.UpdatedAt = int64(indicator.updated_at)

	global.GoBackTestFinished(&i)
}

//export OnAccountStatus
func OnAccountStatus(accountStatus *C.AccountStatus) {
	var a AccountStatus
	a.AccountId = C.GoString(&accountStatus.account_id[0])
	a.AccountName = C.GoString(&accountStatus.account_name[0])
	a.State = int32(accountStatus.state)
	a.ErrorCode = int32(accountStatus.error_code)
	a.ErrorMsg = C.GoString(&accountStatus.error_msg[0])

	global.GoAccountStatus(&a)
}

//export OnError
func OnError(errorCode C.int, errorMsg *C.char) {
	global.GoError(int(errorCode), C.GoString(errorMsg))

}

//export OnStop
func OnStop() {
	global.GoStop()
}

//export OnMarketDataConnected
func OnMarketDataConnected() {
	global.GoMarketDataConnected()
}

//export OnTradeDataConnected
func OnTradeDataConnected() {
	global.GoTradeDataConnected()
}

//export OnMarketDataDisconnected
func OnMarketDataDisconnected() {
	global.GoMarketDataDisconnected()
}

//export OnTradeDataDisconnected
func OnTradeDataDisconnected() {
	global.GoTradeDataDisconnected()
}

func CreateInstance() {
	C.create_instance()
}

type GoStrategy interface {
	GoInit()
	GoTick(tick *Tick)
	GoBar(bar *Bar)
	Gol2Transaction(l2transaction *L2Transaction)
	Gol2Order(l2order *L2Order)
	GoL2OrderQueue(l2queue *L2OrderQueue)
	GoOrderStatus(order *Order)
	GoExecutionReport(rpt *ExecRpt)
	GoAlgoOrderStatus(order *AlgoOrder)
	GoCash(cash *Cash)
	GoPosition(position *Position)
	GoParameter(param *Parameter)
	GoSchedule(dataRule string, timeRule string)
	GoBackTestFinished(indicator *Indicator)
	GoAccountStatus(accountStatus *AccountStatus)
	GoError(errorCode int, errorMsg string)
	GoStop()
	GoMarketDataConnected()
	GoTradeDataConnected()
	GoMarketDataDisconnected()
	GoTradeDataDisconnected()
}

func RegisterStrategy(strategy GoStrategy) {
	global = strategy
}

var _ GoStrategy = &Strategy{}

type Strategy struct {
	GoStrategy
}

func NewStrategy() *Strategy {
	return &Strategy{}
}

// Run 运行策略
func (s *Strategy) Run() {
	C.emc_run()
}

// Stop 停止策略
func (s *Strategy) Stop() {
	C.emc_stop()
}

// SetStrategyId 设置策略ID
func (s *Strategy) SetStrategyId(strategyId string) {
	C.emc_set_strategy_id(C.CString(strategyId))
}

// SetToken 设置用户token
func (s *Strategy) SetToken(token string) {
	C.emc_set_token(C.CString(token))
}

// SetMode 设置策略运行模式
func (s *Strategy) SetMode(mode int) {
	C.emc_set_mode(C.int(mode))
}

// Schedule 设置定时任务
func (s *Strategy) Schedule(dataRule string, timeRule string) int32 {
	return int32(C.emc_schedule(C.CString(dataRule), C.CString(timeRule)))
}

// Now 获取当前时间戳
func (s *Strategy) Now() float64 {
	return float64(C.emc_now())
}

// SetBackTestConfig 设置回测参数
func (s *Strategy) SetBackTestConfig(startTime, endTime string, initialCash, transactionRatio, commissionRatio, slippageRatio float64, adjust, checkCache int) {
	C.emc_set_backtest_config(
		C.CString(startTime),
		C.CString(endTime),
		C.double(initialCash),
		C.double(transactionRatio),
		C.double(commissionRatio),
		C.double(slippageRatio),
		C.int(adjust),
		C.int(checkCache))
}

// Subscribe 订阅行情
func (s *Strategy) Subscribe(symbol string, frequency string, unsubscribePrevious bool) int32 {
	return int32(C.emc_subscribe(C.CString(symbol), C.CString(frequency), C.bool(unsubscribePrevious)))
}

// Unsubscribe 退订行情
func (s *Strategy) Unsubscribe(symbol string, frequency string) int32 {
	return int32(C.emc_unsubscribe(C.CString(symbol), C.CString(frequency)))
}

// GetAccounts 查询交易账号
func (s *Strategy) GetAccounts() string {
	var length, errorCode C.int
	cstr := C.emc_get_accounts(&length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// GetAccountStatus 查询指定交易账号状态
func (s *Strategy) GetAccountStatus(account string) string {
	var length, errorCode C.int
	cstr := C.emc_get_account_status(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// GetAllAccountStatus 查询所有交易账号状态
func (s *Strategy) GetAllAccountStatus() string {
	var length, errorCode C.int
	cstr := C.emc_get_all_account_status(&length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// EntrustedOrder 委托下单
func (s *Strategy) EntrustedOrder(param *EntrustParam) string {
	var length, errorCode C.int
	var cstr *C.char
	switch param.EntrustType {
	case FixedVolumeEntrust: // 定量委托
		cstr = C.emc_order_volume(C.CString(param.Symbol), C.int(param.Volume), C.int(param.Side), C.int(param.OrderType), C.int(param.PositionEffect), C.double(param.Price), C.CString(param.Account), &length, &errorCode)
	case FixedPriceEntrust: // 定价委托
		cstr = C.emc_order_value(C.CString(param.Symbol), C.double(param.Value), C.int(param.Side), C.int(param.OrderType), C.int(param.PositionEffect), C.double(param.Price), C.CString(param.Account), &length, &errorCode)
	case PercentEntrust: // 按总资产比例委托
		cstr = C.emc_order_percent(C.CString(param.Symbol), C.double(param.Percent), C.int(param.Side), C.int(param.OrderType), C.int(param.PositionEffect), C.double(param.Price), C.CString(param.Account), &length, &errorCode)
	}
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// AdjustPosition 调整持仓
func (s *Strategy) AdjustPosition(param *AdjustPositionParam) string {
	var length, errorCode C.int
	var cstr *C.char

	switch param.AdjustType {
	case AdjustPositionToVolume: // 调仓到目标持仓量
		cstr = C.emc_order_target_volume(C.CString(param.Symbol), C.int(param.Volume), C.int(param.PositionSide), C.int(param.OrderType), C.double(param.Price), C.CString(param.Account), &length, &errorCode)
	case AdjustPositionToAmount: // 调仓到目标持仓额度
		cstr = C.emc_order_target_value(C.CString(param.Symbol), C.double(param.Value), C.int(param.PositionSide), C.int(param.OrderType), C.double(param.Price), C.CString(param.Account), &length, &errorCode)
	case AdjustPositionToPercent: // 调仓到目标持仓比例（总资产的比例）
		cstr = C.emc_order_target_percent(C.CString(param.Symbol), C.double(param.Percent), C.int(param.PositionSide), C.int(param.OrderType), C.double(param.Price), C.CString(param.Account), &length, &errorCode)
	}
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// OrderCloseAll 撤销所有委托
func (s *Strategy) OrderCloseAll() string {
	var length, errorCode C.int
	cstr := C.emc_order_close_all(&length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// OrderCancel 委托撤单
func (s *Strategy) OrderCancel(clOrdId string, account string) int32 {
	return int32(C.emc_order_cancel(C.CString(clOrdId), C.CString(account)))
}

// PlaceOrder 委托下单
func (s *Strategy) PlaceOrder(param *PlaceOrderParam) string {
	var length, errorCode C.int
	cstr := C.emc_place_order(C.CString(param.Symbol), C.int(param.Volume), C.int(param.Side), C.int(param.OrderType), C.int(param.PositionEffect), C.double(param.Price),
		C.int(param.OrderDuration), C.int(param.OrderQualifier), C.double(param.StopPrice),
		C.int(param.OrderBusiness), C.CString(param.Account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// OrderAfterHour 盘后定价交易
func (s *Strategy) OrderAfterHour(symbol, account string, volume int32, side int32, price float64) string {
	var length, errorCode C.int
	cstr := C.emc_order_after_hour(C.CString(symbol), C.int(volume), C.int(side), C.double(price), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// GetOrders 查询所有委托
func (s *Strategy) GetOrders(account string) string {
	var length, errorCode C.int
	cstr := C.emc_get_orders(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// GetUnfinishedOrders 查询未结委托
func (s *Strategy) GetUnfinishedOrders(account string) string {
	var length, errorCode C.int
	cstr := C.emc_get_unfinished_orders(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// GetExecutionReports 查询成交
func (s *Strategy) GetExecutionReports(account string) string {
	var length, errorCode C.int
	cstr := C.emc_get_execution_reports(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// GetCash 查询资金
func (s *Strategy) GetCash(account string) string {
	var length, errorCode C.int
	cstr := C.emc_get_cash(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// GetPosition 查询持仓
func (s *Strategy) GetPosition(account string) []*Position {
	var length, errorCode C.int
	cstr := C.emc_get_position(C.CString(account), &length, &errorCode)
	if int32(length) <= 0 && int32(errorCode) == 0 {
		return nil
	}
	str := C.GoString(cstr)
	C.free(unsafe.Pointer(cstr))

	var pos []*Position
	err := sonic.Unmarshal([]byte(str), &pos)
	if err != nil {
		return nil
	}
	return pos
}

// OrderAlgo 委托算法单
func (s *Strategy) OrderAlgo(param *AlgoOrderParam) string {
	var length, errorCode C.int
	cstr := C.emc_order_algo(C.CString(param.Symbol), C.int(param.Volume), C.int(param.PositionEffect), C.int(param.Side),
		C.int(param.OrderType), C.double(param.Price), C.CString(param.AlgoName), C.CString(param.AlgoParam), C.CString(param.Account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// AlgoOrderCancel 撤销算法单
func (s *Strategy) AlgoOrderCancel(clOrdId string, account string) int32 {
	return int32(C.emc_algo_order_cancel(C.CString(clOrdId), C.CString(account)))
}

// AlgoOrderPause 暂停或恢复算法单
func (s *Strategy) AlgoOrderPause(clOrdId string, status int32, account string) int32 {
	return int32(C.emc_algo_order_pause(C.CString(clOrdId), C.int(status), C.CString(account)))
}

// GetAlgoOrders 查询算法委托
func (s *Strategy) GetAlgoOrders(account string) string {
	var length, errorCode C.int
	cstr := C.emc_get_algo_orders(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// GetAlgoChildOrders 查询算法委托
func (s *Strategy) GetAlgoChildOrders(clOrdId, account string) string {
	var length, errorCode C.int
	cstr := C.emc_get_algo_child_orders(C.CString(clOrdId), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// RawFunc 功能号调用
func (s *Strategy) RawFunc(account, funcId, funcArgs string) string {
	var rsp *C.char
	result := C.emc_raw_func(C.CString(account), C.CString(funcId), C.CString(funcArgs), &rsp)
	if result != 0 || rsp == nil {
		return ""
	}
	str := C.GoString(rsp)
	C.free(unsafe.Pointer(rsp))
	return str
}

func (s *Strategy) CreditBusiness(param *CreditParam) string {
	var length, errorCode C.int
	var cstr *C.char
	switch param.CreditType {
	case CreditBuying: // 融资买入
		cstr = C.emc_credit_buying_on_margin(C.int(param.PositionSrc), C.CString(param.Symbol), C.int(param.Volume), C.double(param.Price), C.int(param.OrderType), C.int(param.OrderDuration), C.int(param.OrderQualifier), C.CString(param.Account), &length, &errorCode)
	case CreditSelling: // 融卷卖出
		cstr = C.emc_credit_short_selling(C.int(param.PositionSrc), C.CString(param.Symbol), C.int(param.Volume), C.double(param.Price), C.int(param.OrderType), C.int(param.OrderDuration), C.int(param.OrderQualifier), C.CString(param.Account), &length, &errorCode)
	case BuyingShareRepayShare: // 买券还券
		cstr = C.emc_credit_repay_share_by_buying_share(C.CString(param.Symbol), C.int(param.Volume), C.double(param.Price), C.int(param.OrderType), C.int(param.OrderDuration), C.int(param.OrderQualifier), C.CString(param.Account), &length, &errorCode)
	case SellingShareRepay: // 卖券还款
		cstr = C.emc_credit_repay_cash_by_selling_share(C.CString(param.Symbol), C.int(param.Volume), C.double(param.Price), C.int(param.OrderType), C.int(param.OrderDuration), C.int(param.OrderQualifier), C.CString(param.Account), &length, &errorCode)
	case CollateralBuying: // 担保品买入
		cstr = C.emc_credit_buying_on_collateral(C.CString(param.Symbol), C.int(param.Volume), C.double(param.Price), C.int(param.OrderType), C.int(param.OrderDuration), C.int(param.OrderQualifier), C.CString(param.Account), &length, &errorCode)
	case CollateralSelling: // 担保品卖出
		C.emc_credit_selling_on_collateral(C.CString(param.Symbol), C.int(param.Volume), C.double(param.Price), C.int(param.OrderType), C.int(param.OrderDuration), C.int(param.OrderQualifier), C.CString(param.Account), &length, &errorCode)
	}
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// CreditRepayShareDirectly 直接还券
func (s *Strategy) CreditRepayShareDirectly(symbol, account string, volume int32) string {
	var length, errorCode C.int
	cstr := C.emc_credit_repay_share_directly(C.CString(symbol), C.int(volume), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// CreditRepayCashDirectly 直接还款
func (s *Strategy) CreditRepayCashDirectly(amount float64, account string) (err string, actualRepayAmount float64) {
	var actual C.double
	buf := (*C.char)(C.malloc(1024))
	result := C.emc_credit_repay_cash_directly(C.double(amount), C.CString(account), &actual, buf, 1024)
	if int32(result) != 0 {
		errMsg := C.GoString(buf)
		C.free(unsafe.Pointer(buf))
		return errMsg, 0.0
	}
	return "", float64(actual)
}

// CreditCollateralIn 担保品转入
func (s *Strategy) CreditCollateralIn(symbol, account string, volume int32) string {
	var length, errorCode C.int
	cstr := C.emc_credit_collateral_in(C.CString(symbol), C.int(volume), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// CreditCollateralOut 担保品转出
func (s *Strategy) CreditCollateralOut(symbol, account string, volume int32) string {
	var length, errorCode C.int
	cstr := C.emc_credit_collateral_out(C.CString(symbol), C.int(volume), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// CreditGetCollateralInstruments 查询担保证券
func (s *Strategy) CreditGetCollateralInstruments(account string) string {
	var length, errorCode C.int
	cstr := C.emc_credit_get_collateral_instruments(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// CreditGetBorrowableInstruments 查询标的证券，可做融券标的股票列表
func (s *Strategy) CreditGetBorrowableInstruments(positionSrc int32, account string) string {
	var length, errorCode C.int
	cstr := C.emc_credit_get_borrowable_instruments(C.int(positionSrc), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

func (s *Strategy) CreditGetBorrowableInstrumentsPositions(positionSrc int32, account string) string {
	var length, errorCode C.int
	cstr := C.emc_credit_get_borrowable_instruments_positions(C.int(positionSrc), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// CreditGetContracts 查询融资融券合约
func (s *Strategy) CreditGetContracts(positionSrc int32, account string) string {
	var length, errorCode C.int
	cstr := C.emc_credit_get_contracts(C.int(positionSrc), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// CreditGetCash 查询融资融券资金
func (s *Strategy) CreditGetCash(account string) string {
	var length, errorCode C.int
	cstr := C.emc_credit_get_cash(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// IpoBuy 新股新债申购
func (s *Strategy) IpoBuy(symbol, account string, volume int32, price float64) string {
	var length, errorCode C.int
	cstr := C.emc_ipo_buy(C.CString(symbol), C.int(volume), C.double(price), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// IpoGetQuota 查询客户新股新债申购额度
func (s *Strategy) IpoGetQuota(account string) string {
	var length, errorCode C.int
	cstr := C.emc_ipo_get_quota(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// IpoGetInstruments 查询当日新股新债清单
func (s *Strategy) IpoGetInstruments(securityType int32, account string) string {
	var length, errorCode C.int
	cstr := C.emc_ipo_get_instruments(C.int(securityType), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// ipoGetMatchNumber 配号查询
func (s *Strategy) ipoGetMatchNumber(account string) string {
	var length, errorCode C.int
	cstr := C.emc_ipo_get_match_number(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// IpoGetLotInfo 中签查询
func (s *Strategy) IpoGetLotInfo(account string) string {
	var length, errorCode C.int
	cstr := C.emc_ipo_get_lot_info(C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// FundEtfBuy ETF申购
func (s *Strategy) FundEtfBuy(symbol, account string, volume int32, price float64) string {
	var length, errorCode C.int
	cstr := C.emc_fund_etf_buy(C.CString(symbol), C.int(volume), C.double(price), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// FundEtfRedemption ETF赎回
func (s *Strategy) FundEtfRedemption(symbol, account string, volume int32, price float64) string {
	var length, errorCode C.int
	cstr := C.emc_fund_etf_redemption(C.CString(symbol), C.int(volume), C.double(price), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// FundSubscribing 基金认购
func (s *Strategy) FundSubscribing(symbol, account string, volume int32) string {
	var length, errorCode C.int
	cstr := C.emc_fund_subscribing(C.CString(symbol), C.int(volume), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// FundBuy 基金申购
func (s *Strategy) FundBuy(symbol, account string, volume int32) string {
	var length, errorCode C.int
	cstr := C.emc_fund_buy(C.CString(symbol), C.int(volume), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// FundRedemption 基金赎回
func (s *Strategy) FundRedemption(symbol, account string, volume int32) string {
	var length, errorCode C.int
	cstr := C.emc_fund_redemption(C.CString(symbol), C.int(volume), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// BondReverseRepurchaseAgreement 国债逆回购
func (s *Strategy) BondReverseRepurchaseAgreement(param *NationalDebtParam) string {
	var length, errorCode C.int
	cstr := C.emc_bond_reverse_repurchase_agreement(C.CString(param.Symbol), C.int(param.Volume), C.double(param.Price),
		C.int(param.OrderType), C.int(param.OrderDuration), C.int(param.OrderQualifier), C.CString(param.Account),
		&length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// BondConvertibleCall 可转债转股
func (s *Strategy) BondConvertibleCall(symbol, account string, volume int32, price float64) string {
	var length, errorCode C.int
	cstr := C.emc_bond_convertible_call(C.CString(symbol), C.int(volume), C.double(price), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// BondConvertiblePut 可转债回售
func (s *Strategy) BondConvertiblePut(symbol, account string, volume int32, price float64) string {
	var length, errorCode C.int
	cstr := C.emc_bond_convertible_put(C.CString(symbol), C.int(volume), C.double(price), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// BondConvertiblePutCancel 可转债回售撤销
func (s *Strategy) BondConvertiblePutCancel(symbol, account string, volume int32) string {
	var length, errorCode C.int
	cstr := C.emc_bond_convertible_put_cancel(C.CString(symbol), C.int(volume), C.CString(account), &length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// AddParameters 添加参数
func (s *Strategy) AddParameters(param *Parameter) int32 {
	var p C.Parameter
	p.value = C.double(param.Value)
	p.min = C.double(param.Min)
	p.max = C.double(param.Max)
	p.readonly = C.bool(param.Readonly)
	C.memcpy(unsafe.Pointer(&p.key), unsafe.Pointer(C.CString(param.Key)), C.size_t(len(param.Key)+1))
	C.memcpy(unsafe.Pointer(&p.name), unsafe.Pointer(C.CString(param.Name)), C.size_t(len(param.Name)+1))
	C.memcpy(unsafe.Pointer(&p.intro), unsafe.Pointer(C.CString(param.Intro)), C.size_t(len(param.Intro)+1))
	C.memcpy(unsafe.Pointer(&p.group), unsafe.Pointer(C.CString(param.Group)), C.size_t(len(param.Group)+1))
	return int32(C.emc_add_parameters(&p, C.int(1)))
}

// DelParameters 删除参数
func (s *Strategy) DelParameters(keys string) int32 {
	return int32(C.emc_del_parameters(C.CString(keys)))
}

// SetParameters 设置参数
func (s *Strategy) SetParameters(param *Parameter) int32 {
	var p C.Parameter
	p.value = C.double(param.Value)
	p.min = C.double(param.Min)
	p.max = C.double(param.Max)
	p.readonly = C.bool(param.Readonly)
	C.memcpy(unsafe.Pointer(&p.key), unsafe.Pointer(C.CString(param.Key)), C.size_t(len(param.Key)+1))
	C.memcpy(unsafe.Pointer(&p.name), unsafe.Pointer(C.CString(param.Name)), C.size_t(len(param.Name)+1))
	C.memcpy(unsafe.Pointer(&p.intro), unsafe.Pointer(C.CString(param.Intro)), C.size_t(len(param.Intro)+1))
	C.memcpy(unsafe.Pointer(&p.group), unsafe.Pointer(C.CString(param.Group)), C.size_t(len(param.Group)+1))
	return int32(C.emc_set_parameters(&p, C.int(1)))
}

func (s *Strategy) GetParameters() string {
	var length, errorCode C.int
	cstr := C.emc_get_parameters(&length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

func (s *Strategy) SetSymbols(symbols string) int32 {
	return int32(C.emc_set_symbols(C.CString(symbols)))
}

func (s *Strategy) GetSymbols() string {
	var length, errorCode C.int
	cstr := C.emc_get_symbols(&length, &errorCode)
	if int32(length) > 0 && int32(errorCode) == 0 {
		str := C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		return str
	}
	return ""
}

// GetVersion 获取sdk版本号
func GetVersion() string {
	return C.GoString(C.emc_get_version())
}

// SetToken 设置token
func SetToken(token string) {
	C.emc_set_token(C.CString(token))
}

// SetServerAddress 自定义服务地址
func SetServerAddress(address string) {
	C.emc_set_serv_addr(C.CString(address))
}

// SetMfp 第三方系统设置留痕信息
func SetMfp(mfp string) {
	C.emc_set_mfp(C.CString(mfp))
}

// Current 查询当前行情快照
func (s *Strategy) Current(symbols string) []*Tick {
	var length, errorCode C.int
	ticks := C.emc_current(C.CString(symbols), &length, &errorCode)
	var ts []*Tick
	var size C.Tick
	for i := 0; i < int(length); i++ {
		var t Tick
		tick := *(**C.Tick)(unsafe.Pointer(uintptr(unsafe.Pointer(ticks)) + (unsafe.Sizeof(size) * uintptr(i))))
		t.Symbol = C.GoString(&tick.symbol[0])
		t.CreatedAt = float64(tick.created_at)
		t.Price = float32(tick.price)
		t.Open = float32(tick.open)
		t.High = float32(tick.high)
		t.Low = float32(tick.low)
		t.CumVolume = float64(tick.cum_volume)
		t.CumAmount = float64(tick.cum_amount)
		t.CumPosition = int64(tick.cum_position)
		t.LastAmount = float64(tick.last_amount)
		t.LastVolume = int32(tick.last_volume)
		t.TradeType = int32(tick.trade_type)
		for j := 0; j < 10; j++ {
			var q Quote
			q.BidPrice = float32(tick.quotes[j].bid_price)
			q.BidVolume = int64(tick.quotes[j].bid_volume)
			q.AskPrice = float32(tick.quotes[j].ask_price)
			q.AskVolume = int64(tick.quotes[j].ask_volume)
			t.Quotes = append(t.Quotes, &q)
		}
		ts = append(ts, &t)
	}
	C.emc_ticks_free(ticks, C.int(length))
	return ts
}

func (s *Strategy) HistoryTicks(param *HistoryTickParam) []*Tick {
	var length, errorCode C.int
	var ticks **C.Tick
	if param.N == 0 {
		ticks = C.emc_history_ticks(C.CString(param.Symbols), C.CString(param.StartTime), C.CString(param.EndTime), C.int(param.Adjust),
			C.CString(param.AdjustEndTime), C.bool(param.SkipSuspended), C.CString(param.FillMissing), &length, &errorCode)
	} else {
		ticks = C.emc_history_ticks_n(C.CString(param.Symbols), C.int(param.N), C.CString(param.EndTime), C.int(param.Adjust), C.CString(param.AdjustEndTime), C.bool(param.SkipSuspended), C.CString(param.FillMissing), &length, &errorCode)
	}

	var ts []*Tick
	var size C.Tick
	for i := 0; i < int(length); i++ {
		var t Tick
		tick := *(**C.Tick)(unsafe.Pointer(uintptr(unsafe.Pointer(ticks)) + (unsafe.Sizeof(size) * uintptr(i))))
		t.Symbol = C.GoString(&tick.symbol[0])
		t.CreatedAt = float64(tick.created_at)
		t.Price = float32(tick.price)
		t.Open = float32(tick.open)
		t.High = float32(tick.high)
		t.Low = float32(tick.low)
		t.CumVolume = float64(tick.cum_volume)
		t.CumAmount = float64(tick.cum_amount)
		t.CumPosition = int64(tick.cum_position)
		t.LastAmount = float64(tick.last_amount)
		t.LastVolume = int32(tick.last_volume)
		t.TradeType = int32(tick.trade_type)
		for j := 0; j < 10; j++ {
			var q Quote
			q.BidPrice = float32(tick.quotes[j].bid_price)
			q.BidVolume = int64(tick.quotes[j].bid_volume)
			q.AskPrice = float32(tick.quotes[j].ask_price)
			q.AskVolume = int64(tick.quotes[j].ask_volume)
			t.Quotes = append(t.Quotes, &q)
		}
		ts = append(ts, &t)
	}
	C.emc_ticks_free(ticks, C.int(length))
	return ts
}

func (s *Strategy) HistoryBars(param *HistoryBarParam) []*Bar {
	var length, errorCode C.int
	var bars **C.Bar
	if param.N == 0 {
		bars = C.emc_history_bars(C.CString(param.Symbols), C.CString(param.Frequency), C.CString(param.StartTime), C.CString(param.EndTime), C.int(param.Adjust),
			C.CString(param.AdjustEndTime), C.bool(param.SkipSuspended), C.CString(param.FillMissing), &length, &errorCode)
	} else {
		bars = C.emc_history_bars_n(C.CString(param.Symbols), C.CString(param.Frequency), C.int(param.N), C.CString(param.EndTime), C.int(param.Adjust),
			C.CString(param.AdjustEndTime), C.bool(param.SkipSuspended), C.CString(param.FillMissing), &length, &errorCode)
	}

	var bs []*Bar
	var size C.Bar
	for i := 0; i < int(length); i++ {
		var b Bar
		bar := *(**C.Bar)(unsafe.Pointer(uintptr(unsafe.Pointer(bars)) + (unsafe.Sizeof(size) * uintptr(i))))
		b.Symbol = C.GoString(&bar.symbol[0])
		b.Bob = float64(bar.bob)
		b.Eob = float64(bar.eob)
		b.Open = float32(bar.open)
		b.Close = float32(bar.close)
		b.High = float32(bar.high)
		b.Low = float32(bar.low)
		b.Volume = float64(bar.volume)
		b.Amount = float64(bar.amount)
		b.PreClose = float32(bar.pre_close)
		b.Position = int64(bar.position)
		b.Frequency = C.GoString(&bar.frequency[0])
		bs = append(bs, &b)
	}
	C.emc_bars_free(bars, C.int(length))
	return bs
}

func (s *Strategy) HistoryL2Ticks(param *HistoryTickParam) []*Tick {
	var length, errorCode C.int
	var ticks **C.Tick
	ticks = C.emc_history_l2ticks(C.CString(param.Symbols), C.CString(param.StartTime), C.CString(param.EndTime), C.int(param.Adjust),
		C.CString(param.AdjustEndTime), C.bool(param.SkipSuspended), C.CString(param.FillMissing), &length, &errorCode)

	var ts []*Tick
	var size C.Tick
	for i := 0; i < int(length); i++ {
		var t Tick
		tick := *(**C.Tick)(unsafe.Pointer(uintptr(unsafe.Pointer(ticks)) + (unsafe.Sizeof(size) * uintptr(i))))
		t.Symbol = C.GoString(&tick.symbol[0])
		t.CreatedAt = float64(tick.created_at)
		t.Price = float32(tick.price)
		t.Open = float32(tick.open)
		t.High = float32(tick.high)
		t.Low = float32(tick.low)
		t.CumVolume = float64(tick.cum_volume)
		t.CumAmount = float64(tick.cum_amount)
		t.CumPosition = int64(tick.cum_position)
		t.LastAmount = float64(tick.last_amount)
		t.LastVolume = int32(tick.last_volume)
		t.TradeType = int32(tick.trade_type)
		for j := 0; j < 10; j++ {
			var q Quote
			q.BidPrice = float32(tick.quotes[j].bid_price)
			q.BidVolume = int64(tick.quotes[j].bid_volume)
			q.AskPrice = float32(tick.quotes[j].ask_price)
			q.AskVolume = int64(tick.quotes[j].ask_volume)
			t.Quotes = append(t.Quotes, &q)
		}
		ts = append(ts, &t)
	}
	C.emc_ticks_free(ticks, C.int(length))
	return ts
}

func (s *Strategy) HistoryL2Bars(param *HistoryBarParam) []*Bar {
	var length, errorCode C.int
	var bars **C.Bar
	bars = C.emc_history_bars(C.CString(param.Symbols), C.CString(param.Frequency), C.CString(param.StartTime), C.CString(param.EndTime), C.int(param.Adjust),
		C.CString(param.AdjustEndTime), C.bool(param.SkipSuspended), C.CString(param.FillMissing), &length, &errorCode)

	var bs []*Bar
	var size C.Bar
	for i := 0; i < int(length); i++ {
		var b Bar
		bar := *(**C.Bar)(unsafe.Pointer(uintptr(unsafe.Pointer(bars)) + (unsafe.Sizeof(size) * uintptr(i))))
		b.Symbol = C.GoString(&bar.symbol[0])
		b.Bob = float64(bar.bob)
		b.Eob = float64(bar.eob)
		b.Open = float32(bar.open)
		b.Close = float32(bar.close)
		b.High = float32(bar.high)
		b.Low = float32(bar.low)
		b.Volume = float64(bar.volume)
		b.Amount = float64(bar.amount)
		b.PreClose = float32(bar.pre_close)
		b.Position = int64(bar.position)
		b.Frequency = C.GoString(&bar.frequency[0])
		bs = append(bs, &b)
	}
	C.emc_bars_free(bars, C.int(length))
	return bs
}

func (s *Strategy) HistoryL2Transactions(symbols, startTime, endTime string) []*L2Transaction {
	var length, errorCode C.int
	var transactions **C.L2Transaction
	transactions = C.emc_history_l2transactions(C.CString(symbols), C.CString(startTime), C.CString(endTime), &length, &errorCode)

	var l2Transactions []*L2Transaction
	var size C.L2Transaction
	for i := 0; i < int(length); i++ {
		var l2 L2Transaction
		l2transaction := *(**C.L2Transaction)(unsafe.Pointer(uintptr(unsafe.Pointer(transactions)) + (unsafe.Sizeof(size) * uintptr(i))))
		l2.Symbol = C.GoString(&l2transaction.symbol[0])
		l2.CreatedAt = float64(l2transaction.created_at)
		l2.Price = float32(l2transaction.price)
		l2.Volume = int64(l2transaction.volume)
		l2.Side = string(l2transaction.side)
		l2.ExecType = string(l2transaction.exec_type)
		l2.ExecIndex = int64(l2transaction.exec_index)
		l2.AskOrderIndex = int64(l2transaction.ask_order_index)
		l2.BidOrderIndex = int64(l2transaction.bid_order_index)
		l2Transactions = append(l2Transactions, &l2)
	}
	C.emc_l2transactions_free(transactions, C.int(length))
	return l2Transactions
}

func (s *Strategy) HistoryL2Orders(symbols, startTime, endTime string) []*L2Order {
	var length, errorCode C.int
	var orders **C.L2Order
	orders = C.emc_history_l2orders(C.CString(symbols), C.CString(startTime), C.CString(endTime), &length, &errorCode)

	var l2Orders []*L2Order
	var size C.L2Order
	for i := 0; i < int(length); i++ {
		var l2 L2Order
		l2order := *(**C.L2Order)(unsafe.Pointer(uintptr(unsafe.Pointer(orders)) + (unsafe.Sizeof(size) * uintptr(i))))
		l2.Symbol = C.GoString(&l2order.symbol[0])
		l2.CreatedAt = float64(l2order.created_at)
		l2.Price = float32(l2order.price)
		l2.Volume = int64(l2order.volume)
		l2.Side = string(l2order.side)
		l2.OrderType = string(l2order.order_type)
		l2.OrderIndex = int64(l2order.order_index)
		l2Orders = append(l2Orders, &l2)
	}
	C.emc_l2orders_free(orders, C.int(length))
	return l2Orders
}

func (s *Strategy) HistoryL2OrdersQueue(symbols, startTime, endTime string) []*L2OrderQueue {
	var length, errorCode C.int
	var orders **C.L2OrderQueue
	orders = C.emc_history_l2orders_queue(C.CString(symbols), C.CString(startTime), C.CString(endTime), &length, &errorCode)

	var l2Orders []*L2OrderQueue
	var size C.L2OrderQueue
	for i := 0; i < int(length); i++ {
		var l2 L2OrderQueue
		l2queue := *(**C.L2OrderQueue)(unsafe.Pointer(uintptr(unsafe.Pointer(orders)) + (unsafe.Sizeof(size) * uintptr(i))))
		l2.Symbol = C.GoString(&l2queue.symbol[0])
		l2.CreatedAt = float64(l2queue.created_at)
		l2.Price = float32(l2queue.price)
		l2.Volume = int64(l2queue.volume)
		l2.Side = string(l2queue.side)
		l2.QueueOrders = int32(l2queue.queue_orders)
		for j := 0; j < int(l2.QueueOrders); j++ {
			l2.QueueVolumes = append(l2.QueueVolumes, int32(l2queue.queue_volumes[j]))
		}
	}
	C.emc_l2orders_queue_free(orders, C.int(length))
	return l2Orders
}

type ObjectId int64

func (o *ObjectId) FreeDataSet() {
	C.emc_free_dataset(C.longlong(*o))
}

func (o *ObjectId) DataSetStatus() int32 {
	return int32(C.emc_dataset_status(C.longlong(*o)))
}

func (o *ObjectId) DataSetNext() {
	C.emc_dataset_next(C.longlong(*o))
}

func (o *ObjectId) DataSetGetInteger(key string) int32 {
	return int32(C.emc_dataset_get_integer(C.longlong(*o), C.CString(key)))
}

func (o *ObjectId) DataSetGetLongInteger(key string) int64 {
	return int64(C.emc_dataset_get_long_integer(C.longlong(*o), C.CString(key)))
}

func (o *ObjectId) DataSetGetReal(key string) float64 {
	return float64(C.emc_dataset_get_real(C.longlong(*o), C.CString(key)))
}

func (o *ObjectId) DataSetGetString(key string) string {
	return C.GoString(C.emc_dataset_get_string(C.longlong(*o), C.CString(key)))
}

func (o *ObjectId) DataSetGetDebugString() string {
	return C.GoString(C.emc_dataset_debug_string(C.longlong(*o)))
}

func (s *Strategy) GetFundamentals(param *FundamentalParam) ObjectId {
	var errorCode C.int
	var oid C.longlong
	if param.N == 0 {
		C.emc_get_fundamentals(C.CString(param.Table), C.CString(param.Symbols), C.CString(param.StartDate), C.CString(param.EndDate), C.CString(param.Fields),
			C.CString(param.Filter), C.CString(param.OrderBy), C.int(param.Limit), &oid, &errorCode)
	} else {
		C.emc_get_fundamentals_n(C.CString(param.Table), C.CString(param.Symbols), C.CString(param.EndDate), C.CString(param.Fields),
			C.int(param.N), C.CString(param.Filter), C.CString(param.OrderBy), &oid, &errorCode)
	}
	if int32(errorCode) != 0 {
		return ObjectId(errorCode)
	}
	return ObjectId(oid)
}

// GetInstruments 查询最新交易标的信息
func (s *Strategy) GetInstruments(exchanges, secTypes, fields string) ObjectId {
	var errorCode C.int
	var oid C.longlong
	C.emc_get_instruments(C.CString(exchanges), C.CString(secTypes), C.CString(fields), &oid, &errorCode)

	if int32(errorCode) != 0 {
		return ObjectId(errorCode)
	}
	return ObjectId(oid)
}

// GetHistoryInstruments 查询交易标的历史数据
func (s *Strategy) GetHistoryInstruments(symbols, startDate, endDate, fields string) ObjectId {
	var errorCode C.int
	var oid C.longlong
	C.emc_get_history_instruments(C.CString(symbols), C.CString(startDate), C.CString(endDate), C.CString(fields), &oid, &errorCode)

	if int32(errorCode) != 0 {
		return ObjectId(errorCode)
	}
	return ObjectId(oid)
}

// GetInstrumentInfos 查询交易标的基本信息
func (s *Strategy) GetInstrumentInfos(symbols, exchanges, secTypes, names, fields string) ObjectId {
	var errorCode C.int
	var oid C.longlong
	C.emc_get_instrumentinfos(C.CString(symbols), C.CString(exchanges), C.CString(secTypes), C.CString(names), C.CString(fields), &oid, &errorCode)

	if int32(errorCode) != 0 {
		return ObjectId(0)
	}
	return ObjectId(errorCode)
}

// GetConstituents 查询指数成份股
func (s *Strategy) GetConstituents(index, tradeDate string) ObjectId {
	var errorCode C.int
	var oid C.longlong
	C.emc_get_constituents(C.CString(index), C.CString(tradeDate), &oid, &errorCode)

	if int32(errorCode) != 0 {
		return ObjectId(errorCode)
	}
	return ObjectId(oid)
}

// GetIndustry 查询行业股票列表
func (s *Strategy) GetIndustry(code string) string {
	var errorCode, length C.int
	cstr := C.emc_get_industry(C.CString(code), &length, &errorCode)

	if int32(errorCode) != 0 || length == 0 {
		return ""
	}
	return C.GoString(cstr)
}

// GetConcept 查询概念板块股票列表
func (s *Strategy) GetConcept(code string) string {
	var errorCode, length C.int
	cstr := C.emc_get_concept(C.CString(code), &length, &errorCode)

	if int32(errorCode) != 0 || length == 0 {
		return ""
	}
	return C.GoString(cstr)
}

// GetTradingDates 查询交易日列表
func (s *Strategy) GetTradingDates(exchange, startDate, endDate string) string {
	var errorCode, length C.int
	cstr := C.emc_get_trading_dates(C.CString(exchange), C.CString(startDate), C.CString(endDate), &length, &errorCode)

	if int32(errorCode) != 0 || length == 0 {
		return ""
	}
	return C.GoString(cstr)
}

// GetPreviousTradingDate 返回指定日期的上一个交易日
func (s *Strategy) GetPreviousTradingDate(exchange, date string) string {
	outDate := (*C.char)(C.malloc(1024))
	r := C.emc_get_previous_trading_date(C.CString(exchange), C.CString(date), outDate)

	if int32(r) != 0 {
		C.free(unsafe.Pointer(outDate))
		return ""
	}
	str := C.GoString(outDate)
	C.free(unsafe.Pointer(outDate))
	return str
}

// GetNextTradingDate 返回指定日期的上一个交易日
func (s *Strategy) GetNextTradingDate(exchange, date string) string {
	outDate := (*C.char)(C.malloc(1024))
	r := C.emc_get_next_trading_date(C.CString(exchange), C.CString(date), outDate)
	if int32(r) != 0 {
		C.free(unsafe.Pointer(outDate))
		return ""
	}
	str := C.GoString(outDate)
	C.free(unsafe.Pointer(outDate))
	return str
}

// GetDividend 查询分红送配
func (s *Strategy) GetDividend(symbol, startDate, endDate string) ObjectId {
	var errorCode C.int
	var oid C.longlong
	C.emc_get_dividend(C.CString(symbol), C.CString(startDate), C.CString(endDate), &oid, &errorCode)

	if int32(errorCode) != 0 {
		return ObjectId(errorCode)
	}
	return ObjectId(oid)
}

// GetContinuousContracts 获取连续合约
func (s *Strategy) GetContinuousContracts(symbol, startDate, endDate string) ObjectId {
	var errorCode C.int
	var oid C.longlong
	C.emc_get_continuous_contracts(C.CString(symbol), C.CString(startDate), C.CString(endDate), &oid, &errorCode)

	if int32(errorCode) != 0 {
		return ObjectId(errorCode)
	}
	return ObjectId(oid)
}
