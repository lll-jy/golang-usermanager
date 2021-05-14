// Package paths handles file paths that are used throughout the app.
package paths

var PlaceholderPath string
//var FileBasePath2 string // May set to some cloud space
var FileBasePath string
var TempPath string // On local device

func SetupPaths(src string) {
	switch src {
	case "main":
		PlaceholderPath = "assets/placeholder.jpeg"
		//FileBasePath2 = "../../../Desktop/EntryTask/entry-task/test/data/upload"
		FileBasePath = "test/data/upload"
		//FileBasePath2 = ""
		//FileBasePath = "../../../../tmp"
		TempPath = "test/data/temp"
		break
	case "test":
		PlaceholderPath = "../assets/placeholder.jpeg"
		//FileBasePath2 = "../../../../Desktop/EntryTask/entry-task/test/data/test_upload"
		FileBasePath = "data/test_upload"
		//FileBasePath2 = ""
		//FileBasePath = "../../../../tmp"
		TempPath = "data/test_temp"
		break
	}
}
