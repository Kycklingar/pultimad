package yp

import (
	"fmt"
	"log"

	"github.com/kycklingar/pultimad/daemon"
	"github.com/kycklingar/pultimad/yiff.party/db"
	yp "github.com/kycklingar/pultimad/yiff.party/parser"
)

type sfTask struct {
	creator *yp.Creator
	db      *db.DB
}

func (t sfTask) Domain() string { return ypDomain }

func (t sfTask) Description() string { return fmt.Sprintf("%s - Grabbing shared files", t.creator.Name) }

func (t *sfTask) Do() []daemon.Taskif {
	sharedFiles, err := t.creator.SharedFiles()
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, sf := range sharedFiles {
		err = t.db.StoreSharedFile(sf)
		if err != nil {
			log.Println(err)
			return nil
		}
	}

	err = t.db.CheckCreator(t.creator)
	if err != nil {
		log.Println(err)
	}

	return nil
}
