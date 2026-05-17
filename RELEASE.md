# Releases

## v0.20.0 — 2026-05-25

### 🎨 Style
- Reduce the get button for small screens [L004]

### 🏗 Build / Migrations
- Add automatic branch id generator [L001]

### 🐛 Fix
- Handle empty argImageSource for SplitSeq

### 🧩 Refactor
- Rename static paths

---

## v0.19.1 — 2026-05-15

### 🐛 Fix
- Add the creation of necessary directories to dockerfile

---

## v0.19.0 — 2026-05-15

### ✨ Features
- Added thumbnail support for media files in lists (#291)
- Implemented download directory sharding for better storage scalability
- Added prometheus db metrics
- Added media descriptions (incl. watch page integration)
- Added PWA share target URL import handling
- Introduced configurable short link TTL via SHORT_LINK_TTL_DAYS

### 🎨 Style
- Added global UI blocking overlay for menu interactions
- Implemented backdrop (screen dimming) when blocking menus are open
- Added scroll lock behavior while blocking menu is active
- Added persistent notification area at the bottom of the screen
- Introduced unified notification styling for light and dark themes
- Added description to watch page media
- Replaced global user-select disabling with targeted touch interaction styles
- Improved mobile layout: avatar positioning, title alignment, and user menu placement
- Improved infinite scroll preload trigger behavior on mobile
- Maskable icon visual updates

### 🖥️ UI / Frontend
- Added persistent notification area at the bottom of the screen
- Implemented smooth animated notifications with auto-dismiss
- Added support for multiple notification types: success, error, warning, and info
- Introduced unified notification styling with proper contrast for both light and dark themes
- Added a description to the watch

### ⚡ Performance
- Added SVG rendering cache to reduce recomputation overhead
- Introduced in-memory cache for thumbnail file storage

### 🧩 Refactor
- Introduce streaming full names iterator in repository
- Reorganized UI folder into composition-based structure
- Renamed REST handlers for improved consistency
- Refactored HTML template data handling
- Refactored downloader migrations and internal data flow

### 🧪 Test
- Added admin server (http://127.0.0.1:6060/.../...)
  - /debug/pprof/
  - /metrics
  - /healthz
- Added fasthttp prometheus middleware

### 🏗 Build / Migrations
- Added macOS amd64 and arm64 release targets
- Introduced migrations:
  - media info backfill
  - media description backfill
  - thumbnail generation pipeline
- Downloader-related migration refactor and data normalization updates

### 🐛 Fix
- Fixed Prometheus metrics interval update behavior
- Fixed SSE reconnection logic in frontend
- Fixed metric path normalization issues
- Fixed PWA manifest enctype warning

---

## v0.18.6 — 2026-04-29

### Feature
- Added support for `<a>` elements in menu actions
- Added "Copy short link" menu action

### Refactor
- Introduced explicit item types for menu rendering (link / action / divider)

---

## v0.18.5 — 2026-04-20

### Style
- Remove the default text from the Share link

### Build
- Upgrade go to 1.26.2

---

## v0.18.4 — 2026-04-19

### Fix
- Corrected Open Graph meta tags for proper link preview generation in Telegram

---

## v0.18.3 — 2026-04-16

### Refactor
- Add exceptions

---

## v0.18.2 — 2026-04-15

### Feature
- Introduce EventKey-based client routing

### Refactor
- Create a basic exception

---

## v0.18.1 — 2026-04-15

### Fix
- Improve session cookie invalidation

---

## v0.18.0 — 2026-04-13

### Feature
- Add support for generating short URLs for downloaded media. Closes #257
- Add a media viewing page
- Add an Error Page
- Add hash when opening the playback window
- Add a menu to the row
- Implement a share link using standard OS functions

### Style
- Update play button icon

### Refactor
- Redesign main.go
- Simplify errorx wrapping logic, expand functions, add unit tests
- Optimize and extend the methods in the errorx package

### Fix
- Remove the spaces in from the authorization data
- Refine the definition of the file extension
- Exclude technical errors from public

---

## v0.17.1 — 2026-04-02

### Docs
- Add ip address to log

---

## v0.17.0 — 2026-04-02

### Feature
- Add authorization for admins and privileged rights. Closes #187

### Style
- Improve the design of active elements

### Refactor
- Restructure templates into layouts, pages, and components
- Restructure errorx package

---

## v0.16.2 — 2026-03-23

### Style
- Add link to version

### Chore
- Bump CSS/JS file versions

---

## v0.16.1 — 2026-03-23

### Feature
- Add demo mode
- Add error notification

### Refactor
- Work pools

---

## v0.16.0 — 2026-03-22

### Feature
- Add name search to the download history. Closes #184
- Add a status bar at the bottom. Closes #156
- Add Server status

### Style
- Improved display on mobile devices
- Add smooth tooltip animations

### Perfomance
- Update download progress via SSE

### Refactor
- CSS class names

---

## v0.15.1 — 2026-03-17

### Fix
- Preload warning for font

### Perfomance
- Optimize infinite scroll history loading

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
