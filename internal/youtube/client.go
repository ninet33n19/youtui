package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

type Client struct {
	maxResults int
}

func NewClient(maxResults int) *Client {
	return &Client{maxResults: maxResults}
}

func (c *Client) Search(query string) ([]Video, error) {
	cmd := exec.Command("yt-dlp",
		"--dump-json",
		"--flat-playlist",
		"--no-warnings",
		"--ignore-errors",
		"--skip-download",
		fmt.Sprintf("ytsearch%d:%s", c.maxResults, query))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute yt-dlp: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var videos []Video

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var video Video
		if err := json.Unmarshal([]byte(line), &video); err == nil {
			if video.ID != "" {
				video.ThumbnailURL = video.GetThumbnailURL()
				videos = append(videos, video)
			}
		}
	}

	if len(videos) == 0 {
		return nil, fmt.Errorf("no videos found for query: %s", query)
	}

	return videos, nil
}

func (c *Client) Download(videoID, outputDir string) error {
	outputTemplate := fmt.Sprintf("%s/%%(title)s.%%(ext)s", outputDir)

	cmd := exec.Command("yt-dlp",
		"-f", "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
		"--merge-output-format", "mp4",
		"-o", outputTemplate,
		videoID)

	_, err := cmd.CombinedOutput()
	return err
}

func DownloadThumbnail(url, fallbackURL string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Try fallback
		resp, err = http.Get(fallbackURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("could not download thumbnail")
		}
	}

	return io.ReadAll(resp.Body)
}

func RenderThumbnail(imagePath string, width, height int) (string, error) {
	cmd := exec.Command("chafa",
		"--format", "kitty",
		"--size", fmt.Sprintf("%dx%d", width, height),
		imagePath)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("chafa failed: %w", err)
	}

	return string(output), nil
}
