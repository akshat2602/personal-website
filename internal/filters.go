package internal

import (
	"fmt"
	"strings"

	"github.com/flosch/pongo2/v6"
	"github.com/gosimple/slug"
)

func init() {
	// slugify filter: converts a string to a URL-safe slug
	pongo2.RegisterFilter("slugify", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		tag := in.String()
		return pongo2.AsValue(makeTagSlug(tag)), nil
	})

	// markdownify filter: renders markdown inline
	pongo2.RegisterFilter("markdownify", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		rendered := renderInlineMarkdown(in.String())
		return pongo2.AsSafeValue(rendered), nil
	})

	// slugpost filter: slug for a post filename
	pongo2.RegisterFilter("slugpost", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		return pongo2.AsValue(slug.Make(in.String())), nil
	})

	// joinstr filter with separator
	pongo2.RegisterFilter("joinstr", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		sep := param.String()
		if in.CanSlice() {
			parts := make([]string, 0, in.Len())
			for i := 0; i < in.Len(); i++ {
				parts = append(parts, in.Index(i).String())
			}
			return pongo2.AsValue(strings.Join(parts, sep)), nil
		}
		return in, nil
	})

	// twitter share URL
	pongo2.RegisterFilter("twittershare", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		title := in.String()
		url := param.String()
		shareURL := fmt.Sprintf("https://twitter.com/intent/tweet?text=%s&url=%s",
			strings.ReplaceAll(title, " ", "%20"),
			strings.ReplaceAll(url, "/", "%2F"))
		return pongo2.AsValue(shareURL), nil
	})
}
