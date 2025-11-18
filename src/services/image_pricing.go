package services

import (
	"fmt"
	"maps"
	"os"

	"gopkg.in/yaml.v3"
)

// imagePricing manages image generation pricing
type imagePricing struct {
	pricingMap map[ImageModel]map[ImageQuality]map[ImageSize]float64
}

// ImageQuality represents the quality level of image generation
type ImageQuality string

const (
	QualityLow      ImageQuality = "low"
	QualityMedium   ImageQuality = "medium"
	QualityHigh     ImageQuality = "high"
	QualityStandard ImageQuality = "standard"
	QualityHD       ImageQuality = "hd"
)

// ImageSize represents the dimensions of the generated image
type ImageSize string

const (
	Size256x256   ImageSize = "256x256"
	Size512x512   ImageSize = "512x512"
	Size1024x1024 ImageSize = "1024x1024"
	Size1024x1536 ImageSize = "1024x1536"
	Size1536x1024 ImageSize = "1536x1024"
	Size1024x1792 ImageSize = "1024x1792"
	Size1792x1024 ImageSize = "1792x1024"
)

// ImageModel represents the image generation model
type ImageModel string

const (
	ModelGPTImage1     ImageModel = "gpt-image-1"
	ModelGPTImage1Mini ImageModel = "gpt-image-1-mini"
	ModelDALLE3        ImageModel = "dall-e-3"
	ModelDALLE2        ImageModel = "dall-e-2"
)

// PricingConfig represents the pricing configuration structure
type PricingConfig struct {
	Version     string        `yaml:"version"`
	LastUpdated string        `yaml:"last_updated"`
	Models      []ModelConfig `yaml:"models"`
}

// ModelConfig represents a model's pricing configuration
type ModelConfig struct {
	Model       string          `yaml:"model"`
	Description string          `yaml:"description"`
	Enabled     bool            `yaml:"enabled"`
	Qualities   []QualityConfig `yaml:"qualities"`
}

// QualityConfig represents a quality level's pricing configuration
type QualityConfig struct {
	Quality     string       `yaml:"quality"`
	Description string       `yaml:"description"`
	Sizes       []SizeConfig `yaml:"sizes"`
}

// SizeConfig represents a size's pricing configuration
type SizeConfig struct {
	Size        string  `yaml:"size"`
	Price       float64 `yaml:"price"`
	Description string  `yaml:"description"`
}

// getDefaultPricing returns the default pricing matrix
func getDefaultPricing() map[ImageModel]map[ImageQuality]map[ImageSize]float64 {
	return map[ImageModel]map[ImageQuality]map[ImageSize]float64{
		ModelGPTImage1: {
			QualityLow: {
				Size1024x1024: 0.011,
				Size1024x1536: 0.016,
				Size1536x1024: 0.016,
			},
			QualityMedium: {
				Size1024x1024: 0.042,
				Size1024x1536: 0.063,
				Size1536x1024: 0.063,
			},
			QualityHigh: {
				Size1024x1024: 0.167,
				Size1024x1536: 0.25,
				Size1536x1024: 0.25,
			},
		},
		ModelGPTImage1Mini: {
			QualityLow: {
				Size1024x1024: 0.005,
				Size1024x1536: 0.006,
				Size1536x1024: 0.006,
			},
			QualityMedium: {
				Size1024x1024: 0.011,
				Size1024x1536: 0.015,
				Size1536x1024: 0.015,
			},
			QualityHigh: {
				Size1024x1024: 0.036,
				Size1024x1536: 0.052,
				Size1536x1024: 0.052,
			},
		},
		ModelDALLE3: {
			QualityStandard: {
				Size1024x1024: 0.04,
				Size1024x1792: 0.08,
				Size1792x1024: 0.08,
			},
			QualityHD: {
				Size1024x1024: 0.08,
				Size1024x1792: 0.12,
				Size1792x1024: 0.12,
			},
		},
		ModelDALLE2: {
			QualityStandard: {
				Size256x256:   0.016,
				Size512x512:   0.018,
				Size1024x1024: 0.02,
			},
		},
	}
}

// NewImagePricing creates a new Pricing instance with default pricing
func NewImagePricing() *imagePricing {
	return &imagePricing{
		pricingMap: getDefaultPricing(),
	}
}

