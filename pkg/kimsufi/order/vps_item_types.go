package order

import "time"

// VPSItemInfo represents VPS item information for ordering
type VPSItemInfo struct {
	PlanCode    string             `json:"planCode"`
	ProductName string             `json:"productName"`
	ProductType string             `json:"productType"`
	Prices      []VPSItemInfoPrice `json:"prices"`
}

// VPSItemInfos is a collection of VPS item information
type VPSItemInfos []VPSItemInfo

// VPSItemInfoPrice represents pricing information for VPS items
type VPSItemInfoPrice struct {
	Capacities      []string            `json:"capacities"`
	Description     string              `json:"description"`
	Duration        string              `json:"duration"`
	Interval        int                 `json:"interval"`
	MaximumQuantity int                 `json:"maximumQuantity"`
	MaximumRepeat   *int                `json:"maximumRepeat"`
	MinimumQuantity int                 `json:"minimumQuantity"`
	MinimumRepeat   int                 `json:"minimumRepeat"`
	Price           VPSPrice            `json:"price"`
	PriceInUcents   int64               `json:"priceInUcents"`
	PricingMode     string              `json:"pricingMode"`
	PricingType     string              `json:"pricingType"`
	Promotions      []VPSPricePromotion `json:"promotions,omitempty"`
}

// VPSPrice represents price information
type VPSPrice struct {
	CurrencyCode string  `json:"currencyCode"`
	Text         string  `json:"text"`
	Value        float64 `json:"value"`
}

// VPSPricePromotion represents promotional pricing
type VPSPricePromotion struct {
	Description string     `json:"description"`
	Discount    VPSPrice   `json:"discount"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Name        string     `json:"name"`
	StartDate   *time.Time `json:"startDate,omitempty"`
}

// VPSItemOption represents options available for VPS items
type VPSItemOption struct {
	Exclusive   bool               `json:"exclusive"`
	Family      string             `json:"family"`
	Mandatory   bool               `json:"mandatory"`
	PlanCode    string             `json:"planCode"`
	Prices      []VPSItemInfoPrice `json:"prices"`
	ProductName string             `json:"productName"`
	ProductType string             `json:"productType"`
}

// VPSItemOptions is a collection of VPS item options
type VPSItemOptions []VPSItemOption

// VPSItemPriceConfig defines pricing configuration for VPS items
type VPSItemPriceConfig struct {
	Duration    string
	PricingMode string
}

// VPSItemRequest represents a request to add VPS item to cart
type VPSItemRequest struct {
	Duration    string `json:"duration"`
	PlanCode    string `json:"planCode"`
	PricingMode string `json:"pricingMode"`
	Quantity    int    `json:"quantity"`
}

// VPSItemResponse represents the response after adding VPS item to cart
type VPSItemResponse struct {
	CartID         string             `json:"cartId"`
	Configurations []interface{}      `json:"configurations"`
	ItemID         int64              `json:"itemId"`
	Prices         []VPSItemInfoPrice `json:"prices"`
	ProductName    string             `json:"productName"`
	Settings       VPSItemSettings    `json:"settings"`
}

// VPSItemSettings contains VPS item configuration settings
type VPSItemSettings struct {
	Duration    string `json:"duration"`
	PlanCode    string `json:"planCode"`
	PricingMode string `json:"pricingMode"`
	Quantity    int    `json:"quantity"`
}
