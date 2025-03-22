package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		// TODO find a solution for bad output "../../../../../../../../pkg/nlp/tokenizer.go:229"
		// start := strings.Index(file, "/pkg")
		// if start > -1 {
		// 	return file[start:] + ":" + strconv.Itoa(line)
		// }
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	if err := cli.Command().Execute(); err != nil {
		log.Fatal().Err(err).Msg("error during execution")
	}
}
