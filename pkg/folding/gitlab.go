package folding

import (
	"fmt"
	"math/rand"
	"time"
)

type gitlabFolder struct {
}

func (g gitlabFolder) Start(display string, open bool) string {
	id := fmt.Sprintf("%08x\n", rand.Int31())
	collapsed := "true"
	if open {
		collapsed = "false"
	}
	fmt.Printf("\033[0Ksection_start:%d:%s[collapsed=%s]\r\033[0K%s\n",
		time.Now().Unix(),
		id,
		collapsed,
		display)
	return id
}

func (g gitlabFolder) Stop(id string) {
	fmt.Printf("\033[0Ksection_end:%d:%s\r\033[0K", time.Now().Unix(), id)
}
