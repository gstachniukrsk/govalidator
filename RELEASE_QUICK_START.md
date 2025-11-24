# 🚀 Release Quick Start

## TL;DR - How to Release

```bash
# Just commit and push to master with conventional format
git commit -m "feat: your new feature"
git push origin master

# Release happens automatically in ~2-3 minutes! 🎉
```

## Commit Format Cheat Sheet

| Commit Type | Version Bump | Example |
|-------------|--------------|---------|
| `feat:` | **Minor** (0.x.0) | `feat: add JSONPresenter` |
| `fix:` | **Patch** (0.0.x) | `fix: correct Unicode handling` |
| `perf:` | **Patch** (0.0.x) | `perf: optimize validator` |
| `feat!:` | **Major** (x.0.0) | `feat!: rename validators` |
| `BREAKING CHANGE:` | **Major** (x.0.0) | (in commit body) |
| `docs:`, `style:`, `refactor:`, `test:`, `chore:` | **No release** | Just documentation |

## Quick Examples

### New Feature (Minor Version)
```bash
git commit -m "feat: add CombinedPresenter for error formatting"
# Result: v0.1.0 → v0.2.0
```

### Bug Fix (Patch Version)
```bash
git commit -m "fix: MinLengthValidator now counts runes correctly"
# Result: v0.1.0 → v0.1.1
```

### Breaking Change (Major Version)
```bash
git commit -m "feat!: rename IntValidator to IsIntegerValidator

BREAKING CHANGE: IntValidator renamed to IsIntegerValidator
for consistency with other validators."
# Result: v0.1.0 → v1.0.0
```

### Documentation (No Release)
```bash
git commit -m "docs: update README with examples"
# Result: No new release created
```

## What Happens Automatically

1. ✅ Tests run
2. ✅ Linters run
3. ✅ Version calculated from commits
4. ✅ CHANGELOG.md updated
5. ✅ Git tag created
6. ✅ GitHub release published
7. ✅ Release notes generated

## First Time Setup

```bash
# Create initial tag (only needed once)
git tag v0.1.0
git push origin v0.1.0

# OR just push with feat: and it creates v0.1.0 automatically
git commit -m "feat: initial release"
git push origin master
```

## Pro Tips

💡 **Batch changes**: Squash multiple commits in your PR to create one release with all changes

💡 **Skip CI**: Use `[skip ci]` or `[ci skip]` in commit message to prevent release

💡 **Check before push**: Run tests locally: `go test ./...` and `golangci-lint run`

💡 **View releases**: Go to https://github.com/gstachniukrsk/govalidator/releases

## Troubleshooting

**Release didn't trigger?**
- Check commit message format (must be `type: description`)
- Check GitHub Actions logs
- Verify tests passed

**Wrong version?**
- Check your commit type (feat vs fix)
- Check for `BREAKING CHANGE:` in commit body
- Can't change after release - just make a new release

**Need to delete a release?**
```bash
# Delete locally
git tag -d v1.2.3

# Delete remotely
git push origin :refs/tags/v1.2.3

# Delete GitHub release via web UI
```

## Full Documentation

See [RELEASING.md](RELEASING.md) for complete documentation.
