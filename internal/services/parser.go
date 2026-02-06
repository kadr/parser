package services

import "parser/pkg/sites"

type Site interface {
	Parse(categoryName string, limit *int) ([]sites.Result, error)
}

type ParserService struct {
	site Site
}

func NewService(site Site) ParserService {
	return ParserService{site: site}
}

func (p *ParserService) Parse(categoryName string, limit *int) ([]sites.Result, error) {
	return p.site.Parse(categoryName, limit)
}
