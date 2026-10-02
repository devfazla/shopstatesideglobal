"""cleanhtml.core — HTML to Markdown cleaner.

Behaviour-identical port of the Go implementation (internal/markdown/markdown.go)
and the JavaScript implementation (js/src/clean-html.js). All three ports are
validated against the same canonical vectors in ``testdata/clean-html-cases.json``.
"""

from __future__ import annotations

import html
import re

I = re.IGNORECASE
S = re.DOTALL

_RE_COMMENT = re.compile(r"<!--.*?-->", I | S)
_RE_SCRIPT = re.compile(r"<script\b[^>]*>.*?</script\s*>", I | S)
_RE_STYLE = re.compile(r"<style\b[^>]*>.*?</style\s*>", I | S)
_RE_ALT = re.compile(r"\balt\s*=\s*[\"']([^\"']*)[\"']", I)

_RE_HEADING = re.compile(r"<h([1-6])\b[^>]*>(.*?)</h[1-6]\s*>", I | S)
_RE_IMAGE = re.compile(r"<img\b[^>]*?src\s*=\s*[\"']([^\"']*)[\"'][^>]*?>", I | S)
_RE_LINK = re.compile(r"<a\b[^>]*?href\s*=\s*[\"']([^\"']*)[\"'][^>]*>(.*?)</a\s*>", I | S)

_RE_BOLD = re.compile(r"<(?:strong|b)\b[^>]*>(.*?)</(?:strong|b)\s*>", I | S)
_RE_ITALIC = re.compile(r"<(?:em|i)\b[^>]*>(.*?)</(?:em|i)\s*>", I | S)
_RE_CODE = re.compile(r"<code\b[^>]*>(.*?)</code\s*>", I | S)
_RE_PRE = re.compile(r"<pre\b[^>]*>(.*?)</pre\s*>", I | S)

_RE_QUOTE = re.compile(r"<blockquote\b[^>]*>(.*?)</blockquote\s*>", I | S)
_RE_LIST = re.compile(r"<(ul|ol)\b[^>]*>(.*?)</(?:ul|ol)\s*>", I | S)
_RE_ITEM = re.compile(r"<li\b[^>]*>(.*?)</li\s*>", I | S)

_RE_BR = re.compile(r"<br\b[^>]*/?>", I | S)
_RE_HR = re.compile(r"<hr\b[^>]*/?>", I | S)

_RE_BLOCK = re.compile(
    r"</?(?:p|div|section|article|header|footer|main|aside|nav|figure|figcaption|"
    r"table|thead|tbody|tr|td|th|dl|dt|dd|form|fieldset)\b[^>]*>",
    I | S,
)
_RE_ANYTAG = re.compile(r"</?[a-zA-Z][^>]*>", I)

_RE_MULTI_BLANK = re.compile(r"\n{3,}")
_RE_TRAILING_SP = re.compile(r"[ \t]+$", re.MULTILINE)
_RE_MULTI_SPACES = re.compile(r"[ \t]{2,}")
_RE_PRE_TOKEN = re.compile(r"\x00PRE(\d+)\x00")

# Fenced code blocks extracted before generic tag stripping.
_pre_store: list = []


def _code_text(s: str) -> str:
    s = _RE_BR.sub("\n", s)
    s = _RE_ANYTAG.sub("", s)
    return html.unescape(s)


def _convert_pre(s: str) -> str:
    def repl(m: "re.Match") -> str:
        code = re.sub(r"\n+$", "", _code_text(m.group(1)))
        fence = "```\n" + code + "\n```"
        _pre_store.append(fence)
        return "\n\x00PRE%d\x00\n" % (len(_pre_store) - 1)

    return _RE_PRE.sub(repl, s)


def _restore_pre(s: str) -> str:
    def repl(m: "re.Match") -> str:
        i = int(m.group(1))
        return _pre_store[i] if 0 <= i < len(_pre_store) else ""

    return _RE_PRE_TOKEN.sub(repl, s)


def _inline_text(s: str) -> str:
    s = _RE_BOLD.sub(r"\1", s)
    s = _RE_ITALIC.sub(r"\1", s)
    s = _RE_CODE.sub(r"\1", s)
    s = _RE_LINK.sub(r"\2", s)
    s = _RE_IMAGE.sub("", s)
    s = _RE_BR.sub(" ", s)
    s = _RE_ANYTAG.sub("", s)
    s = html.unescape(s)
    return _RE_MULTI_SPACES.sub(" ", s).strip()


