// Package state holds state
package state

import (
	"github.com/mmycroft/gator/internal/config"
	"github.com/mmycroft/gator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}
