import os
import re
import markdown
import yaml
from pathlib import Path
from datetime import datetime

# Directories
CONTENT_DIR = "content/posts"
TEMPLATE_DIR = "templates"
STATIC_DIR = "static"
OUTPUT_DIR = "output"


def load_template(name):
    return Path(TEMPLATE_DIR, name).read_text(encoding="utf-8")


def save_html(path, content):
    output_path = Path(OUTPUT_DIR, path)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(content, encoding="utf-8")


def parse_front_matter(md_text):
    match = re.match(r"^---\\n(.*?)\\n---\\n(.*)", md_text, re.DOTALL)
    if match:
        front_matter_raw, content = match.groups()
        metadata = yaml.safe_load(front_matter_raw)
    else:
        metadata = {}
        content = md_text
    return metadata, content


def process_sidenotes(md_text):
    pattern = re.compile(
        r"{{< sidenote \\\"(.*?)\\\" \\\"(.*?)\\\" \\\"(.*?)\\\" >}}(.*?){{< /sidenote >}}",
        re.DOTALL,
    )

    def replacer(match):
        side, note_id, label, content = match.groups()
        return f"""<span class="sidenote">
<label class="sidenote-label" for="{note_id}">{label}</label>
<input class="sidenote-checkbox" type="checkbox" id="{note_id}"></input>
<span class="sidenote-content sidenote-{side}">{markdown.markdown(content.strip())}</span>
</span>"""

    return pattern.sub(replacer, md_text)


def build_post(md_path):
    md_text = Path(md_path).read_text(encoding="utf-8")
    metadata, content = parse_front_matter(md_text)

    content = process_sidenotes(content)

    # Convert to HTML
    html_content = markdown.markdown(content, extensions=["extra", "codehilite", "toc"])

    # Apply template
    template = load_template("base.html")

    title = metadata.get("title", Path(md_path).stem.title())
    date = metadata.get("date", "")
    tags = metadata.get("tags", [])
    tags_html = " ".join([f"<span class='tag'>{tag}</span>" for tag in tags])

    final_html = template
    final_html = final_html.replace("{{ title }}", title)
    final_html = final_html.replace("{{ date }}", date)
    final_html = final_html.replace("{{ tags }}", tags_html)
    final_html = final_html.replace("{{ content }}", html_content)
    final_html = final_html.replace("{{ year }}", str(datetime.now().year))

    # Optional: if bsky comments needed
    if "bsky" in metadata:
        bsky_uri = metadata["bsky"]
        comments_html = f'<div id="comments-section" data-bsky-uri="{bsky_uri}"></div>'
        final_html = final_html.replace("{{ comments }}", comments_html)
    else:
        final_html = final_html.replace("{{ comments }}", "")

    output_file = Path(md_path).with_suffix(".html").name
    save_html(output_file, final_html)

    return title, output_file


def build_index(posts):
    template = load_template("index.html")
    posts_html = "\\n".join(
        [f'<li><a href="{link}">{title}</a></li>' for title, link in posts]
    )
    final_html = template.replace("{{ posts }}", posts_html)
    final_html = final_html.replace("{{ title }}", "Akshat Sharma")
    final_html = final_html.replace("{{ year }}", str(datetime.now().year))
    save_html("index.html", final_html)


def copy_static():
    os.system(f"cp -r {STATIC_DIR} {OUTPUT_DIR}")


def main():
    Path(OUTPUT_DIR).mkdir(exist_ok=True)
    posts = []
    for md_file in Path(CONTENT_DIR).glob("*.md"):
        title, link = build_post(md_file)
        posts.append((title, link))
    build_index(posts)
    copy_static()


if __name__ == "__main__":
    main()
