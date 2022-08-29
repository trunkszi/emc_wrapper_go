#ifndef INTERFACE_H
#define INTERFACE_H

#include "third_party/gmsdk/include/gmdef.h"
#include <stdbool.h>


#ifdef __cplusplus
extern "C" {
#endif

#if defined(_WIN32)
#   define EMCAPI         __declspec(dllexport)
#elif defined(__GNUC__) && ((__GNUC__ >= 4) || (__GNUC__ == 3 && __GNUC_MINOR__ >= 3))
#   define EMCAPI         __attribute__((visibility("default")))
#else
#   define EMCAPI
#endif
typedef  long long int64_t;

typedef void(*on_init_callback)();

typedef void(*on_tick_callback)(struct Tick *tick);

typedef void(*on_bar_callback)(struct Bar *bar);

typedef void(*on_l2transaction_callback)(struct L2Transaction *l2transaction);

typedef void(*on_l2order_callback)(struct L2Order *l2order);

typedef void(*on_l2order_queue_callback)(struct L2OrderQueue *l2queue);

typedef void(*on_order_status_callback)(struct Order *order);

typedef void(*on_execution_report_callback)(struct ExecRpt *rpt);

typedef void(*on_algo_order_status_callback)(struct AlgoOrder *order);

typedef void(*on_cash_callback)(struct Cash *cash);

typedef void(*on_position_callback)(struct Position *position);

typedef void(*on_parameter_callback)(struct Parameter *param);

typedef void(*on_schedule_callback)(char *data_rule, char *time_rule);

typedef void(*on_backtest_finished_callback)(struct Indicator *indicator);

typedef void(*on_account_status_callback)(struct AccountStatus *account_status);

typedef void(*on_error_callback)(int error_code, char *error_msg);

typedef void(*on_stop_callback)();

typedef void(*on_market_data_connected_callback)();

typedef void(*on_trade_data_connected_callback)();

typedef void(*on_market_data_disconnected_callback)();

typedef void(*on_trade_data_disconnected_callback)();

EMCAPI void create_instance();

EMCAPI
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
                        on_trade_data_disconnected_callback trade_data_disconnected_callback);

EMCAPI int emc_run();

EMCAPI void emc_stop();

EMCAPI void emc_set_strategy_id(const char *strategy_id);

EMCAPI void emc_set_token(const char *token);

EMCAPI void emc_set_mode(int mode);

EMCAPI int emc_schedule(const char *data_rule, const char *time_rule);

EMCAPI double emc_now();

EMCAPI int emc_set_backtest_config(
        const char *start_time,
        const char *end_time,
        double initial_cash,
        double transaction_ratio,
        double commission_ratio,
        double slippage_ratio,
        int adjust,
        int check_cache
);

EMCAPI int emc_subscribe(const char *symbols, const char *frequency, bool unsubscribe_previous);

EMCAPI int emc_unsubscribe(const char *symbols, const char *frequency);


EMCAPI char *emc_get_accounts(int *length, int *error_code);

EMCAPI char *emc_get_account_status(const char *account, int *length, int *error_code);


EMCAPI char *emc_get_all_account_status(int *length, int *error_code);


EMCAPI  char *
emc_order_volume(const char *symbol, int volume, int side, int order_type, int position_effect, double price,
                 const char *account, int *length, int *error_code);

EMCAPI char *
emc_order_value(const char *symbol, double value, int side, int order_type, int position_effect, double price,
                const char *account, int *length, int *error_code);

EMCAPI char *
emc_order_percent(const char *symbol, double percent, int side, int order_type, int position_effect, double price,
                  const char *account, int *length, int *error_code);

EMCAPI char *
emc_order_target_volume(const char *symbol, int volume, int position_side, int order_type, double price,
                        const char *account, int *length, int *error_code);

EMCAPI char *
emc_order_target_value(const char *symbol, double value, int position_side, int order_type, double price,
                       const char *account, int *length, int *error_code);

EMCAPI char *
emc_order_target_percent(const char *symbol, double percent, int position_side, int order_type, double price,
                         const char *account, int *length, int *error_code);

