package archive

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/holgerhuo/gobackup/config"
	"github.com/holgerhuo/gobackup/helper"
)

// Run archive
func Run(model config.ModelConfig) (err error) {
	if model.Archive == nil {
		return nil
	}

	slog.Info("Starting archive creation", 
		"component", "archive",
		"model", model.Name)

	helper.MkdirP(model.DumpPath)

	includes := model.Archive.GetStringSlice("includes")
	includes = cleanPaths(includes)

	excludes := model.Archive.GetStringSlice("excludes")
	excludes = cleanPaths(excludes)

	if len(includes) == 0 {
		return fmt.Errorf("archive.includes have no config")
	}
	slog.Info("Archive configuration", 
		"component", "archive",
		"model", model.Name,
		"includeRules", len(includes),
		"excludeRules", len(excludes))

	opts := options(model.DumpPath, excludes, includes)
	helper.Exec("tar", opts...)

	slog.Info("Archive creation completed", 
		"component", "archive",
		"model", model.Name,
		"archivePath", filepath.Join(model.DumpPath, "archive.tar"))

	return nil
}

func options(dumpPath string, excludes, includes []string) (opts []string) {
	tarPath := filepath.Join(dumpPath, "archive.tar")
	if helper.IsGnuTar {
		opts = append(opts, "--ignore-failed-read")
	}
	opts = append(opts, "-cPf", tarPath)

	for _, exclude := range excludes {
		opts = append(opts, "--exclude="+filepath.Clean(exclude))
	}

	opts = append(opts, includes...)

	return opts
}

func cleanPaths(paths []string) (results []string) {
	for _, p := range paths {
		results = append(results, filepath.Clean(p))
	}
	return
}
