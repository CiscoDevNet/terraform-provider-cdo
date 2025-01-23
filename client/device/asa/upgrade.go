package asa

import (
	"context"
	"errors"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"strings"
	"time"
)

func ValidateVersionCompatibility(ctx context.Context, client http.Client, deviceUid string, softwareVersion string, asdmVersion string) error {
	compatibilityUrl := url.GetCompatibleAsaVersions(client.BaseUrl(), deviceUid)
	req := client.NewGet(ctx, compatibilityUrl)
	compatibleVersionsResponse := model.CdoListResponse[CompatibleVersion]{}
	if err := req.Send(&compatibleVersionsResponse); err != nil {
		return err
	}

	trimmedSoftwareVersion := strings.TrimSpace(softwareVersion)
	trimmedAsdmVersion := strings.TrimSpace(asdmVersion)
	for _, compatibleVersion := range compatibleVersionsResponse.Items {
		if trimmedSoftwareVersion != "" && trimmedAsdmVersion != "" {
			if compatibleVersion.SoftwareVersion == softwareVersion && compatibleVersion.AsdmVersion == asdmVersion {
				return nil
			}
		} else if trimmedAsdmVersion == "" {
			if compatibleVersion.SoftwareVersion == softwareVersion {
				return nil
			}
		} else {
			if compatibleVersion.AsdmVersion == asdmVersion {
				return nil
			}
		}
	}

	return errors.New(fmt.Sprintf("Device cannot be upgraded to the specified software and ASDM versions.\n%s\n", buildCompatibleVersionsAsString(compatibleVersionsResponse.Items)))
}

func UpgradeAsa(ctx context.Context, client http.Client, deviceUid string, softwareVersion string, asdmVersion string) error {
	upgradeUrl := url.GetUpgradeAsaUrl(client.BaseUrl(), deviceUid)
	transaction, err := publicapi.TriggerTransaction(ctx, client, upgradeUrl, upgradeAsaInput{
		SoftwareVersion: softwareVersion,
		AsdmVersion:     asdmVersion,
	})
	if err != nil {
		return err
	}

	// poll every 30 seconds for up to 30 minutes
	_, err = publicapi.WaitForTransactionToFinish(ctx, client, transaction, retry.NewOptionsBuilder().
		Logger(client.Logger).
		Timeout(30*time.Minute).
		Retries(-1).
		EarlyExitOnError(true).
		Message(fmt.Sprintf("Upgrading ASA device to %s (ASDM version: %s)", softwareVersion, asdmVersion)).
		Delay(30*time.Second).
		Build())
	if err != nil {
		return err
	}

	client.Logger.Println("ASA upgrade successful!")
	return nil
}

type upgradeAsaInput struct {
	SoftwareVersion string `json:"softwareVersion,omitempty"`
	AsdmVersion     string `json:"asdmVersion,omitempty"`
}

func buildCompatibleVersionsAsString(compatibleVersions []CompatibleVersion) string {
	softwareVersions := []string{}
	asdmVersions := []string{}
	softwareVersionSet := make(map[string]struct{})
	asdmVersionSet := make(map[string]struct{})

	for _, version := range compatibleVersions {
		if version.SoftwareVersion != "" {
			if _, exists := softwareVersionSet[version.SoftwareVersion]; !exists {
				softwareVersionSet[version.SoftwareVersion] = struct{}{}
				softwareVersions = append(softwareVersions, version.SoftwareVersion)
			}
		}
		if version.AsdmVersion != "" {
			if _, exists := asdmVersionSet[version.AsdmVersion]; !exists {
				asdmVersionSet[version.AsdmVersion] = struct{}{}
				asdmVersions = append(asdmVersions, version.AsdmVersion)
			}
		}
	}

	return fmt.Sprintf("Compatible ASA versions: %s\nCompatible ASDM versions: %s", strings.Join(softwareVersions, ", "), strings.Join(asdmVersions, ", "))
}
