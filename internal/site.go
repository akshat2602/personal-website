package internal

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/flosch/pongo2/v6"
)

// SocialIcon is a single social link entry
type SocialIcon struct {
	Name string
	URL  string
	SVG  string // raw SVG HTML
}

// SiteConfig holds site-wide configuration
type SiteConfig struct {
	BaseURL        string
	Title          string
	Description    string
	Author         string
	ContentDir     string
	OutputDir      string
	TemplateDir    string
	StaticDir      string
	ProjectRoot    string
	SocialIcons    []SocialIcon
}

// LoadSVGMap reads the PaperModX svg.toml and returns a name->svg map
func LoadSVGMap(svgTomlPath string) (map[string]string, error) {
	var raw map[string]string
	if _, err := toml.DecodeFile(svgTomlPath, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// BuildSocialIcons resolves SVG icons for each social entry
func BuildSocialIcons(names []string, urls []string, svgMap map[string]string) []SocialIcon {
	icons := make([]SocialIcon, 0, len(names))
	for i, name := range names {
		if i >= len(urls) {
			break
		}
		key := strings.ToLower(strings.TrimSpace(name))
		// Cal.com -> calcom
		key = strings.ReplaceAll(key, ".", "")
		key = strings.ReplaceAll(key, " ", "")
		svg := svgMap[key]
		if svg == "" {
			svg = svgMap["default"]
		}
		icons = append(icons, SocialIcon{
			Name: name,
			URL:  strings.TrimSpace(urls[i]),
			SVG:  strings.TrimSpace(svg),
		})
	}
	return icons
}

// PageData is the data passed to templates
type PageData struct {
	Site    *SiteConfig
	Post    *Post
	Posts   []*Post
	Tags    map[string][]*Post
	Tag     string
	TagSlug string
	// Pagination
	CurrentPage int
	TotalPages  int
	HasPrev     bool
	HasNext     bool
	PrevURL     string
	NextURL     string
	PagePosts   []*Post
}

// SearchItem is one entry in the search JSON index
type SearchItem struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Permalink string `json:"permalink"`
	Summary   string `json:"summary"`
}

// RSS structures
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel RSSChannel
}

