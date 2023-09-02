package service

import (
	"strconv"
	"work_report/model"

	"golang.org/x/net/html"
)

type ParseJiraService interface {
	Parse(node *html.Node) map[int]model.Work
}

type parseJiraService struct {
	ParseHTMLService
}

func NewParseJiraService(hs ParseHTMLService) ParseJiraService {
	return &parseJiraService{hs}
}

func (s parseJiraService) Parse(node *html.Node) map[int]model.Work {
	works := map[int]model.Work{}
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
				day := i + 1
				work, ok := works[day]
				if !ok {
					work = model.Work{}
				}
				work.AddContent(title)

				floatHour, _ := strconv.ParseFloat(hour.FirstChild.Data, 64)
				work.AddHour(floatHour)

				works[day] = work
			}
		}
	}

	return works
}
