package main

import (
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/cog"
)

var (
	log    *cog.Logger
	stderr *lorg.Log
)

func init() {
	stderr = lorg.NewLog()
	stderr.SetIndentLines(true)
	stderr.SetFormat(
		lorg.NewFormat("${time} ${level:[%s]:right:short} ${prefix}%s"),
	)

	log = cog.NewLogger(stderr)
}
