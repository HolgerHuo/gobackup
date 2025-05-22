package compressor

import (
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/holgerhuo/gobackup/config"
	"github.com/spf13/viper"
)

// Base compressor
type Base struct {
	name  string
	model config.ModelConfig
	viper *viper.Viper
}

// Context compressor
type Context interface {
	perform() (archivePath string, err error)
}

func (ctx *Base) archiveFilePath(ext string) string {
	return path.Join(ctx.model.TempPath, time.Now().Format("2006.01.02.15.04.05")+ext)
}

func newBase(model config.ModelConfig) (base Base) {
	base = Base{
		name:  model.Name,
		model: model,
		viper: model.CompressWith.Viper,
	}
	return
}

// Run compressor
func Run(model config.ModelConfig) (archivePath string, err error) {
	base := newBase(model)

	var ctx Context
	switch model.CompressWith.Type {
	case "tgz":
		ctx = &Tgz{Base: base}
	case "zstd":
		ctx = &Zstd{Base: base}
	default:
		ctx = &Zstd{Base: base}
	}

	slog.Info("starting compression", 
		"component", "compressor",
		"model", model.Name,
		"type", model.CompressWith.Type)

	// set workdir
	os.Chdir(path.Join(model.DumpPath, "../"))
	archivePath, err = ctx.perform()
	if err != nil {
		return
	}
	slog.Info("compression completed", 
		"component", "compressor",
		"model", model.Name,
		"type", model.CompressWith.Type,
		"archivePath", archivePath)

	return
}
