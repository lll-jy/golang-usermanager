// Package paths handles file paths that are used throughout the app.
package paths

var PlaceholderPath string
var FileBasePath string // May set to some cloud space
var FileBaseRelativePath string
var TempPath string // On local device

func SetupPaths(src string) {
	switch src {
	case "main":
		PlaceholderPath = "assets/placeholder.jpeg"
		FileBasePath = "../../../Desktop/EntryTask/entry-task/test/data/upload"
		FileBaseRelativePath = "test/data/upload"
		TempPath = "test/data/temp"
		break
	case "test":
		PlaceholderPath = "../assets/placeholder.jpeg"
		FileBasePath = "../../../../Desktop/EntryTask/entry-task/test/data/test_upload"
		FileBaseRelativePath = "data/test_upload"
		TempPath = "data/test_temp"
		break
	}
}
