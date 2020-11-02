package web

import (
	"encoding/json"
	"fmt"
	"github.com/dminGod/golook/parsers"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

var app parsers.Application

func getPkgGraph(c *gin.Context) {

	d := c.Query("dir")
	app = parsers.FetchApp(d)

	var pi parsers.PkgImports

	piSt := make(map[string]interface{})
	uqLinks := make(map[string]parsers.Links)

	for i, v := range app.ChildPackages {
		pi.Data = append(pi.Data, parsers.PkgNames{
			ID: i, Value: v.GetLinesInPkg(), Label: v.Name,
			CountFiles: len(v.ChildFiles), CountFuncs: len(v.ChildFuncs), CountImports: len(v.ChildImports),
			CountInterfaces: len(v.ChildInterfaces), CountMethods: len(v.ChildMethods), CountStructs: len(v.ChildStructs),
		})
		piSt[v.Name] = i

		fmt.Printf("Setting '%v' \n", v.Name)
	}

	for i, v := range app.ChildPackages {
		for _, vv := range v.ChildImports {

			_, f := path.Split(strings.Replace(vv, `"`, "", -1))
			if _, ok := piSt[f]; ok {
				nm := fmt.Sprintf("%v-%v", i, piSt[f].(int))
				uqLinks[nm] = parsers.Links{i, piSt[f].(int), "to"}
			}
		}
	}

	for _, v := range uqLinks {
		pi.Links = append(pi.Links, v)
	}

	k, _ := json.Marshal(pi)
	c.Writer.Write(k)
}

func loadDirs(c *gin.Context) {

	var ret string

	for _, v := range parsers.GetLocalRepos() {
		ret += fmt.Sprintf(`<a href="?dir=%v">%v</a><br/>`, v, v)
	}

	c.Writer.Write([]byte(ret))
}

func jsonFS(c *gin.Context) {

	out := make(map[string][]JsonOut)

	out["children"] = []JsonOut{}

	for _, pkgs := range app.ChildPackages {

		var tmOut JsonOut

		tmOut.PkgName = pkgs.Name

		for _, file := range pkgs.ChildFiles {
			tmOut.DataPoints = append(tmOut.DataPoints, FileSize{FileName: fmt.Sprintf("%v-%v",
				file.Name, pkgs.Name), Size: file.NumberLines})
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

		tmOut.DataPoints = append(tmOut.DataPoints, FileSize{FileName: pkgs.Name, Size: lc})

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
			{FileName: pkgs.Name + "-structs", Size: len(pkgs.ChildStructs)},
			{FileName: pkgs.Name + "-funcs", Size: len(pkgs.ChildFuncs)},
			{FileName: pkgs.Name + "-methods", Size: len(pkgs.ChildMethods)},
		}
		out["children"] = append(out["children"], tmOut)
	}

	r, _ := json.Marshal(out)
	c.Writer.Write(r)
}
