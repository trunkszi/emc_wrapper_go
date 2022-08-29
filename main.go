package main

import (
	"fmt"
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"wrapper"
)

var (
	logger *zap.Logger
)

func init() {
	dir, _ := os.Getwd()
	fileName := dir + "/log/emc.log"
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	log := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(log)
	logger = zap.L()
}

type GoStrategy struct {
	wrapper.Strategy
	Weight []float32
	Band   []float32
}

func (s *GoStrategy) GoInit() {
	//s.Weight = []float32{0.5, 0.3, 0.0, 0.3, 0.5}
	//s.Band = []float32{-40.0, -3.0, -2.0, 2.0, 3.0, 40.0}

	s.Subscribe("SHSE.603501", "tick", false)
	//s.Subscribe("SHSE.603501", "60s", false)

	//param := wrapper.HistoryBarParam{
	//	Symbols:       "SHSE.603501",
	//	StartTime:     "2022-01-01 09:00:00",
	//	Frequency:     "60s",
	//	N:             300,
	//	AdjustEndTime: "",
	//	FillMissing:   "Last",
	//	Adjust:        1,
	//	SkipSuspended: false,
	//}
	//var timeSeries []float32
	//var avg float32
	//var count int32
	//bars := s.HistoryBars(&param)
	//for _, bar := range bars {
	//	timeSeries = append(timeSeries, bar.Close)
	//	avg += bar.Close
	//	count++
	//}
	//mean := avg / float32(count)
	//var accum float32 = 0.0
	//
	//for _, d := range timeSeries {
	//	accum += (d - mean) * (d - mean)
	//}
	//dev := math.Sqrt(float64(accum / float32(len(timeSeries)-1)))
	//logger.Info(zap.Float64("dev", dev).String)
	//
	//for i := 0; i < len(s.Band); i++ {
	//	s.Band[i] = mean + s.Band[i]*float32(dev)
	//}
}

func (s *GoStrategy) GoTick(tick *wrapper.Tick) {
	output, _ := sonic.Marshal(tick)
	logger.Info(zap.String("tick", string(output)).String)
}

func (s *GoStrategy) GoBar(bar *wrapper.Bar) {
	//fmt.Println("tick: ", tick)
	fmt.Printf("%+v\n", bar)
	logger.Debug("bar: ", zap.Any("bar", bar))
}
func (s *GoStrategy) Gol2Transaction(l2transaction *wrapper.L2Transaction) {
	fmt.Println("GoL2Transaction")
	fmt.Println("l2transaction: ", l2transaction)
}
func (s *GoStrategy) Gol2Order(l2order *wrapper.L2Order) {
	fmt.Println("GoL2Order")
	fmt.Println("l2order: ", l2order)
}
func (s *GoStrategy) GoL2OrderQueue(l2queue *wrapper.L2OrderQueue) {
	fmt.Println("GoL2OrderQueue")
	fmt.Println("l2queue: ", l2queue)
}
func (s *GoStrategy) GoOrderStatus(order *wrapper.Order) {
	fmt.Println("GoOrderStatus")
	fmt.Println("order: ", order)
}
func (s *GoStrategy) GoExecutionReport(rpt *wrapper.ExecRpt) {
	fmt.Println("GoExecutionReport")
	fmt.Println("rpt: ", rpt)
}
func (s *GoStrategy) GoAlgoOrderStatus(order *wrapper.AlgoOrder) {
	fmt.Println("GoAlgoOrderStatus")
	fmt.Println("order: ", order)
}
func (s *GoStrategy) GoCash(cash *wrapper.Cash) {
	fmt.Println("GoCash")
	fmt.Println("cash: ", cash)
}
func (s *GoStrategy) GoPosition(position *wrapper.Position) {
	fmt.Println("GoPosition")
	fmt.Println("position: ", position)
}
func (s *GoStrategy) GoParameter(param *wrapper.Parameter) {
	fmt.Println("GoParameter")
	fmt.Println("param: ", param)
}
func (s *GoStrategy) GoSchedule(dataRule string, timeRule string) {
	fmt.Println("GoSchedule")
	fmt.Println("dataRule: ", dataRule)
	fmt.Println("timeRule: ", timeRule)
}
func (s *GoStrategy) GoBackTestFinished(indicator *wrapper.Indicator) {
	fmt.Println("GoBackTestFinished")
	fmt.Printf("%+v\n", indicator)
	//fmt.Println("indicator: ", indicator)
}
func (s *GoStrategy) GoAccountStatus(accountStatus *wrapper.AccountStatus) {
	fmt.Println("GoAccountStatus")
	fmt.Println("accountStatus: ", accountStatus)
}
func (s *GoStrategy) GoError(errorCode int, errorMsg string) {
	fmt.Println("GoError")
	fmt.Println("errorCode: ", errorCode)
	fmt.Println("errorMsg: ", errorMsg)
}
func (s *GoStrategy) GoStop() {
	fmt.Println("GoStop")
}

func (s *GoStrategy) GoMarketDataConnected() {
	fmt.Println("GoMarketDataConnected")
}

func (s *GoStrategy) GoTradeDataConnected() {
	fmt.Println("GoTradeDataConnected")
}

func (s *GoStrategy) GoMarketDataDisconnected() {
	fmt.Println("GoMarketDataDisconnected")
}

func (s *GoStrategy) GoTradeDataDisconnected() {
	fmt.Println("GoTradeDataDisconnected")
}

func main() {
	defer logger.Sync()
	s := &GoStrategy{}
	wrapper.RegisterStrategy(s)
	s.SetStrategyId("016fef2e-192f-11ed-9e30-d45d64a6d10d")
	s.SetToken("1a0119ff3c508a2bd60ff086c577cb2c45ffe0d0")
	s.SetMode(2)
	s.SetBackTestConfig("2022-01-01 09:00:00",
		"2022-8-11 16:00:00",
		10000000.0,
		1,
		0.0001,
		0.0001,
		1,
		1)
	s.Run()
	select {}
}
