// Provide utilities for repeated making requests until some conditions are satisified or error out.
// This is the main reason for having replay-able request
// TODO: implement exp backoff
package retry

import (
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
}

// Func is the retryable function for retrying.
// ok: whether to to stop
// err: if not nil, stop retrying
type Func func() (ok bool, err error)

const (
	DefaultTimeout = 3 * time.Minute
	DefaultDelay   = 3 * time.Second
	DefaultRetries = -1
)

var (
	DefaultOpts = Options{
		Timeout: DefaultTimeout,
		Delay:   DefaultDelay,
		Retries: DefaultRetries,

		Logger: cdo.DefaultLogger,
	}
)

func NewOptionsWithLogger(logger *log.Logger) *Options {
	return NewOptions(
		logger,
		DefaultTimeout,
		DefaultDelay,
		DefaultRetries,
	)
}

func NewOptions(logger *log.Logger, timeout time.Duration, delay time.Duration, retries int) *Options {
	return &Options{
		Timeout: timeout,
		Delay:   delay,
		Retries: retries,

		Logger: logger,
	}
}

// Do run retry function until response of request satisfy check function, or ends early according to configuration.
func Do(retryFunc Func, opt Options) error {

	endTime := time.Now().Add(opt.Timeout)
	timeout := func() bool {
		return time.Now().After(endTime)
	}
	maxRetryReached := func(retries int) bool {
		return opt.Retries > 0 && retries >= opt.Retries
	}

	// initial attempt
	if timeout() {
		return fmt.Errorf("timeout")
	}
	ok, err := retryFunc()
	if err != nil {
		return fmt.Errorf("error in retry func, cause=%w", err)
	}
	if ok {
		return nil
	}

	// retry starts
	opt.Logger.Println("[RETRY] starts")

	for retries := 1; ; retries++ {
		// if timeout now or will time out during delay between requests
		if timeout() || time.Now().Add(opt.Delay).After(endTime) {
			return fmt.Errorf("timeout")
		}
		time.Sleep(opt.Delay)

		opt.Logger.Printf("[RETRY] attempt=%v\n", retries)

		ok, err = retryFunc()
		if err != nil {
			return fmt.Errorf("error in retry func, cause=%w", err)
		}
		if ok {
			opt.Logger.Println("[RETRY] success")
			return nil
		}
		opt.Logger.Println("[RETRY] failed")

		if maxRetryReached(retries) {
			return fmt.Errorf("max retry reached")
		}
	}
}
