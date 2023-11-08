package state

type Type string

const (
	DONE  Type = "DONE"
	ERROR Type = "ERROR"

	// asa/ios device only
	BAD_CREDENTIALS                   Type = "BAD_CREDENTIALS"
	WAIT_FOR_USER_TO_UPDATE_CREDS     Type = "WAIT_FOR_USER_TO_UPDATE_CREDS"
	PRE_WAIT_FOR_USER_TO_UPDATE_CREDS Type = "$PRE_WAIT_FOR_USER_TO_UPDATE_CREDS"
	PRE_READ_METADATA                 Type = "$PRE_READ_METADATA"
)
