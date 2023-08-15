package state

const (
	DONE  = "DONE"
	ERROR = "ERROR"

	// asa/ios device only
	BAD_CREDENTIALS                   = "BAD_CREDENTIALS"
	PRE_WAIT_FOR_USER_TO_UPDATE_CREDS = "$PRE_WAIT_FOR_USER_TO_UPDATE_CREDS"
	PRE_READ_METADATA                 = "$PRE_READ_METADATA"
)
