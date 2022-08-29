#include "interface.h"
#include "gmapi.h"
#include "strategy.h"
#include <cstdlib>
#include <memory>
#include <vector>
#include <string>
#include <chrono>

#include <type_traits>
#include <map>
#include "co/json.h"
#include "co/all.h"

bool is_include_chinese(const char *source) {
    char c;
    for (int i = 0; i < strlen(source) + 1; i++) {
        c = source[i];
        if (c == 0) {
            break;
        }
        if (c & 0x80) {
            if (source[i + 1] & 0x80) {
                return true;
            }
        }
    }
    return false;
}

std::string format_time(double duration) {
    std::chrono::duration<double> milliseconds = std::chrono::duration<double>(duration);
    auto sec = std::chrono::duration_cast<std::chrono::seconds>(milliseconds);

    auto tp = std::chrono::time_point<std::chrono::system_clock>(sec); // 将秒数转换为时间点
    std::time_t tm = std::chrono::system_clock::to_time_t(tp); // 将时间点转换为时间戳
    std::stringstream ss;
    ss << std::put_time(std::localtime(&tm), "%F %T");
    return std::move(ss.str());
}



// 摘自 https://programmer.group/c-realizes-the-conversion-of-utf8-and-gbk-coded-strings.html
#ifdef _WIN32

#include <windows.h>

std::string gbk_to_utf8(const char *src_str) {
    if (!is_include_chinese(src_str)) {
        return src_str;
    }
    int len = MultiByteToWideChar(CP_ACP, 0, src_str, -1, NULL, 0);
    wchar_t *wstr = new wchar_t[len + 1];
    memset(wstr, 0, len + 1);
    MultiByteToWideChar(CP_ACP, 0, src_str, -1, wstr, len);
    len = WideCharToMultiByte(CP_UTF8, 0, wstr, -1, NULL, 0, NULL, NULL);
    char *str = new char[len + 1];
    memset(str, 0, len + 1);
    WideCharToMultiByte(CP_UTF8, 0, wstr, -1, str, len, NULL, NULL);
    std::string strTemp = str;
    if (wstr) delete[] wstr;
    if (str) delete[] str;
    return strTemp;
}

std::string utf8_to_gbk(const char *src_str) {
    int len = MultiByteToWideChar(CP_UTF8, 0, src_str, -1, NULL, 0);
    wchar_t *wszGBK = new wchar_t[len + 1];
    memset(wszGBK, 0, len * 2 + 2);
    MultiByteToWideChar(CP_UTF8, 0, src_str, -1, wszGBK, len);
    len = WideCharToMultiByte(CP_ACP, 0, wszGBK, -1, NULL, 0, NULL, NULL);
    char *szGBK = new char[len + 1];
    memset(szGBK, 0, len + 1);
    WideCharToMultiByte(CP_ACP, 0, wszGBK, -1, szGBK, len, NULL, NULL);
    std::string strTemp(szGBK);
    if (wszGBK) delete[] wszGBK;
    if (szGBK) delete[] szGBK;
    return strTemp;
}

#else
#include <iconv.h>
int GbkToUtf8(char *str_str, size_t src_len, char *dst_str, size_t dst_len)
{
    iconv_t cd;
    char **pin = &str_str;
    char **pout = &dst_str;

    cd = iconv_open("utf8", "gbk");
    if (cd == 0)
        return -1;
    memset(dst_str, 0, dst_len);
    if (iconv(cd, pin, &src_len, pout, &dst_len) == -1)
        return -1;
    iconv_close(cd);
    *pout = '\0';

    return 0;
}
#endif

Json order_to_json(Order &order) {
    auto json = json::object();
    json.add_member("strategy_id", gbk_to_utf8(order.strategy_id));
    json.add_member("account_id", gbk_to_utf8(order.account_id));
    json.add_member("account_name", gbk_to_utf8(order.account_name));
    json.add_member("cl_ord_id", gbk_to_utf8(order.cl_ord_id));
    json.add_member("order_id", gbk_to_utf8(order.order_id));
    json.add_member("ex_ord_id", gbk_to_utf8(order.ex_ord_id));
    json.add_member("algo_order_id", order.algo_order_id);
    json.add_member("order_business", order.order_business);

    json.add_member("symbol", gbk_to_utf8(order.symbol));
    json.add_member("side", order.side);
    json.add_member("position_effect", order.position_effect);
    json.add_member("position_side", order.position_side);

    json.add_member("order_type", order.order_type);
    json.add_member("order_duration", order.order_duration);
    json.add_member("order_qualifier", order.order_qualifier);
    json.add_member("order_src", order.order_src);
    json.add_member("position_src", order.position_src);

    json.add_member("status", order.status);
    json.add_member("ord_rej_reason", order.ord_rej_reason);
    json.add_member("ord_rej_reason_detail", gbk_to_utf8(order.ord_rej_reason_detail));

    json.add_member("price", order.price);
    json.add_member("stop_price", order.stop_price);

    json.add_member("order_style", order.order_style);
    json.add_member("volume", order.volume);
    json.add_member("value", order.value);
    json.add_member("percent", order.percent);
    json.add_member("target_volume", order.target_volume);
    json.add_member("target_value", order.target_value);
    json.add_member("target_percent", order.target_percent);

    json.add_member("filled_volume", order.filled_volume);
    json.add_member("filled_vwap", order.filled_vwap);
    json.add_member("filled_amount", order.filled_amount);
    json.add_member("filled_commission", order.filled_commission);

    json.add_member("created_at", format_time(order.created_at));
    json.add_member("updated_at", format_time(order.updated_at));

    return std::move(json);
}


class EmtStrategy : public Strategy {
private:
    on_init_callback init_callback;
    on_tick_callback tick_callback;
    on_bar_callback bar_callback;
    on_l2transaction_callback l2transaction_callback;
    on_l2order_callback l2order_callback;
    on_l2order_queue_callback l2order_queue_callback;
    on_order_status_callback order_status_callback;
    on_execution_report_callback execution_report_callback;
    on_algo_order_status_callback algo_order_status_callback;
    on_cash_callback cash_callback;
    on_position_callback position_callback;
    on_parameter_callback parameter_callback;
    on_schedule_callback schedule_callback;
    on_backtest_finished_callback backtest_finished_callback;
    on_account_status_callback account_status_callback;
    on_error_callback error_callback;
    on_stop_callback stop_callback;
    on_market_data_connected_callback market_data_connected_callback;
    on_trade_data_connected_callback trade_data_connected_callback;
    on_market_data_disconnected_callback market_data_disconnected_callback;
    on_trade_data_disconnected_callback trade_data_disconnected_callback;
public:
    EmtStrategy();

    ~EmtStrategy();

    void
    register_callbacks(on_init_callback init_callback, on_tick_callback tick_callback, on_bar_callback bar_callback,
                       on_l2transaction_callback l2transaction_callback, on_l2order_callback l2order_callback,
                       on_l2order_queue_callback l2order_queue_callback, on_order_status_callback order_status_callback,
                       on_execution_report_callback execution_report_callback,
                       on_algo_order_status_callback algo_order_status_callback, on_cash_callback cash_callback,
                       on_position_callback position_callback, on_parameter_callback parameter_callback,
                       on_schedule_callback schedule_callback, on_backtest_finished_callback backtest_finished_callback,
                       on_account_status_callback account_status_callback, on_error_callback error_callback,
                       on_stop_callback stop_callback, on_market_data_connected_callback market_data_connected_callback,
                       on_trade_data_connected_callback trade_data_connected_callback,
                       on_market_data_disconnected_callback market_data_disconnected_callback,
                       on_trade_data_disconnected_callback trade_data_disconnected_callback) {
        this->init_callback = init_callback;
        this->tick_callback = tick_callback;
        this->bar_callback = bar_callback;
        this->l2transaction_callback = l2transaction_callback;
        this->l2order_callback = l2order_callback;
        this->l2order_queue_callback = l2order_queue_callback;
        this->order_status_callback = order_status_callback;
        this->execution_report_callback = execution_report_callback;
        this->algo_order_status_callback = algo_order_status_callback;
        this->cash_callback = cash_callback;
        this->position_callback = position_callback;
        this->parameter_callback = parameter_callback;
        this->schedule_callback = schedule_callback;
        this->backtest_finished_callback = backtest_finished_callback;
        this->account_status_callback = account_status_callback;
        this->error_callback = error_callback;
        this->stop_callback = stop_callback;
        this->market_data_connected_callback = market_data_connected_callback;
        this->trade_data_connected_callback = trade_data_connected_callback;
        this->market_data_disconnected_callback = market_data_disconnected_callback;
        this->trade_data_disconnected_callback = trade_data_disconnected_callback;
    }

