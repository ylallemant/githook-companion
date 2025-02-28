package options

var (
	Domain  = "githook-companion"
	Current = NewOptions()
)

func NewOptions() *Options {
	options := new(Options)

	return options
}

type Options struct {
	Semver    bool
	Commit    bool
	Separator string
}
