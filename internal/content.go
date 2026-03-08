package internal

import (
	"bytes"
	"fmt"
	"html"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"github.com/gosimple/slug"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// Frontmatter holds the YAML metadata from a post
type Frontmatter struct {
	Title       string      `yaml:"title"`
	Date        interface{} `yaml:"date"`
	Tags        []string    `yaml:"tags"`
	Draft       bool        `yaml:"draft"`
	Bsky        string      `yaml:"bsky"`
	Description string      `yaml:"description"`
	Layout      string      `yaml:"layout"`
	Hidden      bool        `yaml:"hidden"`
	NoComments  bool        `yaml:"noComments"`
}

// TOCItem represents a heading in the table of contents
type TOCItem struct {
	Level    int
	ID       string
	Text     string
	Children []*TOCItem
}

// Post holds all data for a single page/post
type Post struct {
	Title       string
	Date        time.Time
	DateStr     string
	Tags        []string
	Draft       bool
	Bsky        string
	Description string
	Layout      string
	Hidden      bool
	NoComments  bool
	Content     string // Raw markdown
	HTML        string // Rendered HTML
	Slug        string
	Permalink   string
	RelURL      string
	ReadingTime int
	WordCount   int
	TOC         []*TOCItem
	Summary     string
	Filename    string // original filename
	OGImage     bool   // whether an OG image was generated
}

var dateFormats = []string{
	"2006-01-02",
	"02 Jan 2006",
	"January 2, 2006",
	"2 January 2006",
	"15 Feb 2006",
	"2 Jan 2006",
}

func parseDate(raw interface{}) time.Time {
	if raw == nil {
		return time.Time{}
	}

	switch v := raw.(type) {
	case time.Time:
		return v
	case string:
		for _, format := range dateFormats {
			t, err := time.Parse(format, v)
			if err == nil {
				return t
			}
		}
		// Try day-month-year with full month name
		t, err := time.Parse("2 January 2006", v)
		if err == nil {
			return t
		}
		fmt.Fprintf(os.Stderr, "warning: could not parse date %q\n", v)
		return time.Time{}
	}
	return time.Time{}
}

// ParseContent parses a markdown file with YAML frontmatter
func ParseContent(path string) (*Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	content := string(data)
	post := &Post{}
	post.Filename = filepath.Base(path)

	// Extract YAML frontmatter
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content[3:], "---", 2)
		if len(parts) == 2 {
			var fm Frontmatter
			if err := yaml.Unmarshal([]byte(parts[0]), &fm); err != nil {
				return nil, fmt.Errorf("parsing frontmatter in %s: %w", path, err)
			}
			post.Title = fm.Title
			post.Tags = fm.Tags
			post.Draft = fm.Draft
			post.Bsky = fm.Bsky
			post.Description = fm.Description
			post.Layout = fm.Layout
			post.Hidden = fm.Hidden
			post.NoComments = fm.NoComments
			post.Date = parseDate(fm.Date)
			if !post.Date.IsZero() {
				post.DateStr = post.Date.Format("2006-01-02")
			}
			post.Content = strings.TrimSpace(parts[1])
		}
	} else {
		post.Content = content
	}

	// Generate slug from filename (not title, to preserve URL compatibility)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)
	post.Slug = slug.Make(nameWithoutExt)

	// Calculate word count and reading time
	words := strings.Fields(post.Content)
	post.WordCount = len(words)
	post.ReadingTime = int(math.Max(1, math.Ceil(float64(post.WordCount)/200.0)))

	return post, nil
}

// ProcessShortcodes processes Hugo shortcodes in content
func ProcessShortcodes(content string) string {
	// Process {{% sidenote "position" "id" "label" %}}...{{% /sidenote %}} 
	sidenoteRe := regexp.MustCompile(`\{\{%\s*sidenote\s+"(left|right)"\s+"([^"]+)"\s+"([^"]+)"\s*%\}\}([\s\S]*?)\{\{%\s*/sidenote\s*%\}\}`)
	content = sidenoteRe.ReplaceAllStringFunc(content, func(match string) string {
		groups := sidenoteRe.FindStringSubmatch(match)
		if len(groups) < 5 {
			return match
		}
		position := groups[1]
		id := groups[2]
		label := groups[3]
		inner := strings.TrimSpace(groups[4])
		// Render inner markdown
		innerHTML := renderInlineMarkdown(inner)
		return fmt.Sprintf(`<span class="sidenote"><label class="sidenote-label" for="%s">%s</label><input class="sidenote-checkbox" type="checkbox" id="%s"><span class="sidenote-content sidenote-%s">%s</span></span>`,
			id, label, id, position, innerHTML)
	})

	// Process {{< figure src="..." caption="..." align="..." >}}
	figureRe := regexp.MustCompile(`\{\{<\s*figure([^>]*?)>\}\}`)
	content = figureRe.ReplaceAllStringFunc(content, func(match string) string {
		groups := figureRe.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}
		attrs := groups[1]

		src := extractAttr(attrs, "src")
		caption := extractAttr(attrs, "caption")
		align := extractAttr(attrs, "align")
		linkHref := extractAttr(attrs, "link")
		width := extractAttr(attrs, "width")
		height := extractAttr(attrs, "height")
		alt := extractAttr(attrs, "alt")

		if alt == "" && caption != "" {
			alt = caption
		}

		var sb strings.Builder
		sb.WriteString("<figure")
		if align == "center" {
			sb.WriteString(` class="align-center"`)
		}
		sb.WriteString(">")

		if linkHref != "" {
			sb.WriteString(fmt.Sprintf(`<a href="%s">`, linkHref))
		}

		imgTag := fmt.Sprintf(`<img loading="lazy" src="%s"`, src)
		if alt != "" {
			imgTag += fmt.Sprintf(` alt="%s"`, alt)
		}
		if width != "" {
			imgTag += fmt.Sprintf(` width="%s"`, width)
		}
		if height != "" {
			imgTag += fmt.Sprintf(` height="%s"`, height)
		}
		imgTag += ">"
		sb.WriteString(imgTag)

		if linkHref != "" {
			sb.WriteString("</a>")
		}

		if caption != "" {
			sb.WriteString(fmt.Sprintf("<figcaption><p>%s</p></figcaption>", caption))
		}
		sb.WriteString("</figure>")
		return sb.String()
	})

	return content
}

