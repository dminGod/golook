package parsers

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

type PackageInfo struct {
	Name           string
	FolderLocation string
	PackageRef     *ast.Package
	AstFiles       map[string]*ast.File
	ParentApp      *Application
	FileSet        *token.FileSet

	ChildImports    []string
	ChildFiles      []*FileInfo
	ChildTests      []*FileInfo
	ChildStructs    []*StructInfo
	ChildMethods    []*MethodInfo
	ChildInterfaces []*InterfaceInfo
	ChildFuncs      []*FuncInfo

	Structs       []string
	StructMethods []*PkgStructInfo `json:"package_struct_details"`
}

type PkgStructInfo struct {
	StructName     string                 `json:"struct_name"`
	MethodsCount   int                    `json:"methods_count"`
	LinesInMethods int                    `json:"loc_struct_methods"`
	UniqueMethods  map[string]interface{} `json:"unique_methods"`
	MethodToLoc    map[string]int         `json:"methods_loc"`
	MethodInfo     []*MethodInfo          `json:"-"`
}

func (p *PackageInfo) TrackMethod(info *MethodInfo) {

	var found bool

	if info.Name == "" {
		return
	}

	// Loop over each of the existing structs to see if we already have this value
	// if we already are tracking this struct then we will update the info
	for i, v := range p.StructMethods {

		if v.StructName == info.StructName {

			found = true

			v.MethodInfo = append(p.StructMethods[i].MethodInfo, info)
			v.MethodsCount += 1
			v.LinesInMethods += info.NumberLines
			v.UniqueMethods[info.Name] = nil
			v.MethodToLoc[info.Name] = info.NumberLines
		}
	}

	// If we cannot find this struct then we need to make a new one:
	if found == false {

		var pkg PkgStructInfo
		pkg.StructName = info.StructName
		pkg.UniqueMethods = make(map[string]interface{})
		pkg.MethodToLoc = make(map[string]int)

		pkg.MethodInfo = append(pkg.MethodInfo, info)

		pkg.MethodsCount = 1
		pkg.LinesInMethods += info.NumberLines
		pkg.UniqueMethods[info.Name] = nil
		pkg.MethodToLoc[info.Name] = info.NumberLines

		p.StructMethods = append(p.StructMethods, &pkg)
	}
}

func (p *PackageInfo) UniqueImports() int {

	var st map[string]interface{}

	st = make(map[string]interface{})

	for _, v := range p.ChildImports {
		st[v] = nil
	}

	return len(st)
}

func (p *PackageInfo) GetUniqueImports() (uImports []string) {

	var st map[string]interface{}
	st = make(map[string]interface{})

	for _, v := range p.ChildImports {
		st[v] = nil
	}

	for k, _ := range st {
		uImports = append(uImports, k)
	}

	return
}

func (p *PackageInfo) SortFiles() {
	sort.SliceStable(p.ChildFiles, func(i int, j int) bool {
		return p.ChildFiles[i].NumberLines > p.ChildFiles[j].NumberLines
	})
}

func (p *PackageInfo) GetLinesInPkg() (lns int) {

	for _, f := range p.ChildFiles {

		lns += f.NumberLines
	}

	return
}

// Make a package - Loop over each file of the package and call LoadFile on the file
func ParsePackage(a *Application, fileSet *token.FileSet, folder string, pkg *ast.Package) (retPkg *PackageInfo) {

	retPkg = &PackageInfo{
		Name:           pkg.Name,
		FolderLocation: folder,
		AstFiles:       pkg.Files,
		ParentApp:      a,
		FileSet:        fileSet,
		PackageRef:     pkg,
	}

	// Loop over the files
	for _, f := range pkg.Files {

		tmFile := LoadFile(retPkg, f)

		if strings.HasSuffix(tmFile.Name, "_test.go") == false {

			retPkg.ChildFiles = append(retPkg.ChildFiles, tmFile)
		} else {

			retPkg.ChildTests = append(retPkg.ChildFiles, tmFile)
		}
	}

	// TODO: Map the methods and structs here
	return
}
