package storage

import (
	"log/slog"
	"path"

	"github.com/holgerhuo/gobackup/helper"
)

// Local storage
//
// type: local
// path: /data/backups
type Local struct {
	Base
	destPath string
}

func (ctx *Local) open() (err error) {
	ctx.destPath = ctx.model.StoreWith.Viper.GetString("path")
	helper.MkdirP(ctx.destPath)
	return
}

func (ctx *Local) close() {}

func (ctx *Local) upload(fileKey string) (err error) {
	_, err = helper.Exec("cp", ctx.archivePath, ctx.destPath)
	if err != nil {
		slog.Error("Local storage upload failed",
			"component", "storage",
			"type", "local",
			"model", ctx.model.Name,
			"source", ctx.archivePath,
			"destination", ctx.destPath,
			"error", err)
		return err
	}
	
	slog.Info("Local storage upload successful",
		"component", "storage",
		"type", "local",
		"model", ctx.model.Name,
		"destination", ctx.destPath)
	return nil
}

func (ctx *Local) delete(fileKey string) (err error) {
	filePath := path.Join(ctx.destPath, fileKey)
	_, err = helper.Exec("rm", filePath)
	if err != nil {
		slog.Error("Local storage file deletion failed",
			"component", "storage",
			"type", "local",
			"model", ctx.model.Name,
			"file", filePath,
			"error", err)
	} else {
		slog.Debug("Local storage file deleted",
			"component", "storage",
			"type", "local",
			"model", ctx.model.Name,
			"file", filePath)
	}
	return
}
