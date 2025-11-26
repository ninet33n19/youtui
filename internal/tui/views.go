package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ninet33n19/youtui/internal/styles"
	"github.com/ninet33n19/youtui/internal/util"
	"github.com/ninet33n19/youtui/internal/youtube"
)

// View renders the current view
func (m *Model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	switch m.state {
	case StateSearch:
		return m.searchView()
	case StateList:
		return m.listView()
	case StateDetail:
		return m.detailView()
	case StateDownloading:
		return m.downloadView()
	default:
		return m.searchView()
	}
}

func (m *Model) searchView() string {
	renderedLogo := styles.LogoStyle.Render(styles.Logo)
	subtitle := styles.Subtle.Render("Search and download YouTube videos")
	input := styles.InputBox.Render(m.textInput.View())
	help := styles.Help.Render("enter to search • ctrl+c to quit")

	content := lipgloss.JoinVertical(lipgloss.Center,
		renderedLogo,
		"",
		subtitle,
		"",
		"",
		input,
		"",
	)

	if m.loading {
		content += fmt.Sprintf("\n\n%s Searching...", m.spinner.View())
	}

	if m.err != nil {
		content += "\n\n" + styles.Error.Render(fmt.Sprintf("Error: %v", m.err))
	}

	content += "\n\n" + help

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m *Model) listView() string {
	contentWidth := m.CardWidth()

	// Header
	header := lipgloss.JoinHorizontal(lipgloss.Left,
		styles.Title.Render(" YouTUI "),
		"  ",
		styles.Muted.Render("Results for: "),
		lipgloss.NewStyle().Foreground(styles.Malibu).Italic(true).Render(fmt.Sprintf(`"%s"`, m.searchQuery)),
	)

	resultInfo := styles.Muted.Render(fmt.Sprintf("Found %d videos", len(m.videos)))

	var s strings.Builder
	s.WriteString(header + "\n")
	s.WriteString(resultInfo + "\n\n")

	// Download message
	if m.downloadMsg != "" {
		if strings.Contains(m.downloadMsg, "failed") {
			s.WriteString(styles.Error.Render("  "+styles.CrossIcon+" "+m.downloadMsg) + "\n\n")
		} else {
			s.WriteString(styles.Success.Render("  "+styles.CheckIcon+" "+m.downloadMsg) + "\n\n")
		}
	}

	// Video cards
	start, end := m.VisibleRange()
	for i := start; i <= end; i++ {
		if i >= len(m.videos) {
			break
		}
		card := m.renderVideoCard(m.videos[i], i+1, m.cursor == i, contentWidth)
		s.WriteString(card + "\n")
	}

	// Scroll indicator
	if len(m.videos) > (end - start + 1) {
		scrollInfo := styles.Muted.Render(fmt.Sprintf("  %s %d/%d", styles.ArrowIcon, m.cursor+1, len(m.videos)))
		s.WriteString("\n" + scrollInfo)
	}

	// Help
	help := styles.Help.Render("  ↑/k up • ↓/j down • enter select • esc back • q quit")
	s.WriteString("\n" + help)

	return styles.App.Render(s.String())
}

func (m *Model) renderVideoCard(video youtube.Video, index int, isSelected bool, width int) string {
	cStyle := styles.Card.Width(width)
	tStyle := styles.VideoTitle
	iStyle := styles.Index

	if isSelected {
		cStyle = styles.CardSelected.Width(width)
		tStyle = styles.VideoTitleSelected
		iStyle = styles.IndexSelected
	}

	indexStr := iStyle.Render(fmt.Sprintf("%d.", index))

	titleMaxWidth := width - 10
	title := util.Truncate(video.Title, titleMaxWidth)
	titleRender := tStyle.Render(title)

	channelRender := styles.Channel.Render(styles.ChannelIcon + " " + util.Truncate(video.Channel, 30))
	durationRender := styles.Duration.Render(styles.DurationIcon + " " + video.FormatDuration())

	metaLine := lipgloss.JoinHorizontal(lipgloss.Left,
		channelRender,
		styles.Separator(),
		durationRender,
	)

	cardContent := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, indexStr, " ", titleRender),
		lipgloss.NewStyle().PaddingLeft(4).Render(metaLine),
	)

	return cStyle.Render(cardContent)
}

func (m *Model) detailView() string {
	video := m.SelectedVideo()
	if video == nil {
		return "No video selected"
	}

	headerStyle := styles.DetailHeader.Width(m.width - 8)

	header := styles.DetailTitle.Render(video.Title)
	channelInfo := styles.Channel.Render(styles.ChannelIcon + " " + video.Channel)
	durationInfo := styles.Duration.Render(styles.DurationIcon + " " + video.FormatDuration())

	meta := lipgloss.JoinHorizontal(lipgloss.Left,
		channelInfo,
		styles.Separator(),
		durationInfo,
	)

	headerContent := lipgloss.JoinVertical(lipgloss.Left, header, meta)
	headerCard := headerStyle.Render(headerContent)

	var imageArea string
	if m.loading {
		imageArea = fmt.Sprintf("\n%s Loading Thumbnail...\n", m.spinner.View())
	} else if m.renderedImg != "" {
		imageArea = m.renderedImg
	} else {
		imageArea = styles.Muted.Render("No thumbnail available")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		headerCard,
		"",
		imageArea,
		"",
		styles.Help.Render("  enter download (1080p MP4) • esc back"),
	)

	return styles.App.Render(content)
}

func (m *Model) downloadView() string {
	video := m.SelectedVideo()
	if video == nil {
		return "No video selected"
	}

	content := fmt.Sprintf(
		"Downloading:\n\n%s\n\n%s Please wait...",
		lipgloss.NewStyle().Bold(true).Foreground(styles.Dolly).Render(video.Title),
		m.spinner.View(),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, styles.DownloadBox.Render(content))
}
