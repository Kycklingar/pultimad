package yp

import (
	"log"

	"github.com/kycklingar/pultimad/config"
	"github.com/kycklingar/pultimad/daemon"
	"github.com/kycklingar/pultimad/yiff.party/db"
	parser "github.com/kycklingar/pultimad/yiff.party/parser"
)

type Checker struct {
	db *db.DB
}

func (c *Checker) Init(conf config.Config) error {
	cfg, ok := conf.(Config)
	if !ok {
		return config.ErrInvalid
	}

	var err error
	c.db, err = db.Connect(cfg.Connstr)
	if err != nil {
		return err
	}

	err = c.db.Ping()
	if err != nil {
		return err
	}

	creators, err := parser.LoadCreators()
	if err != nil {
		return err
	}

	for _, creator := range creators {
		c.db.StoreCreator(creator)
	}

	return nil
}

func (c *Checker) Check() []daemon.Taskif {
	creators, err := c.db.GetCreators(5)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
