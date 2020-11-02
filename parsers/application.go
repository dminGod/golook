package parsers

import (
	"fmt"
	"github.com/kubernetes/kubernetes/pkg/util/slice"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"sort"
	"strings"
)

type PkgImports struct {
	Data  []PkgNames `json:"data"`
	Links []Links    `json:"links"`
}

type PkgNames struct {
	ID    int    `json:"id"`
	Value int    `json:"value"`
	Label string `json:"label"`

	CountFiles      int `json:"files_count"`
	CountFuncs      int `json:"funcs_count"`
	CountInterfaces int `json:"interface_count"`
	CountStructs    int `json:"structs_count"`
	CountMethods    int `json:"methods_count"`
	CountImports    int `json:"imports_count"`
}

type Links struct {
	From   int    `json:"from"`
	To     int    `json:"to"`
	Arrows string `json:"arrows"`
}

type Application struct {
	Name       string `boltholdKey:"name"`
	BaseFolder string
	SubFolders []string

	LinesInApp    int
	UniqueImports map[string]interface{}
	AllImports    []string
	ImportCounts  map[string]int

	ChildPackages   []*PackageInfo
	ChildFiles      []*FileInfo
	ChildStructs    []*StructInfo
	ChildMethods    []*MethodInfo
	ChildInterfaces []*InterfaceInfo
	ChildFuncs      []*FuncInfo

	GraphData PkgImports
}

// Sort the packages based on the line size of the package
func (a *Application) SortPackages() {

	sort.SliceStable(a.ChildPackages, func(i int, j int) bool {
		return a.ChildPackages[i].GetLinesInPkg() > a.ChildPackages[j].GetLinesInPkg()
	})

	return
}

func (a *Application) CountUniqueImports() int {

	return len(a.UniqueImports)
}

func (a *Application) GetUniqueImports() (uImports []string) {

	for k, _ := range a.UniqueImports {
		uImports = append(uImports, k)
	}

	slice.SortStrings(uImports)

	return
}

// Loop over each subfolder
// Get the packages in the subfolder
// Loop over the packages
// For each package, parse every file and load the details on the file level
// populate the data on the package level
// populate the data of the packages on the application level

func (a *Application) ReadApp() (err error) {

	for _, folder := range a.SubFolders {

		// Create a new fileset
		fileSet := token.NewFileSet()

		// Make a list of packages that are part of this SubFolder
		pkgs := make(map[string]*ast.Package)

		// Get the packages in the directory
		pkgs, err = parser.ParseDir(fileSet, folder, filesFilterFunc, parser.ParseComments)
		if err != nil {
			return
		}

		for _, pkg := range pkgs {

			tmpPkg := ParsePackage(a, fileSet, folder, pkg)

			for _, file := range tmpPkg.ChildFiles {

				a.ChildFiles = append(a.ChildFiles, file)
			}

			for _, str := range tmpPkg.ChildStructs {

				a.ChildStructs = append(a.ChildStructs, str)
			}

			for _, mthd := range tmpPkg.ChildMethods {

				for _, str := range tmpPkg.ChildStructs {

					if mthd.StructName == str.Name {
						mthd.Stuct = str
					}
				}

				a.ChildMethods = append(a.ChildMethods, mthd)
			}

			for _, fun := range tmpPkg.ChildFuncs {
				a.ChildFuncs = append(a.ChildFuncs, fun)
			}

			a.LinesInApp += tmpPkg.GetLinesInPkg()
			a.AllImports = append(a.AllImports, tmpPkg.GetUniqueImports()...)

			for _, v := range tmpPkg.GetUniqueImports() {

				v = strings.Trim(v, `"`)

				if _, ok := a.ImportCounts[v]; ok {
					a.ImportCounts[v] += 1
				} else {
					a.ImportCounts[v] = 1
				}

				a.UniqueImports[v] = nil
			}

			a.ChildPackages = append(a.ChildPackages, tmpPkg)
		}
	}

	var pi PkgImports

	piSt := make(map[string]interface{})
	uqLinks := make(map[string]Links)

	for i, v := range a.ChildPackages {
		pi.Data = append(pi.Data, PkgNames{
			ID:              i,
			Value:           v.GetLinesInPkg(),
			Label:           v.Name,
			CountFiles:      len(v.ChildFiles),
			CountFuncs:      len(v.ChildFuncs),
			CountImports:    len(v.ChildImports),
			CountInterfaces: len(v.ChildInterfaces),
			CountMethods:    len(v.ChildMethods),
			CountStructs:    len(v.ChildStructs)})
		piSt[v.Name] = i

		// fmt.Printf("Setting '%v' \n", v.Name)
	}

	for i, v := range a.ChildPackages {
		for _, vv := range v.ChildImports {

			_, f := path.Split(strings.Replace(vv, `"`, "", -1))
			if _, ok := piSt[f]; ok {
				nm := fmt.Sprintf("%v-%v", i, piSt[f].(int))
				uqLinks[nm] = Links{i, piSt[f].(int), "to"}
			}
		}
	}

	for _, v := range uqLinks {
		pi.Links = append(pi.Links, v)

	}

	a.GraphData = pi

	return
}

// Main function that reads the whole application
func FetchApp(dir string) (RetApp Application) {

	RetApp = NewApplication(dir)
	RetApp.ReadApp()

	//err := addApp(RetApp)
	//if err != nil {
	//	fmt.Println(err)
	//}

	// x := getApp(RetApp.Name)

	// fmt.Println(x, "<-- Parsed application")
	return
}

/*
Set the BaseFolder
Get all the subdirectories in the folder and set them as the SubFolder
Return the application
*/

func NewApplication(baseFolder string) (RetApp Application) {

	RetApp.ImportCounts = make(map[string]int)
	RetApp.UniqueImports = make(map[string]interface{})
	RetApp.BaseFolder = baseFolder
	RetApp.SubFolders = GetAllSubDirs(baseFolder)
	RetApp.Name = baseFolder

	return
}
