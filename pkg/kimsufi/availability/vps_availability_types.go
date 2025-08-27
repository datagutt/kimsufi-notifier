package availability

// VPSAvailabilities represents the response from /vps/order/rule/datacenter
type VPSAvailabilities struct {
	Datacenters []VPSDatacenterAvailability `json:"datacenters"`
}

// VPSDatacenterAvailability represents availability status for a single datacenter
type VPSDatacenterAvailability struct {
	Status        string `json:"status"`
	Datacenter    string `json:"datacenter"`
	Code          string `json:"code"`
	LinuxStatus   string `json:"linuxStatus"`
	WindowsStatus string `json:"windowsStatus"`
}

// GetAvailableDatacenters returns datacenters that are available (not out-of-stock)
func (va VPSAvailabilities) GetAvailableDatacenters() []VPSDatacenterAvailability {
	var available []VPSDatacenterAvailability
	for _, dc := range va.Datacenters {
		if dc.Status == "available" || dc.LinuxStatus == "available" || dc.WindowsStatus == "available" {
			available = append(available, dc)
		}
	}
	return available
}

// GetDatacenterCodes returns all datacenter codes
func (va VPSAvailabilities) GetDatacenterCodes() []string {
	var codes []string
	for _, dc := range va.Datacenters {
		codes = append(codes, dc.Datacenter)
	}
	return codes
}

// GetAvailableDatacenterCodes returns only available datacenter codes
func (va VPSAvailabilities) GetAvailableDatacenterCodes() []string {
	var codes []string
	for _, dc := range va.GetAvailableDatacenters() {
		codes = append(codes, dc.Datacenter)
	}
	return codes
}

// HasAvailability returns true if at least one datacenter is available
func (va VPSAvailabilities) HasAvailability() bool {
	return len(va.GetAvailableDatacenters()) > 0
}

// GetStatus returns overall status - "available" if any datacenter is available, "out-of-stock" otherwise
func (va VPSAvailabilities) GetStatus() string {
	if va.HasAvailability() {
		return "available"
	}
	return "out-of-stock"
}

// IsDatacenterAvailable checks if a specific datacenter is available
func (va VPSAvailabilities) IsDatacenterAvailable(datacenter string) bool {
	for _, dc := range va.Datacenters {
		if dc.Datacenter == datacenter {
			return dc.Status == "available" || dc.LinuxStatus == "available" || dc.WindowsStatus == "available"
		}
	}
	return false
}

// GetDatacenterByCode finds a datacenter by its code (e.g., "BHS", "GRA")
func (va VPSAvailabilities) GetDatacenterByCode(code string) *VPSDatacenterAvailability {
	for _, dc := range va.Datacenters {
		if dc.Datacenter == code {
			return &dc
		}
	}
	return nil
}
