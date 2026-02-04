# Releases

## v0.12.0 — 2026-02-10

### Feature
- 

### Style
- Change the placeholder for the input field. Closes #188

---

## v0.11.1 — 2026-02-04

### Hotfix
- Fix json decoding error in yt-dlp service. Closes #190

---

## v0.11.0 — 2026-02-04

### Feature
- Add “Paste from Clipboard” button to input field. Closes #149
- Update the audio and video info at the end of the download
- Delete duplicates per user
- Add site logos for uploaded videos. Closes #180

### Fix
- Correct video decoding and resolution
- Resolve ffmpeg audio codec handling error
- Auto-detect and convert audio-only media

### Style
- Change paste icon and add title tooltip

### Chore
- Update CSS cache version

---

## v0.10.1 — 2026-01-30

### Build
- Add an update to the assets folder after updating the version. Closes #167

---

## v0.10.0 — 2026-01-30

### Feature
- Add download progress. Closes #24
- Add a Repeat button for the Failed status. Closes #159
- Add a color spinner for different statuses

### Fix
- Add timeout for yt-dlp requests. Closes #152
- Fix release title
- Fix yt-dlp - HTTP Error 403: Forbidden

---

## v0.9.11 — 2026-01-24

### Hotfix
- Avoid creating fasthttp.FS per request

---

## v0.9.10 — 2026-01-22

### Hotfix
- Fixed maintenance interval behavior.

---

## v0.9.9 — 2026-01-21

### Fixes
- Fixed compression when sending files.
- Added timeouts and logging for fetching YouTube channel information.

### Refactor
- Improved channel avatar retrieval logic.
- Refactored logging.

---

## v0.9.8 — 2026-01-21

### Fixes
- Added retry logic when deleting downloaded files.

### Refactor
- Fixed informational log output during server startup.

---

## v0.9.7 — 2026-01-21

### Hotfix
- Fixed Dockerfile build issue.

### Fixes
- Fixed publishing of build artifacts for GitHub releases.
- Added preflight checks for external dependencies (e.g., ffmpeg).

### CI
- Added manual release workflow.
