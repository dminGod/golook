package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func StartUI() {

	r := gin.Default()
	r.GET("/show_app", showApp)
	r.GET("/load_dirs", loadDirs)
	r.GET("/get_package_graph", getPkgGraph)

	r.GET("/json_file_size", jsonFS)
	r.GET("/json_pkg_size", jsonPackageSize)
	r.GET("/json_func_count", jsonMethodFuncs)

	r.Static("/static", "./static")

	r.Run("0.0.0.0:8081") // listen and serve on 0.0.0.0:8080

}

type PkgImports struct {
	Data  []PkgNames `json:"data"`
	Links []Links    `json:"links"`
}

type PkgNames struct {
	ID    int    `json:"id"`
	Value int    `json:"value"`
	Label string `json:"label"`

	CountFiles      int `json:"files_count"`
	CountFuncs      int `json:"funcs_count"`
	CountInterfaces int `json:"interface_count"`
	CountStructs    int `json:"structs_count"`
	CountMethods    int `json:"methods_count"`
	CountImports    int `json:"imports_count"`
}

type Links struct {
	From   int    `json:"from"`
	To     int    `json:"to"`
	Arrows string `json:"arrows"`
}

func getPkgGraph(c *gin.Context) {

	d := c.Query("dir")
	app = FetchApp(d)

	var pi PkgImports

	piSt := make(map[string]interface{})
	uqLinks := make(map[string]Links)

	for i, v := range app.ChildPackages {
		pi.Data = append(pi.Data, PkgNames{
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
				uqLinks[nm] = Links{i, piSt[f].(int), "to"}
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

	for _, v := range GetLocalRepos() {
		ret += fmt.Sprintf(`<a href="?dir=%v">%v</a><br/>`, v, v)
	}

	c.Writer.Write([]byte(ret))
}

func showApp(c *gin.Context) {

	d := c.Query("dir")

	t := time.Now()
	app = FetchApp(d)

	fmt.Println("Time taken: %v", time.Now().Sub(t).Seconds())

	c.Writer.Header().Set("Content-Type", "text/html")
	var ret string
	app.SortPackages()

	for _, pkgs := range app.ChildPackages {

		ret += fmt.Sprintf(`
<h3 class="package_h3">%v</h3>
<div class="package_header">
Folder: %v <br/>
Lines: %v  <br/>
Child Files: %v - 
Funcs : %v, 
Interfaces: %v, 
Structs: %v, 
Methods: %v,
Unique Imports: %v
</div>`,
			pkgs.Name,
			pkgs.FolderLocation,
			pkgs.GetLinesInPkg(),
			len(pkgs.ChildFiles),
			len(pkgs.ChildFuncs),
			len(pkgs.ChildInterfaces),
			len(pkgs.ChildStructs),
			len(pkgs.ChildMethods),
			pkgs.UniqueImports(),
		)

		ret += `<div class="file_details">`

		pkgs.SortFiles()

		ret += `
<div style="float:left;">
<table class="table" style="width: 500px;"><th>
<tr>
<td>File</td>
<td>Action</td>
<td>Loc</td>
<td>Imports</td>
<td>Funcs</td>
<td>Structs</td>
<td>Interfaces</td>
<td>Methods</td>
</tr>
`
		for _, file := range pkgs.ChildFiles {

			var strStructs, strImports, strFuncs string

			for _, v := range file.Structs {
				strStructs += "" + url.QueryEscape(v.Content) + "</br>"
			}

			for _, v := range file.Imports {
				strImports += "" + url.QueryEscape(v) + "</br>"
			}

			for _, v := range file.Funcs {

				strFuncs += "" + url.QueryEscape(v.Name) + " - LOC: " + strconv.Itoa(v.NumberLines) + "</br>"
			}

			mn := ""

			if file.HasMainFunc {
				mn = `<span class="glyphicon glyphicon-star toltip" style="color: orange;">
<span class="tooltiptext">Has a main function</span>
</span>`
			}

			if file.InitFuncCount > 0 {

				mn += fmt.Sprintf(`<span class="glyphicon glyphicon-italic toltip" style="color: blue;">
<span class="tooltiptext">Has init function</span>
</span>`)
			}

			ret += fmt.Sprintf(`<tr> 
<td><b>%v</b></td>
<td class="cnt">%v</td>
<td class="cnt">%v</td>
<td class="cnt"> <div class="importCount" data="%v"> <span style='color: blue;'> %v </span> </div></td> 
<td class="cnt"> <div class="funcsCount" data="%v"> <span style='color: blue'> %v </span> </td> 
<td class="cnt"> <div class="structCount" data="%v"><span style='color: blue;'>%v</span>
</div> </td> <td class="cnt"> %v </td> <td class="cnt"> %v </td></tr>`,
				file.Name,
				mn,
				file.NumberLines,
				strImports, len(file.Imports),
				strFuncs, len(file.Funcs), strStructs, len(file.Structs),
				len(file.Interfaces), len(file.Methods))
		}

		ret += "</table>"

		ret += `
</div>
<div class="file_details_inside" style="float:left; padding: 15px; width: 580px; overflow: scroll;">
</div>
<div style="clear:both"></div>
</div>`
	}

	//	gd, _ := json.Marshal(app.GraphData)

	//	ret += fmt.Sprintf(`<script>
	//	(function a(){
	//	window.graphData = %v;
	//
	//	// create an array with nodes
	//	var nodes = new vis.DataSet(window.graphData.data);
	//	// create an array with edges
	//	var edges = new vis.DataSet(window.graphData.links);
	//
	//	// create a network
	//	var container = document.getElementById("mynetwork");
	//	var data = {
	//		nodes: nodes,
	//		edges: edges
	//	};
	//
	//	var options = {
	//		nodes: {
	//			shape: "dot"
	//		} };
	//
	//	var network = new vis.Network(container, data, options);
	//
	//
	//
	//})();
	//</script>`, string(gd))

	ret += "</body></html>"

	c.Writer.Write([]byte(ret))
}

type JsonOut struct {
	PkgName    string     `json:"name"`
	DataPoints []FileSize `json:"children"`
}

type FileSize struct {
	FileName string `json:"name"`
	Size     int    `json:"size"`
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
