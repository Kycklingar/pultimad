package yp

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type Post struct {
	Creator     *Creator
	ID          int
	Title       string
	Body        string
	FileURL     string
	Attachments []*url.URL
}

func parsePost(node *node) (*Post, error) {
	var post = new(Post)
	err := post.parseID(node)
	if err != nil {
		return nil, err
	}
	err = post.parseTitle(node)
	if err != nil {
		return nil, err
	}
	err = post.parseBody(node)
	if err != nil {
		return nil, err
	}
	err = post.parseFile(node)
	if err != nil {
		return nil, err
	}
	err = post.parseAttachments(node)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *Post) parseID(n *node) error {
	pid := n.getFirstChildNode().getAttr("id")
	var err error
	p.ID, err = strconv.Atoi(pid[1:len(pid)])
	return err
}

func (p *Post) parseTitle(n *node) error {
	cl := n.getClasses("card-reveal")
	if len(cl) <= 0 {
		return nil
	}

	p.Title = cl[0].getFirstChildNode().n.FirstChild.Data
	return nil
}

func (p *Post) parseBody(n *node) error {
	body := &n.getClasses("post-body")[0].n //.FirstChild.FirstChild
	if body == nil {
		return errors.New("Body is nil")
	}
	var wb = new(strings.Builder)
	html.Render(wb, body)
	p.Body = wb.String()
	return nil
}

func (p *Post) parseFile(n *node) error {
	cl := n.getClasses("card-action")
	if len(cl) <= 0 {
		return nil
	}
	var err error
	url, err := url.Parse(cl[0].getFirstChildNode().getAttr("href"))
	p.FileURL = url.String()
	return err
}

func (p *Post) parseAttachments(n *node) error {
	cl := n.getClasses("card-attachments")
	for i, _ := range cl {
		ps := cl[i].getChildrenOfType("p")
		if len(ps) <= 0 {
			return nil
		}

		links := ps[0].getChildrenOfType("a")
		if len(links) <= 0 {
			return nil
		}

		for _, a := range links {
			u, err := url.Parse(a.getAttr("href"))
			if err != nil {
				return err
			}
			if u.Host == "" {
				u.Host = "yiff.party"
				u.Scheme = "https"
			}
			p.Attachments = append(p.Attachments, u)
		}
	}

	return nil
}
