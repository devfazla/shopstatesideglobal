"""Command-line interface mirroring the Go ``--clean-html`` flag.

Usage:
    python -m cleanhtml --clean-html "<p>hi</p>"
    python -m cleanhtml --clean-html-file page.html
    cat page.html | python -m cleanhtml
"""

from __future__ import annotations

import argparse
import sys

from .core import clean_html


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        prog="cleanhtml",
        description="Convert an HTML string to clean Markdown.",
    )
    parser.add_argument(
        "--clean-html",
        dest="html",
        default="",
        help='HTML string to convert, or "-" to read from stdin.',
    )
    parser.add_argument(
        "--clean-html-file",
        dest="html_file",
        default="",
        help="Read HTML from a file and print clean Markdown.",
    )
    args = parser.parse_args(argv)

    if args.html_file:
        with open(args.html_file, "r", encoding="utf-8") as fh:
            data = fh.read()
    elif args.html == "-":
        data = sys.stdin.read()
    elif args.html.strip():
        data = args.html
    elif not sys.stdin.isatty():
        data = sys.stdin.read()
    else:
        parser.print_help(sys.stderr)
        return 1

    out = clean_html(data)
    if out:
        sys.stdout.write(out)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
