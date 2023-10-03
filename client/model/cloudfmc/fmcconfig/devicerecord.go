package fmcconfig

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"

// DeviceRecord schema is from the device tab of <fmc-url-here>/api/api-explorer/
type DeviceRecord struct {
	Id                     string               `json:"id"`
	Type                   string               `json:"type"`
	Links                  internal.Links       `json:"links"`
	Name                   string               `json:"name"`
	Description            string               `json:"description"`
	Model                  string               `json:"model"`
	ModelId                string               `json:"modelId"`
	ModelNumber            string               `json:"modelNumber"`
	ModelType              string               `json:"modelType"`
	HealthStatus           string               `json:"healthStatus"`
	HealthMessage          string               `json:"healthMessage"`
	SwVersion              string               `json:"sw_version"`
	HealthPolicy           internal.NoLinkItem  `json:"healthPolicy"`
	AccessPolicy           internal.NoLinkItem  `json:"accessPolicy"`
	Hostname               string               `json:"hostname"`
	LicenseCaps            []string             `json:"license_caps"` // this is different from the license_caps we have in CDO, e.g. it has ESSENTIALS instead of BASE
	PerformanceTier        string               `json:"performance_tier"`
	KeepLocalEvents        bool                 `json:"keepLocalEvents"`
	ProhibitPacketTransfer bool                 `json:"prohibitPacketTransfer"`
	IsConnected            bool                 `json:"isConnected"`
	FtdMode                string               `json:"ftdMode"` // e.g. ROUTED / TRANSPARENT
	AnalyticsOnly          bool                 `json:"analyticsOnly"`
	SnortEngine            string               `json:"snortEngine"`
	Metadata               DeviceRecordMetadata `json:"metadata"`
	DeploymentStatus       string               `json:"deploymentStatus"`
}

type DeviceRecordMetadata struct {
	ReadOnly                  ReadOnly            `json:"readOnly"`
	InventoryData             InventoryData       `json:"inventoryData"`
	DeviceSerialNumber        string              `json:"deviceSerialNumber"`
	Domain                    internal.NoLinkItem `json:"domain"`
	IsMultiInstance           bool                `json:"isMultiInstance"`
	SnortVersion              string              `json:"snortVersion"`
	VdbVersion                string              `json:"vdbVersion"`
	LspVersion                string              `json:"lspVersion"`
	ClusterBootstrapSupported bool                `json:"clusterBootstrapSupported"`
}

type InventoryData struct {
	CPUCores   string `json:"cpuCores"`
	CPUType    string `json:"cpuType"`
	MemoryInMB string `json:"memoryInMB"`
}
type ReadOnly struct {
	State  bool   `json:"state"`
	Reason string `json:"reason"`
}
