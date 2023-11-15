package sec

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
)

func ComputeEventOnlyBootstrapData(accessToken, tenantName, baseUrl, host string) string {
	bootstrapUrl := fmt.Sprintf("%s/sdc/bootstrap/%s", baseUrl, tenantName)

	rawBootstrapData := fmt.Sprintf("CDO_TOKEN=%q\nCDO_DOMAIN=%q\nCDO_TENANT=%q\nCDO_BOOTSTRAP_URL=%q\nONLY_EVENTING=\"true\"\n", accessToken, host, tenantName, bootstrapUrl)

	bootstrapData := base64.StdEncoding.EncodeToString([]byte(rawBootstrapData))

	return bootstrapData
}

func generateBootstrapData(ctx context.Context, client http.Client) (string, error) {
	userToken, err := user.GetToken(ctx, client, user.NewGetTokenInput())
	if err != nil {
		return "", err
	}

	return ComputeEventOnlyBootstrapData(userToken.AccessToken, userToken.TenantName, client.BaseUrl(), client.Host()), nil
}
