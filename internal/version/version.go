package version

// These values are set via -ldflags at build time.
// Defaults are useful for local dev.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func Get() Info {
	return Info{Version: Version, Commit: Commit, Date: Date}
}
