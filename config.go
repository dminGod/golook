package main

import (
	"flag"
	"os"
	"fmt"
)

var showTestFiles, showImports, showStructs, shortStructs,
showStructMethods, showPublicFuncs, showPrivateFuncs, limitedView, showFuncs bool

// TODO: Add cobra for this and then serve it with that... just looks like a bunch of
	// configuraiotn optoins that are put in as a laundry list.
func getFlags() {

	flag.StringVar(&dir, "dir", "", "The directory of the code")
	flag.BoolVar(&showTestFiles, "use_test_files", false, "Also parses _test.go files if set to true -- its set to false by default")

	flag.BoolVar(&limitedView, "short_view", false, "To show limited results -- Shows only what you specify with other flags (like show_struct, show_public_funcs etc.) -- Set to false by default")
	flag.BoolVar(&showImports, "show_import", false, "Shows only imports used -- Set to true by default")
	flag.BoolVar(&showStructs, "show_struct", false, "Shows only the struct definitions -- Set to true by default")
	flag.BoolVar(&showStructMethods, "show_struct_methods", false, "Shows struct methods  -- Set to true by default")
	flag.BoolVar(&showFuncs, "show_funcs", false, "Shows all functions -- Set to true by default")
	flag.BoolVar(&showPublicFuncs, "show_funcs_public", false, "Shows public functions -- Set to true by default")
	flag.BoolVar(&showPrivateFuncs, "show_funcs_private", false, "Shows private functions -- Set to true by default")
	flag.BoolVar(&shortStructs, "short_struct", false, "Shows only the first line of the struct -- Set to false by default")

	flag.Parse()
}

// Set the defaults if you are running it directly from the IDE
func setDefaults(){

	var err error

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



