package encryptor

import (
	"fmt"
	"github.com/holgerhuo/gobackup/helper"
)

// OpenSSL encryptor for use openssl aes-256-cbc
//
// - base64: false
// - salt: true
// - password:
type OpenSSL struct {
	Base
	salt     bool
	base64   bool
	iter     int
	pbkdf2   bool
	password string
}

func (ctx *OpenSSL) perform() (encryptPath string, err error) {
	sslViper := ctx.viper
	sslViper.SetDefault("salt", true)
	sslViper.SetDefault("base64", false)
    sslViper.SetDefault("iter", 10000)
	sslViper.SetDefault("pbkdf2", true)

	ctx.salt = sslViper.GetBool("salt")
	ctx.base64 = sslViper.GetBool("base64")
	ctx.iter = sslViper.GetInt("iter")
	ctx.pbkdf2 = sslViper.GetBool("pbkdf2")
	ctx.password = sslViper.GetString("password")

	if len(ctx.password) == 0 {
		err = fmt.Errorf("password option is required")
		return
	}

	encryptPath = ctx.archivePath + ".enc"

	opts := ctx.options()
	opts = append(opts, "-in", ctx.archivePath, "-out", encryptPath)
	_, err = helper.Exec("openssl", opts...)
	return
}

func (ctx *OpenSSL) options() (opts []string) {
	opts = append(opts, "aes-256-cbc")
	if ctx.base64 {
		opts = append(opts, "-base64")
	}
	if ctx.salt {
		opts = append(opts, "-salt")
	}
	if ctx.iter > 0 {
		opts = append(opts, "-iter", fmt.Sprintf("%d", ctx.iter))
	}
	if ctx.pbkdf2 {
		opts = append(opts, "-pbkdf2")
	}
	opts = append(opts, `-k`, ctx.password)
	return opts
}
