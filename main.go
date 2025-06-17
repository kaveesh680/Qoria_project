package main

import (
	"abt-dashboard-api/internal"
	_ "net/http/pprof"
)

func main() {

	internal.Init()
}
