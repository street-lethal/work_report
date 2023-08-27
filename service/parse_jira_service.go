package service

import (
	"golang.org/x/net/html"
)

type ParseJiraService interface {
	Parse(node *html.Node) map[int][]string
}

type parseJiraService struct {
	ParseHTMLService
}

func NewParseJiraService(hs ParseHTMLService) ParseJiraService {
	return &parseJiraService{hs}
}

func (s parseJiraService) Parse(node *html.Node) map[int][]string {
	works := map[int][]string{}
	tickets := s.ParseHTMLService.FindAll(node, func(n *html.Node) bool {
		return s.ParseHTMLService.HasClass(n, "fixedDataTableRowLayout_rowWrapper")
	})

	for _, ticket := range tickets {
		titleAndHours := s.ParseHTMLService.FindAll(ticket, func(n *html.Node) bool {
			return s.HasClass(n, "fixedDataTableCellGroupLayout_cellGroupWrapper")
		})

		titleTags := titleAndHours[0]
		hourTags := titleAndHours[1]

		titles := s.ParseHTMLService.FindAll(titleTags, func(n *html.Node) bool {
			return s.HasClass(n, "public_fixedDataTableCell_cellContent")
		})

		title := s.GetAttrValue(titles[1].FirstChild, "title")
		if title == "" {
			continue
		}

		hours := s.FindAll(hourTags, func(n *html.Node) bool {
			return s.HasClass(n, "public_fixedDataTableCell_cellContent")
		})

		for i, hour := range hours {
			if hour.FirstChild != nil {
				if works[i] == nil {
					works[i] = []string{}
				}
				works[i] = append(works[i], title)
			}
		}
	}

	return works
}
