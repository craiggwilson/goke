package build

var base = ""
var mainFile = "build.go"
var buildOutputFile = "build.exe"
var packages = []string{
	"./task",
	"./task/command",
	"./task/internal",
}
