package instance

import (
	"github.com/lcabraja/APP-Project-LukaCabraja/communication/http"
	"github.com/lcabraja/APP-Project-LukaCabraja/configuration"
	"github.com/google/uuid"
	"time"
)

type App struct {
	c  *configuration.Configuration
	hm *http.HttpManager

	started     time.Time
	instances   map[string]*Instance
	instanceIdx []string
}

func NewApp(c *configuration.Configuration, useDefault bool) *App {
	app := &App{
		c:  c,
		hm: http.NewHttpManager(c.Prefix, c.Host()),

		started:     time.Now(),
		instances:   make(map[string]*Instance),
		instanceIdx: make([]string, 0),
	}

	if useDefault {
		app.setupDefaultInstance()
	}

	return app
}

func (a *App) Start() {
	a.hm.ListenAndServe()
}

func (a *App) createInstance(name string) (*Instance, string, error) {
	i, err := NewInstance(name)
	if err != nil {
		return nil, "", err
	}

	u := uuid.NewString()
	a.instances[u] = i
	a.instanceIdx = append(a.instanceIdx, u)
	return i, u, nil
}
