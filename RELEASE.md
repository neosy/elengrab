# Releases

## v0.15.0 — 2026-03-15

### Fix
- Preload warning for font

---

## v0.15.0 — 2026-03-14

### Feature
- Add in-line playback. Closes #224
- Add tap tooltip support for format info on mobile devices

### Style
- Improve row display

### Performance
- Replace row polling with SSE updates. Closes #227

### Fix
- Postprocessor yt-dlp

### Refactor
- Simplify transcoding parameter selection

---

## v0.14.4 — 2026-03-11

### Fix
- Improve yt-dlp format selection logic

---

## v0.14.3 — 2026-03-07

### Fix
- Improve yt-dlp format selection logic

---

## v0.14.2 — 2026-02-21

### Fix
- Add terminated flag to prevent use after Stop

### Chore
- Make it easier to build an image

---

## v0.14.1 — 2026-02-20

### Fix
- Avoid panic on multiple quit channel closes

---

## v0.14.0 — 2026-02-20

### Feature
- Add dynamic worker pool

### Refactor
- Simplify the code and add comments

---

## v0.13.4 — 2026-02-18

### Fix
- Add graceful shutdown delay before force killing

---

## v0.13.3 — 2026-02-17

### Fix
- Improve process kill handling

---

## v0.13.2 — 2026-02-17

### Refactor
- Added more detailed logging

### Perfomance
- Speed up fetching site title

---

## v0.13.1 — 2026-02-17

### Fix
- Correct width and height types for formats

### Refactor
- Added more detailed logging
- Improved yt-dlp service structure

### Docs
- README.md

---

## v0.13.0 — 2026-02-15

### Feature
- Add an icon when compiling an .exe file. Closes #192
- Add authorization parameters for downloading from youtube. Closes #185

### Refactor
-  Improve FetchLogos and FetchBestLogo API

---

## v0.12.0 — 2026-02-12

### Feature
- Add saving the download date. Closes #183

### Style
- Change the placeholder for the input field. Closes #188
- Change the title in the index. Closes #193
- Add tooltips for Quality and Format. Closes #186
- Add tooltips to media logos. Closes #195

### Refactor
- Moved media tables to a separate media.db file. Closes #197

### Fix
- Update delete endpoint path

---

## v0.11.2 — 2026-02-11

### Fix
- Reduce the length of the file name

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
