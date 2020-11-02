package web

import "github.com/gin-gonic/gin"

func StartUI() {

	r := gin.Default()
	r.GET("/show_app", showApp)

	r.GET("/show_app_v2", showAppV2)
	r.GET("/show_app_from_db_v2", showFromDB)

	r.GET("/load_dirs", loadDirs)
	r.GET("/get_package_graph", getPkgGraph)

	r.GET("/json_file_size", jsonFS)
	r.GET("/json_pkg_size", jsonPackageSize)
	r.GET("/json_func_count", jsonMethodFuncs)

	r.Static("/static", "./static")

	r.Run("0.0.0.0:8081") // listen and serve on 0.0.0.0:8080

}

type JsonOut struct {
	PkgName    string     `json:"name"`
	DataPoints []FileSize `json:"children"`
}

type FileSize struct {
	FileName string `json:"name"`
	Size     int    `json:"size"`
}
