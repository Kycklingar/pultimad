package yp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var Creators2 map[int]string

func PopulateCreators(creators []*Creator) {
	Creators2 = make(map[int]string)
	fmt.Println("Populating creators2")
	for _, c := range creators {
		//fmt.Println(c.ID, c.Name)
		Creators2[c.ID] = c.Name
	}
}

func LoadCreators() error {
	type creator struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	type Creators struct {
		Creators []creator `json:"creators"`
	}

	var m Creators

	err := getJson("https://yiff.party/json/creators.json", &m)
	if err != nil {
		return err
	}

	Creators2 = make(map[int]string)

	fmt.Println(m)
	for _, c := range m.Creators {
		fmt.Println(c.Id, c.Name)
		Creators2[c.Id] = c.Name
	}

	return nil
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
