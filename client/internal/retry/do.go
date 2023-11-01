// Provide utilities for repeated making requests until some conditions are satisified or error out.
// This is the main reason for having replay-able request
// TODO: implement exp backoff
package retry

import (
	"context"
	"errors"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/goutil"
	"log"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
)

// Options represents the configuration when retrying.
// Note that there can also be retries and delay within the http.Request.
// Note that it will terminate whichever timeout or max attempts happens first.
// Note that the total attempts made is retries + 1
type Options struct {
	Timeout time.Duration // Timeout is the duration before force terminate
	Delay   time.Duration // Delay is the duration between consecutive requests.
	Retries int           // Retries is the max number of retries before terminating. Negative means no limit.

	Logger *log.Logger

	// EarlyExitOnError will cause Retry to return immediately if error is returned from Func;
	// if false, it will only return when max retries exceeded, which errors are combined and returned together.
	// Note user can decide in the retry function whether to return error when error occur,
	// so if someone does not want to return early, they can just not return the error, so in theory there is no need for this parameter to ignore error returned.
	// But often, if retry times out, we want to see what caused the timeout, we want to see the error that occurred,
	// to do this without this parameter, the retry function needs to manually accumulate the errors if there are any,
	// this is troublesome to implement in the retry function, and most of the time we would like to see the error when this happens.
	// So we have a parameter to handle this.
	EarlyExitOnError bool

	// Message is added to error to tell user what this retry is doing when error occur.
	Message string
}

// Func is the retryable function for retrying.
// ok: whether ok to stop
// error: if not nil, stop retrying and return that error
type Func func() (ok bool, err error)

const (
	DefaultTimeout          = 3 * time.Minute
	DefaultDelay            = 3 * time.Second
	DefaultRetries          = -1
	DefaultEarlyExitOnError = true
)

var (
	DefaultOpts = Options{
		Timeout: DefaultTimeout,
		Delay:   DefaultDelay,
		Retries: DefaultRetries,

		Logger: cdo.DefaultLogger,

		EarlyExitOnError: DefaultEarlyExitOnError,
	}
)

func NewOptions(logger *log.Logger, timeout time.Duration, delay time.Duration, retries int, earlyExitOnError bool) *Options {
	return &Options{
		Timeout: timeout,
		Delay:   delay,
		Retries: retries,

		Logger: logger,

		EarlyExitOnError: earlyExitOnError,
	}
}

// Do run retry function until response of request satisfy check function, or ends early according to configuration.
func Do(ctx context.Context, retryFunc Func, opt Options) error {
	// set up context
	ctxToUse := ctx
	if opt.Timeout >= 0 {
		var cancel context.CancelFunc
		ctxToUse, cancel = context.WithTimeout(ctx, opt.Timeout)
		defer cancel()
	}
	return doInternal(ctxToUse, retryFunc, opt)
}

func willTimeoutAfterDelay(ctx context.Context, delay time.Duration) bool {
	ddl, ok := ctx.Deadline()
	if !ok {
		return false // no deadline is set, so we will never time out after delay
	}
	return time.Now().Add(goutil.Max(delay, 0)).After(ddl)
}

func doInternal(ctx context.Context, retryFunc Func, opt Options) error {
	// setup time
	startTime := time.Now()
	// setup errors
	// retryErrors[i] = nil: no error occur at this attempt
	// retryErrors[i] != nil: error occur at this attempt
	retryErrors := make([]error, goutil.Max(opt.Retries, 0)+1) // +1 because total attempts = retries + 1 initial attempt

	for attempt := 0; attempt <= opt.Retries || opt.Retries < 0; attempt++ {
		select {
		case <-ctx.Done():
			// context timeout/cancelled
			if errors.Is(ctx.Err(), context.Canceled) {
				return newContextCancelledErrorf(opt.Message, "%s at attempt=%d/%d, after=%s, errors:\n%w\n", ctx.Err(), attempt, opt.Retries, time.Since(startTime), errors.Join(retryErrors...))
			} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return newTimeoutErrorf(opt.Message, "%s at attempt=%d/%d, after=%s, errors:\n%w\n", ctx.Err(), attempt, opt.Retries, time.Since(startTime), errors.Join(retryErrors...))
			} else {
				// channel not yet closed, not possible, if it happens, ignore and continue...
			}
		default:
			if attempt > 0 {
				// not the first attempt, this is a retry, so we do delay
				if willTimeoutAfterDelay(ctx, opt.Delay) {
					return newTimeoutErrorf(opt.Message, "at attempt=%d/%d, after=%s, errors:\n%w\n", attempt, opt.Retries, time.Since(startTime), errors.Join(retryErrors...))
				}
				time.Sleep(opt.Delay)
			}
			// do attempt
			ok, err := retryFunc()
			if opt.Logger != nil {
				opt.Logger.Printf("attempt=%d/%d, ok=%t, error=%s\n", attempt, opt.Retries, ok, err)
			}
			if err == nil {
				retryErrors = append(retryErrors, nil)
			} else {
				retryErrors = append(retryErrors, newAttemptErrorf("at attempt=%d/%d, ok=%t, error=%w\n", attempt, opt.Retries, ok, err))
			}

			if err != nil && opt.EarlyExitOnError {
				return newFuncErrorf(opt.Message, "at attempt=%d/%d, after=%s, errors:\n%w\n", attempt, opt.Retries, time.Since(startTime), errors.Join(retryErrors...))
			}
			if ok {
				return nil
			}
		}
	}
	// max retry exceeded
	return newRetriesExceededErrorf(opt.Message, "after %d retries, time taken=%s, errors:\n%w\n", opt.Retries, time.Since(startTime), errors.Join(retryErrors...))
}