def _plain_text(s: str) -> str:
    s = _RE_ANYTAG.sub("", s)
    s = html.unescape(s)
    return _RE_MULTI_SPACES.sub(" ", s).strip()


def _convert_headings(s: str) -> str:
    def repl(m: "re.Match") -> str:
        level = int(m.group(1))
        if level < 1 or level > 6:
            level = 1
        text = _inline_text(m.group(2))
        return ("\n" + "#" * level + " " + text + "\n") if text else ""

    return _RE_HEADING.sub(repl, s)


def _convert_images(s: str) -> str:
    def repl(m: "re.Match") -> str:
        src = m.group(1).strip()
        a = _RE_ALT.search(m.group(0))
        alt = a.group(1).strip() if a else ""
        return ("![" + alt + "](" + src + ")") if src else ""

    return _RE_IMAGE.sub(repl, s)


def _convert_links(s: str) -> str:
    def repl(m: "re.Match") -> str:
        href = m.group(1).strip()
        text = _inline_text(m.group(2))
        if not href:
            return text
        if not text:
            return href
        if text == href:
            return href
        return "[" + text + "](" + href + ")"

    return _RE_LINK.sub(repl, s)


def _convert_lists(s: str) -> str:
    def repl(m: "re.Match") -> str:
        ordered = m.group(1).lower() == "ol"
        state = {"idx": 0}

        def item_repl(im: "re.Match") -> str:
            text = _inline_text(im.group(1))
            if not text:
                return ""
            if ordered:
                state["idx"] += 1
                return "\n%d. %s" % (state["idx"], text)
            return "\n- " + text

        out = _RE_ITEM.sub(item_repl, m.group(2))
        if not out:
            return ""
        return "\n" + out.lstrip("\n") + "\n"

    return _RE_LIST.sub(repl, s)


def _convert_quotes(s: str) -> str:
    def repl(m: "re.Match") -> str:
        text = _inline_text(m.group(1))
        if not text:
            return ""
        return "\n" + "\n".join("> " + line for line in text.split("\n")) + "\n"

    return _RE_QUOTE.sub(repl, s)


def _convert_inline(s: str) -> str:
    def bold_repl(m: "re.Match") -> str:
        text = _plain_text(m.group(1))
        return "**" + text + "**" if text else ""

    def italic_repl(m: "re.Match") -> str:
        text = _plain_text(m.group(1))
        return "_" + text + "_" if text else ""

    def code_repl(m: "re.Match") -> str:
        text = _plain_text(m.group(1))
        return "`" + text + "`" if text else ""

    s = _RE_BOLD.sub(bold_repl, s)
    s = _RE_ITALIC.sub(italic_repl, s)
    s = _RE_CODE.sub(code_repl, s)
    return s


def _normalise(s: str) -> str:
    s = _restore_pre(s)
    s = s.replace("\u00a0", " ")
    s = _RE_TRAILING_SP.sub("", s)
    s = _RE_MULTI_BLANK.sub("\n\n", s)
    return s.strip() + "\n"


def clean_html(s: str) -> str:
    """Convert an HTML string into clean Markdown (with a trailing newline).

    Unknown tags are dropped, whitespace is normalised, and HTML entities are
    decoded. Empty or whitespace-only input returns an empty string.
    """
    if not isinstance(s, str) or s.strip() == "":
        return ""

    global _pre_store
    _pre_store = []

    s = s.replace("\r\n", "\n").replace("\r", "\n")
    s = _RE_COMMENT.sub("", s)
    s = _RE_SCRIPT.sub("", s)
    s = _RE_STYLE.sub("", s)

    s = _convert_pre(s)
    s = _convert_headings(s)
    s = _convert_images(s)
    s = _convert_links(s)
    s = _convert_lists(s)
    s = _convert_quotes(s)
    s = _convert_inline(s)

    s = _RE_BR.sub("\n", s)
    s = _RE_HR.sub("\n---\n", s)
    s = _RE_BLOCK.sub("\n\n", s)

    s = _RE_ANYTAG.sub("", s)
    s = html.unescape(s)

    return _normalise(s)
