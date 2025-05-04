package model

import (
	"log/slog"
	"os"

	"github.com/holgerhuo/gobackup/archive"
	"github.com/holgerhuo/gobackup/compressor"
	"github.com/holgerhuo/gobackup/config"
	"github.com/holgerhuo/gobackup/database"
	"github.com/holgerhuo/gobackup/encryptor"
	"github.com/holgerhuo/gobackup/helper"
	"github.com/holgerhuo/gobackup/notifier"
	"github.com/holgerhuo/gobackup/storage"
)

// Model class
type Model struct {
	Config config.ModelConfig
}

// Perform model
func (ctx Model) Perform() {
	slog.Info("Backup model starting", 
		"component", "model",
		"model", ctx.Config.Name,
		"workDir", ctx.Config.DumpPath)

	var err error

	if len(ctx.Config.BeforeScript) > 0 {
		slog.Info("Executing before script", 
			"component", "model",
			"model", ctx.Config.Name)
		_, err := helper.ExecWithStdio(ctx.Config.BeforeScript, true)
		if err != nil {
			slog.Error("Before script execution failed", 
				"component", "model",
				"model", ctx.Config.Name,
				"error", err)
		}
	}

	defer func() {
		if err != nil {
			notifier.Failure(ctx.Config, err.Error())
		} else {
			notifier.Success(ctx.Config)
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			ctx.cleanup()
		}

		ctx.cleanup()
	}()

	err = database.Run(ctx.Config)
	if err != nil {
		slog.Error("Database backup failed", 
			"component", "model",
			"model", ctx.Config.Name,
			"error", err)
		return
	}

	if ctx.Config.Archive != nil {
		err = archive.Run(ctx.Config)
		if err != nil {
			slog.Error("Archive creation failed", 
				"component", "model",
				"model", ctx.Config.Name,
				"error", err)
			return
		}
	}

	archivePath, err := compressor.Run(ctx.Config)
	if err != nil {
		slog.Error("Compression failed", 
			"component", "model",
			"model", ctx.Config.Name,
			"error", err)
		return
	}

	archivePath, err = encryptor.Run(archivePath, ctx.Config)
	if err != nil {
		slog.Error("Encryption failed", 
			"component", "model",
			"model", ctx.Config.Name,
			"error", err)
		return
	}

	err = storage.Run(ctx.Config, archivePath)
	if err != nil {
		slog.Error("Storage operation failed", 
			"component", "model",
			"model", ctx.Config.Name,
			"error", err)
		return
	}

}

// Cleanup model temp files
func (ctx Model) cleanup() {
	slog.Info("Cleaning up temporary files", 
		"component", "model",
		"model", ctx.Config.Name,
		"tempPath", ctx.Config.TempPath)
	
	err := os.RemoveAll(ctx.Config.TempPath)
	if err != nil {
		slog.Error("Cleanup failed", 
			"component", "model",
			"model", ctx.Config.Name,
			"tempPath", ctx.Config.TempPath,
			"error", err)
	}

	if len(ctx.Config.AfterScript) > 0 {
		slog.Info("Executing after script", 
			"component", "model",
			"model", ctx.Config.Name)
		_, err := helper.ExecWithStdio(ctx.Config.AfterScript, true)
		if err != nil {
			slog.Error("After script execution failed", 
				"component", "model",
				"model", ctx.Config.Name,
				"error", err)
		}
	}

	slog.Info("Backup model completed", 
		"component", "model",
		"model", ctx.Config.Name)
}
