"""Tests for the Python clean-html port.

Validates against the same canonical vectors used by the Go and JavaScript
implementations, so a passing run means all three ports agree on behaviour.
Run with:  python -m pytest  (or)  python tests/test_clean_html.py
"""

import json
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from cleanhtml import clean_html  # noqa: E402

VECTORS_PATH = os.path.join(
    os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))),
    "testdata",
    "clean-html-cases.json",
)


def _trim_last_newline(s: str) -> str:
    return s[:-1] if s.endswith("\n") else s


def _load_cases():
    with open(VECTORS_PATH, "r", encoding="utf-8") as fh:
        return json.load(fh)["cases"]


def _check():
    failures = []
    for case in _load_cases():
        got = _trim_last_newline(clean_html(case["in"]))
        if got != case["want"]:
            failures.append((case["name"], case["in"], got, case["want"]))
    return failures


def test_clean_html_matches_shared_vectors():
    failures = _check()
    assert not failures, "mismatched cases: " + ", ".join(f[0] for f in failures)


if __name__ == "__main__":
    fails = _check()
    for case in _load_cases():
        got = _trim_last_newline(clean_html(case["in"]))
        status = "ok  " if got == case["want"] else "FAIL"
        print(f"{status} - {case['name']}")
        if got != case["want"]:
            print(f"     input: {case['in']!r}")
            print(f"     got:   {got!r}")
            print(f"     want:  {case['want']!r}")
    total = len(_load_cases())
    if fails:
        print(f"\n{len(fails)} case(s) failed")
        sys.exit(1)
    print(f"\nAll {total} cases passed")
