# Beads Workflow Context

> **Context Recovery**: Run `bd prime` after compaction, clear, or new session
> Hooks auto-call this in Claude Code when .beads/ detected

# 🚨 TWO-PHASE WORKFLOW 🚨

**This project uses a split workflow:**

| Phase | Command | Actions |
|-------|---------|---------|
| **1. Implement** | `/start-issue <id>` | Set status → Plan → Implement |
| **2. Finalize** | `/finish-issue` | Commit → Close → Push |

## What this means for you:

### ✅ ALLOWED during implementation:
- Set ticket to `in_progress`
- Read, analyze, plan
- Write and modify code
- Run tests, build, verify

### ❌ FORBIDDEN during implementation:
- `git add` / `git commit` / `git push`
- `bd close`
- `bd sync` (syncs commits)

### When implementation is complete:

**DO NOT** run the old "Session Close Protocol". Instead say:

> "✅ Implementation complete. Files changed: [list files]. Run `/finish-issue` when ready to commit and close."

Then **STOP** and wait for the user.

---

## Core Rules
- Track strategic work in beads (multi-session, dependencies, discovered work)
- Use `bd create` for issues, TodoWrite for simple single-session execution
- When in doubt, prefer bd—persistence you don't need beats lost context
- Session management: check `bd ready` for available work

## Essential Commands

### Finding Work
- `bd ready` - Show issues ready to work (no blockers)
- `bd list --status=open` - All open issues
- `bd list --status=in_progress` - Your active work
- `bd show <id>` - Detailed issue view with dependencies

### Creating & Updating
- `bd create --title="..." --type=task|bug|feature --priority=2` - New issue
  - Priority: 0-4 or P0-P4 (0=critical, 2=medium, 4=backlog). NOT "high"/"medium"/"low"
- `bd update <id> --status=in_progress` - Claim work
- `bd update <id> --assignee=username` - Assign to someone
- `bd update <id> --title/--description/--notes/--design` - Update fields inline
- `bd close <id>` - Mark complete (⚠️ only in /finish-issue!)
- **WARNING**: Do NOT use `bd edit` - it opens $EDITOR which blocks agents

### Dependencies & Blocking
- `bd dep add <issue> <depends-on>` - Add dependency
- `bd blocked` - Show all blocked issues
- `bd show <id>` - See what's blocking/blocked by this issue

### Project Health
- `bd stats` - Project statistics
- `bd doctor` - Check for issues

## Workflow Summary

```
/start-issue <id>
    │
    ├── bd update <id> --status=in_progress
    ├── [Plan mode]
    ├── [Implement]
    └── STOP → "Implementation complete"
    
    ... user reviews ...

/finish-issue
    │
    ├── bd close <id>
    ├── git add -A
    ├── git commit -m "..."
    ├── bd sync
    └── Done!
```
