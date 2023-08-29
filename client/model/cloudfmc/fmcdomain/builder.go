package fmcdomain

type InfoBuilder struct {
	info *Info
}

func NewInfoBuilder() *InfoBuilder {
	info := &Info{}
	b := &InfoBuilder{info: info}
	return b
}

func (b *InfoBuilder) Links(links Links) *InfoBuilder {
	b.info.Links = links
	return b
}

func (b *InfoBuilder) Paging(paging Paging) *InfoBuilder {
	b.info.Paging = paging
	return b
}

func (b *InfoBuilder) Items(items []Item) *InfoBuilder {
	b.info.Items = items
	return b
}

func (b *InfoBuilder) Build() Info {
	return *b.info
}
