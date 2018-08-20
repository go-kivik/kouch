package kouch

// Common command line flags
const (
	// Curl-equivalent flags
	FlagVerbose    = "verbose"
	FlagOutputFile = "output"
	FlagData       = "data"
	FlagHead       = "head"
	FlagDumpHeader = "dump-header"

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
	FlagAutoRev      = "auto-rev"

	// Curl-equivalent short flags
	FlagShortVerbose    = "v"
	FlagShortOutputFile = "o"
	FlagShortData       = "d"
	FlagShortHead       = "I"
	FlagShortDumpHeader = "D"

	// Short versions, custom
	FlagShortServerRoot   = "S"
	FlagShortOutputFormat = "F"
	FlagShortRev          = "r"
	FlagShortAutoRev      = "R"
)
