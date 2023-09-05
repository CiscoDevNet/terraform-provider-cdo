package cloudftd

type UpdateSpecificFtdOutputBuilder struct {
	updateSpecificFtdOutput *UpdateSpecificFtdOutput
}

func NewUpdateSpecificFtdOutputBuilder() *UpdateSpecificFtdOutputBuilder {
	updateSpecificFtdOutput := &UpdateSpecificFtdOutput{}
	b := &UpdateSpecificFtdOutputBuilder{updateSpecificFtdOutput: updateSpecificFtdOutput}
	return b
}

func (b *UpdateSpecificFtdOutputBuilder) SpecificUid(specificUid string) *UpdateSpecificFtdOutputBuilder {
	b.updateSpecificFtdOutput.SpecificUid = specificUid
	return b
}

func (b *UpdateSpecificFtdOutputBuilder) Build() UpdateSpecificFtdOutput {
	return *b.updateSpecificFtdOutput
}
