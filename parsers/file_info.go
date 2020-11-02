package parsers

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

type FileInfo struct {
	Name         string
	PathAndName  string
	FileLocation string
	AstFile      *ast.File
	NumberLines  int
	Size         int

	Structs    []*StructInfo
	Methods    []*MethodInfo
	Funcs      []*FuncInfo
	Interfaces []*InterfaceInfo
	Imports    []string

	HasMainFunc   bool
	InitFuncCount int

	Lines       []string
	PackageInfo *PackageInfo
	AstPackage  *ast.Package
	Application *Application
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

				//tmRole := PriMethodCompType

				if isPublic {
					//	tmRole = PubMethodCompType
				}

				meth := MethodInfo{
					Name:        funcName,
					Application: f.PackageInfo.ParentApp,
					Package:     f.PackageInfo,
					File:        f,

					StructName: strName,

					// TODO : Get the struct reference too -- maybe after parsing everything
					//Stuct:
					//	Role:        tmRole,
					NumberLines: sz,
				}

				f.Methods = append(f.Methods, &meth)
				f.PackageInfo.ChildMethods = append(f.PackageInfo.ChildMethods, &meth)
				f.PackageInfo.TrackMethod(&meth)

			} else {

				//tmRole := PriFunctionCompType

				if isPublic {
					//	tmRole = PubFunctionCompType
				}

				fun := FuncInfo{
					Name:        funcName,
					Application: f.PackageInfo.ParentApp,
					Package:     f.PackageInfo,
					File:        f,
					//	Role:        tmRole,
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

					// role := PriInterfacesCompType

					if isPublic {
						//	role = PubInterfacesCompType
					}

					meth := InterfaceInfo{
						Name:        decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name,
						Application: f.PackageInfo.ParentApp,
						Package:     f.PackageInfo,
						File:        f,

						//	Role:    role,
						NumberLines: sz,
					}

					f.Interfaces = append(f.Interfaces, &meth)
					f.PackageInfo.ChildInterfaces = append(f.PackageInfo.ChildInterfaces, &meth)

				} else {

					//	role := PriStructCompType

					if isPublic {
						//		role = PubStructCompType
					}

					str := StructInfo{
						Name: nm,
						//	Role: role,
						Package:     f.PackageInfo,
						Application: f.PackageInfo.ParentApp,
						File:        f,
						NumberLines: sz,
						Content:     con,
					}

					f.Structs = append(f.Structs, &str)
					f.PackageInfo.ChildStructs = append(f.PackageInfo.ChildStructs, &str)

				}

			case token.IMPORT:
				for _, vv := range decl.(*ast.GenDecl).Specs {

					f.PackageInfo.ChildImports = append(f.PackageInfo.ChildImports, vv.(*ast.ImportSpec).Path.Value)
					f.Imports = append(f.Imports, strings.Trim(vv.(*ast.ImportSpec).Path.Value, `"`))
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
