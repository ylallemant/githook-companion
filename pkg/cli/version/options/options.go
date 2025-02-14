package options

var (
	Domain  = "git-butler"
	Current = NewOptions()
)

func NewOptions() *Options {
	options := new(Options)

	return options
}

type Options struct {
	Semver bool
	Commit bool
}
