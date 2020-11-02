package main

import (
	"github.com/dminGod/golook/parsers"
	"github.com/dminGod/golook/web"
)

var FileList []parsers.FuncInfo

var app parsers.Application

func main() {

	parsers.GetFlags()

	parsers.InitLocalDB()
	//runDebug()
	parsers.SetDefaults()

	// Start the UI frontend
	web.StartUI()
}

func runDebug() {

}
