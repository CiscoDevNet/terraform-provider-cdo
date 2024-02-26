package testing

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/duoadminpanel"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
)

func (m Model) DuoAdminPanelCreateInput() duoadminpanel.CreateInput {
	return duoadminpanel.CreateInput{
		Name:           m.DuoAdminPanelName,
		Host:           m.DuoAdminPanelHost,
		IntegrationKey: "test-int-key",
		SecretKey:      "test-secret-key",
		Labels:         NewTestingLabels(),
	}
}

func (m Model) CreateDuoAdminPanelCreateOutput() duoadminpanel.CreateOutput {
	return m.DuoAdminPanelReadOutput()
}

func (m Model) DuoAdminPanelReadOutput() duoadminpanel.ReadOutput {
	return duoadminpanel.ReadOutput{
		Uid:   m.DuoAdminPanelUid.String(),
		Name:  m.DuoAdminPanelName,
		State: state.DONE,
		Tags:  tags.NewUngrouped(m.DuoAdminPanelLabels...),
	}
}
