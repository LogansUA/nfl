package pagination

import (
	"net/url"
	"os"
	"strconv"
)

type Pagination struct {
	Page, Limit, Offset int
}

func (p *Pagination) Create(page, limit int) {
	l := limit
	if l <= 0 {
		l = 1
	}

	p.Limit = l
	p.Page = page
	p.Offset = (page-1)*l + 1

	return
}

func (p *Pagination) ParseParams(params url.Values) {
	var (
		page  int
		limit int
	)

	l, _ := strconv.Atoi(os.Getenv("PAGINATION_LIMIT"))

	limit = parseLimit(l)

	if pages, ok := params["page"]; ok {
		page, _ = strconv.Atoi(pages[0])
	} else {
		page = 1
	}

	if perPages, ok := params["per_page"]; ok {
		perPage, _ := strconv.Atoi(perPages[0])

		limit = parseLimit(perPage)
	}

	p.Create(page, limit)
}

func parseLimit(limit int) int {
	l := limit

	if l <= 0 {
		l = 1
	}

	return l
}
