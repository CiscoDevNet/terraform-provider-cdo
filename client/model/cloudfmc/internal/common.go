package internal

type Paging struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Pages  int `json:"pages"`
}

func NewPaging(
	count int,
	offset int,
	limit int,
	pages int,
) Paging {
	return Paging{
		Count:  count,
		Offset: offset,
		Limit:  limit,
		Pages:  pages,
	}
}

type Links struct {
	Self string `json:"self"`
}

func NewLinks(self string) Links {
	return Links{
		Self: self,
	}
}
