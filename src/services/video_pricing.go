package services

import (
	"fmt"
	"maps"
	"os"

	"gopkg.in/yaml.v3"
)

// videoPricing manages video generation pricing
type videoPricing struct {
	pricingMap map[VideoModel]map[VideoResolution]float64
}

// VideoResolution represents the output resolution of video generation
type VideoResolution string

const (
	// Portrait resolutions
	ResolutionPortrait720x1280  VideoResolution = "720x1280"
	ResolutionPortrait1024x1792 VideoResolution = "1024x1792"

	// Landscape resolutions
	ResolutionLandscape1280x720  VideoResolution = "1280x720"
	ResolutionLandscape1792x1024 VideoResolution = "1792x1024"
)

// VideoModel represents the video generation model
type VideoModel string

const (
	ModelSora2    VideoModel = "sora-2"
	ModelSora2Pro VideoModel = "sora-2-pro"
)

// VideoPricingConfig represents the video pricing configuration structure
type VideoPricingConfig struct {
	Version     string             `yaml:"version"`
	LastUpdated string             `yaml:"last_updated"`
	Models      []VideoModelConfig `yaml:"models"`
}

// VideoModelConfig represents a model's pricing configuration
type VideoModelConfig struct {
	Model       string                  `yaml:"model"`
	Description string                  `yaml:"description"`
	Enabled     bool                    `yaml:"enabled"`
	Resolutions []VideoResolutionConfig `yaml:"resolutions"`
}

// VideoResolutionConfig represents a resolution's pricing configuration
type VideoResolutionConfig struct {
	Resolution  string  `yaml:"resolution"`
	Price       float64 `yaml:"price_per_second"`
	Description string  `yaml:"description"`
}

// getDefaultVideoPricing returns the default pricing matrix (price per second)
func getDefaultVideoPricing() map[VideoModel]map[VideoResolution]float64 {
	return map[VideoModel]map[VideoResolution]float64{
		ModelSora2: {
			ResolutionPortrait720x1280:  0.10,
			ResolutionLandscape1280x720: 0.10,
		},
		ModelSora2Pro: {
			ResolutionPortrait720x1280:   0.30,
			ResolutionLandscape1280x720:  0.30,
			ResolutionPortrait1024x1792:  0.50,
			ResolutionLandscape1792x1024: 0.50,
		},
	}
}

// NewVideoPricing creates a new videoPricing instance with default pricing
func NewVideoPricing() *videoPricing {
	return &videoPricing{
		pricingMap: getDefaultVideoPricing(),
	}
}

// NewVideoPricingWithConfig creates a new videoPricing instance and loads configuration from file
// If the config file doesn't exist or fails to load, it falls back to default pricing
func NewVideoPricingWithConfig(configPath string) (*videoPricing, error) {
	p := NewVideoPricing()

	// Try to load config file
	if err := p.LoadFromFile(configPath); err != nil {
		// If file doesn't exist, just use default pricing (not an error)
		if os.IsNotExist(err) {
			return p, nil
		}
		// Other errors should be returned
		return nil, fmt.Errorf("failed to load video pricing config: %w", err)
	}

	return p, nil
}

// LoadFromFile loads pricing configuration from a YAML file and merges with default pricing
func (p *videoPricing) LoadFromFile(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config VideoPricingConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse video pricing config: %w", err)
	}

	// Merge config into pricing map
	p.mergeConfig(config)

	return nil
}

// mergeConfig merges the configuration into the pricing map, overriding default values
func (p *videoPricing) mergeConfig(config VideoPricingConfig) {
	for _, modelConfig := range config.Models {
		// Skip disabled models
		if !modelConfig.Enabled {
			continue
		}

		model := VideoModel(modelConfig.Model)

		// Initialize model map if it doesn't exist
		if p.pricingMap[model] == nil {
			p.pricingMap[model] = make(map[VideoResolution]float64)
		}

		for _, resolutionConfig := range modelConfig.Resolutions {
			resolution := VideoResolution(resolutionConfig.Resolution)
			// Override or add the price
			p.pricingMap[model][resolution] = resolutionConfig.Price
		}
	}
}

// GetVideoPrice returns the price per second for generating video with the specified parameters
func (p *videoPricing) GetVideoPrice(model VideoModel, resolution VideoResolution) (float64, bool) {
	if modelPricing, ok := p.pricingMap[model]; ok {
		if price, ok := modelPricing[resolution]; ok {
			return price, true
		}
	}
	return 0, false
}

// CalculateVideoCost calculates the total cost for generating video based on duration in seconds
func (p *videoPricing) CalculateVideoCost(model VideoModel, resolution VideoResolution, durationSeconds float64) (float64, bool) {
	pricePerSecond, ok := p.GetVideoPrice(model, resolution)
	if !ok {
		return 0, false
	}
	return pricePerSecond * durationSeconds, true
}

// GetAllPricing returns a copy of the entire pricing map
func (p *videoPricing) GetAllPricing() map[VideoModel]map[VideoResolution]float64 {
	// Create a deep copy to prevent external modification
	result := make(map[VideoModel]map[VideoResolution]float64)

	for model, resolutions := range p.pricingMap {
		result[model] = make(map[VideoResolution]float64)
		maps.Copy(result[model], resolutions)
	}

	return result
}

// GetSupportedModels returns a list of all supported video models
func (p *videoPricing) GetSupportedModels() []VideoModel {
	models := make([]VideoModel, 0, len(p.pricingMap))
	for model := range p.pricingMap {
		models = append(models, model)
	}
	return models
}

// GetSupportedResolutions returns a list of supported resolutions for a given model
func (p *videoPricing) GetSupportedResolutions(model VideoModel) []VideoResolution {
	if modelPricing, ok := p.pricingMap[model]; ok {
		resolutions := make([]VideoResolution, 0, len(modelPricing))
		for resolution := range modelPricing {
			resolutions = append(resolutions, resolution)
		}
		return resolutions
	}
	return nil
}

var VideoPricing *videoPricing

// InitVideoDefault initializes the global VideoPricing with default pricing
func InitVideoDefault() {
	VideoPricing = NewVideoPricing()
}

// InitVideoPricing initializes the global VideoPricing with pricing from file
func InitVideoPricing(file string) {
	VideoPricing = NewVideoPricing()
	VideoPricing.LoadFromFile(file)
}
