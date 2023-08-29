// central url management, all urls goes into here
package url

import (
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
)

func ReadDevice(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices/%s", baseUrl, uid)
}

func ReadDeviceByNameAndDeviceType(baseUrl string, deviceName string, deviceType string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices?q=name:%s+AND+deviceType:%s", baseUrl, deviceName, deviceType)
}
func ReadAllDevicesByType(baseUrl string, deviceType devicetype.Type) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices?q=deviceType:%s", baseUrl, deviceType)
}

func CreateDevice(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices", baseUrl)
}

func UpdateDevice(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices/%s", baseUrl, uid)
}

func DeleteDevice(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices/%s", baseUrl, uid)
}

func ReadSpecificDevice(baseUrl string, specificUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/device/%s/specific-device", baseUrl, specificUid)
}

func ReadAllConnectors(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies", baseUrl)
}

func ReadAsaConfig(baseUrl string, specificUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/asa/configs/%s", baseUrl, specificUid)
}

func UpdateAsaConfig(baseUrl string, specificUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/asa/configs/%s", baseUrl, specificUid)
}

func ReadConnectorByUid(baseUrl string, connectorUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies/%s", baseUrl, connectorUid)
}

func ReadConnectorByName(baseUrl string, connectorName string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies?q=name:%s", baseUrl, connectorName)
}

func CreateConnector(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies", baseUrl)
}

func UpdateConnector(baseUrl string, connectorUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies/%s", baseUrl, connectorUid)
}

func DeleteConnector(baseUrl string, connectorUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies/%s", baseUrl, connectorUid)
}

func UserToken(baseUrl string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/oauth/token", baseUrl)
}

func ReadSmartLicense(baseUrl string) string {
	return fmt.Sprintf("%s/fmc/api/fmc_platform/v1/license/smartlicenses", baseUrl)
}

func ReadAccessPolicies(baseUrl string, domainUid string, limit int) string {
	return fmt.Sprintf("%s/fmc/api/fmc_config/v1/domain/%s/policy/accesspolicies?limit=%d", baseUrl, domainUid, limit)
}

func CreateUser(baseUrl string, username string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/users/%s", baseUrl, username)
}

func ReadUserByUsername(baseUrl string, username string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/users?q=name:%s", baseUrl, username)
}

func UserByUid(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/users/%s", baseUrl, uid)
}

func GenerateApiToken(baseUrl string, username string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/oauth/token/%s", baseUrl, username)
}

func RevokeApiToken(baseUrl string, tokenId string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/oauth/revoke/%s", baseUrl, tokenId)
}
