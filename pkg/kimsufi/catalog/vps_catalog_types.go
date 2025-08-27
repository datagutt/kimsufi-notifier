package catalog

import "time"

// VPSCatalog represents the OVH VPS catalog.
// Definition can be found at https://eu.api.ovh.com/console/?section=%2Forder&branch=v1#get-/order/catalog/public/vps
type VPSCatalog struct {
	Addons       []VPSAddon      `json:"addons"`
	CatalogID    int             `json:"catalogId"`
	Locale       Locale          `json:"locale"`
	PlanFamilies []VPSPlanFamily `json:"planFamilies"`
	Plans        []VPSPlan       `json:"plans"`
	Products     []VPSProduct    `json:"products"`
}

type VPSAddon struct {
	AddonFamilies            []VPSAddonFamily      `json:"addonFamilies"`
	Blobs                    VPSProductBlob        `json:"blobs"`
	Configurations           []VPSConfiguration    `json:"configurations"`
	ConsumptionConfiguration *VPSConsumptionConfig `json:"consumptionConfiguration,omitempty"`
	Family                   string                `json:"family"`
	InvoiceName              string                `json:"invoiceName"`
	PlanCode                 string                `json:"planCode"`
	PricingType              string                `json:"pricingType"`
	Pricings                 []VPSPricing          `json:"pricings"`
	Product                  string                `json:"product"`
}

type VPSAddonFamily struct {
	Addons    []string `json:"addons"`
	Default   string   `json:"default"`
	Exclusive bool     `json:"exclusive"`
	Mandatory bool     `json:"mandatory"`
	Name      string   `json:"name"`
}

type VPSPlan struct {
	AddonFamilies            []VPSAddonFamily      `json:"addonFamilies"`
	Blobs                    VPSProductBlob        `json:"blobs"`
	Configurations           []VPSConfiguration    `json:"configurations"`
	ConsumptionConfiguration *VPSConsumptionConfig `json:"consumptionConfiguration,omitempty"`
	Family                   string                `json:"family"`
	InvoiceName              string                `json:"invoiceName"`
	PlanCode                 string                `json:"planCode"`
	PricingType              string                `json:"pricingType"`
	Pricings                 []VPSPricing          `json:"pricings"`
	Product                  string                `json:"product"`
}

type VPSPlanFamily struct {
	Name string `json:"name"`
}

type VPSProduct struct {
	Blobs          VPSProductBlob     `json:"blobs"`
	Configurations []VPSConfiguration `json:"configurations"`
	Description    string             `json:"description"`
	Name           string             `json:"name"`
}

type VPSConsumptionConfig struct {
	BillingStrategy string `json:"billingStrategy"`
	PingEndPolicy   string `json:"pingEndPolicy"`
	ProrataUnit     string `json:"prorataUnit"`
}

type VPSConfiguration struct {
	IsCustom    bool     `json:"isCustom"`
	IsMandatory bool     `json:"isMandatory"`
	Name        string   `json:"name"`
	Values      []string `json:"values"`
}

type VPSPricing struct {
	Capacities              []string                    `json:"capacities"`
	Commitment              int                         `json:"commitment"`
	Description             string                      `json:"description"`
	EngagementConfiguration *VPSEngagementConfiguration `json:"engagementConfiguration,omitempty"`
	Interval                int                         `json:"interval"`
	IntervalUnit            string                      `json:"intervalUnit"`
	Mode                    string                      `json:"mode"`
	MustBeCompleted         bool                        `json:"mustBeCompleted"`
	Phase                   int                         `json:"phase"`
	Price                   int                         `json:"price"`
	Promotions              []VPSPromotion              `json:"promotions,omitempty"`
	Quantity                VPSQuantityRange            `json:"quantity"`
	Repeat                  VPSQuantityRange            `json:"repeat"`
	Strategy                string                      `json:"strategy"`
	Tax                     int                         `json:"tax"`
	Type                    string                      `json:"type"`
}

type VPSEngagementConfiguration struct {
	DefaultEndAction string `json:"defaultEndAction"`
	Duration         string `json:"duration"`
	Type             string `json:"type"`
}

type VPSPromotion struct {
	Context                 string     `json:"context"`
	Description             string     `json:"description"`
	Discount                VPSPrice   `json:"discount"`
	Duration                int        `json:"duration"`
	EndDate                 *time.Time `json:"endDate,omitempty"`
	GlobalQuantity          int        `json:"globalQuantity"`
	IsGlobalQuantityLimited bool       `json:"isGlobalQuantityLimited"`
	MinimumDuration         int        `json:"minimumDuration"`
	Name                    string     `json:"name"`
	Quantity                int        `json:"quantity"`
	StartDate               *time.Time `json:"startDate,omitempty"`
	Tags                    []string   `json:"tags,omitempty"`
	Total                   VPSPrice   `json:"total"`
	Type                    string     `json:"type"`
	Value                   int        `json:"value"`
}

type VPSPrice struct {
	Tax   int `json:"tax"`
	Value int `json:"value"`
}

