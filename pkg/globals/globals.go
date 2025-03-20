package globals

var (
	Current = new(Globals)
)

type Globals struct {
	ConfigPath     string
	FallbackConfig bool
}
