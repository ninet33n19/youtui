package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ninet33n19/youtui/internal/util"
)

// Update handles all state updates
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case SearchResultMsg:
		return m.handleSearchResult(msg)

	case ImageRenderedMsg:
		return m.handleImageRendered(msg)

	case DownloadFinishedMsg:
		return m.handleDownloadFinished(msg)

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *Model) handleSearchResult(msg SearchResultMsg) (tea.Model, tea.Cmd) {
	m.loading = false
	if msg.Err != nil {
		m.err = msg.Err
		m.state = StateSearch
	} else {
		m.videos = msg.Videos
		m.state = StateList
		m.cursor = 0
	}
	return m, nil
}

func (m *Model) handleImageRendered(msg ImageRenderedMsg) (tea.Model, tea.Cmd) {
	m.loading = false
	if msg.Err != nil {
		m.renderedImg = "Error loading image: " + msg.Err.Error()
	} else {
		m.renderedImg = msg.Content
	}
	return m, nil
}

func (m *Model) handleDownloadFinished(msg DownloadFinishedMsg) (tea.Model, tea.Cmd) {
	m.loading = false
	m.state = StateList
	if msg.Err != nil {
		m.err = msg.Err
		m.downloadMsg = "Download failed!"
	} else {
		m.downloadMsg = "Download Complete! Saved to ./downloads"
	}
	return m, nil
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keybindings
	if key.Matches(msg, m.keys.ForceQuit) {
		return m, tea.Quit
	}

	if key.Matches(msg, m.keys.Quit) && m.state != StateSearch {
		return m, tea.Quit
	}

	// State-specific handling
	switch m.state {
	case StateSearch:
		return m.updateSearchState(msg)
	case StateList:
		return m.updateListState(msg)
	case StateDetail:
		return m.updateDetailState(msg)
	}

	return m, nil
}

func (m *Model) updateSearchState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) && m.textInput.Value() != "" {
		m.loading = true
		m.err = nil
		m.searchQuery = m.textInput.Value()
		return m, tea.Batch(
			SearchCmd(m.ytClient, m.textInput.Value()),
			m.spinner.Tick,
		)
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *Model) updateListState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.state = StateSearch
		m.textInput.SetValue("")
		m.textInput.Focus()
		m.searchQuery = ""
		return m, nil

	case key.Matches(msg, m.keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case key.Matches(msg, m.keys.Down):
		if m.cursor < len(m.videos)-1 {
			m.cursor++
		}

	case key.Matches(msg, m.keys.Enter):
		if len(m.videos) > 0 {
			m.state = StateDetail
			m.loading = true
			m.renderedImg = ""
			video := m.SelectedVideo()
			return m, tea.Batch(
				RenderThumbnailCmd(m.cache, video, m.width, m.height),
				m.spinner.Tick,
			)
		}
	}

	return m, nil
}

func (m *Model) updateDetailState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.renderedImg = util.ClearKittyGraphics()
		m.state = StateList
		return m, tea.ClearScreen

	case key.Matches(msg, m.keys.Enter):
		m.renderedImg = util.ClearKittyGraphics()
		m.state = StateDownloading
		m.loading = true
		m.downloadMsg = ""
		video := m.SelectedVideo()
		return m, tea.Batch(
			tea.ClearScreen,
			DownloadVideoCmd(m.ytClient, video.ID, m.config.DownloadDir),
			m.spinner.Tick,
		)
	}

	return m, nil
}