type VPSQuantityRange struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

type VPSProductBlob struct {
	Commercial *VPSCommercialBlob `json:"commercial,omitempty"`
	Marketing  *VPSMarketingBlob  `json:"marketing,omitempty"`
	Meta       *VPSMetaBlob       `json:"meta,omitempty"`
	Tags       []string           `json:"tags,omitempty"`
	Technical  *VPSTechnicalBlob  `json:"technical,omitempty"`
	Value      string             `json:"value,omitempty"`
}

type VPSCommercialBlob struct {
	Brick        string           `json:"brick,omitempty"`
	BrickSubtype string           `json:"brickSubtype,omitempty"`
	Connection   *VPSConnection   `json:"connection,omitempty"`
	Features     []VPSFeature     `json:"features,omitempty"`
	Line         string           `json:"line,omitempty"`
	Name         string           `json:"name,omitempty"`
	Price        *VPSDisplayPrice `json:"price,omitempty"`
	Range        string           `json:"range,omitempty"`
}

type VPSMarketingBlob struct {
	Content []VPSMarketingContent `json:"content,omitempty"`
}

type VPSMarketingContent struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VPSMetaBlob struct {
	Configurations []VPSMetaConfiguration `json:"configurations,omitempty"`
}

type VPSMetaConfiguration struct {
	Name   string                      `json:"name"`
	Values []VPSMetaConfigurationValue `json:"values"`
}

type VPSMetaConfigurationValue struct {
	Blobs VPSProductBlob `json:"blobs"`
	Value string         `json:"value"`
}

type VPSTechnicalBlob struct {
	Bandwidth             *VPSNetworkSpec        `json:"bandwidth,omitempty"`
	Connection            *VPSConnection         `json:"connection,omitempty"`
	ConnectionPerSeconds  *VPSConnectionRate     `json:"connectionPerSeconds,omitempty"`
	CPU                   *VPSCPUSpec            `json:"cpu,omitempty"`
	Datacenter            *VPSDatacenter         `json:"datacenter,omitempty"`
	EphemeralLocalStorage *VPSStorageSpec        `json:"ephemeralLocalStorage,omitempty"`
	GPU                   *VPSGPUSpec            `json:"gpu,omitempty"`
	License               *VPSLicenseSpec        `json:"license,omitempty"`
	Memory                *VPSMemorySpec         `json:"memory,omitempty"`
	Name                  string                 `json:"name,omitempty"`
	Nodes                 *VPSNodeSpec           `json:"nodes,omitempty"`
	NVME                  *VPSStorageSpec        `json:"nvme,omitempty"`
	OS                    *VPSOSSpec             `json:"os,omitempty"`
	Provider              *VPSProviderSpec       `json:"provider,omitempty"`
	RequestPerSeconds     *VPSRequestRate        `json:"requestPerSeconds,omitempty"`
	Server                *VPSServerSpec         `json:"server,omitempty"`
	Storage               *VPSStorageSpec        `json:"storage,omitempty"`
	Throughput            *VPSThroughputSpec     `json:"throughput,omitempty"`
	Virtualization        *VPSVirtualizationSpec `json:"virtualization,omitempty"`
	Volume                *VPSVolumeSpec         `json:"volume,omitempty"`
	VRack                 *VPSNetworkSpec        `json:"vrack,omitempty"`
}

type VPSConnection struct {
	Clients *VPSConnectionClients `json:"clients,omitempty"`
	Total   int                   `json:"total"`
}

type VPSConnectionClients struct {
	Concurrency int `json:"concurrency"`
	Number      int `json:"number"`
}

type VPSConnectionRate struct {
	Total int    `json:"total"`
	Unit  string `json:"unit"`
}

type VPSFeature struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type VPSDisplayPrice struct {
	Display   *VPSDisplayValue `json:"display,omitempty"`
	Interval  string           `json:"interval,omitempty"`
	Precision int              `json:"precision"`
	Unit      string           `json:"unit,omitempty"`
}

type VPSDisplayValue struct {
	Value string `json:"value"`
}

type VPSNetworkSpec struct {
	Burst      int    `json:"burst"`
	Capacity   int    `json:"capacity"`
	Guaranteed bool   `json:"guaranteed"`
	Interfaces int    `json:"interfaces"`
	IsMax      bool   `json:"isMax"`
	Level      int    `json:"level"`
	Limit      int    `json:"limit"`
	Max        int    `json:"max"`
	MaxUnit    string `json:"maxUnit,omitempty"`
	Shared     bool   `json:"shared"`
	Traffic    int    `json:"traffic"`
	Unit       string `json:"unit,omitempty"`
	Unlimited  bool   `json:"unlimited"`
}

type VPSCPUSpec struct {
	Boost        int    `json:"boost"`
	Brand        string `json:"brand,omitempty"`
	Cores        int    `json:"cores"`
	Customizable bool   `json:"customizable"`
	Frequency    int    `json:"frequency"`
	MaxFrequency int    `json:"maxFrequency"`
	Model        string `json:"model,omitempty"`
	Number       int    `json:"number"`
	Score        int    `json:"score"`
	Threads      int    `json:"threads"`
	Type         string `json:"type,omitempty"`
}

