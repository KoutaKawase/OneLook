# OneLook

A tiny local web app that looks up an English word across multiple online
dictionaries and language resources. All links are grouped by category and
open in a new tab.

- Go standard library only (no third-party dependencies)
- Templates, CSS, and `links.json` are embedded into the binary via `//go:embed`
- Dark theme, minimal UI
- Runs on `http://localhost:8080`

## Requirements

- Go 1.22 or newer
- (Optional) `tmux` for running the app in the background

## Project layout

```text
onelook/
├── main.go            # HTTP server, embed declarations
├── links.go           # JSON loading, validation, URL building
├── links.json         # Search link configuration (embedded)
├── templates/
│   └── index.html     # HTML template (embedded)
├── static/
│   └── style.css      # Stylesheet (embedded)
├── go.mod
└── README.md
```

## Build

From the project directory:

```bash
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o onelook .
```

- `CGO_ENABLED=0` — pure-Go, static binary (no libc dependency).
- `-trimpath` — removes local build paths from the binary.
- `-ldflags="-s -w"` — strips symbol and debug information to shrink the binary.

The resulting `onelook` binary is self-contained: `links.json`, the HTML
template, and the CSS are all embedded. You can move or copy the binary
anywhere and run it without the source files.

> **Note:** If you edit `links.json`, `templates/index.html`, or
> `static/style.css`, you must rebuild the binary for the changes to take
> effect.

## Run

```bash
./onelook
```

Then open <http://localhost:8080> in your browser.

To stop the server, press `Ctrl+C`.

## Run in the background with tmux

Start a detached session:

```bash
tmux new -d -s onelook '/absolute/path/to/onelook'
```

Check that it is running:

```bash
tmux ls
curl http://localhost:8080
```

View the logs (attach):

```bash
tmux attach -t onelook
```

Detach again without stopping the app: press `Ctrl+b`, release, then `d`.

Stop the app:

```bash
tmux kill-session -t onelook
```

Restart after rebuilding:

```bash
tmux kill-session -t onelook
tmux new -d -s onelook '/absolute/path/to/onelook'
```

Write logs to a file instead of the pane:

```bash
tmux new -d -s onelook '/absolute/path/to/onelook >> /tmp/onelook.log 2>&1'
```

### Notes on tmux

- Sessions do **not** survive a reboot. Use `launchd` (macOS) or a similar
  mechanism if you need automatic startup.
- If the app crashes, the pane closes and the session disappears.
  Run `tmux ls` to check; restart if needed.

## Configuration

All search services are defined in `links.json`. To add, remove, or modify a
link, edit that file and rebuild the binary. No Go or HTML changes are
required.

```json
{
  "links": [
    {
      "name": "Oxford Learner's Dictionaries",
      "category": "Dictionaries",
      "url_template": "https://www.oxfordlearnersdictionaries.com/definition/american_english/{query}",
      "query_encoding": "path"
    }
  ]
}
```

| Field            | Description                                              |
| ---------------- | -------------------------------------------------------- |
| `name`           | Display name shown in the UI                             |
| `category`       | Group heading; arbitrary values are allowed              |
| `url_template`   | URL containing `{query}` exactly once                    |
| `query_encoding` | `path` (`url.PathEscape`) or `query` (`url.QueryEscape`) |

Invalid entries are reported as warnings at startup and skipped. The
remaining valid links still work.

## License

Personal use.
