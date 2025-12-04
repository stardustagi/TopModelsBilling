package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type VideoPricing struct {
	Version     string       `yaml:"version"`
	LastUpdated string       `yaml:"last_updated"`
	Models      []VideoModel `yaml:"models"`
}

type VideoModel struct {
	Model    string      `yaml:"model"`
	Enabled  bool        `yaml:"enabled"`
	Category string      `yaml:"category"`
	Sizes    []VideoSize `yaml:"sizes"`
}

type VideoSize struct {
	Size      string          `yaml:"size"`
	Durations []VideoDuration `yaml:"durations"`
}

type VideoDuration struct {
	Duration float64 `yaml:"duration"` // 固定时长，如5、15、30、60秒，0表示按秒计费
	Price    float64 `yaml:"price"`    // 这个时长的总价格(如果duration=0则是每秒价格)
	Cost     float64 `yaml:"cost"`     // 这个时长的总成本(如果duration=0则是每秒成本)
}

func LoadVideoPricing(filePath string) (*VideoPricing, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config VideoPricing
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (vp *VideoPricing) GetPrice(model, size string, seconds float64) float64 {
	for _, m := range vp.Models {
		if m.Model == model && m.Enabled {
			for _, s := range m.Sizes {
				if s.Size == size {
					// 查找精确匹配的duration
					for _, d := range s.Durations {
						if d.Duration == seconds {
							return d.Price
						}
					}
					// 没找到精确匹配，查找duration=0的default配置
					for _, d := range s.Durations {
						if d.Duration == 0 {
							return d.Price * seconds // 按秒计费
						}
					}
				}
			}
		}
	}
	return 0
}

func (vp *VideoPricing) GetCost(model, size string, seconds float64) float64 {
	for _, m := range vp.Models {
		if m.Model == model && m.Enabled {
			for _, s := range m.Sizes {
				if s.Size == size {
					// 查找精确匹配的duration
					for _, d := range s.Durations {
						if d.Duration == seconds {
							return d.Cost
						}
					}
					// 没找到精确匹配，查找duration=0的default配置
					for _, d := range s.Durations {
						if d.Duration == 0 {
							return d.Cost * seconds // 按秒计费
						}
					}
				}
			}
		}
	}
	return 0
}

func (vp *VideoPricing) GetPriceAndCost(model, size string, seconds float64) (float64, float64) {
	for _, m := range vp.Models {
		if m.Model == model && m.Enabled {
			for _, s := range m.Sizes {
				if s.Size == size {
					// 查找精确匹配的duration
					for _, d := range s.Durations {
						if d.Duration == seconds {
							return d.Price, d.Cost
						}
					}
					// 没找到精确匹配，查找duration=0的default配置
					for _, d := range s.Durations {
						if d.Duration == 0 {
							return d.Price * seconds, d.Cost * seconds // 按秒计费
						}
					}
				}
			}
		}
	}
	return 0, 0
}

func (vp *VideoPricing) IsModelEnabled(model string) bool {
	for _, m := range vp.Models {
		if m.Model == model {
			return m.Enabled
		}
	}
	return false
}
