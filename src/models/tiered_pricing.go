package models

type ModelsTieredPricing struct {
	Id          int64 `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	ModelId     int64 `json:"model_id" xorm:"'model_id' BIGINT(20)"`
	TierStart   int64 `json:"tier_start" xorm:"'tier_start' BIGINT(20)"`
	TierEnd     int64 `json:"tier_end" xorm:"'tier_end' BIGINT(20)"`
	InputPrice  int   `json:"input_price" xorm:"'input_price' INT(10)"`
	OutputPrice int   `json:"output_price" xorm:"'output_price' INT(10)"`
	CachePrice  int   `json:"cache_price" xorm:"'cache_price' INT(10)"`
	CreatedAt   int64 `json:"created_at" xorm:"'created_at' BIGINT(20)"`
	UpdatedAt   int64 `json:"updated_at" xorm:"'updated_at' BIGINT(20)"`
}

func (o *ModelsTieredPricing) TableName() string {
	return "models_tiered_pricing"
}

type UserDiscount struct {
	Id           int64 `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	UserId       int64 `json:"user_id" xorm:"'user_id' BIGINT(20)"`
	DiscountRate int   `json:"discount_rate" xorm:"'discount_rate' INT(10)"`
	CreatedAt    int64 `json:"created_at" xorm:"'created_at' BIGINT(20)"`
	UpdatedAt    int64 `json:"updated_at" xorm:"'updated_at' BIGINT(20)"`
}

func (o *UserDiscount) TableName() string {
	return "user_discount"
}