EMCAPI char *emc_order_close_all(int *length, int *error_code);

EMCAPI int emc_order_cancel(const char *cl_ord_id, const char *account);

EMCAPI int emc_order_cancel_all();

EMCAPI char *
emc_place_order(const char *symbol, int volume, int side, int order_type, int position_effect, double price,
                int order_duration, int order_qualifier, double stop_price, int order_business,
                const char *account, int *length, int *error_code);

EMCAPI char *
emc_order_after_hour(const char *symbol, int volume, int side, double price, const char *account, int *length,
                     int *error_code);

EMCAPI  char *emc_get_orders(const char *account, int *length, int *error_code);

EMCAPI char *emc_get_unfinished_orders(const char *account, int *length, int *error_code);

EMCAPI char *emc_get_execution_reports(const char *account, int *length, int *error_code);

EMCAPI
char *emc_get_cash(const char *accounts, int *length, int *error_code);

EMCAPI
char *emc_get_position(const char *account, int *length, int *error_code);


EMCAPI char *
emc_order_algo(const char *symbol, int volume, int position_effect, int side, int order_type, double price,
               const char *algo_name, const char *algo_param, const char *account, int *length, int *error_code);


EMCAPI
int emc_algo_order_cancel(const char *cl_ord_id, const char *account);

EMCAPI
int emc_algo_order_pause(const char *cl_ord_id, int status, const char *account);


EMCAPI
char *emc_get_algo_orders(const char *account, int *length, int *error_code);

EMCAPI
char *
emc_get_algo_child_orders(const char *cl_ord_id, const char *account, int *length, int *error_code);

EMCAPI
int emc_raw_func(const char *account, const char *func_id, const char *func_args, char **rsp);


EMCAPI char *
emc_credit_buying_on_margin(int position_src, const char *symbol, int volume, double price, int order_type,
                            int order_duration, int order_qualifier, const char *account, int *length, int *error_code);

EMCAPI char *
emc_credit_short_selling(int position_src, const char *symbol, int volume, double price, int order_type,
                         int order_duration, int order_qualifier, const char *account, int *length, int *error_code);

EMCAPI char *
emc_credit_repay_share_by_buying_share(const char *symbol, int volume, double price, int order_type, int order_duration,
                                       int order_qualifier, const char *account, int *length, int *error_code);

EMCAPI char *
emc_credit_repay_cash_by_selling_share(const char *symbol, int volume, double price, int order_type, int order_duration,
                                       int order_qualifier, const char *account, int *length, int *error_code);

EMCAPI char *
emc_credit_buying_on_collateral(const char *symbol, int volume, double price, int order_type,
                                int order_duration, int order_qualifier,
                                const char *account, int *length, int *error_code);

EMCAPI char *
emc_credit_selling_on_collateral(const char *symbol, int volume, double price, int order_type, int order_duration,
                                 int order_qualifier, const char *account, int *length, int *error_code);

EMCAPI char *
emc_credit_repay_share_directly(const char *symbol, int volume, const char *account, int *length, int *error_code);

EMCAPI
int emc_credit_repay_cash_directly(double amount, const char *account, double *actual_repay_amount, char *error_msg_buf,
                                   int buf_len);

EMCAPI char *
emc_credit_collateral_in(const char *symbol, int volume, const char *account, int *length, int *error_code);

EMCAPI char *
emc_credit_collateral_out(const char *symbol, int volume, const char *account, int *length, int *error_code);

EMCAPI char *emc_credit_get_collateral_instruments(const char *account, int *length, int *error_code);


EMCAPI char *
emc_credit_get_borrowable_instruments(int position_src, const char *account, int *length, int *error_code);


EMCAPI char *
emc_credit_get_borrowable_instruments_positions(int position_src, const char *account, int *length, int *error_code);


EMCAPI char *emc_credit_get_contracts(int position_src, const char *account, int *length, int *error_code);

EMCAPI
char *emc_credit_get_cash(const char *account, int *length, int *error_code);


EMCAPI char *
emc_ipo_buy(const char *symbol, int volume, double price, const char *account, int *length, int *error_code);

