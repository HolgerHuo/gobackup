package model

import (
	"os"

	"github.com/holgerhuo/gobackup/archive"
	"github.com/holgerhuo/gobackup/compressor"
	"github.com/holgerhuo/gobackup/config"
	"github.com/holgerhuo/gobackup/database"
	"github.com/holgerhuo/gobackup/encryptor"
	"github.com/holgerhuo/gobackup/helper"
	"github.com/holgerhuo/gobackup/logger"
	"github.com/holgerhuo/gobackup/notifier"
	"github.com/holgerhuo/gobackup/storage"
)

// Model class
type Model struct {
	Config config.ModelConfig
}

// Perform model
func (ctx Model) Perform() {
	logger.Info("======== " + ctx.Config.Name + " ========")
	logger.Info("WorkDir:", ctx.Config.DumpPath+"\n")

	var err error

	if len(ctx.Config.BeforeScript) > 0 {
		logger.Info("Executing before_script...")
		_, err := helper.ExecWithStdio(ctx.Config.BeforeScript, true)
		if err != nil {
			logger.Error(err)
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
		logger.Error(err)
		return
	}

	if ctx.Config.Archive != nil {
		err = archive.Run(ctx.Config)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	archivePath, err := compressor.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	archivePath, err = encryptor.Run(archivePath, ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	err = storage.Run(ctx.Config, archivePath)
	if err != nil {
		logger.Error(err)
		return
	}

}

// Cleanup model temp files
func (ctx Model) cleanup() {
	logger.Info("Cleanup temp: " + ctx.Config.TempPath + "/\n")
	err := os.RemoveAll(ctx.Config.TempPath)
	if err != nil {
		logger.Error("Cleanup temp dir "+ctx.Config.TempPath+" error:", err)
	}

	if len(ctx.Config.AfterScript) > 0 {
		logger.Info("Executing after_script...")
		_, err := helper.ExecWithStdio(ctx.Config.AfterScript, true)
		if err != nil {
			logger.Error(err)
		}
	}

	logger.Info("======= End " + ctx.Config.Name + " =======\n\n")
}
