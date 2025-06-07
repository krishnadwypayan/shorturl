package main

import (
	"github.com/krishnadwypayan/shorturl/internal/encoder"
	"github.com/krishnadwypayan/shorturl/internal/logger"
	"github.com/krishnadwypayan/shorturl/internal/snowflake"
)

func main() {
	generator := snowflake.NewGenerator(1)
	id := generator.NextString()
	logger.Info().Msgf("next ID: %s", id)
	logger.Info().Msgf("next ID decoded: %d", encoder.DecodeBase62(id))
}
