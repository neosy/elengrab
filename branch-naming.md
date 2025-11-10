# Git Branch Naming Convention

This document defines consistent naming rules for branches in the repository.

## 🔹 General Rules

* Use **kebab-case** (lowercase words separated by hyphens).
* Branch names should be **short, clear, and descriptive**.
* Prefix each branch with a keyword that reflects its purpose.
* Avoid using uppercase letters, spaces, or underscores.

**Example:**

```
feature-downloader-ui
fix-api-timeout
refactor-htmx-server
```

---

## 🔹 Branch Prefixes and Their Purpose

| Prefix            | Description                                                      |
| ----------------- | ---------------------------------------------------------------- |
| `feature-<name>`  | Adds new functionality (e.g., new API, module, or UI component). |
| `fix-<name>`      | Fixes a bug or error in existing code.                           |
| `refactor-<name>` | Refactors code without changing external behavior.               |
| `style-<name>`    | Updates styling (CSS, HTML layout, UI design).                   |
| `docs-<name>`     | Documentation updates (README, Swagger, etc.).                   |
| `test-<name>`     | Adds or modifies tests.                                          |
| `perf-<name>`     | Improves performance or efficiency.                              |
| `build-<name>`    | Changes related to build tools, Docker, Makefile, or CI/CD.      |
| `chore-<name>`    | Minor maintenance tasks (e.g., dependency updates).              |
| `hotfix-<name>`   | Urgent fix for a production issue.                               |
| `wip-<name>`      | Work in progress — temporary branch for incomplete work.         |

---

## 🔹 Examples

```
feature-user-authentication
fix-youtube-parser
refactor-rest-server
style-header-layout
docs-update-readme
test-order-repository
perf-cache-optimization
build-github-actions
chore-deps-update
hotfix-production-crash
wip-css-redesign
```

---

## 🔹 Recommended Workflow

1. Create a new branch from `main` or `develop`:

   ```bash
   git checkout -b feature-downloader-ui
   ```
2. Work on your changes.
3. Commit with clear, English messages:

   ```
   Refactor REST API code and improve CSS design
   ```
4. Open a pull request and merge after review.

---

👌 Following this convention keeps your repository clean, organized, and easy to maintain.
