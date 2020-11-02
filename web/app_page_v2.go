package web

import (
	"encoding/json"
	"fmt"
	"github.com/antigloss/go/logger"
	"github.com/dminGod/golook/parsers"
	"github.com/gin-gonic/gin"
	"github.com/timshannon/bolthold"
	"os"
	"time"

	"go.etcd.io/bbolt"
)

func showAppV2(c *gin.Context) {

	d := c.Query("dir")
	app = parsers.FetchApp(d)

	c.Writer.Header().Set("Content-Type", "text/html")
	var ret []byte
	app.SortPackages()

	var retObj AppPage

	retObj.ImportCount = make(map[string]int)
	retObj.ApplicationName = app.Name
	retObj.LinesInApp = app.LinesInApp
	retObj.PackagesCount = len(app.ChildPackages)
	retObj.UniqueImports = app.GetUniqueImports()
	retObj.ImportCount = app.ImportCounts

	for _, pkgs := range app.ChildPackages {

		var pkgFiles []File

		pkgs.SortFiles()

		for _, file := range pkgs.ChildFiles {

			var actionTags, structs, imports []string
			var funcs, methods map[string]int
			funcs = make(map[string]int)
			methods = make(map[string]int)

			if file.InitFuncCount > 0 {
				actionTags = append(actionTags, "has_int_func")
			}
			if file.HasMainFunc {
				actionTags = append(actionTags, "has_main_func")
			}

			for _, v := range file.Funcs {
				funcs[v.Name] = v.NumberLines
			}
			for _, v := range file.Methods {
				mName := fmt.Sprintf("%v.%v", v.StructName, v.Name)
				methods[mName] = v.NumberLines
			}
			for _, v := range file.Structs {
				structs = append(structs, v.Content)
			}
			for _, v := range file.Imports {
				imports = append(imports, v)
			}

			pkgFiles = append(pkgFiles, File{

				FileName:       file.Name,
				ActionTags:     actionTags,
				Loc:            file.NumberLines,
				ImportCount:    len(file.Imports),
				FuncCount:      len(file.Funcs),
				StructsCount:   len(file.Structs),
				InterfaceCount: len(file.Interfaces),
				MethodsCount:   len(file.Methods),
				Funcs:          funcs,
				Methods:        methods,
				Structs:        structs,
				Imports:        imports,
			})
		}

		retObj.AppPackages = append(retObj.AppPackages, AppPagePkg{
			Name:               pkgs.Name,
			FolderLocation:     pkgs.FolderLocation,
			LinesInPkg:         pkgs.GetLinesInPkg(),
			CountChildFiles:    len(pkgs.ChildFiles),
			CountInterfaces:    len(pkgs.ChildInterfaces),
			CountChildStructs:  len(pkgs.ChildStructs),
			CountChildFuncs:    len(pkgs.ChildFuncs),
			CountChildMethods:  len(pkgs.ChildMethods),
			CountUniqueImports: pkgs.UniqueImports(),
			ChildFiles:         pkgFiles,
			structMethods:      pkgs.StructMethods,
		})
	}

	var err error

	err = addApp(retObj)
	if err != nil {
		logger.Error(err.Error())
	}

	bhr := getApp("dmin")

	ret, err = json.Marshal(bhr)
	if err != nil {
		logger.Error("There was an error in Unmarshal JSON - Error: %v", err.Error())
		return
	}

	_, err = c.Writer.Write(ret)
	if err != nil {
		logger.Error("Error writing the output - Error : %v", err.Error())
	}
}

func showFromDB(c *gin.Context) {

	c.Writer.Header().Set("Content-Type", "text/html")
	var ret []byte

	var err error

	bhr := getApp("dmin")

	ret, err = json.Marshal(bhr)
	if err != nil {
		logger.Error("There was an error in Unmarshal JSON - Error: %v", err.Error())
		return
	}

	_, err = c.Writer.Write(ret)
	if err != nil {
		logger.Error("Error writing the output - Error : %v", err.Error())
	}
}

type AppPage struct {
	ApplicationName string         `json:"app_name"`
	LinesInApp      int            `json:"loc_app"`
	PackagesCount   int            `json:"count_child_packages"`
	UniqueImports   []string       `json:"unique_imports"`
	ImportCount     map[string]int `json:"import_counts"`
	AppPackages     []AppPagePkg   `json:"packages"`
}

type AppPagePkg struct {
	Name               string                   `json:"name"`
	FolderLocation     string                   `json:"folder_location"`
	LinesInPkg         int                      `json:"loc_package"`
	CountChildFiles    int                      `json:"count_child_files"`
	CountInterfaces    int                      `json:"count_interfaces"`
	CountChildStructs  int                      `json:"count_structs"`
	CountChildFuncs    int                      `json:"count_funcs"`
	CountChildMethods  int                      `json:"count_methods"`
	CountUniqueImports int                      `json:"count_imports"`
	ChildFiles         []File                   `json:"child_files"`
	structMethods      []*parsers.PkgStructInfo `json:"package_struct_details"`
}

type File struct {
	FileName       string   `json:"file_name"`
	ActionTags     []string `json:"file_action_tags"`
	Loc            int      `json:"loc"`
	ImportCount    int      `json:"import_count"`
	FuncCount      int      `json:"func_count"`
	StructsCount   int      `json:"structs_count"`
	InterfaceCount int      `json:"interfaces_count"`
	MethodsCount   int      `json:"methods_count"`

	Structs []string `json:"structs"`
	Imports []string `json:"imports"`

	Funcs   map[string]int `json:"funcs"`
	Methods map[string]int `json:"methods"`
}

func addApp(a AppPage) (err error) {

	err = store.Upsert("dmin", &a)
	if err != nil {
		fmt.Printf("Got error when trying to upsert record for file read complete reference, File Ref: '%+v', Error: '%v' \n",
			a, err.Error())
	}

	return
}

func getApp(aName string) (a AppPage) {

	err := store.Get(aName, &a)
	if err != nil {
		if err == bolthold.ErrNotFound {
			fmt.Printf("No records found for the inode number : %v \n", aName)
			return
		}
	}

	return
}

func init() {

	InitLocalDB()
}

var store *bolthold.Store

func InitLocalDB() {

	var err error

	store, err = bolthold.Open("local_v2.db", 0666, &bolthold.Options{Options: &bbolt.Options{Timeout: 2 * time.Second}})

	if err != nil {
		fmt.Printf("Error creating the local database file - '%v' -- Will exit application -- Error: %v \n", "local.db", err)
		os.Exit(1)
		return
	}
}

/*
[
{ "name":"v1",
"folder_location":"d:/go_code/src/github.com/kubernetes/kubernetes/staging/src/k8s.io/api/core/v1",
"loc_package":82691,
"count_child_files":14,
"count_interfaces":12,
"count_structs":257,
"count_methods":3583,
"count_imports":15 },

{"name":"v1alpha2",
"folder_location":"d:/go_code/src/github.com/kubernetes/kubernetes/staging/src/k8s.io/cri-api/pkg/apis/runtime/v1alpha2",
"loc_package":33444,
"count_child_files":2,
"count_interfaces":38,
"count_structs":115,
"count_methods":1796,
"count_imports":13},

*/
