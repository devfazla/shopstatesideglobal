# Stateside Global (ShopStatesideGlobal)

Documentation, official reference links, and Go tooling for **Stateside Global**, a membership-based commerce platform for verified global goods.

The `docs/` directory holds the reference documentation. Go code lives under `cmd/` and `internal/` and will grow into the main application.

## Official Links

| Platform | Handle | Link |
|----------|--------|------|
| Website  | -      | https://shopstatesideglobal.com |
| About page | -    | https://www.shopstatesideglobal.com/about |
| Instagram | @shopstatesideglobal | https://www.instagram.com/shopstatesideglobal |
| TikTok   | @shopstatesideglobal | https://www.tiktok.com/@shopstatesideglobal |
| YouTube  | @shopstatesideglobal | https://www.youtube.com/@shopstatesideglobal |
| Facebook | shopstatesideglobal  | https://www.facebook.com/shopstatesideglobal |
| X (Twitter) | @shopstatesideglobal | https://x.com/shopstatesideglobal |
| LinkedIn | shopstatesideglobal | https://www.linkedin.com/company/shopstatesideglobal |
| Threads  | @shopstatesideglobal | https://www.threads.net/@shopstatesideglobal |
| Pinterest | shopstatesideglobal | https://www.pinterest.com/shopstatesideglobal |
| Telegram | @shopstatesideglobal | https://t.me/shopstatesideglobal |
| WhatsApp Channel | shopstatesideglobal | https://whatsapp.com/channel/0029Vb7t55QFcowBwTNSaH1k |
| Support  | -      | support@shopstatesideglobal.com |

See [docs/links.md](docs/links.md) for the full list.

## Documentation

| File | Purpose |
|------|---------|
| [docs/about.md](docs/about.md) | What Stateside Global is, how it works, and why it exists |
| [docs/links.md](docs/links.md) | Website, site sections, and social profile links |
| [docs/contact.md](docs/contact.md) | How to get in touch |
| [structure.md](structure.md) | Repo layout and file descriptions |

## Development

Requires [Go](https://go.dev/dl/) 1.23 or later. The code uses only the standard library (no external modules).

```bash
# Build the CLI
go build -o shopstatesideglobal ./cmd/shopstatesideglobal

# Run the tests
go test ./...

# Convert an HTML string to clean Markdown
./shopstatesideglobal --clean-html "<h1>Title</h1><p>Hello <b>world</b></p>"

# Or read from a file / stdin
./shopstatesideglobal --clean-html-file page.html
cat page.html | ./shopstatesideglobal --clean-html -
```

The cleaner lives in [`internal/markdown`](internal/markdown/markdown.go) as `markdown.CleanHTML(html string) string`, and the `--clean-html` flag is wired up in [`cmd/shopstatesideglobal`](cmd/shopstatesideglobal/main.go).

## Maintaining This Repo

- Use full `https://` URLs for all links in the docs.
- When a link or handle changes, update `README.md` and `docs/links.md` together.
- Follow standard Go layout (`cmd/` for binaries, `internal/` for packages) and run `go test ./...` before pushing.