// NewPricingWithConfig creates a new Pricing instance and loads configuration from file
// If the config file doesn't exist or fails to load, it falls back to default pricing
func NewPricingWithConfig(configPath string) (*imagePricing, error) {
	p := NewImagePricing()

	// Try to load config file
	if err := p.LoadFromFile(configPath); err != nil {
		// If file doesn't exist, just use default pricing (not an error)
		if os.IsNotExist(err) {
			return p, nil
		}
		// Other errors should be returned
		return nil, fmt.Errorf("failed to load pricing config: %w", err)
	}

	return p, nil
}

// LoadFromFile loads pricing configuration from a YAML file and merges with default pricing
func (p *imagePricing) LoadFromFile(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config PricingConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse pricing config: %w", err)
	}

	// Merge config into pricing map
	p.mergeConfig(config)

	return nil
}

// mergeConfig merges the configuration into the pricing map, overriding default values
func (p *imagePricing) mergeConfig(config PricingConfig) {
	for _, modelConfig := range config.Models {
		// Skip disabled models
		if !modelConfig.Enabled {
			continue
		}

		model := ImageModel(modelConfig.Model)

		// Initialize model map if it doesn't exist
		if p.pricingMap[model] == nil {
			p.pricingMap[model] = make(map[ImageQuality]map[ImageSize]float64)
		}

		for _, qualityConfig := range modelConfig.Qualities {
			quality := ImageQuality(qualityConfig.Quality)

			// Initialize quality map if it doesn't exist
			if p.pricingMap[model][quality] == nil {
				p.pricingMap[model][quality] = make(map[ImageSize]float64)
			}

			for _, sizeConfig := range qualityConfig.Sizes {
				size := ImageSize(sizeConfig.Size)
				// Override or add the price
				p.pricingMap[model][quality][size] = sizeConfig.Price
			}
		}
	}
}

// GetImagePrice returns the price for generating an image with the specified parameters
func (p *imagePricing) GetImagePrice(model ImageModel, quality ImageQuality, size ImageSize) (float64, bool) {
	if modelPricing, ok := p.pricingMap[model]; ok {
		if qualityPricing, ok := modelPricing[quality]; ok {
			if price, ok := qualityPricing[size]; ok {
				return price, true
			}
		}
	}
	return 0, false
}

// CalculateImageCost calculates the total cost for generating multiple images
func (p *imagePricing) CalculateImageCost(model ImageModel, quality ImageQuality, size ImageSize, count int) (float64, bool) {
	price, ok := p.GetImagePrice(model, quality, size)
	if !ok {
		return 0, false
	}
	return price * float64(count), true
}

// GetAllPricing returns a copy of the entire pricing map
func (p *imagePricing) GetAllPricing() map[ImageModel]map[ImageQuality]map[ImageSize]float64 {
	// Create a deep copy to prevent external modification
	result := make(map[ImageModel]map[ImageQuality]map[ImageSize]float64)

	for model, qualities := range p.pricingMap {
		result[model] = make(map[ImageQuality]map[ImageSize]float64)
		for quality, sizes := range qualities {
			result[model][quality] = make(map[ImageSize]float64)
			maps.Copy(result[model][quality], sizes)
		}
	}

	return result
}

// GetSupportedModels returns a list of all supported models
func (p *imagePricing) GetSupportedModels() []ImageModel {
	models := make([]ImageModel, 0, len(p.pricingMap))
	for model := range p.pricingMap {
		models = append(models, model)
	}
	return models
}

// GetSupportedQualities returns a list of supported qualities for a given model
func (p *imagePricing) GetSupportedQualities(model ImageModel) []ImageQuality {
	if modelPricing, ok := p.pricingMap[model]; ok {
		qualities := make([]ImageQuality, 0, len(modelPricing))
		for quality := range modelPricing {
			qualities = append(qualities, quality)
		}
		return qualities
	}
	return nil
}

// GetSupportedSizes returns a list of supported sizes for a given model and quality
func (p *imagePricing) GetSupportedSizes(model ImageModel, quality ImageQuality) []ImageSize {
	if modelPricing, ok := p.pricingMap[model]; ok {
		if qualityPricing, ok := modelPricing[quality]; ok {
			sizes := make([]ImageSize, 0, len(qualityPricing))
			for size := range qualityPricing {
				sizes = append(sizes, size)
			}
			return sizes
		}
	}
	return nil
}

var ImagePricing *imagePricing

func InitImageDefault() {
	ImagePricing = NewImagePricing()
}

func InitImagePricing(file string) {
	ImagePricing = NewImagePricing()
	ImagePricing.LoadFromFile(file)
}
