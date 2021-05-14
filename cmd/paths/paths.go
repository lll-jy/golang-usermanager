// Package paths handles file paths that are used throughout the app.
package paths

var PlaceholderPath string
var FileBasePath string
var TempPath string // On local device

func SetupPaths(src string) {
	switch src {
	case "main":
		PlaceholderPath = "assets/placeholder.jpeg"
		FileBasePath = "test/data/upload"
		TempPath = "test/data/temp"
		break
	case "test":
		PlaceholderPath = "../assets/placeholder.jpeg"
		FileBasePath = "data/test_upload"
		TempPath = "data/test_temp"
		break
	}
}
