# Personal CV

Custom Markdown-to-HTML static site generator (Go) producing a personal CV
website at [https://tgustinak.github.io](https://tgustinak.github.io),
hosted free on GitHub Pages.

Markdown pages live in `content/`, the HTML shell in `templates/default.html`,
styling in `public/cv.css`. The build renders each `.md` file to `.html`,
minifies the full page, and generates `sitemap.xml`, `robots.txt`, and `404.html`.

## Development

```shell
make build   # one-shot build into the repo root (this is what Pages serves)
make dev     # rebuild on file changes (content/ + templates/)
make test    # unit tests
make vet     # go vet
make clean   # clear Go test cache
```

Preview locally:

```shell
make dev &
python3 -m http.server 8000
# open http://localhost:8000/
```

CI (`.github/workflows/ci.yml`) runs format check, vet, tests, and a site
build on every push and pull request. Generated files (`index.html`,
`sitemap.xml`, …) are committed — the build output *is* the deployed site.
