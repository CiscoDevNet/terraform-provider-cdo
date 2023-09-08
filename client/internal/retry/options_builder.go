package retry

import (
	"log"
	"time"
)

type OptionsBuilder struct {
	options *Options
}

func NewOptionsBuilder() *OptionsBuilder {
	options := &Options{}
	b := &OptionsBuilder{options: options}
	return b
}

func (b *OptionsBuilder) Timeout(timeout time.Duration) *OptionsBuilder {
	b.options.Timeout = timeout
	return b
}

func (b *OptionsBuilder) Delay(delay time.Duration) *OptionsBuilder {
	b.options.Delay = delay
	return b
}

func (b *OptionsBuilder) Retries(retries int) *OptionsBuilder {
	b.options.Retries = retries
	return b
}

func (b *OptionsBuilder) Logger(logger *log.Logger) *OptionsBuilder {
	b.options.Logger = logger
	return b
}

func (b *OptionsBuilder) EarlyExitOnError(earlyExitOnError bool) *OptionsBuilder {
	b.options.EarlyExitOnError = earlyExitOnError
	return b
}

func (b *OptionsBuilder) Build() Options {
	return *b.options
}
