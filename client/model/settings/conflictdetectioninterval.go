package settings

import (
	"encoding/json"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/sliceutil"
)

type ConflictDetectionInterval string

const (
	ConflictDetectionIntervalEvery10Minutes ConflictDetectionInterval = "EVERY_10_MINUTES"
	ConflictDetectionIntervalEveryHour      ConflictDetectionInterval = "EVERY_HOUR"
	ConflictDetectionIntervalEvery6Hours    ConflictDetectionInterval = "EVERY_6_HOURS"
	ConflictDetectionIntervalEvery24Hours   ConflictDetectionInterval = "EVERY_24_HOURS"
	ConflictDetectionIntervalInvalid        ConflictDetectionInterval = ""
)

func conflictDetectionIntervalValues() []ConflictDetectionInterval {
	return []ConflictDetectionInterval{
		ConflictDetectionIntervalEvery10Minutes,
		ConflictDetectionIntervalEveryHour,
		ConflictDetectionIntervalEvery6Hours,
		ConflictDetectionIntervalEvery24Hours,
	}
}

func ResolveConflictDetectionInterval(input string) ConflictDetectionInterval {
	interval := ConflictDetectionInterval(input)
	if sliceutil.Contains(conflictDetectionIntervalValues(), interval) {
		return interval
	}

	return ConflictDetectionIntervalInvalid
}

func (interval *ConflictDetectionInterval) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*interval = ResolveConflictDetectionInterval(value)
	if *interval == ConflictDetectionIntervalInvalid {
		return fmt.Errorf(fmt.Sprintf("invalid conflict detection interval value, got: %s", *interval))
	}

	return nil
}

func (interval ConflictDetectionInterval) String() string {
	return string(interval)
}
