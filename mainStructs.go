package main

import (
	"go/ast"
	"go/token"
	"sort"
)

type Application struct {

	Name       string
	BaseFolder string
	SubFolders []string

	ChildPackages []*PackageInfo
	ChildFiles    []*FileInfo
	ChildStructs  []*StructInfo
	ChildMethods  []*MethodInfo
	ChildFuncs    []*FuncInfo
}

type PackageInfo struct {

	Name           string
	FolderLocation string
	PackageRef     *ast.Package
	AstFiles       map[string]*ast.File
	ParentApp      *Application
	FileSet        *token.FileSet

	ChildImports []string
	ChildFiles   []*FileInfo
	ChildStructs []*StructInfo
	ChildMethods []*MethodInfo
	ChildFuncs   []*FuncInfo
}

func (p *PackageInfo) SortFiles(){
	sort.SliceStable(p.ChildFiles, func(i int, j int) (bool){
		return p.ChildFiles[i].NumberLines > p.ChildFiles[j].NumberLines
	})
}

type FileInfo struct {

	Name         string
	PathAndName  string
	FileLocation string
	AstFile      *ast.File
	NumberLines  int
	Size         int

	Structs 	[]*StructInfo
	Methods     []*MethodInfo
	Funcs       []*FuncInfo

	Lines 		[]string
	PackageInfo *PackageInfo
	AstPackage  *ast.Package
	Application *Application
}

type MethodInfo struct {

	Name    string
	Content string
	StructName  string

	Role    ComponentType

	NumberLines int
	Application *Application
	Package     *PackageInfo
	Stuct       *StructInfo
	File        *FileInfo
}

type FuncInfo struct {

	Name    string
	Content string
	Role    ComponentType
	NumberLines int

	Application *Application
	Package     *PackageInfo
	File        *FileInfo
}

type StructInfo struct {

	Name    string
	Content string
	Role    ComponentType
	NumberLines int

	Application *Application
	Package     *PackageInfo
	File        *FileInfo
}
