package device

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"

type ReadSpecificOutputBuilder struct {
	readSpecificOutput *ReadSpecificOutput
}

func NewReadSpecificOutputBuilder() *ReadSpecificOutputBuilder {
	readSpecificOutput := &ReadSpecificOutput{}
	b := &ReadSpecificOutputBuilder{readSpecificOutput: readSpecificOutput}
	return b
}

func (b *ReadSpecificOutputBuilder) SpecificUid(specificUid string) *ReadSpecificOutputBuilder {
	b.readSpecificOutput.SpecificUid = specificUid
	return b
}

func (b *ReadSpecificOutputBuilder) State(state state.Type) *ReadSpecificOutputBuilder {
	b.readSpecificOutput.State = state
	return b
}

func (b *ReadSpecificOutputBuilder) Namespace(namespace string) *ReadSpecificOutputBuilder {
	b.readSpecificOutput.Namespace = namespace
	return b
}

func (b *ReadSpecificOutputBuilder) Type(type_ string) *ReadSpecificOutputBuilder {
	b.readSpecificOutput.Type = type_
	return b
}

func (b *ReadSpecificOutputBuilder) Build() ReadSpecificOutput {
	return *b.readSpecificOutput
}
