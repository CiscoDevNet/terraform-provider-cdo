// central url management, all urls goes into here
package url

import (
	"fmt"
)

func ReadDevice(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices/%s", baseUrl, uid)
}

func ReadDeviceByNameAndDeviceType(baseUrl string, deviceName string, deviceType string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices?q=name:%s+AND+deviceType:%s", baseUrl, deviceName, deviceType)
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

func ReadAllSdcs(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies", baseUrl)
}

func ReadAsaConfig(baseUrl string, specificUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/asa/configs/%s", baseUrl, specificUid)
}

func UpdateAsaConfig(baseUrl string, specificUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/asa/configs/%s", baseUrl, specificUid)
}

func ReadSdcByUid(baseUrl string, larUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies/%s", baseUrl, larUid)
}

func ReadSdcByName(baseUrl string, larName string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies?q=name:%s", baseUrl, larName)
}

func CreateSdc(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies", baseUrl)
}

func UserToken(baseUrl string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/oauth/token", baseUrl)
}
