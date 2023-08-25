package statemachine

type Instance struct {
	Actions                       []Action          `json:"actions"`
	ActiveStateMachineContext     map[string]string `json:"activeStateMachineContext"`
	AfterHooks                    []Hook            `json:"afterHooks"`
	BeforeHooks                   []Hook            `json:"beforeHooks"`
	CreatedDate                   int               `json:"createdDate"`
	CurrentState                  string            `json:"currentState"`
	EndDate                       int               `json:"endDate"`
	HasErrors                     bool              `json:"hasErrors"`
	ObjectReference               ObjectReference   `json:"objectReference"`
	StateMachineDetails           Details           `json:"stateMachineDetails"`
	StateMachineIdentifier        string            `json:"stateMachineIdentifier"`
	StateMachineInstanceCondition string            `json:"stateMachineInstanceCondition"`
	StateMachinePriority          string            `json:"stateMachinePriority"`
	StateMachineType              string            `json:"stateMachineType"`
	Status                        string            `json:"status"`
	Uid                           string            `json:"uid"`
}
