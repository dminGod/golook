package parsers

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"path"
	"strings"
)

func LoadFile(pkgInfo *PackageInfo, f *ast.File) (file *FileInfo) {

	fnFull := pkgInfo.FileSet.File(f.Pos()).Name()
	path, fn := path.Split(fnFull)

	// Read the actual file
	b, err := ioutil.ReadFile(fnFull)
	if err != nil {
		fmt.Println("File error", fn, err)
	}

	// Make the lines into []string
	farr := strings.Split(string(b), "\n")

	file = &FileInfo{
		PathAndName:  fnFull,
		FileLocation: path,
		Name:         fn,
		NumberLines:  pkgInfo.FileSet.File(f.Pos()).LineCount(),
		Size:         pkgInfo.FileSet.File(f.Pos()).Size(),
		PackageInfo:  pkgInfo,
		AstPackage:   pkgInfo.PackageRef,
		Application:  pkgInfo.ParentApp,
		Lines:        farr,
		AstFile:      f}

	file.Parse()
	return
}
