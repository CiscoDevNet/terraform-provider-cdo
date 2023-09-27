package acctest

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"os"
	"strconv"
	"strings"
)

type env struct{}

var Env = &env{}

func (e *env) UserDataSourceName() string {
	return e.mustGetString("USER_DATA_SOURCE_NAME")
}

func (e *env) UserDataSourceRole() string {
	return e.mustGetString("USER_DATA_SOURCE_ROLE")
}

func (e *env) UserDataSourceIsApiOnly() bool {
	return e.mustGetBool("USER_DATA_SOURCE_IS_API_ONLY")
}

func (e *env) UserResourceName() string {
	return e.mustGetString("USER_RESOURCE_NAME")
}

func (e *env) UserResourceNewName() string {
	return e.mustGetString("USER_RESOURCE_NEW_NAME")
}

func (e *env) UserResourceIsApiOnly() bool {
	return e.mustGetBool("USER_RESOURCE_IS_API_ONLY")
}

func (e *env) UserResourceRole() string {
	return e.mustGetString("USER_RESOURCE_ROLE")
}

func (e *env) TenantDataSourceName() string {
	return e.mustGetString("TENANT_DATA_SOURCE_NAME")
}

func (e *env) TenantDataSourceHumanReadableName() string {
	return e.mustGetString("TENANT_DATA_SOURCE_HUMAN_READABLE_NAME")
}

func (e *env) TenantDataSourceSubscriptionType() string {
	return e.mustGetString("TENANT_DATA_SOURCE_SUBSCRIPTION_TYPE")
}

func (e *env) IosResourceName() string {
	return e.mustGetString("IOS_RESOURCE_NAME")
}

func (e *env) IosResourceSocketAddress() string {
	return e.mustGetString("IOS_RESOURCE_SOCKET_ADDRESS")
}

func (e *env) IosResourceConnectorType() string {
	return e.mustGetString("IOS_RESOURCE_CONNECTOR_TYPE")
}

func (e *env) IosResourceUsername() string {
	return e.mustGetString("IOS_RESOURCE_USERNAME")
}

func (e *env) IosResourcePassword() string {
	return e.mustGetString("IOS_RESOURCE_PASSWORD")
}

func (e *env) IosResourceConnectorName() string {
	return e.mustGetString("IOS_RESOURCE_CONNECTOR_NAME")
}

func (e *env) IosResourceIgnoreCertificate() string {
	return e.mustGetString("IOS_RESOURCE_IGNORE_CERTIFICATE")
}

func (e *env) IosResourceTags() tags.Type {
	tagsEnv := e.mustGetString("IOS_RESOURCE_TAGS")
	return tags.New(strings.Split(tagsEnv, ",")...)
}

func (e *env) IosResourceHost() string {
	return e.mustGetString("IOS_RESOURCE_HOST")
}

func (e *env) IosResourcePort() int64 {
	return e.mustGetInt("IOS_RESOURCE_PORT")
}

func (e *env) IosResourceNewName() string {
	return e.mustGetString("IOS_RESOURCE_NEW_NAME")
}

func (e *env) IosDataSourceName() string {
	return e.mustGetString("IOS_DATA_SOURCE_NAME")
}

func (e *env) IosDataSourceIgnoreCertificate() string {
	return e.mustGetString("IOS_DATA_SOURCE_IGNORE_CERTIFICATE")
}

func (e *env) IosDataSourceTags() tags.Type {
	tagsEnv := e.mustGetString("IOS_DATA_SOURCE_TAGS")
	return tags.New(strings.Split(tagsEnv, ",")...)
}

func (e *env) FtdResourceName() string {
	return e.mustGetString("FTD_RESOURCE_NAME")
}

func (e *env) FtdResourceAccessPolicyName() string {
	return e.mustGetString("FTD_RESOURCE_ACCESS_POLICY_NAME")
}

func (e *env) FtdResourcePerformanceTier() string {
	return e.mustGetString("FTD_RESOURCE_PERFORMANCE_TIER")
}

func (e *env) FtdResourceVirtual() string {
	return e.mustGetString("FTD_RESOURCE_VIRTUAL")
}

func (e *env) FtdResourceLicenses() string {
	return e.mustGetString("FTD_RESOURCE_LICENSES")
}

func (e *env) FtdResourceTags() tags.Type {
	tagsEnv := e.mustGetString("FTD_RESOURCE_TAGS")
	return tags.New(strings.Split(tagsEnv, ",")...)
}

func (e *env) FtdResourceNewName() string {
	return e.mustGetString("FTD_RESOURCE_NEW_NAME")
}

func (e *env) AsaResourceSdcName() string {
	return e.mustGetString("ASA_RESOURCE_SDC_NAME")
}

func (e *env) AsaResourceSdcSocketAddress() string {
	return e.mustGetString("ASA_RESOURCE_SDC_SOCKET_ADDRESS")
}

func (e *env) AsaResourceSdcConnectorName() string {
	return e.mustGetString("ASA_RESOURCE_SDC_CONNECTOR_NAME")
}

func (e *env) AsaResourceSdcConnectorType() string {
	return e.mustGetString("ASA_RESOURCE_SDC_CONNECTOR_TYPE")
}

func (e *env) AsaResourceSdcUsername() string {
	return e.mustGetString("ASA_RESOURCE_SDC_USERNAME")
}

func (e *env) AsaResourceSdcPassword() string {
	return e.mustGetString("ASA_RESOURCE_SDC_PASSWORD")
}

func (e *env) AsaResourceSdcHost() string {
	return e.mustGetString("ASA_RESOURCE_SDC_HOST")
}

func (e *env) AsaResourceSdcPort() int64 {
	return e.mustGetInt("ASA_RESOURCE_SDC_PORT")
}

