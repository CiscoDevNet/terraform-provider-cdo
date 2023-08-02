package asa

type ReadSpecificOutputBuilder struct {
	output ReadSpecificOutput
}

func NewReadSpecificOutputBuilder() *ReadSpecificOutputBuilder {
	return &ReadSpecificOutputBuilder{
		output: ReadSpecificOutput{
			Namespace: "asa",
			Type:      "configs",
		},
	}
}

func (builder *ReadSpecificOutputBuilder) Build() ReadSpecificOutput {
	return builder.output
}

func (builder *ReadSpecificOutputBuilder) WithSpecificUid(uid string) *ReadSpecificOutputBuilder {
	builder.output.SpecificUid = uid

	return builder
}

func (builder *ReadSpecificOutputBuilder) InDoneState() *ReadSpecificOutputBuilder {
	builder.output.State = "DONE"

	return builder
}
