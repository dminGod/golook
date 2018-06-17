package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
)

var showTestFiles, showImports, showStructs, shortStructs,
showStructMethods, showPublicFuncs, showPrivateFuncs, limitedView, showFuncs bool

var dir string

//var pkgInfo PkgFuncs

var FileList []FuncInfo

var app Application

func main(){

	getFlags()
	//pkgInfo = make(map[string]Funcs)

	setDefaults()
	app = FetchApp( dir )

	StartUI()
}

func StartUI(){

	r := gin.Default()

	// b, _ := json.Marshal(app)

	r.GET("/show_app", showApp)
	r.GET("/json_file_size", jsonFS)
	r.GET("/json_pkg_size", jsonPackageSize)
	r.GET("/json_func_count", jsonMethodFuncs)

	r.Static("/static", "./static")
	r.Run("0.0.0.0:8081") // listen and serve on 0.0.0.0:8080
}

// Set the defaults if you are running it directly from the IDE
func setDefaults(){

	var err error

	dir = "/opt/go_code/src/github.com/docker/cli"
	//dir = "/opt/go_code/src/github.com/kubernetes/kubernetes"
	// dir = "/opt/go_code/src/github.com/dminGod/D30-HectorDA"

	// If no dir passed then use the cwd
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			fmt.Println("No directory specified for --dir and could not get the present working directory")
			os.Exit(1)
		}
		fmt.Printf("No --dir specified, using current directory '%v' to look", dir)
	}

	if limitedView == false {
		showStructs, showStructMethods, showFuncs, showPublicFuncs, showPrivateFuncs, showImports = true, true, true, true, true, true
	} else {
		if (showStructs || showStructMethods ||  showFuncs ||  showPublicFuncs || showPrivateFuncs || shortStructs || showImports) == false {
			fmt.Println("The short_view tag specified but no other filter tags like (show_struct, show_funcs etc.) given")
			os.Exit(1)
		}
	}
}
