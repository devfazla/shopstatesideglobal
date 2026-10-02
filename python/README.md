# cleanhtml-statesideglobal

Convert an HTML string into clean Markdown. Behaviour-identical to the Go and
JavaScript `clean-html` ports in this repository. All three languages are
validated against the same shared vectors in [`../testdata/clean-html-cases.json`](../testdata/clean-html-cases.json).

## Install

```bash
pip install cleanhtml-statesideglobal
```

(From this repo: `pip install ./python`)

## Usage

```python
from cleanhtml import clean_html

md = clean_html("<h1>Title</h1><p>Hello <b>world</b></p>")
# '# Title\n\nHello **world**\n'
```

### CLI

```bash
python -m cleanhtml --clean-html "<p>Hello <b>world</b></p>"
python -m cleanhtml --clean-html-file page.html
cat page.html | python -m cleanhtml
```

## Tests

```bash
python tests/test_clean_html.py
# or, if pytest is installed:
pytest
```

## License

MIT
