package main

var dir string


var FileList []FuncInfo

var app Application

func main(){

	getFlags()

	initLocalDB()
	//runDebug()
	setDefaults()

	// Start the UI frontend
	StartUI()
}

func runDebug(){




}




