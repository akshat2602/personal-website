package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/akshat2602/personal-website/internal"
)

func main() {
	// Determine the project root (two levels up from cmd/build/)
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	_ = exe

	// Use working directory as project root
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Load SVG icon map
	svgTomlPath := filepath.Join(projectRoot, "ssg-static", "svg.toml")
	svgMap, err := internal.LoadSVGMap(svgTomlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load svg.toml: %v\n", err)
		svgMap = map[string]string{}
	}

	// Social icons in order (matching config.yml)
	socialNames := []string{
		"github", "linkedin", "hackernews", "goodreads", "letterboxd",
		"Cal.com", "peerlist", "stackoverflow", "email", "rss",
	}
	socialURLs := []string{
		"https://github.com/akshat2602",
		"https://linkedin.com/in/akshat-sharma-2602/",
		"https://news.ycombinator.com/user?id=akshat2602",
		"https://www.goodreads.com/user/show/54753277-akshat",
		"https://letterboxd.com/akshat2602/",
		"https://cal.com/akshat2602",
		"https://peerlist.io/akshat2602",
		"https://stackoverflow.com/users/16431252/akshat-sharma",
		"mailto:mail@akshatsharma.xyz",
		"/index.xml",
	}

	socialIcons := internal.BuildSocialIcons(socialNames, socialURLs, svgMap)

	config := &internal.SiteConfig{
		BaseURL:     "https://akshatsharma.xyz",
		Title:       "Akshat Sharma",
		Description: "hi, my name is Akshat Sharma",
		Author:      "Akshat",
		ContentDir:  filepath.Join(projectRoot, "content", "posts"),
		OutputDir:   filepath.Join(projectRoot, "public"),
		TemplateDir: filepath.Join(projectRoot, "ssg-templates"),
		StaticDir:   filepath.Join(projectRoot, "ssg-static"),
		ProjectRoot: projectRoot,
		SocialIcons: socialIcons,
	}

	gen, err := internal.NewGenerator(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating generator: %v\n", err)
		os.Exit(1)
	}

	if err := gen.Build(); err != nil {
		fmt.Fprintf(os.Stderr, "build error: %v\n", err)
		os.Exit(1)
	}
}
