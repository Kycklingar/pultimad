package daemon

import (
	"os"
	"time"
)

func NewDaemon() *daemon {
	var d = new(daemon)

	return d
}

type daemon struct {
	domains []domain
}

func (d *daemon) RegisterDomain(name string, checker Checker, sleepTime time.Duration) {
	var dom = domain{
		checker:   checker,
		name:      name,
		queue:     newQueue(),
		sleepTime: sleepTime,
	}

	d.domains = append(d.domains, dom)
}

func (d *daemon) Loop(quit chan os.Signal) {
	for _, dom := range d.domains {
		go dom.process()
	}

	<-quit
	d.cleanup()
}

func (d *daemon) cleanup() {
	for _, dom := range d.domains {
		dom.quit()
	}
}
