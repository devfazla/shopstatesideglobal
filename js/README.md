# @statesideglobal/clean-html

Convert an HTML string into clean Markdown. Behaviour-identical to the Go and
Python `clean-html` ports in this repository. All three languages are validated
against the same shared vectors in [`../testdata/clean-html-cases.json`](../testdata/clean-html-cases.json).

## Install

```bash
npm install @statesideglobal/clean-html
```

## Usage

```js
const { cleanHTML } = require("@statesideglobal/clean-html");

const md = cleanHTML("<h1>Title</h1><p>Hello <b>world</b></p>");
// "# Title\n\nHello **world**\n"
```

### CLI

```bash
node index.js --clean-html "<p>Hello <b>world</b></p>"
node index.js --clean-html-file page.html
cat page.html | node index.js
```

## API

- `cleanHTML(html: string): string` — returns Markdown with a trailing newline.
  Empty/whitespace-only input returns `""`.

## Tests

```bash
npm test
```

## License

MIT
