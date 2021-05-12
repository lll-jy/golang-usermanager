package paths

var PlaceholderPath string
var FileBasePath string // EXTEND: May set to some cloud space
var FileBaseRelativePath string
var TempPath string

func SetupPaths(src string) {
	switch src {
	case "main":
		//PlaceholderPath = "assets/placeholder.jpeg"
		//PlaceholderPath = "bakcupAssets/placeholder.jpeg"
		PlaceholderPath = "test/data/original/placeholder.jpeg"
		FileBasePath = "../../../Desktop/EntryTask/entry-task/test/data/upload"
		FileBaseRelativePath = "test/data/upload"
		TempPath = "test/data/temp"
		break
	case "test":
		PlaceholderPath = "../assets/placeholder.jpeg"
		FileBasePath = "../../../../Desktop/EntryTask/entry-task/test/data/upload"
		FileBaseRelativePath = "data/upload"
		TempPath = "data/temp"
		break
	}
}
