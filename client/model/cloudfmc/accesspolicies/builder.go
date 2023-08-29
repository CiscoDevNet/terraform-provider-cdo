package accesspolicies

type Builder struct {
	accessPolicies *AccessPolicies
}

func NewAccessPoliciesBuilder() *Builder {
	accessPolicies := &AccessPolicies{}
	b := &Builder{accessPolicies: accessPolicies}
	return b
}

func (b *Builder) Items(items []Item) *Builder {
	b.accessPolicies.Items = items
	return b
}

func (b *Builder) Links(links Links) *Builder {
	b.accessPolicies.Links = links
	return b
}

func (b *Builder) Paging(paging Paging) *Builder {
	b.accessPolicies.Paging = paging
	return b
}

func (b *Builder) Build() AccessPolicies {
	return *b.accessPolicies
}
