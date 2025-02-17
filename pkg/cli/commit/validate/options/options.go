package options

var (
	Current = NewOptions()
)

func NewOptions() *Options {
	options := new(Options)

	return options
}

type Options struct {
	Message        string
	OutputFilePath string
}