func (e *env) AsaResourceSdcIgnoreCertificate() bool {
	return e.mustGetBool("ASA_RESOURCE_SDC_IGNORE_CERTIFICATE")
}

func (e *env) AsaResourceSdcTags() tags.Type {
	tagsEnv := e.mustGetString("ASA_RESOURCE_SDC_TAGS")
	return tags.New(strings.Split(tagsEnv, ",")...)
}

func (e *env) AsaResourceAlternativeDeviceLocation() string {
	return e.mustGetString("ASA_RESOURCE_SDC_ALTERNATIVE_DEVICE_LOCATION")
}

func (e *env) AsaResourceSdcNewName() string {
	return e.mustGetString("ASA_RESOURCE_SDC_NEW_NAME")
}

func (e *env) AsaResourceSdcWrongPassword() string {
	return e.mustGetString("ASA_RESOURCE_SDC_WRONG_PASSWORD")
}

func (e *env) AsaResourceCdgName() string {
	return e.mustGetString("ASA_RESOURCE_CDG_NAME")
}

func (e *env) AsaResourceCdgSocketAddress() string {
	return e.mustGetString("ASA_RESOURCE_CDG_SOCKET_ADDRESS")
}

func (e *env) AsaResourceCdgConnectorName() string {
	return e.mustGetString("ASA_RESOURCE_CDG_CONNECTOR_NAME")
}

func (e *env) AsaResourceCdgConnectorType() string {
	return e.mustGetString("ASA_RESOURCE_CDG_CONNECTOR_TYPE")
}

func (e *env) AsaResourceCdgUsername() string {
	return e.mustGetString("ASA_RESOURCE_CDG_USERNAME")
}

func (e *env) AsaResourceCdgPassword() string {
	return e.mustGetString("ASA_RESOURCE_CDG_PASSWORD")
}

func (e *env) AsaResourceCdgIgnoreCertificate() bool {
	return e.mustGetBool("ASA_RESOURCE_CDG_IGNORE_CERTIFICATE")
}

func (e *env) AsaResourceCdgTags() tags.Type {
	tagsEnv := e.mustGetString("ASA_RESOURCE_CDG_TAGS")
	return tags.New(strings.Split(tagsEnv, ",")...)
}

func (e *env) AsaResourceCdgHost() string {
	return e.mustGetString("ASA_RESOURCE_CDG_HOST")
}

func (e *env) AsaResourceCdgPort() int64 {
	return e.mustGetInt("ASA_RESOURCE_CDG_PORT")
}

func (e *env) AsaResourceCdgNewName() string {
	return e.mustGetString("ASA_RESOURCE_CDG_NEW_NAME")
}

func (e *env) AsaResourceCdgWrongPassword() string {
	return e.mustGetString("ASA_RESOURCE_CDG_WRONG_PASSWORD")
}

func (e *env) AsaDataSourceConnectorType() string {
	return e.mustGetString("ASA_DATA_SOURCE_CONNECTOR_TYPE")
}

func (e *env) AsaDataSourceName() string {
	return e.mustGetString("ASA_DATA_SOURCE_NAME")
}

func (e *env) AsaDataSourceSocketAddress() string {
	return e.mustGetString("ASA_DATA_SOURCE_SOCKET_ADDRESS")
}

func (e *env) AsaDataSourceHost() string {
	return e.mustGetString("ASA_DATA_SOURCE_HOST")
}

func (e *env) AsaDataSourcePort() int64 {
	return e.mustGetInt("ASA_DATA_SOURCE_PORT")
}

func (e *env) AsaDataSourceIgnoreCertificate() bool {
	return e.mustGetBool("ASA_DATA_SOURCE_IGNORE_CERTIFICATE")
}

func (e *env) AsaDataSourceTags() tags.Type {
	tagsEnv := e.mustGetString("ASA_DATA_SOURCE_TAGS")
	return tags.New(strings.Split(tagsEnv, ",")...)
}

func (e *env) ConnectorDataSourceName() string {
	return e.mustGetString("CONNECTOR_DATA_SOURCE_NAME")
}

func (e *env) ConnectorResourceName() string {
	return e.mustGetString("CONNECTOR_RESOURCE_NAME")
}

func (e *env) ConnectorResourceNewName() string {
	return e.mustGetString("CONNECTOR_RESOURCE_NEW_NAME")
}

func (e *env) CloudFmcDataSourceHostname() string {
	return e.mustGetString("CLOUD_FMC_HOSTNAME")
}

func (e *env) CloudFmcDataSourceSoftwareVersion() string {
	return e.mustGetString("CLOUD_FMC_SOFTWARE_VERSION")
}

func (e *env) mustGetString(envName string) string {
	value, ok := os.LookupEnv(envName)
	if ok {
		return value
	}
	panic(fmt.Sprintf("acceptance test requires environment variable: %s to be set.", envName))
}

func (e *env) mustGetBool(envName string) bool {
	value, ok := os.LookupEnv(envName)
	if ok {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			panic(fmt.Sprintf("acceptance requires environment variable: %s to be a boolean value.", value))
		}
		return boolValue
	}
	panic(fmt.Sprintf("acceptance test requires environment variable: %s to be set.", envName))
}

func (e *env) mustGetInt(envName string) int64 {
	value, ok := os.LookupEnv(envName)
	base := 10
	bitSize := 64
	if ok {
		intValue, err := strconv.ParseInt(value, base, bitSize)
		if err != nil {
			panic(fmt.Sprintf("acceptance requires environment variable: %s to be a int value (base: %d, bitSize: %d).", value, base, bitSize))
		}
		return intValue
	}
	panic(fmt.Sprintf("acceptance test requires environment variable: %s to be set.", envName))
}
