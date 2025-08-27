package order

import "slices"

// GetByPlanCode returns VPS item info for the given plan code
func (infos VPSItemInfos) GetByPlanCode(planCode string) *VPSItemInfo {
	for _, info := range infos {
		if info.PlanCode == planCode {
			return &info
		}
	}
	return nil
}

// GetPriceByConfig returns the price for the given configuration
func (info *VPSItemInfo) GetPriceByConfig(config VPSItemPriceConfig) *VPSItemInfoPrice {
	for _, price := range info.Prices {
		if price.Duration == config.Duration && price.PricingMode == config.PricingMode {
			return &price
		}
	}
	return nil
}

// GetPriceConfigOrDefault returns the given config if valid, otherwise returns default
func (infos VPSItemInfos) GetPriceConfigOrDefault(planCode string, config VPSItemPriceConfig) VPSItemPriceConfig {
	info := infos.GetByPlanCode(planCode)
	if info == nil {
		return config
	}

	price := info.GetPriceByConfig(config)
	if price != nil {
		return config
	}

	// Return first available price config as default
	if len(info.Prices) > 0 {
		return VPSItemPriceConfig{
			Duration:    info.Prices[0].Duration,
			PricingMode: info.Prices[0].PricingMode,
		}
	}

	return config
}

// GetMandatoryOptions returns mandatory options, optionally filtered
func (options VPSItemOptions) GetMandatoryOptions(filter func(VPSItemOptions, VPSItemOption) bool) VPSItemOptions {
	var mandatory VPSItemOptions
	for _, option := range options {
		if !option.Mandatory {
			continue
		}
		if filter != nil && !filter(mandatory, option) {
			continue
		}
		mandatory = append(mandatory, option)
	}
	return mandatory
}

// Get returns the first option with the given family
func (options VPSItemOptions) Get(family string) *VPSItemOption {
	for _, option := range options {
		if option.Family == family {
			return &option
		}
	}
	return nil
}

// Families returns unique families from the options
func (options VPSItemOptions) Families() []string {
	var families []string
	for _, option := range options {
		if !slices.Contains(families, option.Family) {
			families = append(families, option.Family)
		}
	}
	return families
}

// GetPriceByConfig returns the price for the given configuration
func (option *VPSItemOption) GetPriceByConfig(config VPSItemPriceConfig) *VPSItemInfoPrice {
	for _, price := range option.Prices {
		if price.Duration == config.Duration && price.PricingMode == config.PricingMode {
			return &price
		}
	}
	return nil
}

// ToOptions converts VPS item options to generic options format
func (options VPSItemOptions) ToOptions() Options {
	var result Options
	for _, option := range options {
		result = append(result, Option{
			Family:   option.Family,
			PlanCode: option.PlanCode,
		})
	}
	return result
}
