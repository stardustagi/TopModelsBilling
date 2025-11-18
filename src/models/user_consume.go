package models

// UserConsumeRecord 表示用户消费记录
type UserConsumeRecord struct {
	ID               int64  `xorm:"pk autoincr comment('主键，自增')" json:"id"`                   // 主键，自增
	UserId           int64  `xorm:"user_id int notnull index comment('用户ID')" json:"user_id"` // 用户ID
	NodeId           string `json:"node_id" xorm:"'node_id' VARCHAR(64)"`
	DiscountAmount   int64  `xorm:"bigint default 0 comment('折扣数量')" json:"discount_amount"`     // 折扣数量
	TotalConsumed    int64  `xorm:"bigint default 0 comment('本次使用的币数量')" json:"total_consumed"`  // 本次扣费数量
	Caller           string `xorm:"varchar(64) index comment('调用方')" json:"caller"`              // 调用方
	Model            string `xorm:"varchar(64) comment('模型')" json:"model"`                      // 模型
	ModelId          string `xorm:"varchar(64) comment('模型id')" json:"model_id"`                 // 模型id
	ActualProvider   string `xorm:"varchar(64) comment('服务商')" son:"actual_provider"`            // 实际服务商
	ActualProviderId string `xorm:"varchar(64) comment('服务商id')" json:"actual_provider_id"`      // 实际服务商id
	ConsumeType      string `xorm:"varchar(255) default '' comment('消费类型')" json:"consume_type"` // 消费类型
	CreatedAt        int64  `xorm:"created_at comment('创建时间')" json:"created"`                   // 创建时间
	UpdatedAt        int64  `xorm:"updated_at comment('更新时间')" json:"updated"`                   // 更新时间
}

func (UserConsumeRecord) TableName() string {
	return "user_consume"
}

type UserConsumeDetailText struct {
	ID           int64 `xorm:"pk autoincr comment('主键，自增')" json:"id"`                    // 主键，自增
	ConsumdId    int64 `xorm:"consume_id comment('消费记录id')" json:"consume_id"`            // 消费记录id
	InputTokens  int64 `xorm:"bigint default 0 comment('输入token数')" json:"input_tokens"`  // 输入token数
	OutputTokens int64 `xorm:"bigint default 0 comment('输出token数')" json:"output_tokens"` // 输出token数
	CacheTokens  int64 `xorm:"bigint default 0 comment('缓存token数')" json:"cache_tokens"`  // 缓存token数
	InputPrice   int   `xorm:"int default 0 comment('输入token价格')" json:"input_price"`     // 输入token价格
	OutputPrice  int   `xorm:"int default 0 comment('输出token价格')" json:"output_price"`    // 输出token价格
	CachePrice   int   `xorm:"int default 0 comment('缓存token价格')" json:"cache_price"`     // 缓存token价格
	CreatedAt    int64 `xorm:"created_at comment('创建时间')" json:"created"`                 // 创建时间
}

func (UserConsumeDetailText) TableName() string {
	return "user_consume_detail_text"
}

// UserConsumeDetailImage 图片消费明细，可能是多张
type UserConsumeDetailImage struct {
	ID        int64  `xorm:"pk autoincr comment('主键，自增')" json:"id"`        // 主键，自增
	Quality   string `xorm:"varchar(64) comment('Quality')" json:"quality"` // Quality
	Size      string `xorm:"varchar(64) comment('Size')" json:"size"`       // Size
	CreatedAt int64  `xorm:"created_at comment('创建时间')" json:"created"`     // 创建时间
}

func (UserConsumeDetailImage) TableName() string {
	return "user_consume_detail_image"
}

type UserConsumeDetailVideo struct {
	ID        int64   `xorm:"pk autoincr comment('主键，自增')" json:"id"`         // 主键，自增
	ConsumdId int64   `xorm:"consume_id comment('消费记录id')" json:"consume_id"` // 消费记录id
	Seconds   float64 `xorm:"double comment('Seconds')" json:"Seconds"`       // Seconds
	Size      string  `xorm:"varchar(64) comment('Size')" json:"size"`        // Size
	CreatedAt int64   `xorm:"created_at comment('创建时间')" json:"created"`      // 创建时间
}

func (UserConsumeDetailVideo) TableName() string {
	return "user_consume_detail_video"
}
