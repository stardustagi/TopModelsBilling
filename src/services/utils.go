package services

const MICRO = 1_000_000 // 1 代币 = 1_000_000 微代币

// usedTokens: 实际使用的 token 数
// priceCredits: 价格 coins/1M（例如 20）
// 返回值：微代币（int64）
func CalculateTokenCostMicro(usedTokens int64, priceCredits float64) int64 {
	// 每个 token 的微代币价格
	pricePerTokenMicro := (priceCredits * MICRO) / float64(MICRO)
	totalCost := float64(usedTokens) * pricePerTokenMicro
	return int64(totalCost + 0.5) // 四舍五入
}
