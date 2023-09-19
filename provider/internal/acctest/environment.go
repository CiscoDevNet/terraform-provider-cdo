package acctest

import (
	"fmt"
	"os"
)

type env struct{}

var Env = &env{}

func (e *env) UserDataSourceName() string {
	return e.mustGet("USER_DATA_SOURCE_NAME")
}

func (e *env) UserDataSourceRole() string {
	return e.mustGet("USER_DATA_SOURCE_ROLE")
}

func (e *env) UserDataSourceIsApiOnly() string {
	return e.mustGet("USER_DATA_SOURCE_IS_API_ONLY")
}

func (e *env) UserResourceName() string {
	return e.mustGet("USER_RESOURCE_NAME")
}

func (e *env) UserResourceNewName() string {
	return e.mustGet("USER_RESOURCE_NEW_NAME")
}

func (e *env) UserResourceIsApiOnly() string {
	return e.mustGet("USER_RESOURCE_IS_API_ONLY")
}

func (e *env) UserResourceRole() string {
	return e.mustGet("USER_RESOURCE_ROLE")
}

func (e *env) TenantDataSourceName() string {
	return e.mustGet("TENANT_DATA_SOURCE_NAME")
}

func (e *env) TenantDataSourceHumanReadableName() string {
	return e.mustGet("TENANT_DATA_SOURCE_HUMAN_READABLE_NAME")
}

func (e *env) TenantDataSourceSubscriptionType() string {
	return e.mustGet("TENANT_DATA_SOURCE_SUBSCRIPTION_TYPE")
}

func (e *env) IosResourceName() string {
	return e.mustGet("IOS_RESOURCE_NAME")
}

func (e *env) IosResourceSocketAddress() string {
	return e.mustGet("IOS_RESOURCE_SOCKET_ADDRESS")
}

func (e *env) IosResourceConnectorType() string {
	return e.mustGet("IOS_RESOURCE_CONNECTOR_TYPE")
}

func (e *env) IosResourceUsername() string {
	return e.mustGet("IOS_RESOURCE_USERNAME")
}

func (e *env) IosResourcePassword() string {
	return e.mustGet("IOS_RESOURCE_PASSWORD")
}

func (e *env) IosResourceConnectorName() string {
	return e.mustGet("IOS_RESOURCE_CONNECTOR_NAME")
}

func (e *env) IosResourceIgnoreCertificate() string {
	return e.mustGet("IOS_RESOURCE_IGNORE_CERTIFICATE")
}

func (e *env) IosResourceHost() string {
	return e.mustGet("IOS_RESOURCE_HOST")
}

func (e *env) IosResourcePort() string {
	return e.mustGet("IOS_RESOURCE_PORT")
}

func (e *env) IosResourceNewName() string {
	return e.mustGet("IOS_RESOURCE_NEW_NAME")
}

func (e *env) IosDataSourceId() string {
	return e.mustGet("IOS_DATA_SOURCE_ID")
}

func (e *env) IosDataSourceName() string {
	return e.mustGet("IOS_DATA_SOURCE_NAME")
}

func (e *env) IosDataSourceSocketAddress() string {
	return e.mustGet("IOS_DATA_SOURCE_SOCKET_ADDRESS")
}

func (e *env) IosDataSourceHost() string {
	return e.mustGet("IOS_DATA_SOURCE_HOST")
}

func (e *env) IosDataSourcePort() string {
	return e.mustGet("IOS_DATA_SOURCE_PORT")
}

func (e *env) IosDataSourceIgnoreCertificate() string {
	return e.mustGet("IOS_DATA_SOURCE_IGNORE_CERTIFICATE")
}

func (e *env) FtdResourceName() string {
	return e.mustGet("FTD_RESOURCE_NAME")
}

func (e *env) FtdResourceAccessPolicyName() string {
	return e.mustGet("FTD_RESOURCE_ACCESS_POLICY_NAME")
}

func (e *env) FtdResourcePerformanceTier() string {
	return e.mustGet("FTD_RESOURCE_PERFORMANCE_TIER")
}

func (e *env) FtdResourceVirtual() string {
	return e.mustGet("FTD_RESOURCE_VIRTUAL")
}

func (e *env) FtdResourceLicenses() string {
	return e.mustGet("FTD_RESOURCE_LICENSES")
}

func (e *env) FtdResourceNewName() string {
	return e.mustGet("FTD_RESOURCE_NEW_NAME")
}

func (e *env) AsaResourceSdcName() string {
	return e.mustGet("ASA_RESOURCE_SDC_NAME")
}

func (e *env) AsaResourceSdcSocketAddress() string {
	return e.mustGet("ASA_RESOURCE_SDC_SOCKET_ADDRESS")
}

func (e *env) AsaResourceSdcConnectorName() string {
	return e.mustGet("ASA_RESOURCE_SDC_CONNECTOR_NAME")
}

