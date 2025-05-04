package storage

import (
	"log/slog"

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