// extractAttr extracts an attribute value from a shortcode attribute string
func extractAttr(attrs, name string) string {
	re := regexp.MustCompile(fmt.Sprintf(`%s="([^"]*)"`, regexp.QuoteMeta(name)))
	m := re.FindStringSubmatch(attrs)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

// renderInlineMarkdown renders a small snippet of markdown to inline HTML.
// Strips wrapping <p> tags so the result is safe to embed inside a <span>.
func renderInlineMarkdown(md string) string {
	var buf bytes.Buffer
	if err := goldmark.New(
		goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()),
	).Convert([]byte(md), &buf); err != nil {
		return md
	}
	result := strings.TrimSpace(buf.String())
	// Strip single wrapping <p>...</p> so inline content works inside <span>
	if strings.HasPrefix(result, "<p>") && strings.HasSuffix(result, "</p>") {
		// Only strip if there's exactly one <p> wrapping the whole thing
		inner := result[3 : len(result)-4]
		if !strings.Contains(inner, "<p>") {
			return strings.TrimSpace(inner)
		}
	}
	return result
}

// RenderMarkdown converts markdown to HTML
func RenderMarkdown(content string) (string, error) {
	// Pre-process shortcodes
	content = ProcessShortcodes(content)

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(
					// chroma options
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return "", fmt.Errorf("rendering markdown: %w", err)
	}
	
	html := buf.String()
	html = addLazyLoadingToImages(html)
	return html, nil
}

func addLazyLoadingToImages(html string) string {
	re := regexp.MustCompile(`<img\s+src="([^"]*)"`)
	return re.ReplaceAllString(html, `<img loading="lazy" src="$1"`)
}

// MediaCDN is the CDN URL for media files (set during build)
var MediaCDN = ""

// RewriteMediaURLs rewrites /media/ URLs to use the CDN
func RewriteMediaURLs(htmlContent string) string {
	if MediaCDN == "" {
		return htmlContent
	}
	// Match img src="/media/..." and rewrite to CDN URL
	// Matches: src="/media/", src="/media/foo.png", etc.
	re := regexp.MustCompile(`(src|href)=["'](\/media\/[^"']*)["']`)
	return re.ReplaceAllStringFunc(htmlContent, func(match string) string {
		result := re.ReplaceAllString(match, "$1=\""+MediaCDN+"$2\"")
		return result
	})
}

// GenerateTOC generates a table of contents from HTML
func GenerateTOC(htmlContent string) []*TOCItem {
	headingRe := regexp.MustCompile(`<h([1-6])[^>]*id="([^"]*)"[^>]*>([\s\S]*?)</h[1-6]>`)
	matches := headingRe.FindAllStringSubmatch(htmlContent, -1)

	if len(matches) == 0 {
		return nil
	}

	// Strip HTML tags from heading text
	tagRe := regexp.MustCompile(`<[^>]+>`)

	var items []*TOCItem
	for _, m := range matches {
		level := 0
		fmt.Sscanf(m[1], "%d", &level)
		id := m[2]
		text := tagRe.ReplaceAllString(m[3], "")
		text = strings.TrimSpace(text)
		items = append(items, &TOCItem{
			Level: level,
			ID:    id,
			Text:  text,
		})
	}
	return items
}

// GenerateSummary generates a plain text summary from HTML content
func GenerateSummary(htmlContent string, maxLen int) string {
	tagRe := regexp.MustCompile(`<[^>]+>`)
	plain := tagRe.ReplaceAllString(htmlContent, "")
	// Unescape HTML entities (&ldquo; -> ", &amp; -> &, etc.)
	plain = html.UnescapeString(plain)
	plain = strings.TrimSpace(plain)
	// Collapse whitespace
	spaceRe := regexp.MustCompile(`\s+`)
	plain = spaceRe.ReplaceAllString(plain, " ")
	if len(plain) > maxLen {
		return plain[:maxLen] + "..."
	}
	return plain
}

// LoadPosts loads all posts from a directory
func LoadPosts(dir string) ([]*Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading dir %s: %w", dir, err)
	}

	var posts []*Post
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if entry.Name() == "_index.md" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		post, err := ParseContent(path)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}

		if post.Draft {
			continue
		}

		// Render markdown
		htmlContent, err := RenderMarkdown(post.Content)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", entry.Name(), err)
		}
		// Rewrite /media/ URLs to CDN
		htmlContent = RewriteMediaURLs(htmlContent)
		post.HTML = htmlContent
		post.Summary = GenerateSummary(htmlContent, 200)
		post.TOC = GenerateTOC(htmlContent)

		// Set URLs - posts go under /posts/
		post.RelURL = "/posts/" + post.Slug + "/"
		post.Permalink = "https://akshatsharma.xyz" + post.RelURL

		posts = append(posts, post)
	}

	// Sort by date descending
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

// GetAllTags returns all unique tags across posts
func GetAllTags(posts []*Post) map[string][]*Post {
	tags := make(map[string][]*Post)
	for _, p := range posts {
		for _, t := range p.Tags {
			tags[t] = append(tags[t], p)
		}
	}
	return tags
}


