package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ImagePricing struct {
	Version     string  `yaml:"version"`
	LastUpdated string  `yaml:"last_updated"`
	Models      []Model `yaml:"models"`
}

type Model struct {
	Model     string    `yaml:"model"`
	Enabled   bool      `yaml:"enabled"`
	Category  string    `yaml:"category"`
	Qualities []Quality `yaml:"qualities"`
}

type Quality struct {
	Quality string `yaml:"quality"`
	Sizes   []Size `yaml:"sizes"`
}

type Size struct {
	Size  string  `yaml:"size"`
	Price float64 `yaml:"price"`
	Cost  float64 `yaml:"cost"`
}

func LoadImagePricing(filePath string) (*ImagePricing, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config ImagePricing
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (ip *ImagePricing) GetPrice(model, quality, size string) float64 {
	for _, m := range ip.Models {
		if m.Model == model && m.Enabled {
			for _, q := range m.Qualities {
				if q.Quality == quality {
					for _, s := range q.Sizes {
						if s.Size == size {
							return s.Price
						}
					}
				}
			}
		}
	}
	return 0
}

func (ip *ImagePricing) GetCost(model, quality, size string) float64 {
	for _, m := range ip.Models {
		if m.Model == model && m.Enabled {
			for _, q := range m.Qualities {
				if q.Quality == quality {
					for _, s := range q.Sizes {
						if s.Size == size {
							return s.Cost
						}
					}
				}
			}
		}
	}
	return 0
}

func (ip *ImagePricing) GetPriceAndCost(model, quality, size string) (float64, float64) {
	for _, m := range ip.Models {
		if m.Model == model && m.Enabled {
			for _, q := range m.Qualities {
				if q.Quality == quality {
					for _, s := range q.Sizes {
						if s.Size == size {
							return s.Price, s.Cost
						}
					}
				}
			}
		}
	}
	return 0, 0
}

func (ip *ImagePricing) IsModelEnabled(model string) bool {
	for _, m := range ip.Models {
		if m.Model == model {
			return m.Enabled
		}
	}
	return false
}
