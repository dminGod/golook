package web

import (
	"fmt"
	"github.com/dminGod/golook/parsers"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
)

func showApp(c *gin.Context) {

	d := c.Query("dir")

	//t := time.Now()
	app = parsers.FetchApp(d)

	//fmt.Println("Time taken: %v", time.Now().Sub(t).Seconds())

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

	ret += "</body></html>"

	c.Writer.Write([]byte(ret))
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
