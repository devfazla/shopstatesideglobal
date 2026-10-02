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

The `clean-html` function is available in **three languages**, each with the same behaviour. All three are validated against one shared fixture: [`testdata/clean-html-cases.json`](testdata/clean-html-cases.json).

### Go

Requires [Go](https://go.dev/dl/) 1.23+. Standard library only.

```bash
go build -o shopstatesideglobal ./cmd/shopstatesideglobal
go test ./...
./shopstatesideglobal --clean-html "<h1>Title</h1><p>Hello <b>world</b></p>"
```

Function: `markdown.CleanHTML(html string) string` in [`internal/markdown`](internal/markdown/markdown.go).

### JavaScript / npm

Package: [`@statesideglobal/clean-html`](js/) (source in [`js/`](js/)).

```bash
cd js && npm install && npm test
node index.js --clean-html "<p>Hello <b>world</b></p>"
```

Function: `cleanHTML(html)` exported from [`js/index.js`](js/index.js).

### Python

Package: `cleanhtml-statesideglobal` (source in [`python/`](python/)).

```bash
cd python && pip install .
python tests/test_clean_html.py
python -m cleanhtml --clean-html "<p>Hello <b>world</b></p>"
```

Function: `clean_html(html)` from the `cleanhtml` package.

## Maintaining This Repo

- **Cross-language parity:** when you change the cleaner, update all three ports and add the case to `testdata/clean-html-cases.json` first, then run each language's tests.
- Use full `https://` URLs for all links in the docs.
- When a link or handle changes, update `README.md` and `docs/links.md` together.
- Follow standard Go layout (`cmd/` for binaries, `internal/` for packages).
