# Git Branch Naming Convention

This document defines consistent naming rules for branches in the repository.

## 🔹 General Rules

* Use **kebab-case** (lowercase words separated by hyphens).
* Branch names should be **short, clear, and descriptive**.
* Prefix each branch with a keyword that reflects its purpose.
* Avoid using uppercase letters, spaces, or underscores.

**Example:**

```
feature/downloader-ui
fix/api-timeout
refactor/htmx-server
```

---

## 🔹 Branch Prefixes and Their Purpose

| Prefix            | Description                                                      |
| ----------------- | ---------------------------------------------------------------- |
| `feature/<name>`  | Adds new functionality (e.g., new API, module, or UI component). |
| `fix/<name>`      | Fixes a bug or error in existing code.                           |
| `refactor/<name>` | Refactors code without changing external behavior.               |
| `style/<name>`    | Updates styling (CSS, HTML layout, UI design).                   |
| `docs/<name>`     | Documentation updates (README, Swagger, etc.).                   |
| `test/<name>`     | Adds or modifies tests.                                          |
| `perf/<name>`     | Improves performance or efficiency.                              |
| `build/<name>`    | Changes related to build tools, Docker, Makefile, or CI/CD.      |
| `chore/<name>`    | Minor maintenance tasks (e.g., dependency updates).              |
| `hotfix/<name>`   | Urgent fix for a production issue.                               |
| `wip/<name>`      | Work in progress — temporary branch for incomplete work.         |

---

## 🔹 Examples

```
feature/user-authentication
fix/youtube-parser
refactor/rest-server
style/header-layout
docs/update-readme
test/order-repository
perf/cache-optimization
build/github-actions
chore/deps-update
hotfix/production-crash
wip/css-redesign
```

---

## 🔹 Recommended Workflow

1. Create a new branch from `main` or `develop`:

   ```bash
   git checkout -b feature/downloader-ui
   ```
2. Work on your changes.
3. Commit with clear, English messages:

   ```
   Refactor REST API code and improve CSS design
   ```
4. Open a pull request and merge after review.

---

## 🔹 Commit Rules with Issue Linking

Use clear, descriptive commit messages in English. To link a commit or pull request to an issue, include one of the GitHub keywords followed by the issue number:

```
Fixes #20
Closes #20
Resolves #20
```

### Example Commit Messages

```
perf: Improve yt-dlp response rate (#20)
```

(Does **not** link automatically, used only for clarity.)

```
perf: Improve yt-dlp response rate. Fixes #20
```

(Automatically links and closes Issue #20 when merged into main.)

### Recommended Workflow for Linking Issues

1. Create or locate the issue (e.g., Issue #20)
2. Work in a dedicated branch:

```
perf/20-increase-response-yt-dlp
```

3. Write normal commits without necessarily including the issue number.
4. When creating a Pull Request, include:

```
Fixes #20
```

GitHub will:

* Link the PR to the issue
* Automatically close the issue after merge
* Track the relationship in the UI

---

Following these conventions ensures consistent development workflow, improves collaboration, and provides clear traceability across commits, branches, and issues.
