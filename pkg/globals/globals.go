package globals

import "github.com/rs/zerolog"

var (
	Current = new(Globals)
)

type Globals struct {
	ConfigPath     string
	FallbackConfig bool
	Debug          bool
	LogLevel       string
}

func ProcessGlobals() {
	if Current.Debug {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if Current.LogLevel != "" {
		switch Current.LogLevel {
		default:
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		}
	}
}
