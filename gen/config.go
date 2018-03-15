package gen

// Config for the running of the commands
type Config struct {
	PkgName   string
	OutFolder string
	BaseDir   string
	Tags      []string
	NoTests   bool
	NoHooks   bool
	Wipe      bool
}
