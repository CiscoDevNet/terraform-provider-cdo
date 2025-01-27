package cloudftd

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
)

type UpgradePackage struct {
	UpgradePackageUid string `json:"upgradePackageUid"`
	SoftwareVersion   string `json:"softwareVersion"`
}

func ReadUpgradePackages(ctx context.Context, client http.Client, deviceUid string) (*[]UpgradePackage, error) {
	readUrl := url.GetFtdUpgradePackagesUrl(client.BaseUrl(), deviceUid)
	req := client.NewGet(ctx, readUrl)

	upgradePackageResponse := model.CdoListResponse[UpgradePackage]{}
	if err := req.Send(&upgradePackageResponse); err != nil {
		return nil, err
	}

	return &upgradePackageResponse.Items, nil
}
