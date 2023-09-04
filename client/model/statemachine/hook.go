package statemachine

type Hook struct {
	EndDate        int    `json:"endDate"`
	EndMessage     string `json:"endMessage"`
	HookIdentifier string `json:"hookIdentifier"`
	Identifier     string `json:"identifier"`
	StartDate      int    `json:"startDate"`

	CompleteStackTrace *string `json:"completeStackTrace,omitempty"`
	ErrorMessage       *string `json:"errorMessage,omitempty"`
}
