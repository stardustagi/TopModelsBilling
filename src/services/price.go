package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

type PriceInfo struct {
	InputPrice  int `json:"input_price"`  //输入token计费
	OutputPrice int `json:"output_price"` //输出token计费
	CachePrice  int `json:"cache_price"`  //缓存token计费
}

func (o PriceInfo) String() string {
	return fmt.Sprintf("<Price: input:%d, output:%d>", o.InputPrice, o.OutputPrice)
}

type PriceService struct {
	ctx       context.Context
	xorm      xorm.EngineInterface
	PriceInfo map[string]PriceInfo
	mutex     sync.Mutex
}

func NewPriceService(ctx context.Context, xorm xorm.EngineInterface) *PriceService {
	return &PriceService{
		ctx:       ctx,
		xorm:      xorm,
		PriceInfo: map[string]PriceInfo{},
	}
}

// FetchProviderPrice 根据agentId、providerName、modelName获取价格信息
// 先从本地info查找，找不到再去查询数据库，然后加入本地info
func (m *PriceService) FetchProviderPrice(modelId string) (PriceInfo, bool) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 先从本地缓存查找
	// if priceInfo, ok := m.PriceInfo[modelId]; ok {
	// 	return priceInfo, true
	// }

	// 从数据库中查询
	type ModelWithProvider struct {
		ProviderName string `xorm:"provider_name"`
	}

	var result ModelsInfo
	has, err := m.xorm.Where("model_id = ?", modelId).Get(&result)

	if err != nil || !has {
		logrus.Errorf("failed to fetch price info for model_id %s, error: %v", modelId, err)
		return PriceInfo{}, false
	}

	// // 将查询结果加入本地缓存
	priceInfo := PriceInfo{
		InputPrice:  result.InputPrice,
		OutputPrice: result.OutputPrice,
		CachePrice:  result.CachePrice,
	}
	// m.PriceInfo[modelId] = priceInfo

	return priceInfo, true
}
