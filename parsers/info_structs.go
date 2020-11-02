package parsers

type InterfaceInfo struct {
	Name       string
	Content    string
	StructName string

	//Role    ComponentType

	NumberLines int
	Application *Application
	Package     *PackageInfo
	Stuct       *StructInfo
	File        *FileInfo
}

type MethodInfo struct {
	Name       string
	Content    string
	StructName string

	// Role    ComponentType

	NumberLines int
	Application *Application
	Package     *PackageInfo
	Stuct       *StructInfo
	File        *FileInfo
}

type FuncInfo struct {
	Name    string
	Content string
	//Role    ComponentType
	NumberLines int

	Application *Application
	Package     *PackageInfo
	File        *FileInfo
}

type StructInfo struct {
	Name    string
	Content string
	//Role    ComponentType
	NumberLines int

	Application *Application
	Package     *PackageInfo
	File        *FileInfo

	ChildMethods []*MethodInfo
}
