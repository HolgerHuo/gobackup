package compressor

import (
	"github.com/holgerhuo/gobackup/helper"
)

// Zstd .tar.zst compressor
//
// type: zstd
type Zstd struct {
	Base
}

func (ctx *Zstd) perform() (archivePath string, err error) {
	filePath := ctx.archiveFilePath(".tar.zst")

	opts := ctx.options()
	opts = append(opts, filePath)
	opts = append(opts, ctx.name)

	_, err = helper.Exec("tar", opts...)
	if err == nil {
		archivePath = filePath
		return
	}
	return
}

func (ctx *Zstd) options() (opts []string) {
	if helper.IsGnuTar {
		opts = append(opts, "--ignore-failed-read")
	}
	opts = append(opts, "-acf")

	return
}
