package storage

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/holgerhuo/gobackup/config"
	"github.com/spf13/viper"
)

// Base storage
type Base struct {
	model       config.ModelConfig
	archivePath string
	viper       *viper.Viper
	keep        int
}

// Context storage interface
type Context interface {
	open() error
	close()
	upload(fileKey string) error
	delete(fileKey string) error
}

func newBase(model config.ModelConfig, archivePath string) (base Base) {
	base = Base{
		model:       model,
		archivePath: archivePath,
		viper:       model.StoreWith.Viper,
	}

	if base.viper != nil {
		base.keep = base.viper.GetInt("keep")
	}

	return
}

// Run storage
func Run(model config.ModelConfig, archivePath string) (err error) {
	slog.Info("Starting storage operation", 
		"component", "storage",
		"model", model.Name)
	
	newFileKey := filepath.Base(archivePath)
	base := newBase(model, archivePath)
	var ctx Context
	switch model.StoreWith.Type {
	case "local":
		ctx = &Local{Base: base}
	case "s3":
		ctx = &S3{Base: base}
	default:
		return fmt.Errorf("[%s] storage type has not implement", model.StoreWith.Type)
	}

	slog.Info("Storage operation details", 
		"component", "storage",
		"type", model.StoreWith.Type,
		"model", model.Name,
		"fileKey", newFileKey)
	err = ctx.open()
	if err != nil {
		return err
	}
	defer ctx.close()

	err = ctx.upload(newFileKey)
	if err != nil {
		return err
	}

	cycler := Cycler{}
	cycler.run(model.Name, newFileKey, base.keep, ctx.delete)

	slog.Info("Storage operation completed", 
		"component", "storage",
		"type", model.StoreWith.Type,
		"model", model.Name)
	return nil
}
