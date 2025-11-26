## YouTUI

YouTUI is a Bubble Tea-based terminal UI for searching YouTube, previewing video cards (with inline thumbnails rendered via Kitty graphics), and downloading the selected video in up to 1080p MP4 using `yt-dlp`. It ships with polished Lip Gloss styling, thumbnail caching, and a guided workflow that keeps you on the keyboard the entire time.

### Features
- **Fast search** – leverages `yt-dlp`’s search mode to pull concise result lists (default 10).
- **Rich list view** – scrollable cards with titles, channel names, duration, and status toasts for downloads.
- **Thumbnail previews** – renders cached thumbnails with `chafa` + Kitty graphics when you open a video’s detail view.
- **One-key downloads** – press enter in the detail view to grab the video (best ≤1080p) into `./downloads`.
- **Smart caching** – deduplicates thumbnails and stores rendered Kitty frames to keep the UI responsive.

### Requirements
- Go ≥ 1.25 (see `go.mod`).
- Terminal that understands Kitty graphics (Kitty, WezTerm, iTerm2 3.5 beta, etc.).
- `yt-dlp` on `PATH` (searching & downloading).
- `chafa` on `PATH` (thumbnail rendering).
- `ffmpeg` on `PATH` (needed by `yt-dlp` to mux audio/video).

### Getting Started
```bash
git clone https://github.com/ninet33n19/youtui.git
cd youtui
go run ./cmd/youtui   # or: go build -o youtui ./cmd/youtui
```
Running `go run` (or the compiled binary) drops you into the search screen. Type a query, hit enter, and use the keys below to drive the interface.

### Default Key Bindings
| Key(s)          | Context      | Action                    |
|-----------------|--------------|---------------------------|
| `enter`         | search       | Run query                 |
| `↑`/`k`, `↓`/`j`| list         | Move selection            |
| `enter`         | list         | Open detail + load art    |
| `enter`         | detail       | Download video            |
| `esc`           | list/detail  | Back to previous view     |
| `q`             | list/detail  | Quit gracefully           |
| `ctrl+c`        | anywhere     | Force quit                |

### Workflow Overview
1. **Search** – landing screen with focused text input. Enter submits and triggers a Bubble Tea command that invokes the YouTube client.
2. **List view** – scroll through cards; status bar shows download success/failure messages.
3. **Detail view** – fetches + renders the thumbnail using the cache in `os.TempDir()/youtui/thumbnails`. From here, pressing enter downloads the video (best ≤1080p) into `./downloads` inside your current working directory.
4. **Download state** – spinner stays up until `yt-dlp` finishes; confirmation surfaces back in the list view.

### Configuration Notes
- Defaults live in `internal/config/config.go`. By default:
  - Cache directory: `${TMPDIR}/youtui/thumbnails`
  - Download directory: `./downloads` (created automatically)
  - Max results: `10`
- Adjusting behavior currently means editing `config.Default()`; the project keeps the config struct isolated so you can later wire it to env vars or flags.

### Development
- The Bubble Tea program entrypoint is `cmd/youtui/main.go` → `internal/app`.
- Core UI state, commands, and views live under `internal/tui/`.
- YouTube integration (search, download, thumbnails) sits in `internal/youtube/`.
Use `go test ./...` (once tests are added) and `golangci-lint`/`go fmt` to keep contributions tidy.

### Troubleshooting
- **No thumbnails?** Ensure your terminal supports Kitty graphics and that `chafa` is installed; otherwise you’ll see a plain fallback message.
- **Downloads fail immediately?** Confirm `yt-dlp` and `ffmpeg` are available on `PATH`.
- **Stuck cache** – `internal/cache` exposes helpers; delete `${TMPDIR}/youtui` to force a refresh.

Happy hacking!

