package yp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/html"
)

var creatorsURL = yiffPartyURL + "/creators"

func LoadCreators() ([]Creator, error) {
	res, err := http.Get(creatorsURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var buf bytes.Buffer
		buf.ReadFrom(res.Body)
		return nil, errors.New(buf.String())
	}

	node, err := html.Parse(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var out = make(chan *html.Node)
	go getNodeClassesFull(node, "col s12", out)

	var i int
	for {
		n, done := <-out
		if !done {
			break
		}
		if i == 0 {
			node = n
			break
		}
		i++
	}

	var creators []Creator

	out = make(chan *html.Node)
	defer close(out)

	go getNodeElements(node, "a", out)
	for {
		n := <-out
		if n == nil {
			break
		}

		var c Creator
		for _, a := range n.Attr {
			if a.Key == "href" {
				c.ID, err = strconv.Atoi(a.Val[1:])
				if err != nil {
					return nil, err
				}
			}
		}
		if n.FirstChild != nil {
			c.Name = n.FirstChild.Data
		}

		creators = append(creators, c)
	}

	return creators, nil
}

func getNodeElements(node *html.Node, element string, ch chan *html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == element {
			ch <- n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	ch <- nil
}

func getNodeClassesFull(node *html.Node, class string, ch chan *html.Node) {
	defer close(ch)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "class" {
					if a.Val == class {
						ch <- n
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
}

func getJson(url string, out interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	js, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	if res.StatusCode != 200 {
		log.Println(url)
		log.Println(js)
		return fmt.Errorf(string(js))
	}

	err = json.Unmarshal(js, &out)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
