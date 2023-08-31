package fmcplatform

type ReadDomainInfoInputBuilder struct {
	readDomainInfoInput *ReadDomainInfoInput
}

func NewReadDomainInfoInputBuilder() *ReadDomainInfoInputBuilder {
	readDomainInfoInput := &ReadDomainInfoInput{}
	b := &ReadDomainInfoInputBuilder{readDomainInfoInput: readDomainInfoInput}
	return b
}

func (b *ReadDomainInfoInputBuilder) FmcHost(fmcHost string) *ReadDomainInfoInputBuilder {
	b.readDomainInfoInput.FmcHost = fmcHost
	return b
}

func (b *ReadDomainInfoInputBuilder) Build() ReadDomainInfoInput {
	return *b.readDomainInfoInput
}
