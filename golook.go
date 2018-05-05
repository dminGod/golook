package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"go/ast"
	"io/ioutil"
	"strings"
	"os"
	"flag"
	"sort"
	"path"
)

type FuncInfo struct {

	Name string
	Type string
	Content string
	IsPublic bool
	Filename string
	Package string
	AstFile *ast.File
	IsMethod bool
	StructName string
	NumberLines	int

	FileName string
	FileLocation string
	Size int

	StructsCount       int
	PublicStructCount  int
	PrivateStructCount int

	MethodsCount int
	FunctionsCount int

	PublicFuncCount     int
	PrivateFuncCount    int
	PublicMethodsCount  int
	PrivateMethodsCount int

	ImportsCount int
}

type Funcs struct {

	PublicFunctions []FuncInfo
	PrivateFunctions []FuncInfo
	Structs map[string]FuncInfo
	Methods map[string][]FuncInfo
	Imports []string
}

type ByLineCount []FuncInfo
func (a ByLineCount) Len() int           { return len(a) }
func (a ByLineCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLineCount) Less(i, j int) bool { return a[i].NumberLines > a[j].NumberLines }


// ToDO: Loop over a directory recursively to show package level information
// ToDO: The representation of data and presentation need to be made cleanly separate
// ToDO: File related data needs to be put into a different struct
// TODO: Star for checking method is wrong -- it should be directly on reciever

func (f *Funcs) LenPublicMethods(s string) (i int) {

	for _, v := range f.Methods[s]{
		if v.IsPublic {
			i += 1
		}
	}
	return
}

func (f *Funcs) LenPrivateMethods(s string) (i int) {

	for _, v := range f.Methods[s] {
		if v.IsPublic == false {
			i += 1
		}
	}

	return
}

type PkgFuncs map[string]Funcs


var showTestFiles, showImports, showStructs, shortStructs, showStructMethods, showPublicFuncs, showPrivateFuncs, limitedView, showFuncs bool
var dir string

func main() {

	getFlags()

	var err error

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

	loadFolder()


}


func getFolders(base string) (s []string) {

	l, _ := ioutil.ReadDir( base )
	i := 0

	var t [][]string

	for _, b := range l {

		if b.IsDir() && b.Name()[:1] != "." && b.Name() != "vendor" {

			t = append(t, getFolders( path.Join(base, b.Name()) ))
			i++
		}
	}

	if i == 0 {

		return []string{ base }
	}

	for _, b := range t {

		for _, c := range b {

			s = append(s, c)
		}
	}

	return
}


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

