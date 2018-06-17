package main

import (
	"io/ioutil"
	"path"
	"flag"
	"fmt"
)

func retLines(r []string, par int, end int) (ret []string){

	s := len(r)

	if par <= s && end <= s {
		ret = r[par:end]
	}

	return
}

func retLine(r []string, parm int) (ret string){

	s := len(r)

	if parm <= s {
		ret = r[parm]
	}

	return
}

func GetAllSubDirs(base string) (s []string) {

	// Get
	l, err := ioutil.ReadDir( base )
	if err != nil {
		fmt.Println("Error during getting directory info", err)
		return
	}

	i := 0

	var t [][]string

	for _, b := range l {
		// if its a directory no .<name> dirs and vendor directory
		if b.IsDir() && b.Name()[:1] != "." && b.Name() != "vendor" {
			t = append(t, GetAllSubDirs( path.Join(base, b.Name()) ))
			i++
		}
	}

	// If you dont have anything still just return the base
	if i == 0 {
		return []string{ base }
	}

	// Flatten the array
	for _, b := range t {
		for _, c := range b {
			s = append(s, c)
		}
	}

	s = append(s, base)

	return
}

// Get app passed flags
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