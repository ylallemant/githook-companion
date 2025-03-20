package options

var (
	Current = NewOptions()
)

func NewOptions() *Options {
	options := new(Options)

	return options
}

type Options struct {
	Global              bool
	ReferenceRepository string
	ReferencePath       string
	Minimalistic        bool
}
