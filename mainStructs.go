package main

import (
	"go/ast"
	"go/token"
	"sort"
)

type Application struct {

	Name       string		`boltholdKey:"name"`
	BaseFolder string
	SubFolders []string

	ChildPackages []*PackageInfo
	ChildFiles    []*FileInfo
	ChildStructs  []*StructInfo
	ChildMethods  []*MethodInfo
	ChildInterfaces []*InterfaceInfo
	ChildFuncs    []*FuncInfo

	GraphData     PkgImports
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
	ChildTests   []*FileInfo
	ChildStructs []*StructInfo
	ChildMethods []*MethodInfo
	ChildInterfaces []*InterfaceInfo
	ChildFuncs   []*FuncInfo
}

func (p *PackageInfo) UniqueImports() (int){

	var st map[string]interface{}

	st = make(map[string]interface{})

	for _, v := range p.ChildImports {
		st[v] = nil
	}

	return len(st)
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
	Interfaces  []*InterfaceInfo
	Imports     []string

	HasMainFunc bool
	InitFuncCount   int

	Lines 		[]string
	PackageInfo *PackageInfo
	AstPackage  *ast.Package
	Application *Application
}



type InterfaceInfo struct {

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