EMCAPI
char *emc_ipo_get_quota(const char *account, int *length, int *error_code);


EMCAPI
char *emc_ipo_get_instruments(int security_type, const char *account, int *length, int *error_code);


EMCAPI
char *emc_ipo_get_match_number(const char *account, int *length, int *error_code);


EMCAPI
char *emc_ipo_get_lot_info(const char *account, int *length, int *error_code);


EMCAPI
char *emc_fund_etf_buy(const char *symbol, int volume, double price, const char *account, int *length, int *error_code);

EMCAPI
char *emc_fund_etf_redemption(const char *symbol, int volume, double price, const char *account, int *length,
                              int *error_code);

EMCAPI
char *emc_fund_subscribing(const char *symbol, int volume, const char *account, int *length, int *error_code);

EMCAPI
char *emc_fund_buy(const char *symbol, int volume, const char *account, int *length, int *error_code);

EMCAPI
char *emc_fund_redemption(const char *symbol, int volume, const char *account, int *length, int *error_code);

EMCAPI
char *
emc_bond_reverse_repurchase_agreement(const char *symbol, int volume, double price, int order_type, int order_duration,
                                      int order_qualifier, const char *account, int *length, int *error_code);

EMCAPI
char *emc_bond_convertible_call(const char *symbol, int volume, double price, const char *account, int *length,
                                int *error_code);

EMCAPI
char *emc_bond_convertible_put(const char *symbol, int volume, double price, const char *account, int *length,
                               int *error_code);

EMCAPI
char *
emc_bond_convertible_put_cancel(const char *symbol, int volume, const char *account, int *length, int *error_code);

EMCAPI
int emc_add_parameters(struct Parameter *params, int count);

EMCAPI
int emc_del_parameters(const char *keys);

EMCAPI
int emc_set_parameters(struct Parameter *params, int count);

EMCAPI
char *emc_get_parameters(int *length, int *error_code);


EMCAPI int emc_set_symbols(const char *symbols);

EMCAPI char *emc_get_symbols(int *length, int *error_code);


///////////////////////////////////////////////////////////////////////////////////////////////////////////
// gmapi.h


// 获取sdk版本号
EMCAPI const char *emc_get_version();

// 设置token
EMCAPI void gm_set_token(const char *token);

// 自定义服务地址
EMCAPI void emc_set_serv_addr(const char *addr);

//第三方系统设置留痕信息
EMCAPI void emc_set_mfp(const char *mfp);

// 查询当前行情快照
EMCAPI  struct Tick **
emc_current(const char *symbols, int *length, int *error_code);

EMCAPI void emc_ticks_free(struct Tick **ticks, int length);

// 查询历史Tick行情
EMCAPI  struct Tick **
emc_history_ticks(const char *symbols, const char *start_time, const char *end_time, int adjust,
                  const char *adjust_end_time, bool skip_suspended, const char *fill_missing, int *length,
                  int *error_code);

// 查询历史Bar行情
EMCAPI struct Bar **
emc_history_bars(const char *symbols, const char *frequency, const char *start_time, const char *end_time, int adjust,
                 const char *adjust_end_time, bool skip_suspended, const char *fill_missing, int *length,
                 int *error_code);

EMCAPI void emc_bars_free(struct Bar **bars, int length);


// 查询最新n条Tick行情
EMCAPI  struct Tick **
emc_history_ticks_n(const char *symbols, int n, const char *end_time, int adjust, const char *adjust_end_time,
                    bool skip_suspended, const char *fill_missing, int *length, int *error_code);

//查询最新n条Bar行情
EMCAPI struct Bar **
emc_history_bars_n(const char *symbols, const char *frequency, int n, const char *end_time, int adjust,
                   const char *adjust_end_time, bool skip_suspended, const char *fill_missing, int *length,
                   int *error_code);

// 查询历史L2 Tick行情
EMCAPI  struct Tick **
emc_history_l2ticks(const char *symbols, const char *start_time, const char *end_time, int adjust,
                    const char *adjust_end_time, bool skip_suspended, const char *fill_missing, int *length,
                    int *error_code);

