package parsers

import (
	"flag"
	"fmt"
	"os"
)

var showTestFiles bool

var dir string

func GetFlags() {

	flag.StringVar(&dir, "dir", "", "The directory of the code")
	flag.BoolVar(&showTestFiles, "use_test_files", false, "Also parses _test.go files if set to true -- its set to false by default")

	flag.Parse()
}

func SetDefaults() {

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

}
