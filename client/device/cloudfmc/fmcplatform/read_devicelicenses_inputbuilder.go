package fmcplatform

type ReadDeviceLicensesInputBuilder struct {
	readDeviceLicensesInput *ReadDeviceLicensesInput
}

func NewReadDeviceLicensesInputBuilder() *ReadDeviceLicensesInputBuilder {
	readDeviceLicensesInput := &ReadDeviceLicensesInput{}
	b := &ReadDeviceLicensesInputBuilder{readDeviceLicensesInput: readDeviceLicensesInput}
	return b
}

func (b *ReadDeviceLicensesInputBuilder) FmcHost(fmcHost string) *ReadDeviceLicensesInputBuilder {
	b.readDeviceLicensesInput.FmcHost = fmcHost
	return b
}

func (b *ReadDeviceLicensesInputBuilder) Build() ReadDeviceLicensesInput {
	return *b.readDeviceLicensesInput
}
