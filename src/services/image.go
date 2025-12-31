package services

import (
	"fmt"

	"github.com/deepissue/fee_server/config"
)

// loadImagePricing 加载图片定价配置
func (m *FeeService) loadImagePricing() (*config.ImagePricing, error) {
	return config.LoadImagePricing("config/image_pricing.yml")
}

// loadVideoPricing 加载视频定价配置
func (m *FeeService) loadVideoPricing() (*config.VideoPricing, error) {
	return config.LoadVideoPricing("config/video_pricing.yml")
}

// calculateImageActualCost 计算图片生成的实际成本（用于成本统计）
func (m *FeeService) calculateImageActualCost(inst FeeInstance) (float64, float64, error) {
	totalPrice := 0.0
	totalCost := 0.0
	usage, err := inst.ImageUsage()
	if err != nil {
		return 0, 0, err
	}
	// 如果有图片定价配置，使用cost字段
	for _, imageUsage := range usage {
		price, cost, err := m.getImagePriceAndCost(inst.data.Model, imageUsage.Quality, imageUsage.Size)
		if err != nil {
			return 0, 0, err
		}
		totalCost += cost
		totalPrice += price
	}

	return totalPrice, totalCost, nil
}

// getImagePriceAndCost 同时获取图片价格和成本
func (m *FeeService) getImagePriceAndCost(model, quality, size string) (float64, float64, error) {
	// 加载图片定价配置
	imagePricing, err := m.loadImagePricing()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to load image pricing: %v", err)
	}

	price, cost := imagePricing.GetPriceAndCost(model, quality, size)

	if price == 0 {
		price, cost = imagePricing.GetPriceAndCost(model, "default", size)
	}
	if price == 0 {
		price, cost = imagePricing.GetPriceAndCost(model, "default", "default")
	}
	return price, cost, nil
}
