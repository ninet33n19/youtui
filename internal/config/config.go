package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	CacheDir    string
	DownloadDir string
	MaxResults  int
}

func Default() *Config {
	cacheDir := filepath.Join(os.TempDir(), "youtui", "thumbnails")
	currentWorkingDir, _ := os.Getwd() // fallback to current working directory
	downloadDir := filepath.Join(currentWorkingDir, "downloads")

	return &Config{
		CacheDir:    cacheDir,
		DownloadDir: downloadDir,
		MaxResults:  10,
	}
}

func (c *Config) EnsureDirs() error {
	if err := os.MkdirAll(c.CacheDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(c.DownloadDir, os.ModePerm); err != nil {
		return err
	}

	return nil
}
