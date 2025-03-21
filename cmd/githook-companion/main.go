package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ylallemant/githook-companion/pkg/cli"
)

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}

	zerolog.SetGlobalLevel(zerolog.FatalLevel)

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	if err := cli.Command().Execute(); err != nil {
		log.Fatal().Err(err).Msg("error during execution")
	}
}
