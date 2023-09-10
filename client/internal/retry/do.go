// Provide utilities for repeated making requests until some conditions are satisified or error out.
// This is the main reason for having replay-able request
// TODO: implement exp backoff
package retry

import (
	"errors"
	"fmt"
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

	EarlyExitOnError bool // EarlyExitOnError will cause Retry to return immediately if error is returned from Func; if false, it will only return when max retries exceeded, which errors are combined and returned together.
}

// Func is the retryable function for retrying.
// ok: whether ok to stop
// err: if not nil, stop retrying and return that error
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

func NewOptionsWithLogger(logger *log.Logger) *Options {
	return NewOptions(
		logger,
		DefaultTimeout,
		DefaultDelay,
		DefaultRetries,
		DefaultEarlyExitOnError,
	)
}

func NewOptionsWithLoggerAndRetries(logger *log.Logger, retries int) *Options {
	return NewOptions(
		logger,
		DefaultTimeout,
		DefaultDelay,
		retries,
		DefaultEarlyExitOnError,
	)
}

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
func Do(retryFunc Func, opt Options) error {

	startTime := time.Now()
	endTime := startTime.Add(opt.Timeout)
	timeout := func() bool {
		return time.Now().After(endTime)
	}
	willTimeoutAfterDelay := func() bool {
		return time.Now().Add(opt.Delay).After(endTime)
	}
	accumulatedErrs := make([]error, 0, opt.Retries+1)

	// check if it is already timeout
	if timeout() {
		return fmt.Errorf("retry func timeout before any attempt")
	}

	// attempt once before starting retries
	ok, err := retryFunc()
	opt.Logger.Printf("err=%+v", err)
	accumulatedErrs = append(accumulatedErrs, err)
	if err != nil && opt.EarlyExitOnError {
		return fmt.Errorf("error in retry func, cause=%w", err)
	}
	if ok {
		return nil
	}

	// retry starts
	for retries := 1; opt.Retries < 0 || retries <= opt.Retries; retries++ {

		// if timeout now or will time out during delay between requests
		if timeout() || willTimeoutAfterDelay() {
			return fmt.Errorf("timeout at attempt %d/%d, after %s", retries, opt.Retries, time.Now().Sub(startTime))
		}
		time.Sleep(opt.Delay)

		opt.Logger.Printf("[RETRY] attempt=%d/%d\n", retries, opt.Retries)

		// attempt
		ok, err = retryFunc()
		opt.Logger.Printf("err=%+v", err)
		accumulatedErrs = append(accumulatedErrs, err)
		if err != nil && opt.EarlyExitOnError {
			return fmt.Errorf("error in retry func, cause=%w", err)
		}
		if ok {
			opt.Logger.Println("[RETRY] success")
			return nil
		}
		opt.Logger.Println("[RETRY] failed")
	}

	return fmt.Errorf(
		"max retry reached, retries=%d, time taken=%s, accumulated errors=%w",
		opt.Retries,
		time.Now().Sub(startTime),
		errors.Join(accumulatedErrs...),
	)
}
