package cloudftd

import (
	"context"
	"errors"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type FtdUpgradeInput struct {
	Uid             string `json:"uid"`
	SoftwareVersion string `json:"softwareVersion"`
}

type FtdUpgradeService interface {
	Upgrade(uid string, softwareVersion string) (*FtdDevice, error)
}

type ftdUpgradeService struct {
	Ctx    context.Context
	Client *http.Client
}

func NewFtdUpgradeService(ctx context.Context, client *http.Client) FtdUpgradeService {
	return &ftdUpgradeService{
		Ctx:    ctx,
		Client: client,
	}
}

func (f *ftdUpgradeService) Upgrade(uid string, softwareVersion string) (*FtdDevice, error) {
	ftdDevice, err := ReadByUid(f.Ctx, *f.Client, ReadByUidInput{Uid: uid})
	if err != nil {
		return nil, err
	}
	tflog.Debug(f.Ctx, fmt.Sprintf("FTD device found: %v", ftdDevice))

	tflog.Debug(f.Ctx, "Validating if FTD device is suitable for upgrade...")
	err = f.validateDeviceType(ftdDevice)
	if err != nil {
		return nil, err
	}
	err = f.validateConnectivityState(ftdDevice)
	if err != nil {
		return nil, err
	}
	err = f.validateFtdVersion(ftdDevice, softwareVersion)
	if err != nil {
		return nil, err
	}

	return ftdDevice, nil
}

func (f *ftdUpgradeService) validateDeviceType(ftdDevice *FtdDevice) error {
	if ftdDevice.DeviceType != "FTDC" {
		return errors.New("this resource only supports cdFMC managed FTDs")
	}

	return nil
}

func (f *ftdUpgradeService) validateConnectivityState(ftdDevice *FtdDevice) error {
	if ftdDevice.ConnectivityState != 1 {
		return errors.New("FTD device connectivity state is not ONLINE. Only ONLINE devices can be upgraded")
	}

	return nil
}

func (f *ftdUpgradeService) validateFtdVersion(ftdDevice *FtdDevice, softwareVersionToUpgradeToStr string) error {
	versionOnDevice, err := ftd.NewVersion(ftdDevice.SoftwareVersion)
	if err != nil {
		f.Client.Logger.Printf("error parsing software version %s on device\n", ftdDevice.SoftwareVersion)
		return err
	}
	versionToUpgradeTo, err := ftd.NewVersion(softwareVersionToUpgradeToStr)
	if err != nil {
		f.Client.Logger.Printf("error parsing software version %s to upgrade to\n", softwareVersionToUpgradeToStr)
		return err
	}

	if versionOnDevice.GreaterThan(versionToUpgradeTo) {
		return errors.New(fmt.Sprintf("FTD device is on version %s, which is newer than the"+
			" version to upgrade to: %s", ftdDevice.SoftwareVersion, softwareVersionToUpgradeToStr))
	}
	if versionOnDevice.LessThan(versionToUpgradeTo) {
		err = f.validateUpgradePathExistsTo(ftdDevice, versionToUpgradeTo)
		if err != nil {
			return err
		}
		return errors.New("upgrade implementation coming soon")
	}

	return nil
}

func (f *ftdUpgradeService) validateUpgradePathExistsTo(ftdDevice *FtdDevice, toVersion *ftd.Version) error {
	upgradePackages, err := ReadUpgradePackages(f.Ctx, *f.Client, ftdDevice.Uid)
	if err != nil {
		return err
	}
	for _, upgradePackage := range *upgradePackages {
		tflog.Debug(f.Ctx, fmt.Sprintf("Checking upgrade package: %s", upgradePackage.SoftwareVersion))
		softwareVersion, err := ftd.NewVersion(upgradePackage.SoftwareVersion)
		if err != nil {
			f.Client.Logger.Printf("error parsing software version %s in upgrade package\n", upgradePackage.SoftwareVersion)
			return err
		}
		if softwareVersion.Equal(toVersion) {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("%s is not a valid version to upgrade FTD device %s to", toVersion.String(), ftdDevice.Name))
}
