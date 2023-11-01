package statemachine

type Details struct {
	CurrentDataRequirements       *[]string `json:"currentDataRequirements,omitempty"`
	Identifier                    string    `json:"identifier"`
	LastActiveDate                int       `json:"lastActiveDate"`
	LastStep                      string    `json:"lastStep"`
	Priority                      string    `json:"priority"`
	StateMachineInstanceCondition string    `json:"stateMachineInstanceCondition"`
	StateMachineType              string    `json:"stateMachineType"`
	Uid                           string    `json:"uid"`
	LastError                     *Error    `json:"lastError,omitempty"`
}

type Error struct {
	StateMachineType       string  `json:"stateMachineType"`
	ErrorMessage           string  `json:"errorMessage"`
	StateMachineIdentifier string  `json:"stateMachineIdentifier"`
	ActionIdentifier       string  `json:"actionIdentifier"`
	EndState               *string `json:"endState,omitempty"`
	ErrorCode              *string `json:"errorCode,omitempty"`
	ErrorDate              int64   `json:"errorDate"`
}
