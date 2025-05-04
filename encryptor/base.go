package encryptor

import (
	"log/slog"

	"github.com/holgerhuo/gobackup/config"
	"github.com/spf13/viper"
)

// Base encryptor
type Base struct {
	model       config.ModelConfig
	viper       *viper.Viper
	archivePath string
}

// Context encryptor interface
type Context interface {
	perform() (encryptPath string, err error)
}

func newBase(archivePath string, model config.ModelConfig) (base Base) {
	base = Base{
		archivePath: archivePath,
		model:       model,
		viper:       model.EncryptWith.Viper,
	}
	return
}

// Run compressor
func Run(archivePath string, model config.ModelConfig) (encryptPath string, err error) {
	base := newBase(archivePath, model)
	var ctx Context
	switch model.EncryptWith.Type {
	case "openssl":
		ctx = &OpenSSL{Base: base}
	default:
		encryptPath = archivePath
		return
	}

	slog.Info("Starting encryption", 
		"component", "encryptor",
		"model", model.Name,
		"type", model.EncryptWith.Type,
		"sourcePath", archivePath)
	encryptPath, err = ctx.perform()
	if err != nil {
		return
	}
	slog.Info("Encryption completed", 
		"component", "encryptor",
		"model", model.Name,
		"type", model.EncryptWith.Type,
		"encryptPath", encryptPath)

	return
}
