#export GOBIN="/usr/local/bin/"
#go build golook.go structs.go utils.go
#go install golook.go structs.go utils.go

go build -o golook.exe application.go config.go file_info.go golook.go info_structs.go load_data.go local_db.go package_info.go parse_file.go utils.go web_show_app.go