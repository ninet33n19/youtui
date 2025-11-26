package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ninet33n19/youtui/internal/cache"
	"github.com/ninet33n19/youtui/internal/config"
	"github.com/ninet33n19/youtui/internal/styles"
	"github.com/ninet33n19/youtui/internal/youtube"
)

type State int

const (
	StateSearch = iota
	StateList
	StateDetail
	StateDownloading
)

type Model struct {
	state       State
	videos      []youtube.Video
	cursor      int
	searchQuery string

	textInput textinput.Model
	spinner   spinner.Model
	keys      KeyMap

	width  int
	height int

	loading     bool
	err         error
	downloadMsg string
	renderedImg string

	config   *config.Config
	cache    *cache.ThumbnailCache
	ytClient *youtube.Client
}

func NewModel() *Model {
	cfg := config.Default()
	cfg.EnsureDirs()

	ti := textinput.New()
	ti.Placeholder = "search videos..."
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 50
	ti.PromptStyle = styles.InputPrompt
	ti.TextStyle = styles.InputText

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Spinner

	return &Model{
		state:     StateSearch,
		videos:    []youtube.Video{},
		cursor:    0,
		textInput: ti,
		spinner:   s,
		keys:      DefaultKeyMap(),
		config:    cfg,
		cache:     cache.New(cfg.CacheDir),
		ytClient:  youtube.NewClient(cfg.MaxResults),
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) SelectedVideo() *youtube.Video {
	if m.cursor >= 0 && m.cursor < len(m.videos) {
		return &m.videos[m.cursor]
	}

	return nil
}

func (m *Model) CardWidth() int {
	width := m.width - 6
	if width > 80 {
		width = 80
	}

	return width
}

func (m *Model) VisibleRange() (start, end int) {
	cardHeight := 4
	perPage := (m.height - 12) / cardHeight
	if perPage < 1 {
		perPage = 1
	}

	start = 0
	end = len(m.videos) - 1

	if len(m.videos) > perPage {
		if m.cursor < perPage/2 {
			start = 0
			end = perPage - 1
		} else if m.cursor >= len(m.videos)-perPage/2 {
			start = len(m.videos) - perPage
			end = len(m.videos) - 1
		} else {
			start = m.cursor - perPage/2
			end = m.cursor + perPage/2
		}
	}
	return start, end
}
