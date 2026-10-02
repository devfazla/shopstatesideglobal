// Package markdown provides helpers for turning raw HTML into clean,
// readable Markdown. It is intentionally dependency-free (standard library
// only) so it can be reused across the growing codebase without pulling in
// external modules.
package markdown

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

// Regular expressions for the block and inline elements we recognise.
// Non-greedy matches keep conversion scoped to the nearest closing tag.
var (
	reComment = regexp.MustCompile(`(?s)<!--.*?-->`)
	reScript  = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script\s*>`)
	reStyle   = regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style\s*>`)

	reAlt = regexp.MustCompile(`(?is)\balt\s*=\s*["']([^"']*)["']`)

	reHeading = regexp.MustCompile(`(?is)<h([1-6])\b[^>]*>(.*?)</h[1-6]\s*>`)
	reImage   = regexp.MustCompile(`(?is)<img\b[^>]*?src\s*=\s*["']([^"']*)["'][^>]*?>`)
	reLink    = regexp.MustCompile(`(?is)<a\b[^>]*?href\s*=\s*["']([^"']*)["'][^>]*>(.*?)</a\s*>`)

	reBold   = regexp.MustCompile(`(?is)<(strong|b)\b[^>]*>(.*?)</(?:strong|b)\s*>`)
	reItalic = regexp.MustCompile(`(?is)<(em|i)\b[^>]*>(.*?)</(?:em|i)\s*>`)
	reCode   = regexp.MustCompile(`(?is)<code\b[^>]*>(.*?)</code\s*>`)
	rePre    = regexp.MustCompile(`(?is)<pre\b[^>]*>(.*?)</pre\s*>`)

	reQuote = regexp.MustCompile(`(?is)<blockquote\b[^>]*>(.*?)</blockquote\s*>`)
	reList  = regexp.MustCompile(`(?is)<(ul|ol)\b[^>]*>(.*?)</(?:ul|ol)\s*>`)
	reItem  = regexp.MustCompile(`(?is)<li\b[^>]*>(.*?)</li\s*>`)

	reBreak = regexp.MustCompile(`(?is)<br\b[^>]*/?>`)
	reHRule = regexp.MustCompile(`(?is)<hr\b[^>]*/?>`)

	// Block-level tags become paragraph breaks.
	reBlock = regexp.MustCompile(`(?is)</?(p|div|section|article|header|footer|main|aside|nav|figure|figcaption|table|thead|tbody|tr|td|th|dl|dt|dd|form|fieldset)\b[^>]*>`)
	// Anything left over is stripped entirely.
	reAnyTag = regexp.MustCompile(`</?[a-zA-Z][^>]*>`)

	reMultiBlank  = regexp.MustCompile(`\n{3,}`)
	reTrailingSP  = regexp.MustCompile(`(?m)[ \t]+$`)
	reMultiSpaces = regexp.MustCompile(`[ \t]{2,}`)
)

// CleanHTML takes a string containing HTML and returns a cleaned Markdown
// representation of it. Unknown tags are dropped, whitespace is normalised,
// and HTML entities are decoded in the final output.
func CleanHTML(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}

	// Normalise line endings and remove non-content.
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	preStore = nil
	s = reComment.ReplaceAllString(s, "")
	s = reScript.ReplaceAllString(s, "")
	s = reStyle.ReplaceAllString(s, "")

	s = convertPre(s)
	s = convertHeadings(s)
	s = convertImages(s)
	s = convertLinks(s)
	s = convertLists(s)
	s = convertQuotes(s)
	s = convertInline(s)

	s = reBreak.ReplaceAllString(s, "\n")
	s = reHRule.ReplaceAllString(s, "\n---\n")
	s = reBlock.ReplaceAllString(s, "\n\n")

	// Drop any tags we did not explicitly handle, then decode entities so
	// that escaped markup (e.g. &lt;div&gt;) survives as literal text.
	s = reAnyTag.ReplaceAllString(s, "")
	s = html.UnescapeString(s)

	return normalise(s)
}

// convertHeadings rewrites <h1>..<h6> into "#"..###### headings.
func convertHeadings(s string) string {
	return reHeading.ReplaceAllStringFunc(s, func(m string) string {
		sub := reHeading.FindStringSubmatch(m)
		level, err := strconv.Atoi(sub[1])
		if err != nil || level < 1 || level > 6 {
			level = 1
		}
		text := inlineText(sub[2])
		if text == "" {
			return ""
		}
		return "\n" + strings.Repeat("#", level) + " " + text + "\n"
	})
}

// convertImages rewrites <img src alt> into ![alt](src).
func convertImages(s string) string {
	return reImage.ReplaceAllStringFunc(s, func(m string) string {
		sub := reImage.FindStringSubmatch(m)
		src := strings.TrimSpace(sub[1])
		alt := ""
		if a := reAlt.FindStringSubmatch(m); len(a) == 2 {
			alt = strings.TrimSpace(a[1])
		}
		if src == "" {
			return ""
		}
		return "![" + alt + "](" + src + ")"
	})
}

