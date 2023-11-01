package statemachine

// see also: https://github.com/cisco-lockhart/cdo-jvm-modules/blob/master/libs/lh-model/lh-model/src/main/java/com/cisco/lockhart/model/statemachine/model/StateMachineInstance.java
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

type InstanceBuilder struct {
	instance *Instance
}

func NewInstanceBuilder() *InstanceBuilder {
	instance := &Instance{}
	b := &InstanceBuilder{instance: instance}
	return b
}

func (b *InstanceBuilder) Actions(actions []Action) *InstanceBuilder {
	b.instance.Actions = actions
	return b
}

func (b *InstanceBuilder) ActiveStateMachineContext(activeStateMachineContext map[string]string) *InstanceBuilder {
	b.instance.ActiveStateMachineContext = activeStateMachineContext
	return b
}

func (b *InstanceBuilder) AfterHooks(afterHooks []Hook) *InstanceBuilder {
	b.instance.AfterHooks = afterHooks
	return b
}

func (b *InstanceBuilder) BeforeHooks(beforeHooks []Hook) *InstanceBuilder {
	b.instance.BeforeHooks = beforeHooks
	return b
}

func (b *InstanceBuilder) CreatedDate(createdDate int) *InstanceBuilder {
	b.instance.CreatedDate = createdDate
	return b
}

func (b *InstanceBuilder) CurrentState(currentState string) *InstanceBuilder {
	b.instance.CurrentState = currentState
	return b
}

func (b *InstanceBuilder) EndDate(endDate int) *InstanceBuilder {
	b.instance.EndDate = endDate
	return b
}

func (b *InstanceBuilder) HasErrors(hasErrors bool) *InstanceBuilder {
	b.instance.HasErrors = hasErrors
	return b
}

func (b *InstanceBuilder) ObjectReference(objectReference ObjectReference) *InstanceBuilder {
	b.instance.ObjectReference = objectReference
	return b
}

func (b *InstanceBuilder) StateMachineDetails(stateMachineDetails Details) *InstanceBuilder {
	b.instance.StateMachineDetails = stateMachineDetails
	return b
}

func (b *InstanceBuilder) StateMachineIdentifier(stateMachineIdentifier string) *InstanceBuilder {
	b.instance.StateMachineIdentifier = stateMachineIdentifier
	return b
}

func (b *InstanceBuilder) StateMachineInstanceCondition(stateMachineInstanceCondition string) *InstanceBuilder {
	b.instance.StateMachineInstanceCondition = stateMachineInstanceCondition
	return b
}

func (b *InstanceBuilder) StateMachinePriority(stateMachinePriority string) *InstanceBuilder {
	b.instance.StateMachinePriority = stateMachinePriority
	return b
}

func (b *InstanceBuilder) StateMachineType(stateMachineType string) *InstanceBuilder {
	b.instance.StateMachineType = stateMachineType
	return b
}

func (b *InstanceBuilder) Status(status string) *InstanceBuilder {
	b.instance.Status = status
	return b
}

func (b *InstanceBuilder) Uid(uid string) *InstanceBuilder {
	b.instance.Uid = uid
	return b
}

func (b *InstanceBuilder) Build() Instance {
	return *b.instance
}
