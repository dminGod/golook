package main

import (
	"fmt"
	"go/ast"
	"path"
	"reflect"
	"strings"
	"go/token"
	"io/ioutil"
	"go/parser"
	"sort"
)


/*
Set the BaseFolder
Get all the subdirectories in the folder and set them as the SubFolder
Return the application
 */

func NewApplication(baseFolder string) (RetApp Application) {

	RetApp.BaseFolder = baseFolder
	RetApp.SubFolders = GetAllSubDirs(baseFolder)

	return
}

// Make a package - Loop over each file of the package and call LoadFile on the file

func ParsePackage(a *Application, fileSet *token.FileSet, folder string, pkg *ast.Package) (retPkg *PackageInfo) {

	retPkg = &PackageInfo {
		Name:           pkg.Name,
		FolderLocation: folder,
		AstFiles:       pkg.Files,
		ParentApp: a,
		FileSet: fileSet,
		PackageRef: pkg,
	}

	// Loop over the files
	for _, f := range pkg.Files {

		tmFile := LoadFile( retPkg, f )

		if strings.HasSuffix(tmFile.Name, "_test.go") == false {

			retPkg.ChildFiles = append(retPkg.ChildFiles, tmFile)
		} else {

			retPkg.ChildTests = append(retPkg.ChildFiles, tmFile)
		}
	}

	// TODO: Map the methods and structs here
	return
}

// Sort the packages based on the line size of the package

func (a *Application) SortPackages() {

	sort.SliceStable(a.ChildPackages, func(i int, j int) (bool){
		return a.ChildPackages[i].GetLinesInPkg() > a.ChildPackages[j].GetLinesInPkg()
	})

	return
}

// Main function that reads the whole application

