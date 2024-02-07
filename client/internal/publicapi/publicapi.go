package publicapi

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
)

func PostForTransaction(ctx context.Context, client http.Client, url string, body any, options retry.Options) (transaction.Type, error) {

	req := client.NewPost(ctx, url, body)

	t, err := sendAndCheckForError(req)
	if err != nil {
		return t, err
	}

	return poll(ctx, client, options, t.TransactionPollingUrl)
}

func poll(ctx context.Context, client http.Client, options retry.Options, pollingUrl string) (transaction.Type, error) {
	var output transaction.Type
	err := retry.Do(ctx, untilDoneOrError(ctx, client, pollingUrl, &output), options)
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
		var out transaction.Type
		if err := req.Send(&out); err != nil {
			return false, err
		}
		*trans = out
		client.Logger.Printf("status=%s\n", out.Status)
		if out.Status == transactionstatus.DONE || out.Status == transactionstatus.ERROR {
			return true, nil
		} else {
			return false, nil
		}
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
