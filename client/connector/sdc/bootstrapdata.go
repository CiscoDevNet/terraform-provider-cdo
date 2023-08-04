package sdc

import (
	"encoding/base64"
	"fmt"
)

func computeBootstrapData(sdcName, accessToken, tenantName, baseUrl, host string) string {
	bootstrapUrl := fmt.Sprintf("%s/sdc/bootstrap/%s/%s", baseUrl, tenantName, sdcName)

	rawBootstrapData := fmt.Sprintf("CDO_TOKEN=%q\nCDO_DOMAIN=%q\nCDO_TENANT=%q\nCDO_BOOTSTRAP_URL=%q\n", accessToken, host, tenantName, bootstrapUrl)

	bootstrapData := base64.StdEncoding.EncodeToString([]byte(rawBootstrapData))

	return bootstrapData
}
