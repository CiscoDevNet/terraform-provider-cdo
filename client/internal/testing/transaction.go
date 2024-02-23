package testing

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	"time"
)

func (m Model) CreateDoneTransaction(entityUid string, type_ transactiontype.Type) transaction.Type {
	return transaction.Type{
		TransactionUid:  m.TransactionUid.String(),
		TenantUid:       m.TenantUid.String(),
		EntityUid:       entityUid,
		EntityUrl:       fmt.Sprintf("%s/v1/inventory/devices/%s", m.BaseUrl, entityUid),
		PollingUrl:      fmt.Sprintf("%s/v1/transactions/%s", m.BaseUrl, m.TransactionUid),
		SubmissionTime:  m.TransactionSubmissionTime.Format(time.RFC3339),
		LastUpdatedTime: time.Now().Format(time.RFC3339),
		Type:            type_,
		Status:          transactionstatus.DONE,
	}
}

func (m Model) CreateErrorTransaction(entityUid string, type_ transactiontype.Type) transaction.Type {
	return transaction.Type{
		TransactionUid:  m.TransactionUid.String(),
		TenantUid:       m.TenantUid.String(),
		EntityUid:       entityUid,
		EntityUrl:       fmt.Sprintf("%s/v1/inventory/devices/%s", m.BaseUrl, entityUid),
		PollingUrl:      fmt.Sprintf("%s/v1/transactions/%s", m.BaseUrl, m.TransactionUid),
		SubmissionTime:  m.TransactionSubmissionTime.Format(time.RFC3339),
		LastUpdatedTime: time.Now().Format(time.RFC3339),
		Type:            type_,
		Status:          transactionstatus.ERROR,
		ErrorMessage:    randString("test-error-message"),
		ErrorDetails: map[string]string{
			randString("test-error-detail-key"): randString("test-error-detail-value"),
		},
	}
}
