package daemon

import (
	"fmt"
	"log"
	"sync"
	"time"

	yp "github.com/kycklingar/pultimad/parser"
	db "github.com/kycklingar/pultimad/storage"
)

const downloadFiles = true

func NewDaemon(connstr string) (*daemon, error) {
	var d = new(daemon)
	var err error
	d.b, err = db.Connect(connstr)
	if err != nil {
		return nil, err
	}

	d.q = new(queue)
	d.q.l = new(sync.Mutex)

	creators, err := d.b.LoadCreators()
	if err != nil {
		return nil, err
	}

	yp.PopulateCreators(creators)

	return d, nil
}

type daemon struct {
	q *queue
	b *db.DB
}

func (d *daemon) ReloadCreators() error {
	err := yp.LoadCreators()
	if err != nil {
		return err
	}

	for k, v := range yp.Creators2 {
		var c = new(yp.Creator)
		c.ID = k
		c.Name = v
		err = d.b.StoreCreator(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *daemon) QueueCreator(id int) {
	d.q.push(task{creatorTask, id})
}

func (d *daemon) checkForNewJobs() {
	if downloadFiles {
		files, err := d.b.GetFiles(5)
		if err != nil {
			log.Println(err)
		} else {
			for _, file := range files {
				d.q.push(task{fileTask, file})
			}
		}
	}

	creators, err := d.b.GetCreators(5)
	if err != nil {
		log.Println(err)
	} else {
		for _, creator := range creators {
			d.q.push(task{postsTask, creator})
		}
	}

}

func (d *daemon) Loop() {
	for {
		time.Sleep(time.Second * 2)
		tas := d.q.pop()
		switch tas.ttype {
		case sharedFilesTask:
			fmt.Print("Grabing shared files from ")
			creator, ok := tas.data.(*yp.Creator)
			if !ok {
				log.Println("task data not of type *yp.Creator", tas, tas.data)
				break
			}

			fmt.Println(creator.Name)
			files, err := creator.SharedFiles()
			if err != nil {
				log.Println(err)
				break
			}

			for _, file := range files {
				if err = d.b.StoreSharedFile(file); err != nil {
					log.Println(err)
				}
			}

			err = d.b.CheckCreator(creator)
			if err != nil {
				log.Println(err)
			}
		case postsTask:
			fmt.Print("Grabing posts from ")
			creator, ok := tas.data.(*yp.Creator)
			if !ok {
				log.Println("task data not of type *yp.Creator", tas, tas.data)
				break
			}

			fmt.Println(creator.Name)

			posts, err := creator.Next()
			if err != nil {
				log.Println(err)
				break
			}

			if posts == nil {
				d.q.push(task{sharedFilesTask, creator})
				break
			}

			for _, post := range posts {
				if err = d.b.StorePost(post); err != nil {
					log.Println(err)
				}
			}
			d.q.push(tas)
		case fileTask:
			fmt.Print("Downloading file ")
			file, ok := tas.data.(db.File)
			if !ok {
				log.Println("task data not of type db.File")
				break
			}
			fmt.Printf("[%s] %s\n", file.Creator, file.FileURL)
			if err := d.download(file); err != nil {
				d.b.Tried(file)
				log.Println(err)
			}
		default:
			d.checkForNewJobs()
			time.Sleep(time.Second * 5)
		}
	}
}
