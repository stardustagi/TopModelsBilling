package services

import "fmt"

// calculateVideoCost 计算视频生成成本
// calculateImageActualCost 计算图片生成的实际成本（用于成本统计）
func (m *FeeService) calculateVideoActualCost(inst FeeInstance) (float64, float64, error) {
	totalPrice := 0.0
	totalCost := 0.0
	usage, err := inst.VideoUsage()
	if err != nil {
		return 0, 0, err
	}
	// 如果有图片定价配置，使用cost字段
	price, cost, err := m.getVideoPriceAndCost(inst.data.Model, usage.Size, int(usage.Seconds))
	if err != nil {
		return 0, 0, err
	}
	totalCost += cost
	totalPrice += price

	return totalPrice, totalCost, nil
}

// getImagePriceAndCost 同时获取图片价格和成本
func (m *FeeService) getVideoPriceAndCost(model, size string, seconds int) (float64, float64, error) {
	// 加载图片定价配置
	videoPricing, err := m.loadVideoPricing()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to load image pricing: %v", err)
	}

	price, cost := videoPricing.GetPriceAndCost(model, size, seconds)

	if price == 0 {
		price, cost = videoPricing.GetPriceAndCost(model, "default", seconds)
	}
	if price == 0 {
		return 0, 0, fmt.Errorf("No pricing found for model %s, size %s, seconds %d", model, size, seconds)
	}

	return price, cost, nil
}