func FetchApp(dir string) (RetApp Application) {

	RetApp = NewApplication(dir)
	RetApp.ReadApp()


	err := addApp(RetApp)
	if err != nil {
		fmt.Println(err)
	}

	x := getApp(RetApp.Name)

	fmt.Println(x, "<-- Parsed application")
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

				a.ChildMethods = append( a.ChildMethods, mthd )
			}

			for _, fun := range tmpPkg.ChildFuncs {
				a.ChildFuncs = append(a.ChildFuncs, fun)
			}

			a.ChildPackages = append(a.ChildPackages, tmpPkg)
		}
	}


	var pi PkgImports

	piSt := make(map[string]interface{})
	uqLinks := make(map[string]Links)

	for i, v := range a.ChildPackages {
		pi.Data = append(pi.Data, PkgNames{ID: i, Value: v.GetLinesInPkg(), Label: v.Name,
			CountFiles: len(v.ChildFiles), CountFuncs: len(v.ChildFuncs), CountImports: len(v.ChildImports),
			CountInterfaces: len(v.ChildInterfaces), CountMethods: len(v.ChildMethods), CountStructs: len(v.ChildStructs)})
		piSt[v.Name] = i

		fmt.Printf("Setting '%v' \n", v.Name)
	}

	for i, v := range a.ChildPackages {
		for _, vv := range v.ChildImports {

			_, f := path.Split(strings.Replace(vv,`"`, "", -1))
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

func (p *PackageInfo) GetLinesInPkg() (lns int){

	for _, f := range p.ChildFiles {

		lns += f.NumberLines
	}

	return
}

func LoadFile(pkgInfo *PackageInfo, f *ast.File) (file *FileInfo) {

	fnFull := pkgInfo.FileSet.File(f.Pos()).Name()
	path, fn := path.Split(fnFull)

	// Read the actual file
	b, err := ioutil.ReadFile( fnFull )
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
		Lines: 		  farr,
		AstFile:      f}

	file.Parse()
	return
}

func (f *FileInfo) Parse() {

	for _, decl := range f.AstFile.Decls {

		switch decl.(type) {

		case *ast.FuncDecl:

			var isMethod bool
			//var strName string

			st := f.PackageInfo.FileSet.File(decl.(*ast.FuncDecl).Pos()).Line(decl.(*ast.FuncDecl).Pos())
			en := f.PackageInfo.FileSet.File(decl.(*ast.FuncDecl).Pos()).Line(decl.(*ast.FuncDecl).End())
			sz := (en - st) + 1

			fd := decl.(*ast.FuncDecl)
			fc := fd.Name.Name[:1]
			funcName := fd.Name.Name


			if funcName == "main" {
				f.HasMainFunc = true
			}

			if funcName == "init" {
				f.InitFuncCount += 1
			}

			isPublic := fc == strings.ToUpper(fc)

			strName := ""

			if fd.Recv != nil && len(fd.Recv.List) > 0 {

				if _, ok := (fd.Recv.List[0].Type).(*ast.StarExpr); ok {
					strName = fmt.Sprintf("%v", (fd.Recv.List[0].Type).(*ast.StarExpr).X)
				} else {
					strName = fmt.Sprintf((fd.Recv.List[0].Type).(*ast.Ident).Name)
				}

				isMethod = true
			}

			//for _, b := range fd.Body.List {
			//	_ = b
			//	//fmt.Println("This is b", b)
			//}

			if isMethod {

				tmRole := PriMethodCompType

				if isPublic {
					tmRole = PubMethodCompType
				}

				meth := MethodInfo{
					Name:        funcName,
					Application: f.PackageInfo.ParentApp,
					Package:     f.PackageInfo,
					File:        f,

					StructName: strName,

					// TODO : Get the struct reference too -- maybe after parsing everything
					//Stuct:
					Role:        tmRole,
					NumberLines: sz,
				}

				f.Methods = append(f.Methods, &meth)
				f.PackageInfo.ChildMethods = append(f.PackageInfo.ChildMethods, &meth)

			}  else {

				tmRole := PriFunctionCompType

				if isPublic {
					tmRole = PubFunctionCompType
				}

				fun := FuncInfo{
					Name:        funcName,
					Application: f.PackageInfo.ParentApp,
					Package:     f.PackageInfo,
					File:        f,
					Role:        tmRole,
					NumberLines: sz,
				}

				f.Funcs = append(f.Funcs, &fun)
				f.PackageInfo.ChildFuncs = append(f.PackageInfo.ChildFuncs, &fun)
			}

		case *ast.GenDecl:

			switch decl.(*ast.GenDecl).Tok {

			case token.TYPE:
				nm := decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name

				var st, en int
				var con string

				st = (f.PackageInfo.FileSet.Position(decl.Pos()).Line - 1)
				en = (f.PackageInfo.FileSet.Position(decl.End()).Line)
				sz := en - st

				if decl.Pos() != decl.End() {
					con = strings.Join(retLines(f.Lines, st, en), "\n")
				} else {
					st = (f.PackageInfo.FileSet.Position(decl.Pos()).Line - 1)
					con = retLine(f.Lines, st)
				}


				//switch decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(type) {
				//
				//case *ast.StructType:
				//	fmt.Printf("Struct Type %+v  -- Type %+v - Name : %v - File : %v \n",
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec),
				//		reflect.TypeOf(decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type).String(),
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name,
				//		f.Name,
				//		)
				//
				//case *ast.InterfaceType:
				//	//fmt.Printf("Interface Type %+v  -- Type %+v - Name : %v - File : %v  \n",
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec),
				//		reflect.TypeOf(decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type).String(),
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name,
				//		f.Name,
				//		)
				//
				//case *ast.ArrayType:
				//	fmt.Printf("Array Type %+v  -- Type %+v - Name : %v - File : %v  \n",
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec),
				//		reflect.TypeOf(decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type).String(),
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name,
				//		f.Name,
				//		)
				//
				//case *ast.Ident:
				//	fmt.Printf("Identity Type %+v  -- Type %+v - Name : %v - File : %v  \n",
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec),
				//		reflect.TypeOf(decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type).String(),
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name,
				//		f.Name,
				//		)
				//
				//case *ast.MapType:
				//	fmt.Printf("MapType Type %+v  -- Type %+v - Name : %v - File : %v  \n",
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec),
				//		reflect.TypeOf(decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type).String(),
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name,
				//		f.Name,
				//		)
				//
				//case *ast.FuncType:
				//	fmt.Printf("Func Type %+v  -- Type %+v - Name : %v - File : %v  \n",
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec),
				//		reflect.TypeOf(decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type).String(),
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name,
				//		f.Name,
				//	)
				//
				//default:
				//	fmt.Printf("********* Unknow Type %+v  -- Type %+v - Name : %v - File : %v  \n",
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec),
				//		reflect.TypeOf(decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type).String(),
				//		decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name,
				//		f.Name,
				//	)
				//
				//}


				isPublic := nm[:1] == strings.ToUpper(nm[:1])

				if reflect.TypeOf(decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type).String() == "*ast.InterfaceType" {

					role := PriInterfacesCompType

					if isPublic {
						role = PubInterfacesCompType
					}

					meth := InterfaceInfo{
						Name:        decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name,
						Application: f.PackageInfo.ParentApp,
						Package:     f.PackageInfo,
						File:        f,

						Role:    role,
						NumberLines: sz,
					}

					f.Interfaces = append(f.Interfaces, &meth)
					f.PackageInfo.ChildInterfaces = append(f.PackageInfo.ChildInterfaces, &meth)

				} else {

					role := PriStructCompType

					if isPublic {
						role = PubStructCompType
					}

					str := StructInfo{
						Name: nm,
						Role: role,
						Package: f.PackageInfo,
						Application: f.PackageInfo.ParentApp,
						File: f,
						NumberLines: sz,
						Content: con,
						}

					f.Structs = append(f.Structs, &str)
					f.PackageInfo.ChildStructs = append(f.PackageInfo.ChildStructs, &str)

				}


			case token.IMPORT:
				for _, vv := range decl.(*ast.GenDecl).Specs {

					f.PackageInfo.ChildImports = append(f.PackageInfo.ChildImports, vv.(*ast.ImportSpec).Path.Value)
					f.Imports = append(f.Imports, vv.(*ast.ImportSpec).Path.Value)
				}

			case token.VAR:
				for _, vv := range decl.(*ast.GenDecl).Specs {
					j := vv.(*ast.ValueSpec)
					_ = j
					//fmt.Println(j)
				}

			case token.CONST:

			default:
				fmt.Println("General Declaration : ", decl.(*ast.GenDecl).Tok)
			}

		default:
			fmt.Println("V.type : ", reflect.TypeOf(decl))
		}
	}

	return
}
