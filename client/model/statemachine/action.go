package statemachine

type Action struct {
	ActionIdentifier string `json:"actionIdentifier"`
	EndDate          int    `json:"endDate"`
	EndMessage       string `json:"endMessage"`
	EndMessageString string `json:"endMessageString"`
	EndState         string `json:"endState"`
	Identifier       string `json:"identifier"`
	StartDate        int    `json:"startDate"`
	StartState       string `json:"startState"`

	CompleteStackTrace *string `json:"completeStackTrace,omitempty"`
	ErrorCode          *string `json:"errorCode,omitempty"`
	ErrorMessage       *string `json:"errorMessage,omitempty"`
	FurtherDetails     *string `json:"furtherDetails,omitempty"`
}
