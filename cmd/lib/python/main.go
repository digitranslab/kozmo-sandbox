package main

import (
	"github.com/digitranslab/kozmo-sandbox/internal/core/lib/python"
)
import "C"

//export KozmoSeccomp
func KozmoSeccomp(uid int, gid int, enable_network bool) {
	python.InitSeccomp(uid, gid, enable_network)
}

func main() {}
