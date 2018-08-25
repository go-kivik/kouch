package kouch

// Common command line flags
const (
	// Curl-equivalent flags
	FlagVerbose    = "verbose"
	FlagOutputFile = "output"
	FlagData       = "data"
	FlagHead       = "head"

	// Custom flags
	FlagClobber      = "force"
	FlagConfigFile   = "kouchconfig"
	FlagServerRoot   = "root"
	FlagDataJSON     = "data-json"
	FlagDataYAML     = "data-yaml"
	FlagOutputFormat = "output-format"
	FlagFilename     = "filename"
	FlagDocument     = "id"
	FlagDatabase     = "database"
	FlagFullCommit   = "full-commit"
	FlagIfNoneMatch  = "if-none-match"
	FlagRev          = "rev"

	// Curl-equivalent short flags
	FlagShortVerbose    = "v"
	FlagShortOutputFile = "o"
	FlagShortData       = "d"
	FlagShortHead       = "I"

	// Short versions, custom
	FlagShortServerRoot   = "R"
	FlagShortOutputFormat = "F"
	FlagShortRev          = "r"
)