func loadFolder() () {

	cwd, _ := os.Getwd()

	f := []string{cwd}

	for _, v := range getFolders(cwd) {

		f = append(f, v)
	}

	var totLOC int


	for _, a := range f {

		var FileList []FuncInfo

		fset := token.NewFileSet()

		var pkgInfo PkgFuncs
		pkgInfo = make(map[string]Funcs)

		filter := func(f os.FileInfo) (ret bool) {

			ret = true

			if showTestFiles == false {
				ret = strings.HasSuffix(f.Name(), "_test.go") == false
			}

			return
		}

		var recur []map[string]*ast.Package

		pkgs, err := parser.ParseDir(fset, a, filter, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		recur = append(recur, pkgs)

		for pkgn, pkg := range pkgs {

			var tmFunc Funcs
			tmFunc.Methods = make(map[string][]FuncInfo)
			tmFunc.Structs = make(map[string]FuncInfo)
			tmImports := make(map[string]struct{})

			for _, vv := range pkg.Files {

				fnFull := fset.File(vv.Pos()).Name()
				pt, fn := path.Split(fnFull)
				lc := fset.File(vv.Pos()).LineCount()
				sz := fset.File(vv.Pos()).Size()

				curFile := FuncInfo{
					Type: "file",
					NumberLines: lc,
					Size : sz,
					Filename: fn,
					Name: fnFull,
					FileLocation: pt,
					Package: pkg.Name,
				}

				b, err := ioutil.ReadFile(fnFull)
				if err != nil {
					fmt.Println("File error", pt, fn, err)
				}

				farr := strings.Split(string(b), "\n")

				for _, v := range vv.Decls {
					switch v.(type) {

					case *ast.FuncDecl:

						var isMethod bool
						var strName string

						st := fset.File(v.(*ast.FuncDecl).Pos()).Line(v.(*ast.FuncDecl).Pos())
						en := fset.File(v.(*ast.FuncDecl).Pos()).Line(v.(*ast.FuncDecl).End())
						sz := (en - st) + 1

						fd := v.(*ast.FuncDecl)
						fc := fd.Name.Name[:1]
						funcName := fd.Name.Name
						isPublic := fc == strings.ToUpper(fc)

						if fd.Recv != nil && len(fd.Recv.List) > 0 {

							if _, ok := (fd.Recv.List[0].Type).(*ast.StarExpr); ok {
								strName = fmt.Sprintf("%v", (fd.Recv.List[0].Type).(*ast.StarExpr).X)
							} else {
								strName = fmt.Sprintf((fd.Recv.List[0].Type).(*ast.Ident).Name)
							}

							isMethod = true
						}

						fn := FuncInfo{
							Name:       funcName,
							Type:       "function",
							Content:    retLine(farr, fset.Position(v.Pos()).Line-1),
							IsPublic:   isPublic,
							Package:    pkg.Name,
							Filename:   fn,
							IsMethod:   isMethod,
							StructName: strName,
							NumberLines: sz	}

						if isMethod {

							tmFunc.Methods[strName] = append(tmFunc.Methods[strName], fn)
							curFile.MethodsCount += 1

							if isPublic {
								curFile.PublicMethodsCount += 1
							} else {
								curFile.PrivateMethodsCount += 1
							}

						} else {

							curFile.FunctionsCount += 1

							if isPublic {
								curFile.PublicFuncCount += 1
								tmFunc.PublicFunctions = append(tmFunc.PublicFunctions, fn)
							} else {
								curFile.PrivateFuncCount += 1
								tmFunc.PrivateFunctions = append(tmFunc.PrivateFunctions, fn)
							}
						}

					case *ast.GenDecl:

						switch v.(*ast.GenDecl).Tok {

						case token.TYPE:

							nm := v.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name

							var st, en int
							var con string

							st = (fset.Position(v.Pos()).Line - 1)
							en = (fset.Position(v.End()).Line)
							sz = en - st

							if v.Pos() != v.End() {

								con = strings.Join(retLines(farr, st, en), "\n")
							} else {
								st = (fset.Position(v.Pos()).Line - 1)
								con = retLine(farr, st)
							}

							isPublic := nm[:1] == strings.ToUpper(nm[:1])

							if isPublic {
								curFile.PublicStructCount += 1
							} else {
								curFile.PrivateStructCount += 1
							}

							fn := FuncInfo{
								Name:     nm,
								Type:     "struct",
								Content:  con,
								IsPublic: isPublic,
								Package:  pkg.Name,
								Filename: fn,
								NumberLines: sz,
								//IsMethod: isMethod,
								//StructName : strName }
							}

							tmFunc.Structs[nm] = fn
							curFile.StructsCount += 1

						case token.IMPORT:

							for _, vv := range v.(*ast.GenDecl).Specs {
								tmImports[vv.(*ast.ImportSpec).Path.Value] = struct{}{}
								curFile.ImportsCount += 1
							}
						}
					}
				}

				FileList = append(FileList, curFile)
			}

			var gh, ngh []string

			for kk, _ := range tmImports {

				if strings.Contains(kk, "github.com") {
					gh = append(gh, strings.TrimSpace(kk))
				} else {
					ngh = append(ngh, strings.TrimSpace(kk))
				}
			}

			sort.Strings(gh)
			sort.Strings(ngh)

			for _, v := range gh {
				tmFunc.Imports = append(tmFunc.Imports, v)
			}

			for _, v := range ngh {
				tmFunc.Imports = append(tmFunc.Imports, v)
			}

			sort.Sort(ByLineCount(tmFunc.PublicFunctions))
			sort.Sort(ByLineCount(tmFunc.PrivateFunctions))
			sort.Sort(ByLineCount(FileList))

			pkgInfo[ pkgn ] = tmFunc
		}


		for _, v := range FileList {

			totLOC += v.NumberLines
		}


		for pkg, fInfo := range pkgInfo {

			fmt.Sprintf("\n\n\n\n\n\n\n\n---------- Package : %v -------------------------------------------------------------------\n", pkg)



			if showImports {
				fmt.Printf("\n-- File Details -- \n ")

				for _, v := range FileList {

					fmt.Sprintf(`
	%v --- %v
	  Lines: %v, Size: %v,
	  Funcs (Tot, Pub, Priv): (%v, %v, %v)
	  Structs (Tot, Pub, Pri): (%v, %v, %v),
	  Methods (Tot, Pub, Pri): (%v, %v, %v)
	`, v.Filename, v.Name, v.NumberLines, v.Size, v.FunctionsCount, v.PublicFuncCount, v.PrivateFuncCount,
						v.StructsCount, v.PublicStructCount, v.PrivateStructCount, v.MethodsCount, v.PublicMethodsCount, v.PrivateMethodsCount )
				}
			}

			if showImports {
				fmt.Printf("\n-- List of Imports -- \n%v \n ", strings.Join(fInfo.Imports, "\n"))
			}
			ii := 0
			for _, v := range fInfo.Structs {

				showStruct(v, fInfo, true, ii)
				ii++
			}

			ii = 0
			for _, v := range fInfo.Structs {

				showStruct(v, fInfo, false, ii)
				ii++
			}



			if showPublicFuncs || showFuncs {
				fmt.Println("\n---   Public Functions ---------------------------------------------------")

				for _, v := range fInfo.PublicFunctions {

					fmt.Printf("  %v (%v)\t\t%v", v.Name, v.NumberLines, v.Content)
				}
			}

			if showPrivateFuncs || showFuncs {
				fmt.Println("\n---   Private Functions ---------------------------------------------------")

				for _, v := range fInfo.PrivateFunctions {
					fmt.Printf("  %v (%v)\t\t%v\n", v.Name, v.NumberLines, v.Content)
				}
			}
		}


		fmt.Println("Total Lines of Code : ", totLOC)
	}



}

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

func showStruct(v FuncInfo, fInfo Funcs, publicLoop bool, ii int) {

	if ((publicLoop && v.IsPublic) || (publicLoop == false && v.IsPublic == false)) && (showStructs || shortStructs) {

		if ii == 0 && v.IsPublic {

			fmt.Println("\n---   Public Structs ---------------------------------------------------")
		}

		if ii == 0 && v.IsPublic == false {

			fmt.Println("\n---   Private Structs ---------------------------------------------------")
		}

		strCon := v.Content

		if shortStructs {
			strCon = strings.Split(v.Content, "\n")[0]
		}

		fmt.Printf("\n---- %v (%v) ---- \n %v \n", v.Name, v.NumberLines, strCon)

		for i, vv := range fInfo.Methods[v.Name] {

			if i == 0 && fInfo.LenPublicMethods(v.Name) > 0 {
				fmt.Printf("\n\tPublic Methods: \n")
			}

			if vv.IsPublic {

				fmt.Printf("\t  %v (%v)\t\t%v \n", vv.Name, vv.NumberLines, vv.Content)
			}
		}

		for i, vv := range fInfo.Methods[v.Name] {

			if i == 0 && fInfo.LenPrivateMethods(v.Name) > 0 {
				fmt.Printf("\n\tPrivate Methods: \n")
			}

			if vv.IsPublic == false {

				fmt.Printf("\t  %v (%v)\t\t%v \n", vv.Name, vv.NumberLines, vv.Content)
			}
		}

		fmt.Printf("\n")
	}
}