# Repository Structure

Layout of the `shopstatesideglobal` repository. Documentation is Markdown (`.md`); the implementation is Go.

```
shopstatesideglobal/
├── README.md
├── structure.md
├── LICENSE
├── .gitignore
├── go.mod
├── cmd/
│   └── shopstatesideglobal/
│       └── main.go
├── internal/
│   └── markdown/
│       ├── markdown.go
│       └── markdown_test.go
└── docs/
    ├── about.md
    ├── contact.md
    └── links.md
```

## File Descriptions

| File | Purpose |
|------|---------|
| `README.md` | Repo overview, official links table, and index of the docs |
| `structure.md` | This file: the repo layout and what each file is for |
| `LICENSE` | MIT license for the repository contents |
| `.gitignore` | Ignores OS, editor, log, and secret files |
| `go.mod` | Go module definition (`github.com/devfazla/shopstatesideglobal`) |
| `cmd/shopstatesideglobal/main.go` | CLI entry point; exposes the `--clean-html` flag |
| `internal/markdown/markdown.go` | `CleanHTML` converts an HTML string into clean Markdown |
| `internal/markdown/markdown_test.go` | Unit tests for the HTML-to-Markdown cleaner |
| `docs/about.md` | What Stateside Global is, how it works, and why it exists |
| `docs/contact.md` | Support email, official channels, newsletter, and repo maintainer |
| `docs/links.md` | Website, site sections, and social profile links |

## Conventions

- Go code lives under `cmd/` (entry points) and `internal/` (packages); documentation stays under `docs/`.
- Add new documents inside `docs/` and list them in the table above and in `README.md`.
- Update this file whenever a file is added, renamed, or removed.