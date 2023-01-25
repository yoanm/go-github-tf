package core

import (
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func configTrace(header string, c *GhRepoConfig) {
	if zerolog.GlobalLevel() == zerolog.TraceLevel {
		encoded, encodeError := yaml.Marshal(*c)
		if encodeError == nil {
			log.Trace().Msgf("%s\n%s", header, string(encoded))
		} else {
			log.Trace().Msgf("%s Error %s", header, encodeError)
		}
	}
}
