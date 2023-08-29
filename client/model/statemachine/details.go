package statemachine

type Details struct {
	CurrentDataRequirements       *string `json:"currentDataRequirements,omitempty"`
	Identifier                    string  `json:"identifier"`
	LastActiveDate                int     `json:"lastActiveDate"`
	LastStep                      string  `json:"lastStep"`
	Priority                      string  `json:"priority"`
	StateMachineInstanceCondition string  `json:"stateMachineInstanceCondition"`
	StateMachineType              string  `json:"stateMachineType"`
	Uid                           string  `json:"uid"`
}
