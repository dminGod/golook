	package main

import (
	"io/ioutil"
	"path"
	"fmt"
)


// Recursively go to all sub-directories that are not starting with .
// and return back an array of string with the location of the folders

func GetAllSubDirs(base string) (s []string) {

	var foundCounter int

	// Get all the files from the base given
	l, err := ioutil.ReadDir( base )
	if err != nil {
		fmt.Println("Error during getting directory info", err)
		return
	}

	var t [][]string

	for _, b := range l {
		// if its a directory no .<name> dirs and vendor directory
		if b.IsDir() && b.Name()[:1] != "." && b.Name() != "vendor" {

			t = append(t, GetAllSubDirs( path.Join(base, b.Name()) ))
			foundCounter++
		}
	}

	// This is the last layer you need to return back the function
	if foundCounter == 0 {
		return []string{ base }
	}

	// Collect the directories from the [][]dirs
	for _, b := range t {
		for _, c := range b {
			s = append(s, c)
		}
	}

	s = append(s, base)
	return
}

// These functions may be useful for getting data out in some place
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


