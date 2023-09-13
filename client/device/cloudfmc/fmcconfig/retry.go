package fmcconfig

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcconfig/fmctaskstatus"
	"gopkg.in/errgo.v2/errors"
	"time"
)

var (
	TaskFailedError        = fmt.Errorf("task failed")
	UnknownTaskStatusError = fmt.Errorf("unknown task status")
)

func UntilTaskStatusSuccess(ctx context.Context, client http.Client, readInp ReadTaskStatusInput) retry.Func {
	return func() (bool, error) {
		readTaskOutp, err := ReadTaskStatus(ctx, client, readInp)
		if err != nil {
			return false, err
		}
		if readTaskOutp.Status == fmctaskstatus.Success {
			return true, nil
		}
		if readTaskOutp.Status == fmctaskstatus.Running || readTaskOutp.Status == fmctaskstatus.Pending {
			return false, nil
		}
		if readTaskOutp.Status == fmctaskstatus.Failed {
			return false, errors.Because(TaskFailedError, nil, readTaskOutp.Message)
		}
		return false, errors.Because(UnknownTaskStatusError, nil, fmt.Sprintf("task status = %s", readTaskOutp.Status))
	}
}

func UntilCreateDeviceRecordSuccess(ctx context.Context, client http.Client, createDeviceRecordInput CreateDeviceRecordInput, output *CreateDeviceRecordOutput) retry.Func {
	return func() (bool, error) {
		createDeviceOutp, err := CreateDeviceRecord(ctx, client, createDeviceRecordInput)
		if err != nil {
			return false, err
		}
		*output = *createDeviceOutp
		if err != nil {
			return false, err
		}
		readInp := NewReadTaskStatusInput(createDeviceRecordInput.FmcDomainUid, createDeviceOutp.Metadata.Task.Id, createDeviceRecordInput.FmcHostname)
		err = retry.Do(
			UntilTaskStatusSuccess(ctx, client, readInp),
			retry.NewOptionsBuilder().
				Retries(-1).
				Logger(client.Logger).
				Timeout(30*time.Minute). // usually takes ~5 minutes
				EarlyExitOnError(true).
				Delay(3*time.Second).
				Build(),
		)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}
