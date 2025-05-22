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
	"github.com/holgerhuo/gobackup/storage"
)

// Model represents a backup model with its configuration.
type Model struct {
	Config config.ModelConfig
}

// Perform executes the backup process for the model.
func (m *Model) Perform() {
	slog.Info("Backup model starting",
		"component", "model",
		"model", m.Config.Name,
		"workDir", m.Config.DumpPath,
	)

	// Ensure cleanup is always called, even on panic.
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic occurred during backup model execution",
				"component", "model",
				"model", m.Config.Name,
				"error", r,
			)
		}
		m.cleanup()
	}()

	if err := m.runScript(m.Config.BeforeScript, "before"); err != nil {
		slog.Error("Before script execution failed",
			"component", "model",
			"model", m.Config.Name,
			"error", err,
		)
	}

	if err := database.Run(m.Config); err != nil {
		slog.Error("Database backup failed",
			"component", "model",
			"model", m.Config.Name,
			"error", err,
		)
		return
	}

	if m.Config.Archive != nil {
		if err := archive.Run(m.Config); err != nil {
			slog.Error("Archive creation failed",
				"component", "model",
				"model", m.Config.Name,
				"error", err,
			)
			return
		}
	}

	archivePath, err := compressor.Run(m.Config)
	if err != nil {
		slog.Error("Compression failed",
			"component", "model",
			"model", m.Config.Name,
			"error", err,
		)
		return
	}

	archivePath, err = encryptor.Run(archivePath, m.Config)
	if err != nil {
		slog.Error("Encryption failed",
			"component", "model",
			"model", m.Config.Name,
			"error", err,
		)
		return
	}

	if err := storage.Run(m.Config, archivePath); err != nil {
		slog.Error("Storage operation failed",
			"component", "model",
			"model", m.Config.Name,
			"error", err,
		)
		return
	}
}

// runScript executes a shell script if provided.
func (m *Model) runScript(script string, stage string) error {
	if len(script) == 0 {
		return nil
	}
	slog.Info("Executing "+stage+" script",
		"component", "model",
		"model", m.Config.Name,
	)
	_, err := helper.ExecWithStdio(script, true)
	return err
}

// cleanup removes temporary files and runs the after script.
func (m *Model) cleanup() {
	slog.Info("Cleaning up temporary files",
		"component", "model",
		"model", m.Config.Name,
		"tempPath", m.Config.TempPath,
	)

	if err := os.RemoveAll(m.Config.TempPath); err != nil {
		slog.Error("Cleanup failed",
			"component", "model",
			"model", m.Config.Name,
			"tempPath", m.Config.TempPath,
			"error", err,
		)
	}

	if err := m.runScript(m.Config.AfterScript, "after"); err != nil {
		slog.Error("After script execution failed",
			"component", "model",
			"model", m.Config.Name,
			"error", err,
		)
	}

	slog.Info("Backup model completed",
		"component", "model",
		"model", m.Config.Name,
	)
}
