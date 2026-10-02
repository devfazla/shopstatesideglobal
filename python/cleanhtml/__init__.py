"""cleanhtml — Convert HTML strings into clean Markdown.

Behaviour-identical to the Go and JavaScript ``clean-html`` ports from the
Stateside Global project. See ``testdata/clean-html-cases.json`` for the shared
test vectors that all three languages validate against.
"""

from .core import clean_html

__all__ = ["clean_html"]
__version__ = "0.1.0"
