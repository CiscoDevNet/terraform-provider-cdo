package cloudfmc

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

func (b *ReadSpecificOutputBuilder) DomainUid(domainUid string) *ReadSpecificOutputBuilder {
	b.readSpecificOutput.DomainUid = domainUid
	return b
}

func (b *ReadSpecificOutputBuilder) State(state string) *ReadSpecificOutputBuilder {
	b.readSpecificOutput.State = state
	return b
}

func (b *ReadSpecificOutputBuilder) Status(status string) *ReadSpecificOutputBuilder {
	b.readSpecificOutput.Status = status
	return b
}

func (b *ReadSpecificOutputBuilder) Build() ReadSpecificOutput {
	return *b.readSpecificOutput
}
