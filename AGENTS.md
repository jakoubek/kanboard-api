# Agent Instructions

## Plans

- Make the plan extremely concise. Sacrifice grammar for the sake of concision.
- At the end of each plan, give me a list of unresolved questions to answer, if any.

## Main Workflow

The operator drives the work directly via prompts — there is no issue tracker.

### 1. Task

The operator assigns a piece of work in the conversation. Use plan mode to lay out
the approach, then implement the plan.

### 2. Operator tests manually

After the implementation is finished, the operator tests the solution themselves.

### 3. Commit

Only **after** the operator's manual test do you create the commit(s). Implementing
and committing are separate steps — do not commit while still implementing.

## Commits

This project follows the [Conventional Commits](https://www.conventionalcommits.org)
specification. Use the `/commit` command to create commits.

- **Types**: `feat:`, `fix:`, `docs:`, `style:`, `refactor:`, `perf:`, `test:`,
  `build:`, `ci:`, `chore:`, `revert:`. Optional scope, e.g. `feat(auth):`.
  Breaking changes get a `!`, e.g. `feat!:`.
- **English** commit messages.
- **One commit per concern.** Do not bundle unrelated changes. If a single commit
  message would need an "and" or a list of unrelated bullet points, split it —
  at file granularity (`git add <paths>`) or hunk granularity (`git add -p`).
- Only when all changes belong to one coherent topic is a single commit correct.
- Never add co-authors.

## CHANGELOG

Every user-facing change gets a CHANGELOG entry. Maintenance is done with the
`/update-changelog` command.

- Format: [Keep a Changelog](https://keepachangelog.com/) + Semantic Versioning.
- **English** entries (even when commit messages are German), one line per change,
  no commit hashes.
- Categorization:
  - `feat:` → **Added**
  - `fix:` → **Fixed**
  - `refactor:` / `perf:` → **Changed**
  - `BREAKING` → **Changed** (mark as BREAKING)
  - `remove` / `deprecate` → **Removed**
  - Skip `chore:`, `ci:`, `test:`, `style:` unless user-facing.

## Release / Versioning

Semantic Versioning with a `v` prefix (current: `v1.5.0`). Tags are created with
the `/tag-version` command; the changelog is updated with `/update-changelog`.

Version bump rule (Conventional Commits since the last tag):

- `BREAKING CHANGE` or a type ending with `!` (e.g. `feat!:`) → **MAJOR**
- `feat:` / `feat(` present → **MINOR**
- otherwise (`fix`, `docs`, `refactor`, `perf`, `style`, …) → **PATCH**
