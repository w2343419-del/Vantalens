# Vantalens (Hugo Blog + TalentWriter)

[![Hugo](https://img.shields.io/badge/Hugo-Extended-blueviolet?style=flat-square)](https://gohugo.io/)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)

Vantalens is a bilingual Hugo blog project with a local management tool, TalentWriter (Go). The daily workflow uses one unified entry point: `web.exe`. It provides both the control view and the writer view. If you need to debug, you can still run the control and writer services separately.

## Quick Start

### 1) Preview the site

From the repository root:

```bash
hugo server
```

Open http://localhost:1313/VantalensWeb/.

### 2) Run the unified launcher

Enter the backend directory:

```bash
cd TalentWriter
```

Build and run the launcher:

```bash
go build -o web.exe ./cmd/server
./web.exe
```

Windows standalone example:

```powershell
$env:TALENTWRITER_APP_MODE="standalone"
$env:TALENTWRITER_AUTOSTART_HUGO="false"
./web.exe
```

`web.exe` includes both the control and writer views.

### 3) Optional debug services

If you need to debug a single service, run them separately:

```bash
go run ./cmd/control
go run ./cmd/writer
```

## Core Features

- Bilingual content management (ZH/EN)
- Local visual editing and publishing workflow
- Comment moderation, bulk actions, and export
- Admin-only analytics and visitor IP statistics

## Deployment

1. Build the static site with Hugo
2. Push to the GitHub repository
3. Serve it with GitHub Pages

## References

- Comment config: [config/_default/params.toml](config/_default/params.toml)
- Comment settings: [config/comment_settings.json](config/comment_settings.json)
- Analytics guide: [BUSUANZI_SETUP.md](BUSUANZI_SETUP.md)

## License

MIT License. See [LICENSE](LICENSE).
