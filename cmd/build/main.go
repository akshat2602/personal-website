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
