package options

var (
	Current = NewOptions()
)

func NewOptions() *Options {
	options := new(Options)

	return options
}

type Options struct {
	Global           bool
	ParentRepository string
	ParentPath       string
	Minimalistic     bool
}
