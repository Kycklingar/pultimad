package yp

import (
	"fmt"
	"log"

	"github.com/kycklingar/pultimad/daemon"
	"github.com/kycklingar/pultimad/yiff.party/db"
	yp "github.com/kycklingar/pultimad/yiff.party/parser"
)

// Creator posts task
type cpTask struct {
	creator *yp.Creator
	db      *db.DB
}

func (t cpTask) Domain() string { return ypDomain }

func (t cpTask) Description() string { return fmt.Sprintf("%s - Grabbing posts", t.creator.Name) }

func (t *cpTask) Do() []daemon.Taskif {
	posts, err := t.creator.Next()
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, post := range posts {
		err = t.db.StorePost(post)
		if err != nil {
			log.Println(err)
			return nil
		}
	}

	if len(posts) <= 0 {
		task := &sfTask{
			t.creator,
			t.db,
		}
		return []daemon.Taskif{task}
	}

	return []daemon.Taskif{t}

}
