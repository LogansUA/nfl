package pagination

import (
	"net/url"
	"os"
	"strconv"
)

type Pagination struct {
	Page, Limit, Offset int
}

func New(params url.Values) Pagination {
	page, limit := parseParams(params)

	pagination := Pagination{}

	pagination.create(page, limit)

	return pagination
}

func (p *Pagination) create(page, limit int) {
	l := limit
	if l <= 0 {
		l = 1
	}

	p.Limit = l
	p.Page = page
	p.Offset = (page-1)*l + 1

	return
}

func parseParams(params url.Values) (page, limit int) {
	envLimit, _ := strconv.Atoi(os.Getenv("PAGINATION_LIMIT"))

	limit = parseLimit(envLimit)

	if pages, ok := params["page"]; ok {
		page, _ = strconv.Atoi(pages[0])
	} else {
		page = 1
	}

	if perPages, ok := params["per_page"]; ok {
		perPage, _ := strconv.Atoi(perPages[0])

		limit = parseLimit(perPage)
	}

	return
}

func parseLimit(limit int) int {
	l := limit

	if l <= 0 {
		l = 1
	}

	return l
}
