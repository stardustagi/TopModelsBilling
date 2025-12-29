package models

// ProviderConsumeSummary 供应商消费汇总表
type ProviderConsumeSummary struct {
	ID               int64 `xorm:"pk autoincr comment('主键，自增')" json:"id"`
	ActualProviderId int   `xorm:"int notnull unique comment('实际服务商ID')" json:"actual_provider_id"`
	TotalConsumed    int64 `xorm:"bigint default 0 comment('总消费金额')" json:"total_consumed"`
	TotalCost        int64 `xorm:"bigint default 0 comment('总成本')" json:"total_cost"`
	UpdatedAt        int64 `xorm:"updated_at comment('更新时间')" json:"updated_at"`
}

func (ProviderConsumeSummary) TableName() string {
	return "provider_consume_summary"
}

// ProviderModelDailySummary 供应商模型日消费汇总表
type ProviderModelDailySummary struct {
	ID               int64  `xorm:"pk autoincr comment('主键，自增')" json:"id"`
	UserId           int64  `xorm:"bigint notnull comment('用户ID')" json:"user_id"`
	ActualProviderId int    `xorm:"int notnull comment('实际服务商ID')" json:"actual_provider_id"`
	ModelId          int    `xorm:"int notnull comment('模型ID')" json:"model_id"`
	ConsumeType      string `xorm:"varchar(32) notnull comment('消费类型')" json:"consume_type"`
	Date             string `xorm:"varchar(10) notnull comment('日期YYYY-MM-DD')" json:"date"`
	TotalConsumed    int64  `xorm:"bigint default 0 comment('总消费金额')" json:"total_consumed"`
	TotalCost        int64  `xorm:"bigint default 0 comment('总成本')" json:"total_cost"`
	UpdatedAt        int64  `xorm:"updated_at comment('更新时间')" json:"updated_at"`
}

func (ProviderModelDailySummary) TableName() string {
	return "provider_model_daily_summary"
}
