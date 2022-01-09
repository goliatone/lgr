package version

var (
	//Tag is the tagged version e.g. v0.0.1
	Tag = "dev"
	//Time of build
	Time string
	//User that built the package
	User string
	//Commit hash
	Commit string
)

//GetVersion returns version string
func GetVersion() string {
	return Tag + "-" + Time + ":" + User
}
