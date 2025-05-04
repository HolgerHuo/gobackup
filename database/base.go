package database

import (
	"fmt"
	"log/slog"
	"path"

	"github.com/holgerhuo/gobackup/config"
	"github.com/holgerhuo/gobackup/helper"
	"github.com/spf13/viper"
)

// Base database
type Base struct {
	model    config.ModelConfig
	dbConfig config.SubConfig
	viper    *viper.Viper
	name     string
	dumpPath string
}

// Context database interface
type Context interface {
	perform() error
}

func newBase(model config.ModelConfig, dbConfig config.SubConfig) (base Base) {
	base = Base{
		model:    model,
		dbConfig: dbConfig,
		viper:    dbConfig.Viper,
		name:     dbConfig.Name,
	}
	base.dumpPath = path.Join(model.DumpPath, dbConfig.Type, base.name)
	helper.MkdirP(base.dumpPath)
	return
}

// New - initialize Database
func runModel(model config.ModelConfig, dbConfig config.SubConfig) (err error) {
	base := newBase(model, dbConfig)
	var ctx Context
	switch dbConfig.Type {
	case "mysql":
		ctx = &MySQL{Base: base}
	case "redis":
		ctx = &Redis{Base: base}
	case "postgresql":
		ctx = &PostgreSQL{Base: base}
	default:
		err = fmt.Errorf("model: %s databases.%s config `type: %s`, but is not implement", model.Name, dbConfig.Name, dbConfig.Type)
		slog.Warn("Unsupported database type", 
			"component", "database",
			"model", model.Name,
			"database", dbConfig.Name,
			"type", dbConfig.Type,
			"error", err)
		return
	}

	slog.Info("Database operation starting", 
		"component", "database",
		slog.Group("database", 
			"type", dbConfig.Type,
			"name", base.name,
		),
		"model", model.Name)

	// perform
	err = ctx.perform()
	if err != nil {
		return err
	}
	// Log successful completion
	slog.Debug("Database operation completed", 
		"component", "database",
		"type", dbConfig.Type,
		"name", base.name,
		"model", model.Name)

	return
}

// Run databases
func Run(model config.ModelConfig) error {
	if len(model.Databases) == 0 {
		return nil
	}

	slog.Info("Starting database backups", 
		"component", "database",
		"model", model.Name,
		"count", len(model.Databases))
	for _, dbCfg := range model.Databases {
		err := runModel(model, dbCfg)
		if err != nil {
			return err
		}
	}
	slog.Info("Database backups completed", 
		"component", "database",
		"model", model.Name)

	return nil
}
