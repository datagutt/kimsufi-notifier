package catalog

// GetVPSPlan returns the VPS plan with the given plan code.
func (c VPSCatalog) GetVPSPlan(planCode string) *VPSPlan {
	for _, plan := range c.Plans {
		if plan.PlanCode == planCode {
			return &plan
		}
	}
	return nil
}

// GetVPSProduct returns the VPS product with the given product name.
func (c VPSCatalog) GetVPSProduct(productName string) *VPSProduct {
	for _, product := range c.Products {
		if product.Name == productName {
			return &product
		}
	}
	return nil
}

// GetCategory returns the category/family of the VPS plan.
func (p VPSPlan) GetCategory() string {
	if p.Family != "" {
		return p.Family
	}
	if p.Blobs.Commercial != nil && p.Blobs.Commercial.Range != "" {
		return p.Blobs.Commercial.Range
	}
	return "vps"
}

// GetFirstPrice returns the first suitable pricing entry for the VPS plan.
// It follows similar logic to the regular servers, looking for monthly rental prices
// and avoiding installation fees.
func (p VPSPlan) GetFirstPrice() VPSPricing {
	if len(p.Pricings) == 0 {
		return VPSPricing{}
	}

	// Look for a monthly rental price, avoiding installation fees
	for _, pricing := range p.Pricings {
		// Look for monthly rentals with reasonable intervals
		if pricing.IntervalUnit == "month" &&
			pricing.Interval >= 1 &&
			pricing.Type == "rental" &&
			pricing.Price > 0 {
			// Avoid installation fees (usually have "installation" capacity)
			hasInstallationOnly := len(pricing.Capacities) == 1 &&
				contains(pricing.Capacities, "installation")
			if !hasInstallationOnly {
				return pricing
			}
		}
	}

	// Fallback: find any non-zero price that's not installation-only
	for _, pricing := range p.Pricings {
		if pricing.Price > 0 {
			hasInstallationOnly := len(pricing.Capacities) == 1 &&
				contains(pricing.Capacities, "installation")
			if !hasInstallationOnly {
				return pricing
			}
		}
	}

	// Last resort: return the first pricing
	return p.Pricings[0]
}

// GetPrice returns the price as a float64, using the same divider as regular servers.
// OVH stores prices as integers multiplied by 100,000,000 (8 decimal places).
func (vp VPSPricing) GetPrice() float64 {
	return float64(vp.Price) / 100000000.0
}

// contains is a helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// HasDatacenterInfo checks if the VPS plan has datacenter information in its technical specs.
func (p VPSPlan) HasDatacenterInfo() bool {
	return p.Blobs.Technical != nil && p.Blobs.Technical.Datacenter != nil
}

// GetDatacenterInfo returns datacenter information if available.
func (p VPSPlan) GetDatacenterInfo() *VPSDatacenter {
	if p.HasDatacenterInfo() {
		return p.Blobs.Technical.Datacenter
	}
	return nil
}

// GetCPUInfo returns CPU specification information if available.
func (p VPSPlan) GetCPUInfo() *VPSCPUSpec {
	if p.Blobs.Technical != nil && p.Blobs.Technical.CPU != nil {
		return p.Blobs.Technical.CPU
	}
	return nil
}

// GetMemoryInfo returns memory specification information if available.
func (p VPSPlan) GetMemoryInfo() *VPSMemorySpec {
	if p.Blobs.Technical != nil && p.Blobs.Technical.Memory != nil {
		return p.Blobs.Technical.Memory
	}
	return nil
}

// GetStorageInfo returns storage specification information if available.
func (p VPSPlan) GetStorageInfo() *VPSStorageSpec {
	if p.Blobs.Technical != nil && p.Blobs.Technical.Storage != nil {
		return p.Blobs.Technical.Storage
	}
	return nil
}
