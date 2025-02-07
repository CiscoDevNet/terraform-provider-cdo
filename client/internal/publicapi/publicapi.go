package publicapi

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"time"
)

func TriggerTransaction(ctx context.Context, client http.Client, url string, body any) (transaction.Type, error) {
	req := client.NewPost(ctx, url, body)

	return sendAndCheckForError(req)
}

func WaitForTransactionToFinishWithDefaults(ctx context.Context, client http.Client, t transaction.Type, msg string) (transaction.Type, error) {
	return WaitForTransactionToFinish(ctx, client, t, retry.NewOptionsBuilder().
		Logger(client.Logger).
		Timeout(5*time.Minute).
		Retries(-1).
		EarlyExitOnError(true).
		Message(msg).
		Delay(1*time.Second).
		Build())

}

func WaitForTransactionToFinish(ctx context.Context, client http.Client, t transaction.Type, options retry.Options) (transaction.Type, error) {
	if isDone(t) {
		return t, nil
	} else if err := checkForError(t); err != nil {
		return t, err
	} else {
		return pollTransaction(ctx, client, t, options)
	}
}

func pollTransaction(ctx context.Context, client http.Client, t transaction.Type, options retry.Options) (transaction.Type, error) {
	var output transaction.Type
	err := retry.Do(ctx, untilDoneOrError(ctx, client, t.PollingUrl, &output), options)
	if err != nil {
		return transaction.Type{}, err
	}
	if err := checkForError(output); err != nil {
		return output, err
	}

	return output, nil
}

func untilDoneOrError(ctx context.Context, client http.Client, transactionPollingUrl string, trans *transaction.Type) retry.Func {
	return func() (bool, error) {
		req := client.NewGet(ctx, transactionPollingUrl)
		var t transaction.Type
		if err := req.Send(&t); err != nil {
			return false, err
		}
		*trans = t
		client.Logger.Printf("Polled transaction:\ntransactionUid=%s,\nstatus=%s,\ndetails%s\n", t.TransactionUid, t.Status, t.Details)
		return isDoneOrError(t), nil
	}
}

func sendAndCheckForError(req *http.Request) (transaction.Type, error) {
	var t transaction.Type

	if err := req.Send(&t); err != nil {
		return transaction.Type{}, err
	}
	if err := checkForError(t); err != nil {
		return transaction.Type{}, err
	}

	return t, nil
}

func checkForError(transaction transaction.Type) error {
	if transaction.Status == transactionstatus.ERROR {
		return NewTransactionErrorFromTransaction(transaction)
	} else {
		return nil
	}
}

func isDoneOrError(transaction transaction.Type) bool {
	return transaction.Status == transactionstatus.DONE || transaction.Status == transactionstatus.ERROR
}

func isDone(transaction transaction.Type) bool {
	return transaction.Status == transactionstatus.DONE
}