func (e *env) AsaResourceSdcConnectorType() string {
	return e.mustGet("ASA_RESOURCE_SDC_CONNECTOR_TYPE")
}

func (e *env) AsaResourceSdcUsername() string {
	return e.mustGet("ASA_RESOURCE_SDC_USERNAME")
}

func (e *env) AsaResourceSdcPassword() string {
	return e.mustGet("ASA_RESOURCE_SDC_PASSWORD")
}

func (e *env) AsaResourceSdcIgnoreCertificate() string {
	return e.mustGet("ASA_RESOURCE_SDC_IGNORE_CERTIFICATE")
}

func (e *env) AsaResourceSdcHost() string {
	return e.mustGet("ASA_RESOURCE_SDC_HOST")
}

func (e *env) AsaResourceSdcPort() string {
	return e.mustGet("ASA_RESOURCE_SDC_PORT")
}

func (e *env) AsaResourceAlternativeDeviceLocation() string {
	return e.mustGet("ASA_RESOURCE_SDC_ALTERNATIVE_DEVICE_LOCATION")
}

func (e *env) AsaResourceSdcNewName() string {
	return e.mustGet("ASA_RESOURCE_SDC_NEW_NAME")
}

func (e *env) AsaResourceSdcWrongPassword() string {
	return e.mustGet("ASA_RESOURCE_SDC_WRONG_PASSWORD")
}

func (e *env) AsaResourceCdgName() string {
	return e.mustGet("ASA_RESOURCE_CDG_NAME")
}

func (e *env) AsaResourceCdgSocketAddress() string {
	return e.mustGet("ASA_RESOURCE_CDG_SOCKET_ADDRESS")
}

func (e *env) AsaResourceCdgConnectorName() string {
	return e.mustGet("ASA_RESOURCE_CDG_CONNECTOR_NAME")
}

func (e *env) AsaResourceCdgConnectorType() string {
	return e.mustGet("ASA_RESOURCE_CDG_CONNECTOR_TYPE")
}

func (e *env) AsaResourceCdgUsername() string {
	return e.mustGet("ASA_RESOURCE_CDG_USERNAME")
}

func (e *env) AsaResourceCdgPassword() string {
	return e.mustGet("ASA_RESOURCE_CDG_PASSWORD")
}

func (e *env) AsaResourceCdgIgnoreCertificate() string {
	return e.mustGet("ASA_RESOURCE_CDG_IGNORE_CERTIFICATE")
}

func (e *env) AsaResourceCdgHost() string {
	return e.mustGet("ASA_RESOURCE_CDG_HOST")
}

func (e *env) AsaResourceCdgPort() string {
	return e.mustGet("ASA_RESOURCE_CDG_PORT")
}

func (e *env) AsaResourceCdgNewName() string {
	return e.mustGet("ASA_RESOURCE_CDG_NEW_NAME")
}

func (e *env) AsaResourceCdgWrongPassword() string {
	return e.mustGet("ASA_RESOURCE_CDG_WRONG_PASSWORD")
}

func (e *env) AsaDataSourceConnectorType() string {
	return e.mustGet("ASA_DATA_SOURCE_CONNECTOR_TYPE")
}

func (e *env) AsaDataSourceName() string {
	return e.mustGet("ASA_DATA_SOURCE_NAME")
}

func (e *env) AsaDataSourceSocketAddress() string {
	return e.mustGet("ASA_DATA_SOURCE_SOCKET_ADDRESS")
}

func (e *env) AsaDataSourceHost() string {
	return e.mustGet("ASA_DATA_SOURCE_HOST")
}

func (e *env) AsaDataSourcePort() string {
	return e.mustGet("ASA_DATA_SOURCE_PORT")
}

func (e *env) AsaDataSourceIgnoreCertificate() string {
	return e.mustGet("ASA_DATA_SOURCE_IGNORE_CERTIFICATE")
}

func (e *env) ConnectorDataSourceName() string {
	return e.mustGet("CONNECTOR_DATA_SOURCE_NAME")
}

func (e *env) ConnectorResourceName() string {
	return e.mustGet("CONNECTOR_RESOURCE_NAME")
}

func (e *env) ConnectorResourceNewName() string {
	return e.mustGet("CONNECTOR_RESOURCE_NEW_NAME")
}

func (e *env) CloudFmcDataSourceHostname() string {
	return e.mustGet("CLOUD_FMC_HOSTNAME")
}

func (e *env) CloudFmcDataSourceSoftwareVersion() string {
	return e.mustGet("CLOUD_FMC_SOFTWARE_VERSION")
}

func (e *env) mustGet(envName string) string {
	value, ok := os.LookupEnv(envName)
	if ok {
		return value
	}
	panic(fmt.Sprintf("acceptance test requires environment variable: %s to be set.", envName))
}
