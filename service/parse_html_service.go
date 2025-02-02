package service

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type ParseHTMLService interface {
	Parse(htm string) (*html.Node, error)
	FindAll(
		node *html.Node,
		f func(node *html.Node) bool,
	) []*html.Node
	FindFirst(
		node *html.Node,
		f func(node *html.Node) bool,
	) *html.Node
	HasClass(node *html.Node, name string) bool
	IsTag(node *html.Node, tagName string) bool
	AttrValRegExp(node *html.Node, attrName, attrValRegExp string) bool
	GetAttrValue(node *html.Node, attrName string) string
}

type parseHTMLService struct{}

func NewParseHTMLService() ParseHTMLService {
	return &parseHTMLService{}
}

func (s parseHTMLService) Parse(htm string) (*html.Node, error) {
	reader := strings.NewReader(htm)
	return html.Parse(reader)
}

func (s parseHTMLService) FindFirst(
	node *html.Node,
	f func(node *html.Node) bool,
) *html.Node {
	if node == nil {
		return nil
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if f(child) {
			return child
		}
		found := s.FindFirst(child, f)
		if found != nil {
			return found
		}
	}

	return nil
}

func (s parseHTMLService) FindAll(
	node *html.Node,
	f func(node *html.Node) bool,
) (found []*html.Node) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if f(child) {
			found = append(found, child)
		}
		foundChildren := s.FindAll(child, f)
		found = append(found, foundChildren...)
	}

	return
}
func (s parseHTMLService) HasClass(node *html.Node, name string) bool {
	attrs := node.Attr
	for _, attr := range attrs {
		if attr.Key == "class" {
			classes := strings.Split(attr.Val, " ")
			for _, class := range classes {
				if class == name {
					return true
				}
			}
		}
	}
	return false
}

func (s parseHTMLService) IsTag(node *html.Node, tagName string) bool {
	return node.Data == tagName
}

func (s parseHTMLService) AttrValRegExp(node *html.Node, attrName, attrValRegExp string) bool {
	attrs := node.Attr
	for _, attr := range attrs {
		if attr.Key == attrName {
			matched, err := regexp.MatchString(attrValRegExp, attr.Val)
			if err != nil {
				return false
			}

			return matched
		}
	}
	return false
}

func (s parseHTMLService) GetAttrValue(node *html.Node, attrName string) string {
	attrs := node.Attr
	for _, attr := range attrs {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}
