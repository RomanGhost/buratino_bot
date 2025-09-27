package pagination

type Pagination struct {
	Page  int
	Limit int
}

func (p *Pagination) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return (p.Page - 1) * p.Limit
}
