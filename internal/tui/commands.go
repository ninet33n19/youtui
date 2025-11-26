package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ninet33n19/youtui/internal/cache"
	"github.com/ninet33n19/youtui/internal/youtube"
)

type SearchResultMsg struct {
	Videos []youtube.Video
	Err    error
}

type ImageRenderedMsg struct {
	Content string
	Err     error
}

type DownloadFinishedMsg struct {
	Filename string
	Err      error
}

func SearchCmd(client *youtube.Client, query string) tea.Cmd {
	return func() tea.Msg {
		videos, err := client.Search(query)

		return SearchResultMsg{
			Videos: videos,
			Err:    err,
		}
	}
}

func RenderThumbnailCmd(c *cache.ThumbnailCache, video *youtube.Video, width int, height int) tea.Cmd {
	return func() tea.Msg {
		displayHeight := max(height/2, 10)

		if rendered, ok := c.GetRendered(video.ID, width, displayHeight); ok {
			return ImageRenderedMsg{Content: rendered}
		}

		thumbnailPath := c.GetPath(video.ID)

		if !c.Exists(video.ID) {
			data, err := youtube.DownloadThumbnail(
				video.GetThumbnailURL(),
				video.GetFallbackThumbnailURL(),
			)
			if err != nil {
				return ImageRenderedMsg{Err: err}
			}

			if err := c.Save(video.ID, data); err != nil {
				return ImageRenderedMsg{Err: err}
			}
		}

		rendered, err := youtube.RenderThumbnail(thumbnailPath, width-10, displayHeight)
		if err != nil {
			return ImageRenderedMsg{Err: err}
		}

		c.SaveRendered(video.ID, width, displayHeight, rendered)

		return ImageRenderedMsg{Content: rendered}
	}
}

func DownloadVideoCmd(client *youtube.Client, videoID, outputDir string) tea.Cmd {
	return func() tea.Msg {
		err := client.Download(videoID, outputDir)

		return DownloadFinishedMsg{Filename: videoID, Err: err}
	}
}