type RSSChannel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	LastBuild   string    `xml:"lastBuildDate"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

const postsPerPage = 5

// Generator is the main site generator
type Generator struct {
	Config    *SiteConfig
	Templates *pongo2.TemplateSet
	Posts     []*Post
	Tags      map[string][]*Post
}

// NewGenerator creates a new site generator
func NewGenerator(config *SiteConfig) (*Generator, error) {
	loader, err := pongo2.NewLocalFileSystemLoader(config.TemplateDir)
	if err != nil {
		return nil, fmt.Errorf("creating template loader: %w", err)
	}
	ts := pongo2.NewSet("site", loader)

	return &Generator{
		Config:    config,
		Templates: ts,
	}, nil
}

// Build runs the full site build
func (g *Generator) Build() error {
	fmt.Println("Loading posts...")
	posts, err := LoadPosts(g.Config.ContentDir)
	if err != nil {
		return fmt.Errorf("loading posts: %w", err)
	}
	g.Posts = posts
	g.Tags = GetAllTags(posts)

	fmt.Printf("Found %d posts\n", len(posts))

	// Ensure output directory exists
	if err := os.MkdirAll(g.Config.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// Copy static files
	fmt.Println("Copying static files...")
	if err := g.copyStatic(); err != nil {
		return fmt.Errorf("copying static: %w", err)
	}

	// Generate pages
	fmt.Println("Generating home page...")
	if err := g.generateHome(); err != nil {
		return fmt.Errorf("generating home: %w", err)
	}

	fmt.Println("Generating OG images...")
	if err := g.generateOGImages(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: OG image generation: %v\n", err)
	}

	fmt.Println("Generating post pages...")
	if err := g.generatePosts(); err != nil {
		return fmt.Errorf("generating posts: %w", err)
	}

	fmt.Println("Generating archives page...")
	if err := g.generateArchives(); err != nil {
		return fmt.Errorf("generating archives: %w", err)
	}

	fmt.Println("Generating tag pages...")
	if err := g.generateTags(); err != nil {
		return fmt.Errorf("generating tags: %w", err)
	}

	fmt.Println("Generating search page...")
	if err := g.generateSearch(); err != nil {
		return fmt.Errorf("generating search: %w", err)
	}

	fmt.Println("Generating special pages (now, wakatime)...")
	if err := g.generateSpecialPages(); err != nil {
		return fmt.Errorf("generating special pages: %w", err)
	}

	fmt.Println("Generating RSS feed...")
	if err := g.generateRSS(); err != nil {
		return fmt.Errorf("generating RSS: %w", err)
	}

	fmt.Println("Generating search index...")
	if err := g.generateSearchIndex(); err != nil {
		return fmt.Errorf("generating search index: %w", err)
	}

	fmt.Println("Generating 404 page...")
	if err := g.generate404(); err != nil {
		return fmt.Errorf("generating 404: %w", err)
	}

	fmt.Printf("Build complete! Output in %s\n", g.Config.OutputDir)
	return nil
}

func (g *Generator) render(templateName string, data pongo2.Context, outPath string) error {
	t, err := g.Templates.FromFile(templateName)
	if err != nil {
		return fmt.Errorf("loading template %s: %w", templateName, err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("creating dir for %s: %w", outPath, err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", outPath, err)
	}
	defer f.Close()

	return t.ExecuteWriter(data, f)
}

func (g *Generator) baseContext() pongo2.Context {
	return pongo2.Context{
		"site":         g.Config,
		"year":         time.Now().Year(),
		"social_icons": g.Config.SocialIcons,
	}
}

func (g *Generator) generateHome() error {
	// Paginate posts
	visiblePosts := filterHidden(g.Posts)
	totalPages := (len(visiblePosts) + postsPerPage - 1) / postsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	for page := 1; page <= totalPages; page++ {
		start := (page - 1) * postsPerPage
		end := start + postsPerPage
		if end > len(visiblePosts) {
			end = len(visiblePosts)
		}
		pagePosts := visiblePosts[start:end]

		ctx := g.baseContext()
		ctx["posts"] = pagePosts
		ctx["all_posts"] = visiblePosts
		ctx["current_page"] = page
		ctx["total_pages"] = totalPages
		ctx["has_prev"] = page > 1
		ctx["has_next"] = page < totalPages
		ctx["page_title"] = g.Config.Title

		if page > 1 {
			ctx["prev_url"] = fmt.Sprintf("/page/%d/", page-1)
		}
		if page < totalPages {
			ctx["next_url"] = fmt.Sprintf("/page/%d/", page+1)
		}

		var outPath string
		if page == 1 {
			outPath = filepath.Join(g.Config.OutputDir, "index.html")
		} else {
			outPath = filepath.Join(g.Config.OutputDir, "page", fmt.Sprintf("%d", page), "index.html")
			ctx["prev_url"] = fmt.Sprintf("/page/%d/", page-1)
			if page == 2 {
				ctx["prev_url"] = "/"
			}
		}

		if err := g.render("home.html", ctx, outPath); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) generatePosts() error {
	for _, post := range g.Posts {
		ctx := g.baseContext()
		ctx["post"] = post
		ctx["posts"] = g.Posts
		ctx["page_title"] = post.Title + " | " + g.Config.Title
		// Compute description once for use in meta tags
		desc := post.Description
		if desc == "" {
			desc = post.Summary
		}
		ctx["post_description"] = desc

		// Find prev/next posts
		for i, p := range g.Posts {
			if p.Slug == post.Slug {
				if i > 0 {
					ctx["next_post"] = g.Posts[i-1] // newer
				}
				if i < len(g.Posts)-1 {
					ctx["prev_post"] = g.Posts[i+1] // older
				}
				break
			}
		}

		outPath := filepath.Join(g.Config.OutputDir, "posts", post.Slug, "index.html")
		if err := g.render("post.html", ctx, outPath); err != nil {
			return fmt.Errorf("rendering post %s: %w", post.Slug, err)
		}
	}
	return nil
}

func (g *Generator) generateArchives() error {
	// Group posts by year and month
	type MonthGroup struct {
		Month string
		Posts []*Post
	}
	type YearGroup struct {
		Year   string
		Count  int
		Months []MonthGroup
	}

	yearMap := make(map[string]map[string][]*Post)
	for _, p := range filterHidden(g.Posts) {
		year := p.Date.Format("2006")
		month := p.Date.Format("January")
		if yearMap[year] == nil {
			yearMap[year] = make(map[string][]*Post)
		}
		yearMap[year][month] = append(yearMap[year][month], p)
	}

	// Sort years descending
	var years []string
	for y := range yearMap {
		years = append(years, y)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(years)))

	var yearGroups []YearGroup
	for _, year := range years {
		months := yearMap[year]
		// Sort months
		monthOrder := []string{"January", "February", "March", "April", "May", "June",
			"July", "August", "September", "October", "November", "December"}
		var monthGroups []MonthGroup
		for i := len(monthOrder) - 1; i >= 0; i-- {
			m := monthOrder[i]
			if ps, ok := months[m]; ok {
				monthGroups = append(monthGroups, MonthGroup{Month: m, Posts: ps})
			}
		}
		total := 0
		for _, mg := range monthGroups {
			total += len(mg.Posts)
		}
		yearGroups = append(yearGroups, YearGroup{Year: year, Count: total, Months: monthGroups})
	}

	ctx := g.baseContext()
	ctx["year_groups"] = yearGroups
	ctx["page_title"] = "Posts | " + g.Config.Title

	return g.render("archives.html", ctx, filepath.Join(g.Config.OutputDir, "posts", "index.html"))
}

func (g *Generator) generateTags() error {
	// Tags index page
	type TagInfo struct {
		Name  string
		Slug  string
		Count int
	}
	var tagList []TagInfo
	for tag, posts := range g.Tags {
		tagList = append(tagList, TagInfo{
			Name:  tag,
			Slug:  makeTagSlug(tag),
			Count: len(posts),
		})
	}
	sort.Slice(tagList, func(i, j int) bool {
		return tagList[i].Name < tagList[j].Name
	})

	ctx := g.baseContext()
	ctx["tags"] = tagList
	ctx["page_title"] = "Tags | " + g.Config.Title

	if err := g.render("tags.html", ctx, filepath.Join(g.Config.OutputDir, "tags", "index.html")); err != nil {
		return err
	}

	// Individual tag pages
	for tag, posts := range g.Tags {
		tagSlug := makeTagSlug(tag)
		sortedPosts := make([]*Post, len(posts))
		copy(sortedPosts, posts)
		sort.Slice(sortedPosts, func(i, j int) bool {
			return sortedPosts[i].Date.After(sortedPosts[j].Date)
		})

		ctx := g.baseContext()
		ctx["tag"] = tag
		ctx["tag_slug"] = tagSlug
		ctx["posts"] = sortedPosts
		ctx["page_title"] = "Tag: " + tag + " | " + g.Config.Title

		outPath := filepath.Join(g.Config.OutputDir, "tags", tagSlug, "index.html")
		if err := g.render("tag.html", ctx, outPath); err != nil {
			return fmt.Errorf("rendering tag %s: %w", tag, err)
		}
	}
	return nil
}

func (g *Generator) generateSearch() error {
	ctx := g.baseContext()
	ctx["page_title"] = "Search | " + g.Config.Title
	return g.render("search.html", ctx, filepath.Join(g.Config.OutputDir, "search", "index.html"))
}

func (g *Generator) generateSpecialPages() error {
	contentRoot := filepath.Join(g.Config.ProjectRoot, "content")

	// Now page
	nowContent, err := ParseContent(filepath.Join(contentRoot, "now.md"))
	if err == nil {
		renderedHTML, err := RenderMarkdown(nowContent.Content)
		if err == nil {
			nowContent.HTML = renderedHTML
		}
		ctx := g.baseContext()
		ctx["post"] = nowContent
		ctx["page_title"] = "Now | " + g.Config.Title
		if err := g.render("now.html", ctx, filepath.Join(g.Config.OutputDir, "now", "index.html")); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "warning: could not parse now.md: %v\n", err)
	}

	// Wakatime page
	wakaContent, err := ParseContent(filepath.Join(contentRoot, "wakatime.md"))
	if err == nil {
		renderedHTML, err := RenderMarkdown(wakaContent.Content)
		if err == nil {
			wakaContent.HTML = renderedHTML
		}
		ctx := g.baseContext()
		ctx["post"] = wakaContent
		ctx["page_title"] = "Wakatime | " + g.Config.Title
		if err := g.render("now.html", ctx, filepath.Join(g.Config.OutputDir, "wakatime", "index.html")); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateOGImages() error {
	baseImage := filepath.Join(g.Config.StaticDir, "images", "og_base.png")
	for _, post := range g.Posts {
		if err := GenerateOGImage(post, baseImage, g.Config.OutputDir); err != nil {
			fmt.Fprintf(os.Stderr, "warning: OG image for %s: %v\n", post.Slug, err)
		} else {
			post.OGImage = true
		}
	}
	return nil
}

func (g *Generator) generateRSS() error {
	var items []RSSItem
	for i, p := range g.Posts {
		if i >= 20 {
			break
		}
		items = append(items, RSSItem{
			Title:       p.Title,
			Link:        g.Config.BaseURL + p.RelURL,
			Description: p.Summary,
			PubDate:     p.Date.Format(time.RFC1123Z),
			GUID:        g.Config.BaseURL + p.RelURL,
		})
	}

	rss := RSS{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       g.Config.Title,
			Link:        g.Config.BaseURL,
			Description: g.Config.Description,
			Language:    "en",
			LastBuild:   time.Now().Format(time.RFC1123Z),
			Items:       items,
		},
	}

	data, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return err
	}

	outPath := filepath.Join(g.Config.OutputDir, "index.xml")
	return os.WriteFile(outPath, append([]byte(xml.Header), data...), 0644)
}

func (g *Generator) generateSearchIndex() error {
	var items []SearchItem
	for _, p := range g.Posts {
		items = append(items, SearchItem{
			Title:     p.Title,
			Content:   p.Summary,
			Permalink: p.Permalink,
			Summary:   p.Summary,
		})
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(g.Config.OutputDir, "index.json"), data, 0644)
}

func (g *Generator) generate404() error {
	ctx := g.baseContext()
	ctx["page_title"] = "404 | " + g.Config.Title
	return g.render("404.html", ctx, filepath.Join(g.Config.OutputDir, "404.html"))
}

func (g *Generator) copyStatic() error {
	return copyDir(g.Config.StaticDir, g.Config.OutputDir)
}

// copyDir recursively copies src directory to dst
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}

func filterHidden(posts []*Post) []*Post {
	var out []*Post
	for _, p := range posts {
		if !p.Hidden {
			out = append(out, p)
		}
	}
	return out
}

func makeTagSlug(tag string) string {
	// Replace / with - for tags like "subject/distributed-system"
	tag = strings.ReplaceAll(tag, "/", "-")
	s := ""
	for _, c := range strings.ToLower(tag) {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' {
			s += string(c)
		} else if c == ' ' || c == '_' {
			s += "-"
		}
	}
	return strings.Trim(s, "-")
}
