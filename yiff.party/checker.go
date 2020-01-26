package yp

import (
	"github.com/kycklingar/pultimad/daemon"
	"github.com/kycklingar/pultimad/yiff.party/db"
	yp "github.com/kycklingar/pultimad/yiff.party/parser"
)

type Checker struct {
	db *db.DB
}

func (c *Checker) Init(config string) error {
	var err error

	c.db, err = db.Connect(config)
	if err != nil {
		return err
	}

	if err := yp.LoadCreators(); err != nil {
		return err
	}

	return nil
}

func (c *Checker) Check() []daemon.Taskif {
	creators, err := c.db.GetCreators(5)
	if err != nil {
		return nil
	}

	var tasks []daemon.Taskif

	for _, creator := range creators {
		var t = &cpTask{
			creator: creator,
			db:      c.db,
		}

		tasks = append(tasks, t)
	}

	files, err := c.db.GetFiles(20)
	if err != nil {
		return nil
	}

	for _, file := range files {
		var t = &dlTask{
			file: file,
			db:   c.db,
		}

		tasks = append(tasks, t)
	}

	return tasks
}

func (c *Checker) Quit() {
	c.db.Close()
}
