package folding

import (
	log "github.com/sirupsen/logrus"
)

type structuredFolder struct {
	entry *log.Entry
}

func (s structuredFolder) Start(id string) {
	s.entry = log.WithFields(log.Fields{"id": id})
}

func (s structuredFolder) Stop() {
	s.entry = nil
}

func (s structuredFolder) Log() {
	// TODO
	s.entry.Info("log")
}
