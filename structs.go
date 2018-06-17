package main

type Component interface {
	GetRole() ComponentType
	GetStats() map[ComponentType]int
}

type ComponentType int

const (
	_                                 = iota
	ApplicationCompType ComponentType = iota
	PackageCompType
	FileCompType

	PubFunctionCompType
	PubInterfacesCompType
	PubStructCompType
	PubMethodCompType
	PubVariableCompType
	PubConstCompType

	PriFunctionCompType
	PriInterfacesCompType
	PriStructCompType
	PriMethodCompType
	PriVariableCompType
	PriConstCompType

	FunctionCall
	MethodCall

	ImportCompType
)

func (a *PackageInfo) ParseFolder() {

}

func (a *Application) GetRole() (ct ComponentType) {

	ct = ApplicationCompType
	return
}

func (a *Application) GetStats() (ret map[ComponentType]int) {

	// TODO
	return
}

func (a *PackageInfo) GetRole() (ct ComponentType) {
	ct = PackageCompType
	return
}

func (p *PackageInfo) GetStats() (ret map[ComponentType]int) {

	// TODO
	return
}

func (a *FileInfo) GetRole() (ct ComponentType) {

	ct = FileCompType
	return
}

func (p *FileInfo) GetStats() (ret map[ComponentType]int) {

	// TODO

	return
}

func (a *MethodInfo) GetRole() (ct ComponentType) {

	ct = FileCompType
	return
}

func (p *MethodInfo) GetStats() (ret map[ComponentType]int) {

	// TODO
	return
}

func (a *FuncInfo) GetRole() (ct ComponentType) {

	ct = FileCompType
	return
}

func (p *FuncInfo) GetStats() (ret map[ComponentType]int) {

	// TODO

	return
}

func IsPublic(c Component) (b bool) {

	r := c.GetRole()

	for _, value := range []ComponentType{PubFunctionCompType, PubMethodCompType, PubConstCompType,
		PubStructCompType, PubInterfacesCompType, PubVariableCompType} {

		if r == value {
			b = true
		}
	}

	return
}

func GetStats(c Component) (res map[ComponentType]int) {

	res = c.GetStats()
	return
}



//
//type Funcs struct {
//	PublicFunctions  []FuncInfo
//	PrivateFunctions []FuncInfo
//	Structs          map[string]FuncInfo
//	Methods          map[string][]FuncInfo
//	Imports          []string
//}
//
//type ByLineCount []FuncInfo
//
//func (a ByLineCount) Len() int           { return len(a) }
//func (a ByLineCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a ByLineCount) Less(i, j int) bool { return a[i].NumberLines > a[j].NumberLines }

// ToDO: Loop over a directory recursively to show package level information
// ToDO: The representation of data and presentation need to be made cleanly separate
// ToDO: File related data needs to be put into a different struct
// TODO: Star for checking method is wrong -- it should be directly on reciever
//
//func (f *Funcs) LenPublicMethods(s string) (i int) {
//
//	for _, v := range f.Methods[s] {
//		if v.IsPublic {
//			i += 1
//		}
//	}
//	return
//}
//
//func (f *Funcs) LenPrivateMethods(s string) (i int) {
//
//	for _, v := range f.Methods[s] {
//		if v.IsPublic == false {
//			i += 1
//		}
//	}
//
//	return
//}
//
//func NewFunc() (ret Funcs) {
//
//	ret.Methods = make(map[string][]FuncInfo)
//	ret.Structs = make(map[string]FuncInfo)
//
//	return
//}
//
//type PkgFuncs map[string]Funcs
