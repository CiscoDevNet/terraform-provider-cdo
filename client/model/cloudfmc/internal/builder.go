package internal

type Builder struct {
	response *Response
}

func NewResponseBuilder() *Builder {
	b := &Builder{response: &Response{}}
	return b
}

func (b *Builder) Items(items []Item) *Builder {
	b.response.Items = items
	return b
}

func (b *Builder) Links(links Links) *Builder {
	b.response.Links = links
	return b
}

func (b *Builder) Paging(paging Paging) *Builder {
	b.response.Paging = paging
	return b
}

func (b *Builder) Build() Response {
	return *b.response
}
