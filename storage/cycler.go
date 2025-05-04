package storage

import (
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/holgerhuo/gobackup/config"
	"github.com/holgerhuo/gobackup/helper"
)

type PackageList []Package

type Package struct {
	FileKey   string    `json:"file_key"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	cyclerPath = path.Join(config.HomeDir, ".gobackup/cycler")
)

type Cycler struct {
	packages PackageList
	isLoaded bool
}

func (c *Cycler) add(fileKey string) {
	c.packages = append(c.packages, Package{
		FileKey:   fileKey,
		CreatedAt: time.Now(),
	})
}

func (c *Cycler) shiftByKeep(keep int) (first *Package) {
	total := len(c.packages)
	if total <= keep {
		return nil
	}

	first, c.packages = &c.packages[0], c.packages[1:]
	return
}

func (c *Cycler) run(model string, fileKey string, keep int, deletePackage func(fileKey string) error) {
	cyclerFileName := path.Join(cyclerPath, model+".json")

	c.load(cyclerFileName)
	c.add(fileKey)
	defer c.save(cyclerFileName)

	if keep == 0 {
		return
	}

	for {
		pkg := c.shiftByKeep(keep)
		if pkg == nil {
			break
		}

		err := deletePackage(pkg.FileKey)
		if err != nil {
			slog.Warn("Package removal failed",
				"component", "storage.cycler",
				"model", model,
				"fileKey", pkg.FileKey,
				"error", err)
		}
	}
}

func (c *Cycler) load(cyclerFileName string) {
	helper.MkdirP(cyclerPath)

	// write example JSON if not exist
	if !helper.IsExistsPath(cyclerFileName) {
		ioutil.WriteFile(cyclerFileName, []byte("[{}]"), os.ModePerm)
	}

	f, err := ioutil.ReadFile(cyclerFileName)
	if err != nil {
		slog.Error("Failed to load cycler file",
			"component", "storage.cycler",
			"file", cyclerFileName,
			"error", err)
		return
	}
	
	err = json.Unmarshal(f, &c.packages)
	if err != nil {
		slog.Error("Failed to unmarshal cycler data",
			"component", "storage.cycler",
			"file", cyclerFileName,
			"error", err)
	}
	c.isLoaded = true
}

func (c *Cycler) save(cyclerFileName string) {
	if !c.isLoaded {
		slog.Warn("Skipping cycler save - not loaded",
			"component", "storage.cycler",
			"file", cyclerFileName)
		return
	}

	data, err := json.Marshal(&c.packages)
	if err != nil {
		slog.Error("Failed to marshal cycler data",
			"component", "storage.cycler",
			"file", cyclerFileName,
			"error", err)
		return
	}

	err = ioutil.WriteFile(cyclerFileName, data, os.ModePerm)
	if err != nil {
		slog.Error("Failed to save cycler file",
			"component", "storage.cycler",
			"file", cyclerFileName,
			"error", err)
		return
	}
	
	slog.Debug("Cycler file saved successfully",
		"component", "storage.cycler",
		"file", cyclerFileName,
		"packageCount", len(c.packages))
}
