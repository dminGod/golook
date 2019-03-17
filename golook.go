package main


var dir string

//var pkgInfo PkgFuncs

var FileList []FuncInfo

var app Application

func main(){

	getFlags()
	setDefaults()

	app = FetchApp( dir )

	// Start the UI frontend
	StartUI()
}

func collectStats(){




}




