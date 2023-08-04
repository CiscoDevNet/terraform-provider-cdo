package sdc

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/user"
)

func generateBootstrapData(ctx context.Context, client http.Client, sdcName string) (string, error) {
	userToken, err := user.GetToken(ctx, client, user.NewGetTokenInput())
	if err != nil {
		return "", err
	}

	host, err := client.Host()
	if err != nil {
		return "", err
	}

	return ComputeBootstrapData(
		sdcName, userToken.AccessToken, userToken.TenantName, client.BaseUrl(), host,
	), nil
}

func ComputeBootstrapData(sdcName, accessToken, tenantName, baseUrl, host string) string {
	bootstrapUrl := fmt.Sprintf("%s/sdc/bootstrap/%s/%s", baseUrl, tenantName, sdcName)

	rawBootstrapData := fmt.Sprintf("CDO_TOKEN=%q\nCDO_DOMAIN=%q\nCDO_TENANT=%q\nCDO_BOOTSTRAP_URL=%q\n", accessToken, host, tenantName, bootstrapUrl)

	bootstrapData := base64.StdEncoding.EncodeToString([]byte(rawBootstrapData))

	return bootstrapData
}
