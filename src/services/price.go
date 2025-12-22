package services

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/deepissue/fee_server/models"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

type PriceInfo struct {
	InputPrice  int `json:"input_price"`  //输入token计费
	OutputPrice int `json:"output_price"` //输出token计费
	CachePrice  int `json:"cache_price"`  //缓存token计费

	CostInputPrice  int `json:"cost_input_price"`  //输入token计费cost
	CostOutputPrice int `json:"cost_output_price"` //输出token计费cost
	CostCachePrice  int `json:"cose_cache_price"`  //缓存token计费cost
}

type TieredPriceInfo struct {
	InputPrice  int `json:"input_price"`
	OutputPrice int `json:"output_price"`
	CachePrice  int `json:"cache_price"`
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
func (m *PriceService) FetchProviderPrice(modelId int, providerId string) (PriceInfo, bool) {

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

	var result models.ModelsInfo
	has, err := m.xorm.Where("id = ?", modelId).Get(&result)

	if err != nil {
		logrus.Errorf("failed to fetch price info for model_id %d, error: %v", modelId, err)
		return PriceInfo{}, false
	}
	if !has {
		logrus.Errorf("failed to fetch price info for model_id %d not found", modelId)
		return PriceInfo{}, false
	}
	var provider models.ModelsProvider
	id, _ := strconv.Atoi(providerId)
	m.xorm.Where("id = ?", id).Get(&provider)

	// // 将查询结果加入本地缓存
	priceInfo := PriceInfo{
		InputPrice:      result.InputPrice,
		OutputPrice:     result.OutputPrice,
		CachePrice:      result.CachePrice,
		CostInputPrice:  provider.InputPrice,
		CostOutputPrice: provider.OutputPrice,
		CostCachePrice:  provider.CachePrice,
	}
	// m.PriceInfo[modelId] = priceInfo
	return priceInfo, true
}

// FetchTieredPrice 根据模型ID和总token数获取阶梯价格
func (m *PriceService) FetchTieredPrice(modelId int, totalTokens int64) (*TieredPriceInfo, bool) {
	var tier models.ModelsTieredPricing
	has, err := m.xorm.Where("model_id = ? AND tier_start <= ? AND (tier_end >= ? OR tier_end = -1)", modelId, totalTokens, totalTokens).
		OrderBy("tier_start DESC").
		Get(&tier)
	if err != nil {
		logrus.Errorf("failed to fetch tiered price for model_id %d, tokens %d: %v", modelId, totalTokens, err)
		return nil, false
	}
	if !has {
		return nil, false
	}
	return &TieredPriceInfo{
		InputPrice:  tier.InputPrice,
		OutputPrice: tier.OutputPrice,
		CachePrice:  tier.CachePrice,
	}, true
}

// FetchUserDiscount 获取用户折扣率，返回折扣率（100表示无折扣）
func (m *PriceService) FetchUserDiscount(userId int64) int {
	var discount models.UserDiscount
	has, err := m.xorm.Where("user_id = ?", userId).Get(&discount)
	if err != nil {
		logrus.Errorf("failed to fetch user discount for user_id %d: %v", userId, err)
		return 100
	}
	if !has {
		return 100
	}
	return discount.DiscountRate
}