type VPSDatacenter struct {
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	Name        string `json:"name,omitempty"`
	Region      string `json:"region,omitempty"`
}

type VPSStorageSpec struct {
	Disks       []VPSDiskSpec   `json:"disks,omitempty"`
	HotSwap     bool            `json:"hotSwap,omitempty"`
	Raid        string          `json:"raid,omitempty"`
	RaidDetails *VPSRaidDetails `json:"raidDetails,omitempty"`
}

type VPSDiskSpec struct {
	Capacity        int    `json:"capacity"`
	Interface       string `json:"interface,omitempty"`
	IOPS            int    `json:"iops"`
	MaximumCapacity int    `json:"maximumCapacity"`
	Number          int    `json:"number"`
	SizeUnit        string `json:"sizeUnit,omitempty"`
	Specs           string `json:"specs,omitempty"`
	Technology      string `json:"technology,omitempty"`
	Usage           string `json:"usage,omitempty"`
}

type VPSRaidDetails struct {
	CardModel string `json:"cardModel,omitempty"`
	CardSize  string `json:"cardSize,omitempty"`
	Type      string `json:"type,omitempty"`
}

type VPSGPUSpec struct {
	Brand       string         `json:"brand,omitempty"`
	Memory      *VPSMemorySpec `json:"memory,omitempty"`
	Model       string         `json:"model,omitempty"`
	Number      int            `json:"number"`
	Performance int            `json:"performance"`
}

type VPSLicenseSpec struct {
	Application  string       `json:"application,omitempty"`
	Cores        *VPSCoreSpec `json:"cores,omitempty"`
	CPU          *VPSCPUSpec  `json:"cpu,omitempty"`
	Distribution string       `json:"distribution,omitempty"`
	Edition      string       `json:"edition,omitempty"`
	Family       string       `json:"family,omitempty"`
	Feature      string       `json:"feature,omitempty"`
	Flavor       string       `json:"flavor,omitempty"`
	Images       []string     `json:"images,omitempty"`
	NbOfAccount  int          `json:"nbOfAccount"`
	Package      string       `json:"package,omitempty"`
	Version      string       `json:"version,omitempty"`
}

type VPSCoreSpec struct {
	Number int `json:"number"`
	Total  int `json:"total"`
}

type VPSMemorySpec struct {
	Customizable bool   `json:"customizable"`
	ECC          bool   `json:"ecc"`
	Frequency    int    `json:"frequency"`
	Interface    string `json:"interface,omitempty"`
	RamType      string `json:"ramType,omitempty"`
	Size         int    `json:"size"`
	SizeUnit     string `json:"sizeUnit,omitempty"`
}

type VPSNodeSpec struct {
	Number int `json:"number"`
}

type VPSOSSpec struct {
	Distribution string `json:"distribution,omitempty"`
	Edition      string `json:"edition,omitempty"`
	Family       string `json:"family,omitempty"`
	Version      string `json:"version,omitempty"`
}

type VPSProviderSpec struct {
	PointsOfPresence int  `json:"pointsOfPresence"`
	Reference        bool `json:"reference"`
}

type VPSRequestRate struct {
	Total int    `json:"total"`
	Unit  string `json:"unit,omitempty"`
}

type VPSServerSpec struct {
	CPU      *VPSCPUSpec      `json:"cpu,omitempty"`
	Frame    *VPSFrameSpec    `json:"frame,omitempty"`
	Network  *VPSNetworkSpec  `json:"network,omitempty"`
	Range    string           `json:"range,omitempty"`
	Services *VPSServicesSpec `json:"services,omitempty"`
}

type VPSFrameSpec struct {
	DualPowerSupply bool   `json:"dualPowerSupply"`
	Model           string `json:"model,omitempty"`
	Size            string `json:"size,omitempty"`
}

type VPSServicesSpec struct {
	AntiDDoS       string `json:"antiddos,omitempty"`
	IncludedBackup int    `json:"includedBackup"`
	SLA            int    `json:"sla"`
}

type VPSThroughputSpec struct {
	Level int `json:"level"`
}

type VPSVirtualizationSpec struct {
	Hypervisor string `json:"hypervisor,omitempty"`
}

type VPSVolumeSpec struct {
	Capacity *VPSVolumeCapacity `json:"capacity,omitempty"`
	IOPS     *VPSVolumeIOPS     `json:"iops,omitempty"`
}

type VPSVolumeCapacity struct {
	Max int `json:"max"`
}

type VPSVolumeIOPS struct {
	Guaranteed bool   `json:"guaranteed"`
	Level      int    `json:"level"`
	Max        int    `json:"max"`
	MaxUnit    string `json:"maxUnit,omitempty"`
	Unit       string `json:"unit,omitempty"`
}
