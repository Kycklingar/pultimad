package yp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const (
	yiffPartyURL   = "https://yiff.party"
	renderPostsURL = yiffPartyURL + "/render_posts"
)

func NewCreator(id int) (*Creator, error) {
	var c = new(Creator)
	c.ID = id

	if Creators2 == nil {
		return nil, errors.New("There are no creator names. Did you forget to LoadCreators()?")
	}

	c.Name = Creators2[id]

	return c, nil
}

type Creator struct {
	ID   int
	Name string

	currentPage int
	pageCount   int
}

func (c *Creator) Next() ([]*Post, error) {
	if c.currentPage > c.pageCount {
		return nil, nil
	}
	var v = url.Values{}
	v.Set("s", "patreon")
	v.Set("c", strconv.Itoa(c.ID))
	v.Set("p", strconv.Itoa(c.currentPage))
	c.currentPage++

	rc, err := c.downloadPage(renderPostsURL, v)
	if err != nil {
		c.currentPage--
		return nil, err
	}
	defer rc.Close()

	posts, err := c.parsePosts(rc)

	return posts, err
}

func (c Creator) SharedFiles() ([]Shared, error) {
	jsonStream, err := c.downloadPage(fmt.Sprint(yiffPartyURL, "/", c.ID, ".json"), nil)
	if err != nil {
		return nil, err
	}
	defer jsonStream.Close()
	files, err := parseSharedFilesJson(jsonStream)
	for i, _ := range files {
		files[i].CreatorID = c.ID
	}

	return files, err
}

func (c *Creator) parsePosts(body io.Reader) ([]*Post, error) {
	n, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	var posts []*Post

	nd := &node{*n}
	pgnode := nd.getClasses("paginate-count")
	if len(pgnode) > 0 {
		pgn := strings.Split(pgnode[0].n.FirstChild.Data, "/")
		if len(pgn) < 2 {
			return nil, errors.New("could not find pageinate-count")
		}

		c.pageCount, err = strconv.Atoi(strings.Trim(pgn[1], " "))
		if err != nil {
			return nil, err
		}
	}

	cc := nd.getClasses("row yp-posts-row")
	if len(cc) <= 0 {
		return nil, nil
	}

	for child := cc[0].n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			post, err := parsePost(&node{*child})
			if err != nil {
				log.Println(err)
				continue
			}
			post.Creator = c
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	}

	return posts, nil
}

var client *http.Client

func (c *Creator) downloadPage(page string, v url.Values) (io.ReadCloser, error) {
	if client == nil {
		client = new(http.Client)
	}

	query := ""
	if v != nil {
		query = "?" + v.Encode()
	}

	resp, err := http.Get(page + query)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		var buf = new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		return nil, errors.New(buf.String())
	}

	return resp.Body, nil
}