    //初始化完成
    virtual void on_init() override {
        init_callback();
    }

    //收到Tick行情
    virtual void on_tick(Tick *tick) override {
        if (tick == nullptr) {
            return;
        }
        Tick *t = (Tick *) malloc(sizeof(Tick));
        memcpy(t, tick, sizeof(Tick));
        tick_callback(tick);
    }

    //收到bar行情
    virtual void on_bar(Bar *bar) override {
        if (bar == nullptr) {
            return;
        }
        Bar *b = (Bar *) malloc(sizeof(Bar));
        memcpy(b, bar, sizeof(Bar));
        bar_callback(bar);
    }

    //收到逐笔成交（L2行情时有效）
    virtual void on_l2transaction(L2Transaction *l2transaction) override {
        if (l2transaction == nullptr) {
            return;
        }
        L2Transaction *t = (L2Transaction *) malloc(sizeof(L2Transaction));
        memcpy(t, l2transaction, sizeof(L2Transaction));
        l2transaction_callback(l2transaction);
    }

    //收到逐笔委托（深交所L2行情时有效）
    virtual void on_l2order(L2Order *l2order) override {
        if (l2order == nullptr) {
            return;
        }
        L2Order *o = (L2Order *) malloc(sizeof(L2Order));
        memcpy(o, l2order, sizeof(L2Order));
        l2order_callback(l2order);
    }

    //收到委托队列（上交所L2行情时有效）
    virtual void on_l2order_queue(L2OrderQueue *l2queue) override {
        if (l2queue == nullptr) {
            return;
        }
        L2OrderQueue *q = (L2OrderQueue *) malloc(sizeof(L2OrderQueue));
        memcpy(q, l2queue, sizeof(L2OrderQueue));
        l2order_queue_callback(l2queue);
    }

    //委托变化
    virtual void on_order_status(Order *order) override {
        if (order == nullptr) {
            return;
        }
        Order *o = (Order *) malloc(sizeof(Order));
        memcpy(o, order, sizeof(Order));
        order_status_callback(order);
    }

    //执行回报
    virtual void on_execution_report(ExecRpt *rpt) override {
        if (rpt == nullptr) {
            return;
        }
        ExecRpt *r = (ExecRpt *) malloc(sizeof(ExecRpt));
        memcpy(r, rpt, sizeof(ExecRpt));
        execution_report_callback(rpt);
    }

    //算法委托变化
    virtual void on_algo_order_status(AlgoOrder *order) override {
        if (order == nullptr) {
            return;
        }
        AlgoOrder *o = (AlgoOrder *) malloc(sizeof(AlgoOrder));
        memcpy(o, order, sizeof(AlgoOrder));
        algo_order_status_callback(order);
    }

    //资金推送
    virtual void on_cash(Cash *cash) override {
        if (cash == nullptr) {
            return;
        }
        Cash *c = (Cash *) malloc(sizeof(Cash));
        memcpy(c, cash, sizeof(Cash));
        cash_callback(cash);
    }

    //持仓推送
    virtual void on_position(Position *position) override {
        if (position == nullptr) {
            return;
        }
        Position *p = (Position *) malloc(sizeof(Position));
        memcpy(p, position, sizeof(Position));
        position_callback(position);
    }

    //参数变化
    virtual void on_parameter(Parameter *param) override {
        if (param == nullptr) {
            return;
        }
        Parameter *p = (Parameter *) malloc(sizeof(Parameter));
        memcpy(p, param, sizeof(Parameter));
        parameter_callback(param);
    }

    //定时任务触发
    virtual void on_schedule(const char *data_rule, const char *time_rule) override {
        schedule_callback((char *) data_rule, (char *) time_rule);
    }

    //回测完成后收到绩效报告
    virtual void on_backtest_finished(Indicator *indicator) override {
        if (indicator == nullptr) {
            return;
        }
        Indicator *i = (Indicator *) malloc(sizeof(Indicator));
        memcpy(i, indicator, sizeof(Indicator));
        backtest_finished_callback(indicator);
    }

    //实盘账号状态变化
    virtual void on_account_status(AccountStatus *account_status) override {
        if (account_status == nullptr) {
            return;
        }
        AccountStatus *a = (AccountStatus *) malloc(sizeof(AccountStatus));
        memcpy(a, account_status, sizeof(AccountStatus));
        account_status_callback(account_status);
    }

    //错误产生
    virtual void on_error(int error_code, const char *error_msg) override {
        error_callback(error_code, (char *) error_msg);
    }

    //收到策略停止信号
    virtual void on_stop() override {
        stop_callback();
    }

    //数据已经连接上
    virtual void on_market_data_connected() override {
        market_data_connected_callback();
    }

    //交易已经连接上
    virtual void on_trade_data_connected() override {
        trade_data_connected_callback();
    }

    //数据连接断开了
    virtual void on_market_data_disconnected() override {
        market_data_disconnected_callback();
    }

    //交易连接断开了
    virtual void on_trade_data_disconnected() override {
        trade_data_disconnected_callback();
    }
};

EmtStrategy::EmtStrategy() {}

EmtStrategy::~EmtStrategy() {}


std::shared_ptr<EmtStrategy> s;

void register_callbacks(on_init_callback init_callback,
                        on_tick_callback tick_callback,
                        on_bar_callback bar_callback,
                        on_l2transaction_callback l2transaction_callback,
                        on_l2order_callback l2order_callback,
                        on_l2order_queue_callback l2order_queue_callback,
                        on_order_status_callback order_status_callback,
                        on_execution_report_callback execution_report_callback,
                        on_algo_order_status_callback algo_order_status_callback,
                        on_cash_callback cash_callback,
                        on_position_callback position_callback,
                        on_parameter_callback parameter_callback,
                        on_schedule_callback schedule_callback,
                        on_backtest_finished_callback backtest_finished_callback,
                        on_account_status_callback account_status_callback,
                        on_error_callback error_callback,
                        on_stop_callback stop_callback,
                        on_market_data_connected_callback market_data_connected_callback,
                        on_trade_data_connected_callback trade_data_connected_callback,
                        on_market_data_disconnected_callback market_data_disconnected_callback,
                        on_trade_data_disconnected_callback trade_data_disconnected_callback) {

    s->register_callbacks(
            init_callback,
            tick_callback,
            bar_callback,
            l2transaction_callback,
            l2order_callback,
            l2order_queue_callback,
            order_status_callback,
            execution_report_callback,
            algo_order_status_callback,
            cash_callback,
            position_callback,
            parameter_callback,
            schedule_callback,
            backtest_finished_callback,
            account_status_callback,
            error_callback,
            stop_callback,
            market_data_connected_callback,
            trade_data_connected_callback,
            market_data_disconnected_callback,
            trade_data_disconnected_callback);
}


int emc_run() {
    return s->run();
}

void emc_stop() {
    s->on_stop();
}

void emc_set_strategy_id(const char *strategy_id) {
    s->set_strategy_id(strategy_id);
}

void emc_set_token(const char *token) {
    s->set_token(token);
}

void emc_set_mode(int mode) {
    s->set_mode(mode);
}

int emc_schedule(const char *data_rule, const char *time_rule) {
    return s->schedule(data_rule, time_rule);
}

double emc_now() {
    return s->now();
}

int emc_set_backtest_config(const char *start_time, const char *end_time, double initial_cash, double transaction_ratio,
                            double commission_ratio, double slippage_ratio, int adjust, int check_cache) {
    return s->set_backtest_config(start_time, end_time, initial_cash, transaction_ratio, commission_ratio,
                                  slippage_ratio,
                                  adjust, check_cache);
}

int emc_subscribe(const char *symbols, const char *frequency, bool unsubscribe_previous) {
    return s->subscribe(symbols, frequency, unsubscribe_previous);
}

int emc_unsubscribe(const char *symbols, const char *frequency) {
    return s->unsubscribe(symbols, frequency);
}

char *emc_get_accounts(int *length, int *error_code) {
    auto array = s->get_accounts();
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *error_code = array->status();
        *length = 0;
        return nullptr;
    }
    *error_code = array->status();

    auto json = json::array();
    for (int i = 0; i < array->count(); i++) {
        auto a = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("account_id", gbk_to_utf8(a.account_id));
        j.add_member("account_id", gbk_to_utf8(a.account_name));
        j.add_member("title", gbk_to_utf8(a.title));
        j.add_member("intro", gbk_to_utf8(a.intro));
        j.add_member("comment", gbk_to_utf8(a.comment));
    }
    *length = array->count();
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());

    array->release();
    return buf;
}

