// Provide utilities for repeated making requests until some conditions are satisified or error out.
// This is the main reason for having replay-able request
// TODO: implement exp backoff
package retry

import (
	"context"
	"errors"
	"fmt"
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
		accumulatedErrs = append(accumulatedErrs, err)
		if err != nil && opt.EarlyExitOnError {
			return fmt.Errorf("error in retry func, cause=%w", err)
		}
		if ok {
			opt.Logger.Println("[RETRY] success")
			return nil
		}
		opt.Logger.Println("[RETRY] failed")
		opt.Logger.Printf("[RETRY] reason: ok=%t err=%+v", ok, err)
	}

	return fmt.Errorf(
		"max retry reached, retries=%d, time taken=%s, accumulated errors=%w",
		opt.Retries,
		time.Now().Sub(startTime),
		errors.Join(accumulatedErrs...),
	)
}

func Do2(ctx context.Context, retryFunc Func, opt Options) error {
	// set up context
	ctxToUse := ctx
	if opt.Timeout > 0 {
		var cancel context.CancelFunc
		ctxToUse, cancel = context.WithTimeout(ctx, opt.Timeout)
		defer cancel()
	}

	// set up channels
	retryChan := make(chan retryInfo)
	go retry(ctxToUse, retryChan, opt.Delay, opt.Retries)
	return doInternal(ctxToUse, retryFunc, retryChan, opt)
}

type retryInfo struct {
	attempt int
	timeout bool
}

// retry sends a number to the given channel according to given delay and retries,
// it will send retries+1 numbers to the channel, the number sent is the current retry attempt,
// 0 indicate the first attempt, -1 indicate an early timeout during future delay.
func retry(ctx context.Context, c chan<- retryInfo, delay time.Duration, retries int) {
	c <- retryInfo{attempt: 0, timeout: false}
	for attempt := 1; attempt <= retries || retries < 0; attempt++ {
		if willTimeoutAfterDelay(ctx, delay) {
			c <- retryInfo{attempt: attempt, timeout: true}
		}
		time.Sleep(delay)
		c <- retryInfo{attempt: attempt, timeout: false}
	}
	close(c)
}

func willTimeoutAfterDelay(ctx context.Context, delay time.Duration) bool {
	ddl, ok := ctx.Deadline()
	if !ok {
		return false // no deadline is set, so we will never time out after delay
	}
	return time.Now().Add(goutil.Max(delay, 0)).After(ddl)
}

func doInternal(ctx context.Context, retryFunc Func, retryChan <-chan retryInfo, opt Options) error {
	// setup time
	startTime := time.Now()
	// setup errors
	// retryErrors[i] = nil: no error occur at this attempt
	// retryErrors[i] != nil: error occur at this attempt
	retryErrors := make([]error, goutil.Max(opt.Retries, 0)+1)

	// setup retry info
	var retryInfo retryInfo
	var retryOk = true

	for {
		select {
		case <-ctx.Done():
			// context timeout
			return fmt.Errorf("%w at attempt=%d/%d, after=%s, retry errors=%w", ctx.Err(), retryInfo.attempt, opt.Retries, time.Since(startTime), errors.Join(retryErrors...))
		case retryInfo, retryOk = <-retryChan:
			// retry signal received from retry channel
			if !retryOk {
				// max retry exceeded
				return fmt.Errorf("failed after %d retries, time taken=%s, retry errors=%w", opt.Retries, time.Since(startTime), errors.Join(retryErrors...))
			}
			if retryInfo.timeout {
				// early timeout, indicates we foresee that it will time out during the delay,
				return fmt.Errorf("timeout at attempt=%d/%d, after=%s, retry errors:\n%w\n", retryInfo.attempt, opt.Retries, time.Since(startTime), errors.Join(retryErrors...))
			}
			// do attempt
			ok, err := retryFunc()
			retryErrors = append(retryErrors, fmt.Errorf("error at attempt %d/%d: %w", retryInfo.attempt, opt.Retries, err))
			if err != nil && opt.EarlyExitOnError {
				return fmt.Errorf("retry early exited on error=%w", err)
			}
			if ok {
				return nil
			}
		}
	}
}
