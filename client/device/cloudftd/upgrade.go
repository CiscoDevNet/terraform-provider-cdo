package cloudftd

import (
	"context"
	"errors"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

type FtdUpgradeInput struct {
	UpgradePackageUid string `json:"upgradePackageUid"`
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

func (f *ftdUpgradeService) Upgrade(uid string, softwareVersionStr string) (*FtdDevice, error) {
	var upgradePackage *UpgradePackage
	var newSoftwareVersion *ftd.Version
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
	newSoftwareVersion, err = f.validateFtdVersion(ftdDevice, softwareVersionStr)
	if err != nil {
		return nil, err
	}
	if newSoftwareVersion == nil {
		tflog.Debug(f.Ctx, "New software version is the same as the current software version. No upgrade needed.")
		return ftdDevice, nil
	}
	upgradePackage, err = f.validateUpgradePathExistsTo(ftdDevice, newSoftwareVersion)
	if err != nil {
		return nil, err
	}

	tflog.Debug(f.Ctx, "Triggering the FTD device upgrade.")
	return f.doUpgrade(upgradePackage, ftdDevice)
}

func (f *ftdUpgradeService) doUpgrade(upgradePackage *UpgradePackage, ftdDevice *FtdDevice) (*FtdDevice, error) {
	upgradeUrl := url.GetFtdUpgradeUrl(f.Client.BaseUrl(), ftdDevice.Uid)
	transaction, err := publicapi.TriggerTransaction(f.Ctx, *f.Client, upgradeUrl, FtdUpgradeInput{
		UpgradePackageUid: upgradePackage.UpgradePackageUid,
	})
	if err != nil {
		return nil, err
	}

	// poll every 30 seconds for up to 60 minutes
	_, err = publicapi.WaitForTransactionToFinish(f.Ctx, *f.Client, transaction, retry.NewOptionsBuilder().
		Logger(f.Client.Logger).
		Timeout(60*time.Minute).
		Retries(-1).
		EarlyExitOnError(true).
		Message(fmt.Sprintf("Upgrading FTD device to version: %s)", upgradePackage.SoftwareVersion)).
		Delay(30*time.Second).
		Build())
	if err != nil {
		return nil, err
	}

	f.Client.Logger.Println("FTD upgrade successful.")
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

func (f *ftdUpgradeService) validateFtdVersion(ftdDevice *FtdDevice, softwareVersionToUpgradeToStr string) (*ftd.Version, error) {
	versionOnDevice, err := ftd.NewVersion(ftdDevice.SoftwareVersion)
	if err != nil {
		f.Client.Logger.Printf("error parsing software version %s on device\n", ftdDevice.SoftwareVersion)
		return nil, err
	}
	versionToUpgradeTo, err := ftd.NewVersion(softwareVersionToUpgradeToStr)
	if err != nil {
		f.Client.Logger.Printf("error parsing software version %s to upgrade to\n", softwareVersionToUpgradeToStr)
		return nil, err
	}
	if versionOnDevice.GreaterThan(versionToUpgradeTo) {
		return nil, errors.New(fmt.Sprintf("FTD device is on version %s, which is newer than the"+
			" version to upgrade to: %s", ftdDevice.SoftwareVersion, softwareVersionToUpgradeToStr))
	}
	if versionOnDevice.LessThan(versionToUpgradeTo) {
		return versionToUpgradeTo, nil
	}

	return nil, nil
}

func (f *ftdUpgradeService) validateUpgradePathExistsTo(ftdDevice *FtdDevice, toVersion *ftd.Version) (*UpgradePackage, error) {
	upgradePackages, err := ReadUpgradePackages(f.Ctx, *f.Client, ftdDevice.Uid)
	if err != nil {
		return nil, err
	}
	for _, upgradePackage := range *upgradePackages {
		tflog.Debug(f.Ctx, fmt.Sprintf("Checking upgrade package: %s", upgradePackage.SoftwareVersion))
		softwareVersion, err := ftd.NewVersion(upgradePackage.SoftwareVersion)
		if err != nil {
			f.Client.Logger.Printf("error parsing software version %s in upgrade package\n", upgradePackage.SoftwareVersion)
			return nil, err
		}
		if softwareVersion.Equal(toVersion) {
			return &upgradePackage, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("%s is not a valid version to upgrade FTD device %s to", toVersion.String(), ftdDevice.Name))
}
