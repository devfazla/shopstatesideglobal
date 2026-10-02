package markdown

import "testing"

func TestCleanHTML(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "empty",
			in:   "   ",
			want: "",
		},
		{
			name: "paragraph with bold",
			in:   "<p>Hello <b>world</b></p>",
			want: "Hello **world**",
		},
		{
			name: "heading",
			in:   "<h1>Title</h1>",
			want: "# Title",
		},
		{
			name: "link",
			in:   `<a href="https://example.com">Example</a>`,
			want: "[Example](https://example.com)",
		},
		{
			name: "bare link collapses",
			in:   `<a href="https://example.com">https://example.com</a>`,
			want: "https://example.com",
		},
		{
			name: "unordered list",
			in:   "<ul><li>a</li><li>b</li></ul>",
			want: "- a\n- b",
		},
		{
			name: "ordered list",
			in:   "<ol><li>first</li><li>second</li></ol>",
			want: "1. first\n2. second",
		},
		{
			name: "italic and code",
			in:   "<em>hi</em> and <code>x=1</code>",
			want: "_hi_ and `x=1`",
		},
		{
			name: "entity decoding",
			in:   "<p>Tom &amp; Jerry &lt;3</p>",
			want: "Tom & Jerry <3",
		},
		{
			name: "script removed",
			in:   "<script>alert(1)</script><p>ok</p>",
			want: "ok",
		},
		{
			name: "image",
			in:   `<img src="/a.png" alt="logo">`,
			want: "![logo](/a.png)",
		},
		{
			name: "block quote",
			in:   "<blockquote>quoted</blockquote>",
			want: "> quoted",
		},
		{
			name: "line break",
			in:   "one<br>two",
			want: "one\ntwo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanHTML(tt.in)
			// CleanHTML returns a trailing newline; compare trimmed output.
			if trimLast(got) != tt.want {
				t.Errorf("CleanHTML(%q)\n got: %q\nwant: %q", tt.in, got, tt.want+"\n")
			}
		})
	}
}

func trimLast(s string) string {
	for len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return s
}
