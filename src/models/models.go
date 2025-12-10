package models

type ModelsInfo struct {
	Id          int64  `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	Name        string `json:"name" xorm:"'name' comment('模型名') VARCHAR(128)"`
	ApiVersion  string `json:"api_version" xorm:"'api_version' VARCHAR(24)"`
	DeployName  string `json:"deploy_name" xorm:"'deploy_name' VARCHAR(128)"`
	InputPrice  int    `json:"input_price" xorm:"'input_price' INT(10)"`
	OutputPrice int    `json:"output_price" xorm:"'output_price' INT(10)"`
	CachePrice  int    `json:"cache_price" xorm:"'cache_price' INT(10)"`
	Status      string `json:"status" xorm:"'status' comment('模型状态') VARCHAR(12)"`
	LastUpdate  int64  `json:"last_update" xorm:"'last_update' comment('最后更新时间') BIGINT(20)"`
}

func (o *ModelsInfo) TableName() string {
	return "models_info"
}

type ModelsProvider struct {
	Id          int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	OwnerId     int64  `json:"owner_id" xorm:"'owner_id' comment('模型供应商的nodeUserId') BIGINT(12)"`
	ProviderId  string `json:"provider_id" xorm:"'provider_id' VARCHAR(128)"`
	Type        string `json:"type" xorm:"'type' VARCHAR(64)"`
	Name        string `json:"name" xorm:"'name' VARCHAR(128)"`
	Endpoint    string `json:"endpoint" xorm:"'endpoint' VARCHAR(128)"`
	ApiType     string `json:"api_type" xorm:"'api_type' VARCHAR(64)"`
	ModelName   string `json:"model_name" xorm:"'model_name' VARCHAR(64)"`
	InputPrice  int    `json:"input_price" xorm:"'input_price' INT(10)"`
	OutputPrice int    `json:"output_price" xorm:"'output_price' INT(10)"`
	CachePrice  int    `json:"cache_price" xorm:"'cache_price' INT(10)"`
	ApiKeys     string `json:"api_keys" xorm:"'api_keys' comment('apikeys列表') TEXT"`
	Deleted     int64  `json:"deleted" xorm:"'deleted' BIGINT(12)"`
	LastUpdate  int64  `json:"last_update" xorm:"'last_update' BIGINT(12)"`
}

func (o *ModelsProvider) TableName() string {
	return "models_provider"
}
