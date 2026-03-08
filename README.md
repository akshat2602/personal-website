# akshatsharma.xyz

Personal website for [akshatsharma.xyz](https://akshatsharma.xyz), built with a custom Go static site generator.

## Stack

- **Build**: Go (`cmd/build/main.go`) with [Pongo2](https://github.com/flosch/pongo2) (Jinja2) templates
- **Markdown**: [Goldmark](https://github.com/yuin/goldmark) with syntax highlighting via Chroma (Dracula theme)
- **Hosting**: Cloudflare Pages
- **Media**: Cloudflare R2 with CDN (`cdn.akshatsharma.xyz`)

## Features

- Sidenotes (desktop: margin, mobile: toggle)
- Bluesky comments on posts
- Dynamic OG images per post
- Search (Fuse.js)
- Copy code button on code blocks
- RSS feed
- Syntax highlighting (Dracula theme)

## Structure

```
content/posts/     # Markdown posts with YAML frontmatter
templates/         # Jinja2 (Pongo2) HTML templates
static/            # Static assets (CSS, JS, images, svg.toml)
internal/          # Go SSG library (content parsing, site generation)
cmd/build/         # Build entry point
```

## Build

```bash
make build    # build site to public/
make serve    # build + serve on localhost:8080
make test     # build + run Playwright tests
make clean    # remove public/
```

## Deployment

Push to GitHub on `master` branch — Cloudflare Pages auto-deploys.

Build settings in Cloudflare Pages:
- Build command: `go run cmd/build/main.go`
- Build output: `public`