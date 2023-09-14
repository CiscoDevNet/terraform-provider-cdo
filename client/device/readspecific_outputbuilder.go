package device

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

func (b *ReadSpecificOutputBuilder) State(state string) *ReadSpecificOutputBuilder {
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
