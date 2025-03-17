package buildvar

var Tag = "no tag"
var (
	debug = "true"
	Debug = true
)

var (
	isCheatsAvailable = "true"
	IsCheatsAvailable = true
)

func init() {
	if debug != "true" {
		Debug = false
	}

	if isCheatsAvailable != "true" {
		IsCheatsAvailable = false
	}
}
