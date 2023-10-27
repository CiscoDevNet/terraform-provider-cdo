package retry_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestRetryContextOrOptionsErrors(t *testing.T) {
	testCases := []struct {
		testName       string
		contextTimeout time.Duration
		retryOptions   retry.Options
		retryFunc      retry.Func
		expectedError  retry.ErrorType
	}{
		{
			testName:       "Should retries exceed when retries finishes before any timeout occur",
			contextTimeout: time.Second,
			retryOptions: retry.Options{
				Timeout:          time.Second,
				Delay:            2 * time.Millisecond,
				Retries:          2,
				Logger:           log.Default(),
				EarlyExitOnError: false,
			},
			retryFunc: func() (bool, error) {
				return false, nil
			},
			// retries * delay is less than context.Timeout or retryOptions.Timeout
			expectedError: retry.RetriesExceededError,
		},
		{
			testName:       "Should timeout when context timeout before retries or retry options timeout occur",
			contextTimeout: 2 * time.Millisecond,
			retryOptions: retry.Options{
				Timeout:          time.Second,
				Delay:            time.Millisecond,
				Retries:          1000,
				Logger:           log.Default(),
				EarlyExitOnError: false,
			},
			retryFunc: func() (bool, error) {
				return false, nil
			},
			// context timeout is less than retries * delay or retryOptions.Timeout
			expectedError: retry.TimeoutError,
		},
		{
			testName:       "Should timeout when retry options timeout before retries or context retry timeout occur",
			contextTimeout: time.Second,
			retryOptions: retry.Options{
				Timeout:          5 * time.Millisecond,
				Delay:            2 * time.Millisecond,
				Retries:          1000,
				Logger:           log.Default(),
				EarlyExitOnError: false,
			},
			retryFunc: func() (bool, error) {
				return false, nil
			},
			// retryOptions.Timeout is less than retries * delay or context.Timeout
			expectedError: retry.TimeoutError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			retryCtx, cancel := context.WithTimeout(context.Background(), testCase.contextTimeout)
			defer cancel()

			err := retry.Do(retryCtx, testCase.retryFunc, testCase.retryOptions)

			assert.ErrorIs(t, err, testCase.expectedError)
		})
	}
}

func TestRetryFunctionSuccess(t *testing.T) {

	testCases := []int{1, 10, 100} // the varying number of attempts before retry function returns success

	for _, expectedAttempts := range testCases {
		t.Run(
			fmt.Sprintf("Should terminate correctly when it ends successfully after retrying %d times", expectedAttempts),
			func(t *testing.T) {
				var actualAttempt = 0
				err := retry.Do(context.Background(), func() (bool, error) {
					actualAttempt++
					return actualAttempt == expectedAttempts, nil
				}, retry.Options{
					Timeout:          time.Second,
					Delay:            0,
					Retries:          expectedAttempts,
					Logger:           log.Default(),
					EarlyExitOnError: true,
				})

				assert.Nil(t, err)
				assert.Equal(t, expectedAttempts, actualAttempt)
			})
	}
}

func TestRetryFunctionError(t *testing.T) {

	testCases := []int{1, 10, 100} // the varying number of attempts before retry function returns error

	for _, errorAttempt := range testCases {
		t.Run(
			fmt.Sprintf("Should ends correctly when error occur after retrying %d times", errorAttempt),
			func(t *testing.T) {
				var actualAttempt = 0
				err := retry.Do(
					context.Background(),
					func() (bool, error) {
						actualAttempt++
						if actualAttempt < errorAttempt {
							return false, nil
						} else {
							return false, fmt.Errorf("error occured")
						}
					},
					retry.Options{
						Timeout:          time.Second,
						Delay:            0,
						Retries:          errorAttempt,
						Logger:           log.Default(),
						EarlyExitOnError: true,
					})

				assert.NotNil(t, err)
				assert.Equal(t, errorAttempt, actualAttempt)
			})
	}
}

func TestRetryEarlyExitOnError(t *testing.T) {

	// number of times to return errors
	testCases := []int{1, 10, 100}
	// setup errors so that every time the retry func will return different error
	testErrors := make([]error, 100)
	for i := 0; i < cap(testErrors); i++ {
		testErrors[i] = fmt.Errorf("error %d", i)
	}

	for _, attempts := range testCases {
		errorIndex := 0
		t.Run(
			fmt.Sprintf("Should return accumulated errors after retrying %d times", attempts),
			func(t *testing.T) {
				err := retry.Do(
					context.Background(),
					func() (bool, error) {
						defer func() { errorIndex++ }()
						return false, testErrors[errorIndex] // return different errors
					},
					retry.Options{
						Timeout:          time.Second,
						Delay:            0,
						Retries:          attempts - 1,
						Logger:           log.Default(),
						EarlyExitOnError: false,
					})

				assert.NotNil(t, err)
				for i := 0; i < attempts; i++ {
					// assert errors are indeed returned
					assert.ErrorIs(t, err, testErrors[i])
				}
			})
	}
}

func TestRetryShouldAttemptOnceBeforeRetry(t *testing.T) {

	attempts := 0
	err := retry.Do(
		context.Background(),
		func() (bool, error) {
			attempts++
			return false, nil
		},
		retry.Options{
			Timeout:          time.Second,
			Delay:            0,
			Retries:          0, // no retries will happen
			Logger:           log.Default(),
			EarlyExitOnError: false,
		})

	assert.ErrorIs(t, err, retry.RetriesExceededError)
	assert.Equal(t, 1, attempts)
}

func TestRetryShouldTimeoutIfWillTimeoutAfterDelay(t *testing.T) {

	attempts := 0
	startTime := time.Now()
	err := retry.Do(
		context.Background(),
		func() (bool, error) {
			attempts++
			return false, nil
		},
		retry.Options{
			Timeout:          500 * time.Millisecond, // smaller than delay so that it will terminate early
			Delay:            time.Second,
			Retries:          1,
			Logger:           log.Default(),
			EarlyExitOnError: false,
		})

	assert.ErrorIs(t, err, retry.TimeoutError)
	assert.Less(t, time.Since(startTime), 300*time.Millisecond) // time since start should be smaller than delay, to assure that no delay actually took place
	assert.Equal(t, 1, attempts)
}

func TestRetryShouldTimeoutOnContextCancel(t *testing.T) {

	attempts := 0
	startTime := time.Now()
	testContext, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	err := retry.Do(
		testContext,
		func() (bool, error) {
			// cancel on first attempt
			cancel()
			attempts++
			return false, nil
		},
		retry.Options{
			Timeout:          time.Second,
			Delay:            time.Second,
			Retries:          1,
			Logger:           log.Default(),
			EarlyExitOnError: false,
		})

	assert.ErrorIs(t, err, retry.TimeoutError)
	assert.Less(t, time.Since(startTime), 300*time.Millisecond) // time since start should be smaller than delay and retry timeout, to assure that no delay actually took place.
	assert.Equal(t, 1, attempts)                                // only one attempt occurred because context is cancelled during first attempt.
}
