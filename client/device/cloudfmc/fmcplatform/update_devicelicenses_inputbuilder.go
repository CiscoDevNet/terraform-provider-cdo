package fmcplatform

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"

type UpdateDeviceLicensesInputBuilder struct {
	updateDeviceLicensesInput *UpdateDeviceLicensesInput
}

func NewUpdateDeviceLicensesInputBuilder() *UpdateDeviceLicensesInputBuilder {
	updateDeviceLicensesInput := &UpdateDeviceLicensesInput{}
	b := &UpdateDeviceLicensesInputBuilder{updateDeviceLicensesInput: updateDeviceLicensesInput}
	return b
}

func (b *UpdateDeviceLicensesInputBuilder) FmcHost(fmcHost string) *UpdateDeviceLicensesInputBuilder {
	b.updateDeviceLicensesInput.FmcHost = fmcHost
	return b
}

func (b *UpdateDeviceLicensesInputBuilder) LicenseTypes(licenseTypes []license.Type) *UpdateDeviceLicensesInputBuilder {
	b.updateDeviceLicensesInput.LicenseTypes = licenseTypes
	return b
}

func (b *UpdateDeviceLicensesInputBuilder) Build() UpdateDeviceLicensesInput {
	return *b.updateDeviceLicensesInput
}
