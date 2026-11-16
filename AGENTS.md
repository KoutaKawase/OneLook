# Project Specification: Dictionary Search

## 1. Project Overview

Build a simple local web application that allows users to search for an English word across multiple online dictionaries and language-learning websites.

The application should display a list of links for the entered word. Clicking a link opens the corresponding website in a new browser tab.

The application is intended for personal use.

### Core principles

- Keep the implementation as simple as possible.
- Use only the Go standard library.
- Avoid unnecessary abstractions and dependencies.
- Make adding, removing, and modifying search links easy.
- Prefer straightforward, maintainable code.

## 2. Technology Stack

- Language: Go
- Backend: Go standard library only
- HTTP server: `net/http`
- HTML rendering: `html/template`
- JSON parsing: `encoding/json`
- Frontend: HTML, CSS, Vanilla JavaScript
- Database: None
- External dependencies: None

Do not use frontend frameworks, CSS frameworks, JavaScript libraries, ORMs, or third-party Go packages.

## 3. Directory Structure

Use the following structure:

```text
dictionary-search/
├── main.go
├── links.go
├── links.json
├── templates/
│   └── index.html
└── static/
    └── style.css
```

Responsibilities:

- `main.go`: HTTP server initialization and routing.
- `links.go`: JSON configuration loading, data structures, and URL generation.
- `links.json`: Search link configuration.
- `templates/index.html`: HTML template.
- `static/style.css`: Styling and dark theme.

Avoid creating additional files unless there is a clear reason.

## 4. Search Link Configuration

The uploaded `links.json` is the initial source of truth for the search links.

Load this file at application startup using Go's `encoding/json`.

Do not hardcode individual dictionary URLs in Go source code.

### JSON structure

```json
{
  "links": [
    {
      "name": "Oxford Learner's Dictionaries",
      "category": "Dictionaries",
      "url_template": "https://example.com/dictionary/{query}",
      "query_encoding": "path"
    }
  ]
}
```

### Field definitions

| Field | Description |
|---|---|
| name | Display name |
| category | Category used to group links |
| url_template | URL template containing `{query}` |
| query_encoding | URL encoding method |

Supported encoding methods:

- `path`: Use `url.PathEscape`.
- `query`: Use `url.QueryEscape`.

The `{query}` placeholder must be replaced with the appropriately encoded search term.

### Extensibility requirements

Adding a new dictionary or service must only require adding an entry to `links.json`.

No changes to the Go application code or HTML template should be necessary.

The application should support arbitrary categories, not just dictionaries.

Examples:

- Dictionaries
- Collocation
- Pronunciation
- Images
- Thesaurus
- Other

Do not hardcode category names.

Preserve the order of links in the JSON configuration when displaying them.

Group links by category while preserving their original order within each category.

## 5. User Interface

### Main screen

The interface should contain:

1. Application title.
2. Text input field.
3. Search button.
4. Search results grouped by category.

### Search behavior

- The user enters an English word.
- Pressing the Search button submits the word.
- The application displays the corresponding search links.
- The entered word remains in the input field after searching.
- Clicking a link opens the website in a new tab.

Example:

Input:

```text
hello
```

Expected result:

```text
Dictionaries

Oxford Learner's Dictionaries
Cambridge Dictionary
Collins English Dictionary
Merriam Webster
...

Collocation

EnglishCollocation.com
SKELL
Ozdic
...

Pronunciation

YouGlish
PlayPhrase
YARN
...

Images

Google Images
```

This is an illustrative example. The actual links must come from `links.json`.

### Design requirements

- Dark mode by default.
- Minimal and clean interface.
- Center the main content.
- Place the input and Search button on the same row.
- Make category headings visually distinguishable.
- Display links in a simple vertical list.
- Use subtle hover effects.
- Avoid unnecessary animations and visual elements.
- Make the layout usable on desktop and mobile browsers.

Do not implement a light/dark mode switch. Dark mode is sufficient.

## 6. HTTP Specification

Run the HTTP server on:

```text
127.0.0.1:8080
```

Application URL:

```text
http://localhost:8080
```

Use a single endpoint:

| Method | Endpoint | Description |
|---|---|---|
| GET | `/` | Display the search page and results |

Example:

```text
http://localhost:8080/?q=hello
```

Use server-side HTML rendering with `html/template`.

Do not introduce a REST API or asynchronous frontend architecture.

## 7. URL Generation

Implement a common URL generation function.

For example:

```go
func buildURL(link SearchLink, word string) string
```

The function should:

1. Encode the search term according to `query_encoding`.
2. Replace `{query}` in `url_template`.
3. Return the resulting URL.

Do not implement individual URL generation functions for each dictionary.

Validate configuration entries at startup and report malformed entries clearly.

Do not silently modify URLs supplied in `links.json`.

## 8. Error Handling and Input Validation

- Trim leading and trailing whitespace from search input.
- Do not perform a search when the input is empty.
- Handle missing or invalid JSON configuration.
- Handle missing `links.json` with a clear startup error.
- Prevent malformed configuration entries from producing invalid links.
- Use `html/template` for safe HTML rendering.
- Do not allow user input to become an arbitrary redirect destination.
- Apply reasonable input length limits.

## 9. External Links

All search result links must:

- Open in a new browser tab.
- Use `target="_blank"`.
- Include `rel="noopener noreferrer"`.

## 10. Out of Scope

Do not implement:

- Web scraping.
- Direct dictionary definition retrieval.
- Translation functionality.
- Search history.
- Favorites.
- User accounts.
- Authentication.
- Database storage.
- External APIs.
- Background services or scheduled jobs.
- Docker.
- Configuration management interfaces.
- Complex frontend interactions.

## 11. Implementation Guidelines

Prioritize simplicity over abstraction.

- Use Go standard library functionality whenever possible.
- Keep functions small and focused.
- Avoid unnecessary interfaces.
- Avoid excessive package separation.
- Avoid overengineering.
- Keep the application easy to understand and modify.

The JSON configuration should be the only place where individual search services are registered.

## 12. Acceptance Criteria

The implementation is complete when:

1. The application starts on `localhost:8080`.
2. The main page displays a search input and button.
3. Entering a word and pressing Search displays the corresponding links.
4. Links are grouped by category.
5. All links are generated from `links.json`.
6. Clicking a link opens the corresponding website in a new tab.
7. Adding a new service to `links.json` requires no Go or HTML code changes.
8. The interface uses a dark theme.
9. The application works without third-party dependencies.
10. The application remains simple and maintainable.

Before finishing, verify the application with several test inputs, including words containing spaces and punctuation.
