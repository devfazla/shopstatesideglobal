# Repository Structure

Layout of the `shopstatesideglobal` repository. Documentation is Markdown (`.md`); the `clean-html` function is implemented in Go, JavaScript, and Python with identical behaviour.

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
├── js/
│   ├── package.json
│   ├── index.js
│   ├── README.md
│   ├── src/clean-html.js
│   └── test/clean-html.test.js
├── python/
│   ├── pyproject.toml
│   ├── README.md
│   ├── cleanhtml/{__init__.py,core.py,__main__.py}
│   └── tests/test_clean_html.py
├── testdata/
│   └── clean-html-cases.json
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
| `.gitignore` | Ignores OS, editor, log, secret, and build files |
| `go.mod` | Go module definition (`github.com/devfazla/shopstatesideglobal`) |
| `cmd/shopstatesideglobal/main.go` | Go CLI entry point; exposes the `--clean-html` flag |
| `internal/markdown/markdown.go` | Go `CleanHTML` converts an HTML string into clean Markdown |
| `internal/markdown/markdown_test.go` | Go test; reads the shared vectors from `testdata/` |
| `js/` | npm package `@statesideglobal/clean-html` (same function, JS) |
| `python/` | PyPI package `cleanhtml-statesideglobal` (same function, Python) |
| `testdata/clean-html-cases.json` | Canonical cross-language test vectors (single source of truth) |
| `docs/about.md` | What Stateside Global is, how it works, and why it exists |
| `docs/contact.md` | Support email, official channels, newsletter, and repo maintainer |
| `docs/links.md` | Website, site sections, and social profile links |

## Conventions

- The `clean-html` function is ported identically across Go, JavaScript, and Python; all three must keep passing `testdata/clean-html-cases.json`. Add new behaviour as new vectors there first.
- Go code lives under `cmd/` (entry points) and `internal/` (packages); the npm package under `js/`; the Python package under `python/`; documentation under `docs/`.
- Add new documents inside `docs/` and list them in the table above and in `README.md`.
- Update this file whenever a file is added, renamed, or removed.