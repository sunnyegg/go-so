package cron

import "github.com/robfig/cron"

type Cron struct {
	cron *cron.Cron
}

func NewCron() *Cron {
	cron := cron.New()
	return &Cron{
		cron: cron,
	}
}

func (cron *Cron) Start() {
	cron.cron.Start()
}

func (cron *Cron) AddFunc(spec string, cmd func()) {
	cron.cron.AddFunc(spec, cmd)
}
