package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserWallet struct {
	Id            int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	UserId        int64  `json:"user_id" xorm:"'user_id' BIGINT(12)"`
	WalletType    string `json:"wallet_type" xorm:"'wallet_type' VARCHAR(32)"`
	WalletAddress string `json:"wallet_address" xorm:"'wallet_address' VARCHAR(255)"`
	Balance       int64  `json:"balance" xorm:"'balance' BIGINT(12)"`
	CreatedAt     int64  `json:"created_at" xorm:"'created_at' BIGINT(12)"`
	UpdatedAt     int64  `json:"updated_at" xorm:"'updated_at' BIGINT(20)"`
}

func (o *UserWallet) TableName() string {
	return "user_wallet"
}

func (o *UserWallet) GetSliceName(slice string) string {
	return fmt.Sprintf("user_wallet_%s", slice)
}

func (o *UserWallet) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("user_wallet_%d%02d", t.Year(), t.Month())
}

func (o *UserWallet) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("user_wallet_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *UserWallet) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *UserWallet) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *UserWallet) PrimaryKey() interface{} {
	return o.Id
}
