package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"encoding/json"
)

func showApp(c *gin.Context) {

	var ret string

	ret = "<html><body>"

	for _, pkgs := range app.ChildPackages {

	ret += fmt.Sprintf(`
<h3 class="package_h3">%v</h3>
<div class="package_header">
Folder: %v <br/>
Child Files: %v - Funcs : %v - Methods: %v </div>`,
			pkgs.Name,
			pkgs.FolderLocation,
			len(pkgs.ChildFiles),
			len(pkgs.ChildFuncs),
			len(pkgs.ChildMethods),
		)

		ret += `<div class="file_details">`

		for _, file := range pkgs.ChildFiles {
			ret += fmt.Sprintf("<b>%v</b> File: %v </br>", file.Name, file.NumberLines)
		}

		ret += `</div>`
	}

	ret += "</body></html>"

	c.Writer.Write([]byte(ret))
}


type JsonOut struct {

	PkgName	string	`json:"name"`
	DataPoints []FileSize `json:"children"`
}

type FileSize struct {

	FileName 	string `json:"name"`
	Size 		int `json:"size"`
}


func jsonFS(c *gin.Context) {

	out := make(map[string][]JsonOut)

	out["children"] = []JsonOut{}

	for _, pkgs := range app.ChildPackages {

		var tmOut JsonOut

		tmOut.PkgName = pkgs.Name

		for _, file := range pkgs.ChildFiles {
			tmOut.DataPoints = append(tmOut.DataPoints, FileSize{ FileName:
				fmt.Sprintf("%v-%v",
				file.Name, pkgs.Name), Size: file.NumberLines })
		}

		out["children"] = append(out["children"], tmOut)
	}

	r, _ := json.Marshal(out)

	c.Writer.Write(r)
}

func jsonPackageSize(c *gin.Context) {

	out := make(map[string][]JsonOut)

	out["children"] = []JsonOut{}

	for _, pkgs := range app.ChildPackages {

		var tmOut JsonOut
		var lc int

		tmOut.PkgName = pkgs.Name

		for _, file := range pkgs.ChildFiles {
			lc += file.NumberLines
		}

		tmOut.DataPoints = append(tmOut.DataPoints, FileSize{ FileName:
			pkgs.Name, Size: lc })

		out["children"] = append(out["children"], tmOut)
	}

	r, _ := json.Marshal(out)
	c.Writer.Write(r)
}

func jsonMethodFuncs(c *gin.Context) {

	out := make(map[string][]JsonOut)

	out["children"] = []JsonOut{}
	for _, pkgs := range app.ChildPackages {

		var tmOut JsonOut

		tmOut.PkgName = pkgs.Name
		tmOut.DataPoints = []FileSize{
			{ FileName: pkgs.Name + "-structs", Size: len(pkgs.ChildStructs) },
			{ FileName: pkgs.Name + "-funcs", Size: len(pkgs.ChildFuncs) },
			{ FileName: pkgs.Name + "-methods", Size: len(pkgs.ChildMethods) },
		}
		out["children"] = append(out["children"], tmOut)
	}

	r, _ := json.Marshal(out)
	c.Writer.Write(r)
}


