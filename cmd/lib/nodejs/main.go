package main

import "github.com/digitranslab/kozmo-sandbox/internal/core/lib/nodejs"
import "C"

//export KozmoSeccomp
func KozmoSeccomp(uid int, gid int, enable_network bool) {
	nodejs.InitSeccomp(uid, gid, enable_network)
}

func main() {}
