# akshatsharma.xyz

Personal website for [akshatsharma.xyz](https://akshatsharma.xyz), built with a custom Go static site generator.

## Stack

- **Build**: Go (`cmd/build/main.go`) with [Pongo2](https://github.com/flosch/pongo2) (Jinja2) templates
- **Markdown**: [Goldmark](https://github.com/yuin/goldmark) with syntax highlighting via Chroma (Dracula theme)
- **Hosting**: Cloudflare Pages
- **Media**: Cloudflare R2 (`/media/`)

## Structure

```
content/posts/     # Markdown posts with YAML frontmatter
ssg-templates/     # Jinja2 (Pongo2) HTML templates
ssg-static/        # Static assets (CSS, JS, images, svg.toml)
internal/          # Go SSG library (content parsing, site generation)
cmd/build/         # Build entry point
functions/         # Cloudflare Pages Functions (R2 media proxy)
```

## Build

```bash
make build    # build site to public/
make serve    # build + serve on localhost:8080
make test     # build + run Playwright tests
make clean    # remove public/
```

## Deployment

```bash
make build
npx wrangler pages deploy public --project-name=personal-website
```
