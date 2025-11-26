package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

type ThumbnailCache struct {
	cacheDir string
}

func New(cacheDir string) *ThumbnailCache {
	os.MkdirAll(cacheDir, os.ModePerm)

	return &ThumbnailCache{cacheDir: cacheDir}
}

func (c *ThumbnailCache) GetPath(videoID string) string {
	return filepath.Join(c.cacheDir, fmt.Sprintf("%s.jpg", videoID))
}

func (c *ThumbnailCache) Exists(videoID string) bool {
	_, err := os.Stat(c.GetPath(videoID))

	return err == nil
}

func (c *ThumbnailCache) Save(videoID string, data []byte) error {
	return os.WriteFile(c.GetPath(videoID), data, os.ModePerm)
}

func (c *ThumbnailCache) GetRenderedPath(videoID string, width int, height int) string {
	key := fmt.Sprintf("%s_%dx%d.jpg", videoID, width, height)
	hash := sha256.Sum256([]byte(key))
	hashStr := hex.EncodeToString(hash[:8])

	return filepath.Join(c.cacheDir, fmt.Sprintf("%s_rendered.txt", hashStr))
}

func (c *ThumbnailCache) GetRendered(videoID string, width, height int) (string, bool) {
	path := c.GetRenderedPath(videoID, width, height)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func (c *ThumbnailCache) SaveRendered(videoID string, width, height int, content string) error {
	path := c.GetRenderedPath(videoID, width, height)
	return os.WriteFile(path, []byte(content), 0644)
}

func (c *ThumbnailCache) Clear() error {
	entries, err := os.ReadDir(c.cacheDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		os.Remove(filepath.Join(c.cacheDir, entry.Name()))
	}
	return nil
}
