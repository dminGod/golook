package main

import (
	"os"
	"go/ast"
	"sort"
	"go/token"
	"strings"
)


// General filter to exclude some time of files if you want
// Currently removes _test.go files if they are specified

func filesFilterFunc(f os.FileInfo) (ret bool) {

	ret = true

	if showTestFiles == false {

		ret = strings.HasSuffix(f.Name(), "_test.go") == false
	}

	return
}

// Take the imports, sort them and add them as : github then the non github
// Why is this not used?

func SortReturnImports(imports map[string]struct{})(ret []string) {

	var gh, ngh []string

	for kk, _ := range imports {
		if strings.Contains(kk, "github.com") {
			gh = append(gh, strings.TrimSpace(kk))
		} else {
			ngh = append(ngh, strings.TrimSpace(kk))
		}
	}

	sort.Strings( gh )
	sort.Strings( ngh )

	for _, v := range gh {
		ret = append(ret, v)
	}

	for _, v := range ngh {
		ret = append(ret, v)
	}

return
}

// Why is this not used..?

func FindTypes(v ast.Node) ast.Visitor {

	switch v.(type) {
	case *ast.FuncDecl:

		for _, vv := range v.(*ast.FuncDecl).Body.List {

			switch vv.(type) {

			case *ast.ExprStmt:

				switch vv.(*ast.ExprStmt).X.(type) {

				case *ast.CallExpr:

					if _, ok := vv.(*ast.ExprStmt).X.(*ast.CallExpr).Fun.(*ast.SelectorExpr); ok {

						// fmt.Println(vv.(*ast.ExprStmt).X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.Name)
					}
				}
			}
		}

	case *ast.GenDecl:

		switch v.(*ast.GenDecl).Tok {

		case token.TYPE:

		case token.IMPORT:

		case token.VAR:
			//for _, vv := range v.(*ast.GenDecl).Specs {
			//
			//	j := vv.(*ast.ValueSpec)
			//	fmt.Println(j)
			//}

		case token.CONST:

		case token.FUNC:

			//for _, vv := range v.(*ast.GenDecl).Specs {
			//
			//	j := vv.(*ast.ValueSpec)
			//	fmt.Println(j)
			//}

		default:
			// fmt.Println("General Declaration : ", v.(*ast.GenDecl).Tok)
		}

	default:
		// fmt.Println("V.type : ", reflect.TypeOf(v))
	}

	return VisitorFunc(FindTypes)
}

type VisitorFunc func (ast.Node) ast.Visitor

func (f VisitorFunc) Visit(n ast.Node) ast.Visitor { return f(n) }