// convertLinks rewrites <a href>text</a> into [text](href). Links whose label
// equals the href are collapsed to a bare URL for readability.
func convertLinks(s string) string {
	return reLink.ReplaceAllStringFunc(s, func(m string) string {
		sub := reLink.FindStringSubmatch(m)
		href := strings.TrimSpace(sub[1])
		text := inlineText(sub[2])
		if href == "" {
			return text
		}
		if text == "" {
			return href
		}
		if text == href {
			return href
		}
		return "[" + text + "](" + href + ")"
	})
}

// convertLists rewrites <ul>/<ol>/<li> groups into Markdown list items.
func convertLists(s string) string {
	return reList.ReplaceAllStringFunc(s, func(m string) string {
		sub := reList.FindStringSubmatch(m)
		ordered := strings.EqualFold(sub[1], "ol")
		idx := 0
		out := reItem.ReplaceAllStringFunc(sub[2], func(item string) string {
			is := reItem.FindStringSubmatch(item)
			text := inlineText(is[1])
			if text == "" {
				return ""
			}
			if ordered {
				idx++
				return "\n" + strconv.Itoa(idx) + ". " + text
			}
			return "\n- " + text
		})
		if out == "" {
			return ""
		}
		return "\n" + strings.TrimLeft(out, "\n") + "\n"
	})
}

// convertQuotes prefixes blockquote contents with "> ".
func convertQuotes(s string) string {
	return reQuote.ReplaceAllStringFunc(s, func(m string) string {
		sub := reQuote.FindStringSubmatch(m)
		text := inlineText(sub[1])
		if text == "" {
			return ""
		}
		lines := strings.Split(text, "\n")
		for i, ln := range lines {
			lines[i] = "> " + ln
		}
		return "\n" + strings.Join(lines, "\n") + "\n"
	})
}

// convertInline handles bold, italic and inline code. Order matters: inner
// formatting is stripped before the surrounding markers are added.
func convertInline(s string) string {
	s = reBold.ReplaceAllStringFunc(s, func(m string) string {
		text := plainText(reBold.FindStringSubmatch(m)[2])
		if text == "" {
			return ""
		}
		return "**" + text + "**"
	})
	s = reItalic.ReplaceAllStringFunc(s, func(m string) string {
		text := plainText(reItalic.FindStringSubmatch(m)[2])
		if text == "" {
			return ""
		}
		return "_" + text + "_"
	})
	s = reCode.ReplaceAllStringFunc(s, func(m string) string {
		text := plainText(reCode.FindStringSubmatch(m)[1])
		if text == "" {
			return ""
		}
		return "`" + text + "`"
	})
	return s
}

// convertPre rewrites <pre> blocks into fenced code blocks. The inner content
// is protected from later tag stripping via a placeholder token that is
// restored at the end of CleanHTML.
var preStore []string

func convertPre(s string) string {
	return rePre.ReplaceAllStringFunc(s, func(m string) string {
		sub := rePre.FindStringSubmatch(m)
		code := strings.TrimRight(codeText(sub[1]), "\n")
		fence := "```\n" + code + "\n```"
		preStore = append(preStore, fence)
		return "\n\x00PRE" + strconv.Itoa(len(preStore)-1) + "\x00\n"
	})
}

// restorePre swaps placeholder tokens back for their fenced code blocks.
func restorePre(s string) string {
	return regexp.MustCompile("\x00PRE(\\d+)\x00").ReplaceAllStringFunc(s, func(tok string) string {
		num := regexp.MustCompile(`\d+`).FindString(tok)
		i, err := strconv.Atoi(num)
		if err != nil || i < 0 || i >= len(preStore) {
			return ""
		}
		return preStore[i]
	})
}

// inlineText strips tags and collapses whitespace for inline contexts, but is
// applied before generic tag removal so nested formatting survives.
func inlineText(s string) string {
	s = reBold.ReplaceAllString(s, "$2")
	s = reItalic.ReplaceAllString(s, "$2")
	s = reCode.ReplaceAllString(s, "$2")
	s = reLink.ReplaceAllString(s, "$2")
	s = reImage.ReplaceAllString(s, "")
	s = reBreak.ReplaceAllString(s, " ")
	s = reAnyTag.ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	return strings.TrimSpace(reMultiSpaces.ReplaceAllString(s, " "))
}

// plainText removes all tags and collapses runs of whitespace into one space.
func plainText(s string) string {
	s = reAnyTag.ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	return strings.TrimSpace(reMultiSpaces.ReplaceAllString(s, " "))
}

// codeText strips markup inside a <pre> block while preserving newlines, then
// decodes entities so escaped code samples render literally.
func codeText(s string) string {
	s = reBreak.ReplaceAllString(s, "\n")
	s = reAnyTag.ReplaceAllString(s, "")
	return html.UnescapeString(s)
}

// normalise tidies blank lines and trailing whitespace, and restores any
// placeholder code blocks before returning the final Markdown.
func normalise(s string) string {
	s = restorePre(s)
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = reTrailingSP.ReplaceAllString(s, "")
	s = reMultiBlank.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s) + "\n"
}
