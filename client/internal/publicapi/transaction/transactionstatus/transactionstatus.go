package transactionstatus

type Type string

const (
	PENDING     Type = "PENDING"
	IN_PROGRESS Type = "IN_PROGRESS"
	DONE        Type = "DONE"
	ERROR       Type = "ERROR"
)
