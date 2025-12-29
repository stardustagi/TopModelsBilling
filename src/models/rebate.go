package models

// UserRebateConfig 用户消费返点阶梯配置表
type UserRebateConfig struct {
	Id         int64 `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	UserId     int64 `json:"user_id" xorm:"'user_id' BIGINT(20)"`
	TierStart  int64 `json:"tier_start" xorm:"'tier_start' BIGINT(20)"` // 阶梯起始金额
	TierEnd    int64 `json:"tier_end" xorm:"'tier_end' BIGINT(20)"`     // 阶梯结束金额，-1表示无上限
	RebateRate int   `json:"rebate_rate" xorm:"'rebate_rate' INT(10)"`  // 返点比例，如10表示10%
	Status     int   `json:"status" xorm:"'status' INT(10) default 1"`
	CreatedAt  int64 `json:"created_at" xorm:"'created_at' BIGINT(20)"`
	UpdatedAt  int64 `json:"updated_at" xorm:"'updated_at' BIGINT(20)"`
}

func (o *UserRebateConfig) TableName() string {
	return "user_rebate_config"
}

// UserRebateMonthly 用户月度返点记录表
type UserRebateMonthly struct {
	Id            int64  `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	UserId        int64  `json:"user_id" xorm:"'user_id' BIGINT(20)"`
	Month         string `json:"month" xorm:"'month' VARCHAR(7)"`                             // 格式：2024-01
	TotalConsumed int64  `json:"total_consumed" xorm:"'total_consumed' BIGINT(20) default 0"` // 当月消费总额
	RebateAmount  int64  `json:"rebate_amount" xorm:"'rebate_amount' BIGINT(20) default 0"`   // 已返点金额
	RebateUsed    int64  `json:"rebate_used" xorm:"'rebate_used' BIGINT(20) default 0"`       // 已消费返点
	RebateRate    int    `json:"rebate_rate" xorm:"'rebate_rate' INT(10)"`                    // 返点比例快照
	Status        int    `json:"status" xorm:"'status' INT(10) default 0"`                    // 0未返点 1已返点
	CreatedAt     int64  `json:"created_at" xorm:"'created_at' BIGINT(20)"`
	UpdatedAt     int64  `json:"updated_at" xorm:"'updated_at' BIGINT(20)"`
}

func (o *UserRebateMonthly) TableName() string {
	return "user_rebate_monthly"
}