// 查询历史L2 Bar行情
EMCAPI struct Bar **
emc_history_l2bars(const char *symbols, const char *frequency, const char *start_time, const char *end_time, int adjust,
                   const char *adjust_end_time, bool skip_suspended, const char *fill_missing, int *length,
                   int *error_code);

// 查询历史L2 逐笔成交
EMCAPI struct L2Transaction **
emc_history_l2transactions(const char *symbols, const char *start_time, const char *end_time, int *length,
                           int *error_code);
EMCAPI void emc_l2transactions_free(struct L2Transaction **transactions, int length);

// 查询历史L2 逐笔委托
EMCAPI struct L2Order **
emc_history_l2orders(const char *symbols, const char *start_time, const char *end_time, int *length, int *error_code);
EMCAPI void emc_l2orders_free(struct L2Order **orders, int length);

// 查询历史L2 委托队列(最优价最大50笔委托量)
EMCAPI struct L2OrderQueue **
emc_history_l2orders_queue(const char *symbols, const char *start_time, const char *end_time, int *length,
                           int *error_code);

EMCAPI void emc_l2orders_queue_free(struct L2OrderQueue **orders, int length);

EMCAPI void emc_free_dataset(int64_t obj_id);
EMCAPI int emc_dataset_status(int64_t obj_id);
EMCAPI void emc_dataset_next(int64_t obj_id);
EMCAPI int emc_dataset_get_integer(int64_t obj_id, const char *key);
EMCAPI int64_t emc_dataset_get_long_integer(int64_t obj_id, const char *key);
EMCAPI double emc_dataset_get_real(int64_t obj_id, const char *key);
EMCAPI const char *emc_dataset_get_string(int64_t obj_id, const char *key);
EMCAPI const char *emc_dataset_debug_string(int64_t obj_id);

// 查询基本面数据
EMCAPI void emc_get_fundamentals(const char *table, const char *symbols, const char *start_date, const char *end_date,
                                 const char *fields, const char *filter, const char *order_by, int limit,
                                 int64_t *obj_id,
                                 int *error_code);

// 查询基本面数据最新n条
EMCAPI void
emc_get_fundamentals_n(const char *table, const char *symbols, const char *end_date, const char *fields, int n,
                       const char *filter, const char *order_by, int64_t *obj_id, int *error_code);

// 查询最新交易标的信息
EMCAPI void
emc_get_instruments(const char *exchanges, const char *sec_types, const char *fields, int64_t *obj_id, int *error_code);

// 查询交易标的历史数据
EMCAPI void
emc_get_history_instruments(const char *symbols, const char *start_date, const char *end_date, const char *fields,
                            int64_t *obj_id, int *error_code);

// 查询交易标的基本信息
EMCAPI void
emc_get_instrumentinfos(const char *symbols, const char *exchanges, const char *sec_types, const char *names,
                        const char *fields, int64_t *obj_id, int *error_code);

// 查询指数成份股
EMCAPI void emc_get_constituents(const char *index, const char *trade_date, int64_t *obj_id, int *error_code);

// 查询行业股票列表
EMCAPI char *emc_get_industry(const char *code, int *length, int *error_code);

// 查询概念板块股票列表
EMCAPI char *emc_get_concept(const char *code, int *length, int *error_code);

// 查询交易日列表
EMCAPI char *
emc_get_trading_dates(const char *exchange, const char *start_date, const char *end_date, int *length, int *error_code);

// 返回指定日期的上一个交易日
EMCAPI int emc_get_previous_trading_date(const char *exchange, const char *date, char *output_date);

// 返回指定日期的下一个交易日
EMCAPI int emc_get_next_trading_date(const char *exchange, const char *date, char *output_date);

// 查询分红送配
EMCAPI void
emc_get_dividend(const char *symbol, const char *start_date, const char *end_date, int64_t *obj_id, int *error_code);

// 获取连续合约
EMCAPI void
emc_get_continuous_contracts(const char *symbol, const char *start_date, const char *end_date, int64_t *obj_id,
                             int *error_code);


#ifdef __cplusplus
}
#endif

#endif // INTERFACE_H