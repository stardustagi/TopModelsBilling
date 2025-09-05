package models

// UserConsumeRecord 表示用户消费记录
type UserConsumeRecord struct {
	ID               int64  `xorm:"pk autoincr comment('主键，自增')" json:"id"`                   // 主键，自增
	UserId           int64  `xorm:"user_id int notnull index comment('用户ID')" json:"user_id"` // 用户ID
	NodeId           string `json:"node_id" xorm:"'node_id' VARCHAR(64)"`
	TotalConsumed    int64  `xorm:"bigint default 0 comment('本次使用的币数量')" json:"total_consumed"` // 本次扣费数量
	Caller           string `xorm:"varchar(64) index comment('调用方')" json:"caller"`             // 调用方
	Model            string `xorm:"varchar(64) comment('模型')" json:"model"`                     // 模型
	ModelId          string `xorm:"varchar(64) comment('模型id')" json:"model_id"`                // 模型id
	ActualProvider   string `xorm:"varchar(64) comment('服务商')" son:"actual_provider"`           // 实际服务商
	ActualProviderId string `xorm:"varchar(64) comment('服务商id')" json:"actual_provider_id"`     // 实际服务商id
	InputTokens      int64  `xorm:"bigint default 0 comment('输入token数')" json:"input_tokens"`   // 输入token数
	OutputTokens     int64  `xorm:"bigint default 0 comment('输出token数')" json:"output_tokens"`  // 输出token数
	CacheTokens      int64  `xorm:"bigint default 0 comment('缓存token数')" json:"cache_tokens"`   // 缓存token数
	InputPrice       int    `xorm:"int default 0 comment('输入token价格')" json:"input_price"`      // 输入token价格
	OutputPrice      int    `xorm:"int default 0 comment('输出token价格')" json:"output_price"`     // 输出token价格
	CachePrice       int    `xorm:"int default 0 comment('缓存token价格')" json:"cache_price"`      // 缓存token价格
	CreatedAt        int64  `xorm:"created comment('创建时间')" json:"created"`                     // 创建时间
	UpdatedAt        int64  `xorm:"updated comment('更新时间')" json:"updated"`                     // 更新时间
}

func (UserConsumeRecord) TableName() string {
	return "user_consume_record"
}
