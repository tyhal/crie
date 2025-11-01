//go:build linux
// +build linux

// If on linux then try to increase file limits

package cli

import (
	"math"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/tyhal/crie/pkg/errchain"
)

var maxFilesPerRoutine = 5

func convertLimit(limit uint64) int {
	var out int
	if limit > math.MaxInt32 {
		out = math.MaxInt32
	} else {
		out = int(limit)
	}
	return out / maxFilesPerRoutine
}

// MaxConcurrency returns the max concurrency name
func (e *LintCli) MaxConcurrency() int {
	var limit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit)
	if err != nil {
		log.Fatal(errchain.From(err).Link("getting limit for max concurrency"))
	}

	oldLimit := limit.Cur
	limit.Cur = limit.Max
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limit)
	if err != nil {
		log.Debug(err)
		return convertLimit(oldLimit)
	}

	return convertLimit(limit.Cur)
}
