package url

import "fmt"

func ReadTenantSettings(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/settings/tenant", baseUrl)
}

func UpdateTenantSettings(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/settings/tenant", baseUrl)
}
