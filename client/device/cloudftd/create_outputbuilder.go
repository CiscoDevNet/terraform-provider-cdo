package cloudftd

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/publicapilabels"
)

type CreateOutputBuilder struct {
	createOutput *CreateOutput
}

func NewCreateOutputBuilder() *CreateOutputBuilder {
	createOutput := &CreateOutput{}
	b := &CreateOutputBuilder{createOutput: createOutput}
	return b
}

func (b *CreateOutputBuilder) Uid(uid string) *CreateOutputBuilder {
	b.createOutput.Uid = uid
	return b
}

func (b *CreateOutputBuilder) Name(name string) *CreateOutputBuilder {
	b.createOutput.Name = name
	return b
}

func (b *CreateOutputBuilder) Metadata(metadata Metadata) *CreateOutputBuilder {
	b.createOutput.Metadata = metadata
	return b
}

func (b *CreateOutputBuilder) Labels(labels publicapilabels.Type) *CreateOutputBuilder {
	b.createOutput.Labels = labels
	return b
}

func (b *CreateOutputBuilder) Build() CreateOutput {
	return *b.createOutput
}
