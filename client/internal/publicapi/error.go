package publicapi

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
)

type ErrorType interface {
	Error() string
}

type errorType struct{}

func (err errorType) Error() string {
	return "error"
}

var Error = errorType{}

var TransactionError = fmt.Errorf("%w: transaction failed", Error)

func newTransactionErrorf(format string, a ...any) ErrorType {
	return fmt.Errorf("%w, %w", TransactionError, fmt.Errorf(format, a...))
}
func NewTransactionErrorFromTransaction(transaction transaction.Type) ErrorType {
	return newTransactionErrorf("uid=%s, status=%s, message=%s, details=%s", transaction.TransactionUid, transaction.Status, transaction.ErrorMessage, transaction.ErrorDetails)
}
