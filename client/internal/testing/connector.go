package testing

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/status"
)

func (m Model) CdgReadOutput() connector.ReadOutput {
	return connector.ReadOutput{
		Uid:                       m.CdgUid.String(),
		Name:                      m.CdgName,
		DefaultConnector:          false,
		Cdg:                       true,
		TenantUid:                 m.TenantUid.String(),
		PublicKey:                 model.PublicKey{},
		ConnectorStatus:           status.Active,
		IsCommunicationQueueReady: true,
	}
}