char *emc_get_account_status(const char *account, int *length, int *error_code) {
    AccountStatus acc;
    auto result = s->get_account_status(account, acc);
    if (result != 0) {
        *error_code = result;
        return nullptr;
    }
    *error_code = result;

    auto json = json::object();
    json.add_member("account_id", gbk_to_utf8(acc.account_id));
    json.add_member("account_name", gbk_to_utf8(acc.account_name));
    json.add_member("state", acc.state);
    json.add_member("error_code", acc.error_code);
    json.add_member("error_msg", gbk_to_utf8(acc.error_msg));
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_get_all_account_status(int *length, int *error_code) {
    auto array = s->get_all_account_status();
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    *error_code = array->status();

    auto json = json::array();
    for (int i = 0; i < array->count(); i++) {
        auto acc = array->at(i);
        auto j = json.push_object();
        j.add_member("account_id", gbk_to_utf8(acc.account_id));
        j.add_member("account_name", gbk_to_utf8(acc.account_name));
        j.add_member("state", acc.state);
        j.add_member("error_code", acc.error_code);
        j.add_member("error_msg", gbk_to_utf8(acc.error_msg));
    }
    *length = array->count();
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}


char *emc_order_volume(const char *symbol, int volume, int side, int order_type, int position_effect, double price,
                       const char *account, int *length, int *error_code) {
    auto order = s->order_volume(symbol, volume, side, order_type, position_effect, price, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_order_value(const char *symbol, double value, int side, int order_type, int position_effect, double price,
                      const char *account, int *length, int *error_code) {
    auto order = s->order_value(symbol, value, side, order_type, position_effect, price, account);
    *error_code = 0;

    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *
emc_order_percent(const char *symbol, double percent, int side, int order_type, int position_effect, double price,
                  const char *account, int *length, int *error_code) {
    auto order = s->order_percent(symbol, percent, side, order_type, position_effect, price, account);
    *error_code = 0;

    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_order_target_volume(const char *symbol, int volume, int position_side, int order_type, double price,
                              const char *account, int *length, int *error_code) {
    auto order = s->order_target_volume(symbol, volume, position_side, order_type, price, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_order_target_value(const char *symbol, double value, int position_side, int order_type, double price,
                             const char *account, int *length, int *error_code) {
    auto order = s->order_target_value(symbol, value, position_side, order_type, price, account);
    *error_code = 0;

    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_order_target_percent(const char *symbol, double percent, int position_side, int order_type, double price,
                               const char *account, int *length, int *error_code) {
    auto order = s->order_target_percent(symbol, percent, position_side, order_type, price, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_order_close_all(int *length, int *error_code) {
    auto array = s->order_close_all();
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    auto status = array->status();
    if (status != 0) {
        *error_code = status;
        *length = 0;
        return nullptr;
    }
    *error_code = status;
    auto json = json::array();
    for (int i = 0; i < array->count(); i++) {
        auto order = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("strategy_id", gbk_to_utf8(order.strategy_id));
        j.add_member("account_id", gbk_to_utf8(order.account_id));
        j.add_member("account_name", gbk_to_utf8(order.account_name));
        j.add_member("cl_ord_id", gbk_to_utf8(order.cl_ord_id));
        j.add_member("order_id", gbk_to_utf8(order.order_id));
        j.add_member("ex_ord_id", gbk_to_utf8(order.ex_ord_id));
        j.add_member("algo_order_id", order.algo_order_id);
        j.add_member("order_business", order.order_business);

        j.add_member("symbol", gbk_to_utf8(order.symbol));
        j.add_member("side", order.side);
        j.add_member("position_effect", order.position_effect);
        j.add_member("position_side", order.position_side);

        j.add_member("order_type", order.order_type);
        j.add_member("order_duration", order.order_duration);
        j.add_member("order_qualifier", order.order_qualifier);
        j.add_member("order_src", order.order_src);
        j.add_member("position_src", order.position_src);

        j.add_member("status", order.status);
        j.add_member("ord_rej_reason", order.ord_rej_reason);
        j.add_member("ord_rej_reason_detail", gbk_to_utf8(order.ord_rej_reason_detail));

        j.add_member("price", order.price);
        j.add_member("stop_price", order.stop_price);

        j.add_member("order_style", order.order_style);
        j.add_member("volume", order.volume);
        j.add_member("value", order.value);
        j.add_member("percent", order.percent);
        j.add_member("target_volume", order.target_volume);
        j.add_member("target_value", order.target_value);
        j.add_member("target_percent", order.target_percent);

        j.add_member("filled_volume", order.filled_volume);
        j.add_member("filled_vwap", order.filled_vwap);
        j.add_member("filled_amount", order.filled_amount);
        j.add_member("filled_commission", order.filled_commission);

        j.add_member("created_at", format_time(order.created_at));
        j.add_member("updated_at", format_time(order.updated_at));
    }
    *length = array->count();
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

int emc_order_cancel(const char *cl_ord_id, const char *account) {
    return s->order_cancel(cl_ord_id, account);
}

int emc_order_cancel_all() {
    return s->order_cancel_all();
}

char *emc_place_order(const char *symbol, int volume, int side, int order_type, int position_effect, double price,
                      int order_duration, int order_qualifier, double stop_price, int order_business,
                      const char *account, int *length, int *error_code) {
    auto order = s->place_order(symbol, volume, side, order_type, position_effect, price, order_duration,
                                order_qualifier, stop_price, order_business, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_order_after_hour(const char *symbol, int volume, int side, double price, const char *account, int *length,
                           int *error_code) {
    auto order = s->order_after_hour(symbol, volume, side, price, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_get_orders(const char *account, int *length, int *error_code) {
    auto array = s->get_orders(account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    for (int i = 0; i < array->count(); i++) {
        auto order = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("strategy_id", gbk_to_utf8(order.strategy_id));
        j.add_member("account_id", gbk_to_utf8(order.account_id));
        j.add_member("account_name", gbk_to_utf8(order.account_name));
        j.add_member("cl_ord_id", gbk_to_utf8(order.cl_ord_id));
        j.add_member("order_id", gbk_to_utf8(order.order_id));
        j.add_member("ex_ord_id", gbk_to_utf8(order.ex_ord_id));
        j.add_member("algo_order_id", order.algo_order_id);
        j.add_member("order_business", order.order_business);

        j.add_member("symbol", gbk_to_utf8(order.symbol));
        j.add_member("side", order.side);
        j.add_member("position_effect", order.position_effect);
        j.add_member("position_side", order.position_side);

        j.add_member("order_type", order.order_type);
        j.add_member("order_duration", order.order_duration);
        j.add_member("order_qualifier", order.order_qualifier);
        j.add_member("order_src", order.order_src);
        j.add_member("position_src", order.position_src);

        j.add_member("status", order.status);
        j.add_member("ord_rej_reason", order.ord_rej_reason);
        j.add_member("ord_rej_reason_detail", gbk_to_utf8(order.ord_rej_reason_detail));

        j.add_member("price", order.price);
        j.add_member("stop_price", order.stop_price);

        j.add_member("order_style", order.order_style);
        j.add_member("volume", order.volume);
        j.add_member("value", order.value);
        j.add_member("percent", order.percent);
        j.add_member("target_volume", order.target_volume);
        j.add_member("target_value", order.target_value);
        j.add_member("target_percent", order.target_percent);

        j.add_member("filled_volume", order.filled_volume);
        j.add_member("filled_vwap", order.filled_vwap);
        j.add_member("filled_amount", order.filled_amount);
        j.add_member("filled_commission", order.filled_commission);

        j.add_member("created_at", format_time(order.created_at));
        j.add_member("updated_at", format_time(order.updated_at));
    }
    *error_code = array->status();

    *length = array->count();
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

char *emc_get_unfinished_orders(const char *account, int *length, int *error_code) {
    auto array = s->get_unfinished_orders(account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    auto count = array->count();
    for (int i = 0; i < count; i++) {
        auto order = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("strategy_id", gbk_to_utf8(order.strategy_id));
        j.add_member("account_id", gbk_to_utf8(order.account_id));
        j.add_member("account_name", gbk_to_utf8(order.account_name));
        j.add_member("cl_ord_id", gbk_to_utf8(order.cl_ord_id));
        j.add_member("order_id", gbk_to_utf8(order.order_id));
        j.add_member("ex_ord_id", gbk_to_utf8(order.ex_ord_id));
        j.add_member("algo_order_id", order.algo_order_id);
        j.add_member("order_business", order.order_business);

        j.add_member("symbol", gbk_to_utf8(order.symbol));
        j.add_member("side", order.side);
        j.add_member("position_effect", order.position_effect);
        j.add_member("position_side", order.position_side);

        j.add_member("order_type", order.order_type);
        j.add_member("order_duration", order.order_duration);
        j.add_member("order_qualifier", order.order_qualifier);
        j.add_member("order_src", order.order_src);
        j.add_member("position_src", order.position_src);

        j.add_member("status", order.status);
        j.add_member("ord_rej_reason", order.ord_rej_reason);
        j.add_member("ord_rej_reason_detail", gbk_to_utf8(order.ord_rej_reason_detail));

        j.add_member("price", order.price);
        j.add_member("stop_price", order.stop_price);

        j.add_member("order_style", order.order_style);
        j.add_member("volume", order.volume);
        j.add_member("value", order.value);
        j.add_member("percent", order.percent);
        j.add_member("target_volume", order.target_volume);
        j.add_member("target_value", order.target_value);
        j.add_member("target_percent", order.target_percent);

        j.add_member("filled_volume", order.filled_volume);
        j.add_member("filled_vwap", order.filled_vwap);
        j.add_member("filled_amount", order.filled_amount);
        j.add_member("filled_commission", order.filled_commission);

        j.add_member("created_at", format_time(order.created_at));
        j.add_member("updated_at", format_time(order.updated_at));
    }
    *error_code = array->status();
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

char *emc_get_execution_reports(const char *account, int *length, int *error_code) {
    auto array = s->get_execution_reports(account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    auto count = array->count();
    for (int i = 0; i < count; i++) {
        auto rpt = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("strategy_id", gbk_to_utf8(rpt.strategy_id));
        j.add_member("account_id", gbk_to_utf8(rpt.account_id));
        j.add_member("account_name", gbk_to_utf8(rpt.account_name));

        j.add_member("cl_ord_id", gbk_to_utf8(rpt.cl_ord_id));
        j.add_member("order_id", gbk_to_utf8(rpt.order_id));
        j.add_member("exec_id", gbk_to_utf8(rpt.exec_id));

        j.add_member("symbol", gbk_to_utf8(rpt.symbol));

        j.add_member("position_effect", rpt.position_effect);
        j.add_member("side", rpt.side);
        j.add_member("ord_rej_reason", rpt.ord_rej_reason);
        j.add_member("ord_rej_reason_detail", gbk_to_utf8(rpt.ord_rej_reason_detail));
        j.add_member("exec_type", rpt.exec_type);

        j.add_member("price", rpt.price);
        j.add_member("volume", rpt.volume);
        j.add_member("amount", rpt.amount);
        j.add_member("commission", rpt.commission);
        j.add_member("cost", rpt.cost);
        j.add_member("created_at", format_time(rpt.created_at));

    }
    *error_code = array->status();
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

char *emc_get_cash(const char *accounts, int *length, int *error_code) {
    auto array = s->get_cash(accounts);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    auto count = array->count();
    for (int i = 0; i < count; i++) {
        auto cash = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("account_id", gbk_to_utf8(cash.account_id));
        j.add_member("account_name", gbk_to_utf8(cash.account_name));
        j.add_member("currency", cash.currency);

        j.add_member("nav", cash.nav);
        j.add_member("pnl", cash.pnl);
        j.add_member("fpnl", cash.fpnl);
        j.add_member("frozen", cash.frozen);
        j.add_member("order_frozen", cash.order_frozen);
        j.add_member("available", cash.available);
        j.add_member("balance", cash.balance);
        j.add_member("market_value", cash.market_value);
        j.add_member("cum_inout", cash.cum_inout);
        j.add_member("cum_trade", cash.cum_trade);
        j.add_member("cum_pnl", cash.cum_pnl);
        j.add_member("cum_commission", cash.cum_commission);

        j.add_member("last_trade", cash.last_trade);
        j.add_member("last_pnl", cash.last_pnl);
        j.add_member("last_commission", cash.last_commission);
        j.add_member("last_inout", cash.last_inout);
        j.add_member("change_reason", cash.change_reason);
        j.add_member("change_event_id", gbk_to_utf8(cash.change_event_id));

        j.add_member("created_at", format_time(cash.created_at));
        j.add_member("updated_at", format_time(cash.updated_at));
    }
    *error_code = array->status();
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

char *emc_get_position(const char *account, int *length, int *error_code) {
    auto array = s->get_position(strlen(account) == 0 ? nullptr : account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    auto count = array->count();
    for (int i = 0; i < count; i++) {
        auto pos = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("account_id", gbk_to_utf8(pos.account_id));
        j.add_member("account_name", gbk_to_utf8(pos.account_name));

        j.add_member("symbol", gbk_to_utf8(pos.symbol));
        j.add_member("side", pos.side);
        j.add_member("volume", pos.volume);
        j.add_member("volume_today", pos.volume_today);
        j.add_member("vwap", pos.vwap);
        j.add_member("vwap_diluted", pos.vwap_diluted);
        j.add_member("vwap_open", pos.vwap_open);
        j.add_member("amount", pos.amount);

        j.add_member("price", pos.price);
        j.add_member("fpnl", pos.fpnl);
        j.add_member("fpnl_open", pos.fpnl_open);
        j.add_member("cost", pos.cost);
        j.add_member("order_frozen", pos.order_frozen);
        j.add_member("order_frozen_today", pos.order_frozen_today);
        j.add_member("available", pos.available);
        j.add_member("available_today", pos.available_today);
        j.add_member("available_now", pos.available_now);
        j.add_member("market_value", pos.market_value);

        j.add_member("last_price", pos.last_price);
        j.add_member("last_volume", pos.last_volume);
        j.add_member("last_inout", pos.last_inout);
        j.add_member("change_reason", pos.change_reason);
        j.add_member("change_event_id", gbk_to_utf8(pos.change_event_id));

        j.add_member("has_dividend", pos.has_dividend);
        j.add_member("created_at", format_time(pos.created_at));
        j.add_member("updated_at", format_time(pos.updated_at));
    }
    *error_code = array->status();
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}


char *emc_order_algo(const char *symbol, int volume, int position_effect, int side, int order_type, double price,
                     const char *algo_name, const char *algo_param, const char *account, int *length, int *error_code) {
    auto order = s->order_algo(symbol, volume, position_effect, side, order_type, price, algo_name, algo_param,
                               account);
    *error_code = 0;
    auto json = json::object();
    json.add_member("strategy_id", gbk_to_utf8(order.strategy_id));
    json.add_member("account_id", gbk_to_utf8(order.account_id));
    json.add_member("account_name", gbk_to_utf8(order.account_name));

    json.add_member("cl_ord_id", gbk_to_utf8(order.cl_ord_id));
    json.add_member("order_id", gbk_to_utf8(order.order_id));
    json.add_member("ex_ord_id", gbk_to_utf8(order.ex_ord_id));
    json.add_member("order_business", order.order_business);

    json.add_member("symbol", gbk_to_utf8(order.symbol));
    json.add_member("side", order.side);
    json.add_member("position_effect", order.position_effect);
    json.add_member("position_side", order.position_side);

    json.add_member("order_type", order.order_type);
    json.add_member("order_duration", order.order_duration);
    json.add_member("order_qualifier", order.order_qualifier);
    json.add_member("order_src", order.order_src);
    json.add_member("position_src", order.position_src);

    json.add_member("status", order.status);
    json.add_member("ord_rej_reason", order.ord_rej_reason);
    json.add_member("ord_rej_reason_detail", gbk_to_utf8(order.ord_rej_reason_detail));

    json.add_member("price", order.price);
    json.add_member("stop_price", order.stop_price);

    json.add_member("order_style", order.order_style);
    json.add_member("volume", order.volume);
    json.add_member("value", order.value);
    json.add_member("percent", order.percent);
    json.add_member("target_volume", order.target_volume);
    json.add_member("target_value", order.target_value);
    json.add_member("target_percent", order.target_percent);

    json.add_member("filled_volume", order.filled_volume);
    json.add_member("filled_vwap", order.filled_vwap);
    json.add_member("filled_amount", order.filled_amount);
    json.add_member("filled_commission", order.filled_commission);

    json.add_member("algo_name", gbk_to_utf8(order.algo_name));
    json.add_member("algo_param", gbk_to_utf8(order.algo_param));
    json.add_member("algo_status", order.algo_status);
    json.add_member("algo_comment", gbk_to_utf8(order.algo_comment));

    json.add_member("created_at", format_time(order.created_at));
    json.add_member("updated_at", format_time(order.updated_at));

    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

int emc_algo_order_cancel(const char *cl_ord_id, const char *account) {
    return s->algo_order_cancel(cl_ord_id, account);
}

int emc_algo_order_pause(const char *cl_ord_id, int status, const char *account) {
    return s->algo_order_pause(cl_ord_id, status, account);
}

char *emc_get_algo_orders(const char *account, int *length, int *error_code) {
    auto array = s->get_algo_orders(account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    auto count = array->count();
    for (int i = 0; i < count; i++) {
        auto order = array->at(i);
        auto j = json::object();
        j.add_member("index", i);
        j.add_member("strategy_id", gbk_to_utf8(order.strategy_id));
        j.add_member("account_id", gbk_to_utf8(order.account_id));
        j.add_member("account_name", gbk_to_utf8(order.account_name));

        j.add_member("cl_ord_id", gbk_to_utf8(order.cl_ord_id));
        j.add_member("order_id", gbk_to_utf8(order.order_id));
        j.add_member("ex_ord_id", gbk_to_utf8(order.ex_ord_id));
        j.add_member("order_business", order.order_business);

        j.add_member("symbol", gbk_to_utf8(order.symbol));
        j.add_member("side", order.side);
        j.add_member("position_effect", order.position_effect);
        j.add_member("position_side", order.position_side);

        j.add_member("order_type", order.order_type);
        j.add_member("order_duration", order.order_duration);
        j.add_member("order_qualifier", order.order_qualifier);
        j.add_member("order_src", order.order_src);
        j.add_member("position_src", order.position_src);

        j.add_member("status", order.status);
        j.add_member("ord_rej_reason", order.ord_rej_reason);
        j.add_member("ord_rej_reason_detail", gbk_to_utf8(order.ord_rej_reason_detail));

        j.add_member("price", order.price);
        j.add_member("stop_price", order.stop_price);

        j.add_member("order_style", order.order_style);
        j.add_member("volume", order.volume);
        j.add_member("value", order.value);
        j.add_member("percent", order.percent);
        j.add_member("target_volume", order.target_volume);
        j.add_member("target_value", order.target_value);
        j.add_member("target_percent", order.target_percent);

        j.add_member("filled_volume", order.filled_volume);
        j.add_member("filled_vwap", order.filled_vwap);
        j.add_member("filled_amount", order.filled_amount);
        j.add_member("filled_commission", order.filled_commission);

        j.add_member("algo_name", gbk_to_utf8(order.algo_name));
        j.add_member("algo_param", gbk_to_utf8(order.algo_param));
        j.add_member("algo_status", order.algo_status);
        j.add_member("algo_comment", gbk_to_utf8(order.algo_comment));

        j.add_member("created_at", format_time(order.created_at));
        j.add_member("updated_at", format_time(order.updated_at));
    }
    *error_code = array->status();
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

char *emc_get_algo_child_orders(const char *cl_ord_id, const char *account, int *length, int *error_code) {
    auto array = s->get_algo_child_orders(cl_ord_id, account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    auto count = array->count();
    for (int i = 0; i < count; i++) {
        auto order = array->at(i);
        auto j = json::object();
        j.add_member("index", i);
        j.add_member("strategy_id", gbk_to_utf8(order.strategy_id));
        j.add_member("account_id", gbk_to_utf8(order.account_id));
        j.add_member("account_name", gbk_to_utf8(order.account_name));
        j.add_member("cl_ord_id", gbk_to_utf8(order.cl_ord_id));
        j.add_member("order_id", gbk_to_utf8(order.order_id));
        j.add_member("ex_ord_id", gbk_to_utf8(order.ex_ord_id));
        j.add_member("algo_order_id", order.algo_order_id);
        j.add_member("order_business", order.order_business);

        j.add_member("symbol", gbk_to_utf8(order.symbol));
        j.add_member("side", order.side);
        j.add_member("position_effect", order.position_effect);
        j.add_member("position_side", order.position_side);

        j.add_member("order_type", order.order_type);
        j.add_member("order_duration", order.order_duration);
        j.add_member("order_qualifier", order.order_qualifier);
        j.add_member("order_src", order.order_src);
        j.add_member("position_src", order.position_src);

        j.add_member("status", order.status);
        j.add_member("ord_rej_reason", order.ord_rej_reason);
        j.add_member("ord_rej_reason_detail", gbk_to_utf8(order.ord_rej_reason_detail));

        j.add_member("price", order.price);
        j.add_member("stop_price", order.stop_price);

        j.add_member("order_style", order.order_style);
        j.add_member("volume", order.volume);
        j.add_member("value", order.value);
        j.add_member("percent", order.percent);
        j.add_member("target_volume", order.target_volume);
        j.add_member("target_value", order.target_value);
        j.add_member("target_percent", order.target_percent);

        j.add_member("filled_volume", order.filled_volume);
        j.add_member("filled_vwap", order.filled_vwap);
        j.add_member("filled_amount", order.filled_amount);
        j.add_member("filled_commission", order.filled_commission);

        j.add_member("created_at", format_time(order.created_at));
        j.add_member("updated_at", format_time(order.updated_at));
    }
    *error_code = array->status();
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

int emc_raw_func(const char *account, const char *func_id, const char *func_args, char **rsp) {
    auto r = s->raw_func(account, func_id, func_args, *rsp);
    if (r != 0) {
        return r;
    }
    return 0;
}

char *emc_credit_buying_on_margin(int position_src, const char *symbol, int volume, double price, int order_type,
                                  int order_duration, int order_qualifier, const char *account, int *length,
                                  int *error_code) {
    auto order = s->credit_buying_on_margin(position_src, symbol, volume, price, order_type, order_duration,
                                            order_qualifier, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_credit_short_selling(int position_src, const char *symbol, int volume, double price, int order_type,
                               int order_duration, int order_qualifier, const char *account, int *length,
                               int *error_code) {
    auto order = s->credit_short_selling(position_src, symbol, volume, price, order_type, order_duration,
                                         order_qualifier, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *
emc_credit_repay_share_by_buying_share(const char *symbol, int volume, double price, int order_type, int order_duration,
                                       int order_qualifier, const char *account, int *length, int *error_code) {
    auto order = s->credit_repay_share_by_buying_share(symbol, volume, price, order_type, order_duration,
                                                       order_qualifier, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *
emc_credit_repay_cash_by_selling_share(const char *symbol, int volume, double price, int order_type, int order_duration,
                                       int order_qualifier, const char *account, int *length, int *error_code) {
    auto order = s->credit_repay_cash_by_selling_share(symbol, volume, price, order_type, order_duration,
                                                       order_qualifier, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_credit_buying_on_collateral(const char *symbol, int volume, double price, int order_type, int order_duration,
                                      int order_qualifier, const char *account, int *length, int *error_code) {
    auto order = s->credit_buying_on_collateral(symbol, volume, price, order_type, order_duration, order_qualifier,
                                                account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_credit_selling_on_collateral(const char *symbol, int volume, double price, int order_type, int order_duration,
                                       int order_qualifier, const char *account, int *length, int *error_code) {
    auto order = s->credit_selling_on_collateral(symbol, volume, price, order_type, order_duration, order_qualifier,
                                                 account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *
emc_credit_repay_share_directly(const char *symbol, int volume, const char *account, int *length, int *error_code) {
    auto order = s->credit_repay_share_directly(symbol, volume, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

int emc_credit_repay_cash_directly(double amount, const char *account, double *actual_repay_amount, char *error_msg_buf,
                                   int buf_len) {
    return s->credit_repay_cash_directly(amount, account, actual_repay_amount, error_msg_buf, buf_len);
}

char *emc_credit_collateral_in(const char *symbol, int volume, const char *account, int *length, int *error_code) {
    auto order = s->credit_collateral_in(symbol, volume, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_credit_collateral_out(const char *symbol, int volume, const char *account, int *length, int *error_code) {
    auto order = s->credit_collateral_out(symbol, volume, account);
    *error_code = 0;
    auto json = order_to_json(order);
    *length = 1;
    auto count = json.str().size() + 1;
    auto buf = (char *) malloc(count);
    memset(buf, 0, count);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_credit_get_collateral_instruments(const char *account, int *length, int *error_code) {
    auto array = s->credit_get_collateral_instruments(account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto in = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("symbol", gbk_to_utf8(in.symbol));
        j.add_member("name", gbk_to_utf8(in.name));
        j.add_member("pledge_rate", in.pledge_rate);
    }
    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

char *emc_credit_get_borrowable_instruments(int position_src, const char *account, int *length, int *error_code) {
    auto array = s->credit_get_borrowable_instruments(position_src, account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto in = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("symbol", gbk_to_utf8(in.symbol));
        j.add_member("name", gbk_to_utf8(in.name));
        j.add_member("margin_rate_for_cash", in.margin_rate_for_cash);
        j.add_member("margin_rate_for_security", in.margin_rate_for_security);
    }
    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}


char *
emc_credit_get_borrowable_instruments_positions(int position_src, const char *account, int *length, int *error_code) {
    auto array = s->credit_get_borrowable_instruments_positions(position_src, account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto pos = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("symbol", gbk_to_utf8(pos.symbol));
        j.add_member("name", gbk_to_utf8(pos.name));
        j.add_member("balance", pos.balance);
        j.add_member("available", pos.available);
    }

    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}


char *emc_credit_get_contracts(int position_src, const char *account, int *length, int *error_code) {
    auto array = s->credit_get_contracts(position_src, account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto in = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("symbol", gbk_to_utf8(in.symbol));
        j.add_member("name", gbk_to_utf8(in.name));
        j.add_member("orderdate", in.orderdate);
        j.add_member("ordersno", gbk_to_utf8(in.ordersno));
        j.add_member("creditdirect", std::string().append(1, in.creditdirect).c_str());
        j.add_member("orderqty", in.orderqty);
        j.add_member("matchqty", in.matchqty);
        j.add_member("orderamt", in.orderamt);
        j.add_member("orderfrzamt", in.orderfrzamt);
        j.add_member("matchamt", in.matchamt);
        j.add_member("clearamt", in.clearamt);
        j.add_member("lifestatus", std::string().append(1, in.lifestatus).c_str());
        j.add_member("enddate", in.enddate);
        j.add_member("oldenddate", in.oldenddate);
        j.add_member("creditrepay", in.creditrepay);
        j.add_member("creditrepayunfrz", in.creditrepayunfrz);
        j.add_member("fundremain", in.fundremain);
        j.add_member("stkrepay", in.stkrepay);
        j.add_member("stkrepayunfrz", in.stkrepayunfrz);
        j.add_member("stkremain", in.stkremain);
        j.add_member("stkremainvalue", in.stkremainvalue);
        j.add_member("fee", in.fee);
        j.add_member("overduefee", in.overduefee);
        j.add_member("fee_repay", in.fee_repay);
        j.add_member("punifee", in.punifee);
        j.add_member("punifee_repay", in.punifee_repay);
        j.add_member("rights", in.rights);
        j.add_member("overduerights", in.overduerights);
        j.add_member("rights_repay", in.rights_repay);
        j.add_member("lastprice", in.lastprice);
        j.add_member("profitcost", in.profitcost);
        j.add_member("sysdate", in.sysdate);
        j.add_member("sno", gbk_to_utf8(in.sno));
        j.add_member("lastdate", in.lastdate);
        j.add_member("closedate", in.closedate);
        j.add_member("punidebts", in.punidebts);
        j.add_member("punidebts_repay", in.punidebts_repay);
        j.add_member("punidebtsunfrz", in.punidebtsunfrz);
        j.add_member("punifeeunfrz", in.punifeeunfrz);
        j.add_member("punirights", in.punirights);
        j.add_member("punirights_repay", in.punirights_repay);
        j.add_member("punirightsunfrz", in.punirightsunfrz);
        j.add_member("feeunfrz", in.feeunfrz);
        j.add_member("rightsqty", in.rightsqty);
        j.add_member("overduerightsqty", in.overduerightsqty);

    }

    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

char *emc_credit_get_cash(const char *account, int *length, int *error_code) {
    CreditCash c;
    auto r = s->credit_get_cash(c, account);
    if (r != 0) {
        *error_code = r;
        *length = 0;
        return nullptr;
    }
    auto json = json::object();
    json.add_member("fundintrrate", c.fundintrrate);
    json.add_member("stkintrrate", c.stkintrrate);
    json.add_member("punishintrrate", c.punishintrrate);
    json.add_member("creditstatus", std::string().append(1, c.creditstatus).c_str());
    json.add_member("marginrates", c.marginrates);
    json.add_member("realrate", c.realrate);
    json.add_member("asset", c.asset);
    json.add_member("liability", c.liability);
    json.add_member("marginavl", c.marginavl);
    json.add_member("fundbal", c.fundbal);
    json.add_member("fundavl", c.fundavl);
    json.add_member("dsaleamtbal", c.dsaleamtbal);
    json.add_member("guaranteeout", c.guaranteeout);
    json.add_member("gagemktavl", c.gagemktavl);
    json.add_member("fdealavl", c.fdealavl);
    json.add_member("ffee", c.ffee);
    json.add_member("ftotaldebts", c.ftotaldebts);
    json.add_member("dealfmktavl", c.dealfmktavl);
    json.add_member("dfee", c.dfee);
    json.add_member("dtotaldebts", c.dtotaldebts);
    json.add_member("fcreditbal", c.fcreditbal);
    json.add_member("fcreditavl", c.fcreditavl);
    json.add_member("fcreditfrz", c.fcreditfrz);
    json.add_member("dcreditbal", c.dcreditbal);
    json.add_member("dcreditavl", c.dcreditavl);
    json.add_member("dcreditfrz", c.dcreditfrz);
    json.add_member("rights", c.rights);
    json.add_member("serviceuncomerqrights", c.serviceuncomerqrights);
    json.add_member("rightsqty", c.rightsqty);
    json.add_member("serviceuncomerqrightsqty", c.serviceuncomerqrightsqty);
    json.add_member("acreditbal", c.acreditbal);
    json.add_member("acreditavl", c.acreditavl);
    json.add_member("acashcapital", c.acashcapital);
    json.add_member("astkmktvalue", c.astkmktvalue);
    json.add_member("withdrawable", c.withdrawable);
    json.add_member("netcapital", c.netcapital);
    json.add_member("fcreditpnl", c.fcreditpnl);
    json.add_member("dcreditpnl", c.dcreditpnl);
    json.add_member("fcreditmarginoccupied", c.fcreditmarginoccupied);
    json.add_member("dcreditmarginoccupied", c.dcreditmarginoccupied);
    json.add_member("collateralbuyableamt", c.collateralbuyableamt);
    json.add_member("repayableamt", c.repayableamt);
    json.add_member("dcreditcashavl", c.dcreditcashavl);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}


char *emc_ipo_buy(const char *symbol, int volume, double price, const char *account, int *length, int *error_code) {
    auto o = s->ipo_buy(symbol, volume, price, account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_ipo_get_quota(const char *account, int *length, int *error_code) {
    auto array = s->ipo_get_quota(account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto ipo = array->at(i);
        auto j = json::object();
        j.add_member("index", i);
        j.add_member("exchange", gbk_to_utf8(ipo.exchange));
        j.add_member("quota", ipo.quota);
        j.add_member("sse_star_quota", ipo.sse_star_quota);
    }
    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}


char *emc_ipo_get_instruments(int security_type, const char *account, int *length, int *error_code) {
    auto array = s->ipo_get_instruments(security_type, account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto ipo = array->at(i);
        auto j = json::object();
        j.add_member("index", i);
        j.add_member("symbol", gbk_to_utf8(ipo.symbol));
        j.add_member("price", ipo.price);
        j.add_member("min_vol", ipo.min_vol);
        j.add_member("max_vol", ipo.max_vol);
    }

    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}


char *emc_ipo_get_match_number(const char *account, int *length, int *error_code) {
    auto array = s->ipo_get_match_number(account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto ipo = array->at(i);
        auto j = json::object();
        j.add_member("index", i);
        j.add_member("order_id", gbk_to_utf8(ipo.order_id));
        j.add_member("symbol", gbk_to_utf8(ipo.symbol));
        j.add_member("volume", ipo.volume);
        j.add_member("match_number", gbk_to_utf8(ipo.match_number));
        j.add_member("order_at", ipo.order_at);
        j.add_member("match_at", ipo.match_at);
    }
    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}


char *emc_ipo_get_lot_info(const char *account, int *length, int *error_code) {
    auto array = s->ipo_get_lot_info(account);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto ipo = array->at(i);
        auto j = json::object();
        j.add_member("index", i);
        j.add_member("symbol", gbk_to_utf8(ipo.symbol));
        j.add_member("order_at", ipo.order_at);
        j.add_member("lot_at", ipo.lot_at);
        j.add_member("lot_volume", ipo.lot_volume);
        j.add_member("give_up_volume", ipo.give_up_volume);
        j.add_member("price", ipo.price);
        j.add_member("amount", ipo.amount);
        j.add_member("pay_volume", ipo.pay_volume);
        j.add_member("pay_amount", ipo.pay_amount);
    }

    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}


char *
emc_fund_etf_buy(const char *symbol, int volume, double price, const char *account, int *length, int *error_code) {
    auto o = s->fund_etf_buy(symbol, volume, price, account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_fund_etf_redemption(const char *symbol, int volume, double price, const char *account, int *length,
                              int *error_code) {
    auto o = s->fund_etf_redemption(symbol, volume, price, account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_fund_subscribing(const char *symbol, int volume, const char *account, int *length, int *error_code) {
    auto o = s->fund_subscribing(symbol, volume, account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_fund_buy(const char *symbol, int volume, const char *account, int *length, int *error_code) {
    auto o = s->fund_buy(symbol, volume, account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_fund_redemption(const char *symbol, int volume, const char *account, int *length, int *error_code) {
    auto o = s->fund_redemption(symbol, volume, account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *
emc_bond_reverse_repurchase_agreement(const char *symbol, int volume, double price, int order_type, int order_duration,
                                      int order_qualifier, const char *account, int *length, int *error_code) {
    auto o = s->bond_reverse_repurchase_agreement(symbol, volume, price, order_type, order_duration, order_qualifier,
                                                  account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_bond_convertible_call(const char *symbol, int volume, double price, const char *account, int *length,
                                int *error_code) {
    auto o = s->bond_convertible_call(symbol, volume, price, account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *emc_bond_convertible_put(const char *symbol, int volume, double price, const char *account, int *length,
                               int *error_code) {
    auto o = s->bond_convertible_put(symbol, volume, price, account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

char *
emc_bond_convertible_put_cancel(const char *symbol, int volume, const char *account, int *length, int *error_code) {
    auto o = s->bond_convertible_put_cancel(symbol, volume, account);
    auto json = order_to_json(o);
    *error_code = 0;
    *length = 1;
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    return buf;
}

int emc_add_parameters(Parameter *params, int count) {
    return s->add_parameters(params, count);
}

int emc_del_parameters(const char *keys) {
    return s->del_parameters(keys);
}

int emc_set_parameters(Parameter *params, int count) {
    return s->set_parameters(params, count);
}

char *emc_get_parameters(int *length, int *error_code) {
    auto array = s->get_parameters();
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto p = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("key", gbk_to_utf8(p.key));
        j.add_member("value", p.value);
        j.add_member("min", p.min);
        j.add_member("max", p.max);
        j.add_member("name", gbk_to_utf8(p.name));
        j.add_member("intro", gbk_to_utf8(p.intro));
        j.add_member("group", gbk_to_utf8(p.group));
        j.add_member("readonly", p.readonly);
    }
    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}


int emc_set_symbols(const char *symbols) {
    return s->set_symbols(symbols);
}


char *emc_get_symbols(int *length, int *error_code) {
    auto array = s->get_symbols();
    if (array == nullptr) {
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        return nullptr;
    }
    auto json = json::array();
    int count = array->count();
    for (int i = 0; i < count; i++) {
        auto symbol = array->at(i);
        auto j = json.push_object();
        j.add_member("index", i);
        j.add_member("symbol", gbk_to_utf8(symbol));
    }
    *error_code = 0;
    *length = array->count();
    auto len = json.str().size() + 1;
    auto buf = (char *) malloc(len);
    memset(buf, 0, len);
    strcpy(buf, json.str().c_str());
    array->release();
    return buf;
}

void create_instance() {
    s = std::make_shared<EmtStrategy>();
}

const char *emc_get_version() {
    return get_version();
}

void gm_set_token(const char *token) {
    set_token(token);
}

void emc_set_serv_addr(const char *addr) {
    set_serv_addr(addr);
}

void emc_set_mfp(const char *mfp) {
    set_mfp(mfp);
}


Tick **emc_current(const char *symbols, int *length, int *error_code) {
    auto array = current(symbols);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    Tick **ticks = (Tick **) malloc(sizeof(Tick **));
    for (int i = 0; i < array->count(); i++) {
        ticks[i] = (Tick *) malloc(sizeof(Tick));
        memcpy(ticks[i], &array->at(i), sizeof(Tick));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return ticks;
}

void emc_ticks_free(struct Tick **ticks, int length) {
    if (ticks == nullptr) {
        return;
    }

    if (length == 1) {
        free(ticks[0]);
        free(ticks);
    } else {
        for (int i = 0; i < length; i++) {
            free(ticks[i]);
        }
        free(ticks);
    }
}

struct Tick **
emc_history_ticks(const char *symbols, const char *start_time, const char *end_time, int adjust,
                  const char *adjust_end_time, bool skip_suspended, const char *fill_missing, int *length,
                  int *error_code) {
    auto array = history_ticks(symbols, start_time, end_time, adjust, (strlen(adjust_end_time) == 0 ? nullptr : adjust_end_time), skip_suspended, fill_missing);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    Tick **ticks = (Tick **) malloc(sizeof(Tick **));
    for (int i = 0; i < array->count(); i++) {
        ticks[i] = (Tick *) malloc(sizeof(Tick));
        memcpy(ticks[i], &array->at(i), sizeof(Tick));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return ticks;
}

Bar **
emc_history_bars(const char *symbols, const char *frequency, const char *start_time, const char *end_time, int adjust,
                 const char *adjust_end_time, bool skip_suspended, const char *fill_missing, int *length,
                 int *error_code) {
    auto array = history_bars(symbols, frequency, start_time, end_time, adjust, (strlen(adjust_end_time) == 0 ? nullptr : adjust_end_time), skip_suspended, fill_missing);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    Bar **bars = (Bar **) malloc(sizeof(Bar **));
    for (int i = 0; i < array->count(); i++) {
        bars[i] = (Bar *) malloc(sizeof(Bar));
        memcpy(bars[i], &array->at(i), sizeof(Bar));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return bars;
}

void emc_bars_free(Bar **bars, int length) {
    if (bars == nullptr) {
        return;
    }

    if (length == 1) {
        free(bars[0]);
        free(bars);
    } else {
        for (int i = 0; i < length; i++) {
            free(bars[i]);
        }
        free(bars);
    }
}

struct Tick **
emc_history_ticks_n(const char *symbols, int n, const char *end_time, int adjust, const char *adjust_end_time,
                    bool skip_suspended, const char *fill_missing, int *length, int *error_code) {
    auto array = history_ticks_n(symbols, n, end_time, adjust, (strlen(adjust_end_time) == 0 ? nullptr : adjust_end_time), skip_suspended, fill_missing);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    Tick **ticks = (Tick **) malloc(sizeof(Tick **));
    for (int i = 0; i < array->count(); i++) {
        ticks[i] = (Tick *) malloc(sizeof(Tick));
        memcpy(ticks[i], &array->at(i), sizeof(Tick));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return ticks;
}

Bar **emc_history_bars_n(const char *symbols, const char *frequency, int n, const char *end_time, int adjust,
                         const char *adjust_end_time, bool skip_suspended, const char *fill_missing, int *length,
                         int *error_code) {
    auto array = history_bars_n(symbols, frequency, n, end_time, adjust, (strlen(adjust_end_time) == 0 ? nullptr : adjust_end_time), skip_suspended, fill_missing);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    Bar **bars = (Bar **) malloc(sizeof(Bar **));
    for (int i = 0; i < array->count(); i++) {
        bars[i] = (Bar *) malloc(sizeof(Bar));
        memcpy(bars[i], &array->at(i), sizeof(Bar));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return bars;
}

struct Tick **
emc_history_l2ticks(const char *symbols, const char *start_time, const char *end_time, int adjust,
                    const char *adjust_end_time, bool skip_suspended, const char *fill_missing,
                    int *length, int *error_code) {
    auto array = history_l2ticks(symbols, start_time, end_time, adjust, (strlen(adjust_end_time) == 0 ? nullptr : adjust_end_time), skip_suspended, fill_missing);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    Tick **ticks = (Tick **) malloc(sizeof(Tick **));
    for (int i = 0; i < array->count(); i++) {
        ticks[i] = (Tick *) malloc(sizeof(Tick));
        memcpy(ticks[i], &array->at(i), sizeof(Tick));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return ticks;
}

Bar **
emc_history_l2bars(const char *symbols, const char *frequency, const char *start_time, const char *end_time, int adjust,
                   const char *adjust_end_time, bool skip_suspended, const char *fill_missing, int *length,
                   int *error_code) {
    auto array = history_l2bars(symbols, frequency, start_time, end_time, adjust, (strlen(adjust_end_time) == 0 ? nullptr : adjust_end_time), skip_suspended,
                                fill_missing);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    Bar **bars = (Bar **) malloc(sizeof(Bar **));
    for (int i = 0; i < array->count(); i++) {
        bars[i] = (Bar *) malloc(sizeof(Bar));
        memcpy(bars[i], &array->at(i), sizeof(Bar));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return bars;
}

L2Transaction **
emc_history_l2transactions(const char *symbols, const char *start_time, const char *end_time, int *length,
                           int *error_code) {
    auto array = history_l2transactions(symbols, start_time, end_time);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    L2Transaction **transactions = (L2Transaction **) malloc(sizeof(L2Transaction **));
    for (int i = 0; i < array->count(); i++) {
        transactions[i] = (L2Transaction *) malloc(sizeof(L2Transaction));
        memcpy(transactions[i], &array->at(i), sizeof(L2Transaction));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return transactions;
}

struct L2Order **
emc_history_l2orders(const char *symbols, const char *start_time, const char *end_time, int *length, int *error_code) {
    auto array = history_l2orders(symbols, start_time, end_time);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    L2Order **orders = (L2Order **) malloc(sizeof(L2Order **));
    for (int i = 0; i < array->count(); i++) {
        orders[i] = (L2Order *) malloc(sizeof(L2Order));
        memcpy(orders[i], &array->at(i), sizeof(L2Order));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return orders;
}

struct L2OrderQueue **
emc_history_l2orders_queue(const char *symbols, const char *start_time, const char *end_time, int *length,
                           int *error_code) {
    auto array = history_l2orders_queue(symbols, start_time, end_time);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *length = 0;
        *error_code = array->status();
        return nullptr;
    }
    L2OrderQueue **orders = (L2OrderQueue **) malloc(sizeof(L2OrderQueue **));
    for (int i = 0; i < array->count(); i++) {
        orders[i] = (L2OrderQueue *) malloc(sizeof(L2OrderQueue));
        memcpy(orders[i], &array->at(i), sizeof(L2OrderQueue));
    }
    *length = array->count();
    *error_code = 0;
    array->release();
    return orders;
}

std::map<int64_t, DataSet *> data_set;

void emc_free_dataset(int64_t obj_id) {
    auto data = data_set[obj_id];
    data->release();
    data_set.erase(obj_id);
}

int emc_dataset_status(int64_t obj_id) {
    auto data = data_set[obj_id];
    return data->status();
}

bool emc_dataset_is_end(int64_t obj_id) {
    auto data = data_set[obj_id];
    return data->is_end();
}

void emc_dataset_next(int64_t obj_id) {
    auto data = data_set[obj_id];
    data->next();
}

int emc_dataset_get_integer(int64_t obj_id, const char *key) {
    auto data = data_set[obj_id];
    return data->get_integer(key);
}

int64_t emc_dataset_get_long_integer(int64_t obj_id, const char *key) {
    auto data = data_set[obj_id];
    return data->get_long_integer(key);
}

double emc_dataset_get_real(int64_t obj_id, const char *key) {
    auto data = data_set[obj_id];
    return data->get_real(key);
}

const char *emc_dataset_get_string(int64_t obj_id, const char *key) {
    auto data = data_set[obj_id];
    return data->get_string(key);
}

const char *emc_dataset_debug_string(int64_t obj_id) {
    auto data = data_set[obj_id];
    return data->debug_string();
}

void emc_get_fundamentals(const char *table, const char *symbols, const char *start_date, const char *end_date,
                          const char *fields, const char *filter, const char *order_by, int limit, int64_t *obj_id,
                          int *error_code) {

    auto data = get_fundamentals(table, symbols, start_date, end_date, fields, filter, order_by, limit);
    if (data == nullptr) {
        *error_code = -1;
        *obj_id = 0;
        return;
    }
    data_set.insert(std::make_pair(epoch::us(), data));
    *error_code = 0;
}

void emc_get_fundamentals_n(const char *table, const char *symbols, const char *end_date, const char *fields, int n,
                            const char *filter, const char *order_by, int64_t *obj_id, int *error_code) {
    auto data = get_fundamentals_n(table, symbols, end_date, fields, n, filter, order_by);
    if (data == nullptr) {
        *error_code = -1;
        *obj_id = 0;
        return;
    }
    data_set.insert(std::make_pair(epoch::us(), data));
    *error_code = 0;
}

void emc_get_instruments(const char *exchanges, const char *sec_types, const char *fields, int64_t *obj_id,
                         int *error_code) {
    auto data = get_instruments(exchanges, sec_types, fields);
    if (data == nullptr) {
        *error_code = -1;
        *obj_id = 0;
        return;
    }
    data_set.insert(std::make_pair(epoch::us(), data));
    *error_code = 0;
}

void emc_get_history_instruments(const char *symbols, const char *start_date, const char *end_date, const char *fields,
                                 int64_t *obj_id, int *error_code) {
    auto data = get_history_instruments(symbols, start_date, end_date, fields);
    if (data == nullptr) {
        *error_code = -1;
        *obj_id = 0;
        return;
    }
    data_set.insert(std::make_pair(epoch::us(), data));
    *error_code = 0;
}

void emc_get_instrumentinfos(const char *symbols, const char *exchanges, const char *sec_types, const char *names,
                             const char *fields, int64_t *obj_id, int *error_code) {
    auto data = get_instrumentinfos(symbols, exchanges, sec_types, names, fields);
    if (data == nullptr) {
        *error_code = -1;
        *obj_id = 0;
        return;
    }
    data_set.insert(std::make_pair(epoch::us(), data));
    *error_code = 0;
}

void emc_get_constituents(const char *index, const char *trade_date, int64_t *obj_id, int *error_code) {
    auto data = get_constituents(index, trade_date);
    if (data == nullptr) {
        *error_code = -1;
        *obj_id = 0;
        return;
    }
    data_set.insert(std::make_pair(epoch::us(), data));
    *error_code = 0;
}

char *emc_get_industry(const char *code, int *length, int *error_code) {
    auto array = get_industry(code);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *error_code = array->status();
        *length = 0;
        return nullptr;
    }

    auto json = json::array();
    for (int i = 0; i < array->count(); i++) {
        json.push_back(gbk_to_utf8(array->at(i)));
    }
    *length = array->count();
    auto str = (char *) malloc(json.str().size() + 1);
    strcpy(str, json.str().c_str());
    *error_code = 0;
    return str;
}

char *emc_get_concept(const char *code, int *length, int *error_code) {
    auto array = get_concept(code);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *error_code = array->status();
        *length = 0;
        return nullptr;
    }

    auto json = json::array();
    for (int i = 0; i < array->count(); i++) {
        json.push_back(gbk_to_utf8(array->at(i)));
    }
    *length = array->count();
    auto str = (char *) malloc(json.str().size() + 1);
    strcpy(str, json.str().c_str());
    *error_code = 0;
    return str;
}

char *emc_get_trading_dates(const char *exchange, const char *start_date, const char *end_date, int *length,
                            int *error_code) {
    auto array = get_trading_dates(exchange, start_date, end_date);
    if (array == nullptr) {
        *error_code = -1;
        *length = 0;
        return nullptr;
    }
    if (array->status() != 0) {
        *error_code = array->status();
        *length = 0;
        return nullptr;
    }
    auto json = json::array();
    for (int i = 0; i < array->count(); i++) {
        json.push_back(gbk_to_utf8(array->at(i)));
    }
    *length = array->count();
    auto str = (char *) malloc(json.str().size() + 1);
    strcpy(str, json.str().c_str());
    *error_code = 0;
    return str;
}

int emc_get_previous_trading_date(const char *exchange, const char *date, char *output_date) {
    return get_previous_trading_date(exchange, date, output_date);
}

int emc_get_next_trading_date(const char *exchange, const char *date, char *output_date) {
    return get_next_trading_date(exchange, date, output_date);
}

void
emc_get_dividend(const char *symbol, const char *start_date, const char *end_date, int64_t *obj_id, int *error_code) {
    auto data = get_dividend(symbol, start_date, end_date);
    if (data == nullptr) {
        *error_code = -1;
        *obj_id = 0;
        return;
    }
    data_set.insert(std::make_pair(epoch::us(), data));
    *error_code = 0;
}

void emc_get_continuous_contracts(const char *symbol, const char *start_date, const char *end_date, int64_t *obj_id,
                                  int *error_code) {
    auto data = get_continuous_contracts(symbol, start_date, end_date);
    if (data == nullptr) {
        *error_code = -1;
        *obj_id = 0;
        return;
    }
    data_set.insert(std::make_pair(epoch::us(), data));
    *error_code = 0;
}

void emc_l2transactions_free(struct L2Transaction **transactions, int length) {
    if (transactions == nullptr) {
        return;
    }
    if (length == 1) {
        free(transactions[0]);
        free(transactions);
    } else {
        for (int i = 0; i < length; i++) {
            free(transactions[i]);
        }
        free(transactions);
    }
}

void emc_l2orders_free(struct L2Order **orders, int length) {
    if (orders == nullptr) {
        return;
    }
    if (length == 1) {
        free(orders[0]);
        free(orders);
    } else {
        for (int i = 0; i < length; i++) {
            free(orders[i]);
        }
        free(orders);
    }
}

void emc_l2orders_queue_free(struct L2OrderQueue **orders, int length) {
    if (orders == nullptr) {
        return;
    }
    if (length == 1) {
        free(orders[0]);
        free(orders);
    } else {
        for (int i = 0; i < length; i++) {
            free(orders[i]);
        }
        free(orders);
    }
}




