# Release Process

This document describes the automated release process for the govalidator library.

## Overview

The library uses **automated semantic versioning** based on conventional commits. When you push to the `master` branch, GitHub Actions automatically:

1. Analyzes commit messages
2. Determines the next version number
3. Creates a release PR or publishes a release
4. Updates CHANGELOG.md
5. Creates a git tag

## Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) to determine version bumps:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Version Bump Rules

| Commit Type | Version Bump | Example |
|-------------|--------------|---------|
| `feat:` | **Minor** (0.x.0) | `feat: add JSONPresenter for API errors` |
| `fix:` | **Patch** (0.0.x) | `fix: correct rune length calculation in MinLengthValidator` |
| `perf:` | **Patch** (0.0.x) | `perf: optimize path traversal in validator` |
| `feat!:` or `BREAKING CHANGE:` | **Major** (x.0.0) | `feat!: rename IntValidator to IsIntegerValidator` |
| `docs:`, `style:`, `refactor:`, `test:`, `chore:` | **No bump** | `docs: update README with new presenters` |

### Examples

**Minor version bump (new feature):**
```bash
git commit -m "feat: add CombinedPresenter for path+error formatting"
```

**Patch version bump (bug fix):**
```bash
git commit -m "fix: MinLengthValidator now correctly counts Unicode characters"
```

**Major version bump (breaking change):**
```bash
git commit -m "feat!: rename validators for consistency

BREAKING CHANGE: IntValidator renamed to IsIntegerValidator,
StringValidator renamed to IsStringValidator"
```

**No version bump (documentation):**
```bash
git commit -m "docs: add examples for new presenters"
```

## Release Workflow

The project uses **fully automated releases** based on conventional commits.

### How Automatic Releases Work

**File:** `.github/workflows/release.yml`

**Process:**
1. On every push to `master`, GitHub Actions:
   - Runs all tests
   - Runs linters
   - Analyzes commits since last release
   - Determines version bump based on commit types
2. If there are release-worthy commits (feat/fix/etc):
   - Creates a new GitHub release
   - Tags the commit with the version
   - Updates CHANGELOG.md
   - Publishes release notes

**Benefits:**
- ✅ Fully automated - no manual intervention needed
- ✅ Fast releases - happens immediately on push to master
- ✅ Tests and linters run before every release
- ✅ CHANGELOG automatically updated
- ✅ Zero maintenance overhead

### Usage Example

```bash
# Commit with conventional format and push
git add .
git commit -m "feat: add JSONPresenter"
git push origin master

# GitHub Actions immediately:
# 1. Runs tests ✓
# 2. Runs linters ✓
# 3. Creates release v0.2.0
# 4. Updates CHANGELOG.md
# 5. Tags the commit
# 6. Publishes release notes

# Done! Release is live in ~2-3 minutes
```

### Important Notes

**Every push to master can create a release** if it contains:
- `feat:` commits (new features)
- `fix:` commits (bug fixes)
- `perf:` commits (performance improvements)
- Commits with `BREAKING CHANGE:` or `!`

**These do NOT create releases:**
- `docs:` commits
- `style:` commits
- `refactor:` commits
- `test:` commits
- `chore:` commits

**Best Practice:** If you want to batch multiple changes into one release, squash/merge multiple commits in your PR before merging to master.

## Initial Release

To create your first release:

```bash
# Create an initial version tag
git tag v0.1.0
git push origin v0.1.0

# Or let the workflow create it automatically
git commit -m "feat: initial release of govalidator library"
git push origin master
```

## Go Module Versioning

For Go modules, version tags must follow these rules:

- **v0.x.x**: Initial development (pre-v1)
- **v1.x.x**: Stable API, semantic versioning
- **v2.x.x+**: Major versions require module path suffix (e.g., `/v2`)

### Breaking Changes in v2+

If you need to make breaking changes after v1.0.0:

```bash
# Option 1: Stay on v1 and avoid breaking changes (recommended)
git commit -m "feat: add new non-breaking feature"

# Option 2: Move to v2 (requires changing import paths)
# Update go.mod:
# module github.com/gstachniukrsk/govalidator/v2

git commit -m "feat!: major breaking changes for v2"
git push origin master
```

## CHANGELOG.md

The CHANGELOG is automatically generated from commit messages. It follows [Keep a Changelog](https://keepachangelog.com/) format.

Structure:
- **Added**: New features (`feat:`)
- **Changed**: Changes in existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Removed features
- **Fixed**: Bug fixes (`fix:`)
- **Security**: Security fixes

## Manual Release (Emergency)

If automated releases fail, you can create a manual release:

```bash
# 1. Create and push tag
git tag v1.2.3
git push origin v1.2.3

# 2. Create release on GitHub
gh release create v1.2.3 \
  --title "v1.2.3" \
  --notes "Release notes here"

# Or use the GitHub web UI
```

## Troubleshooting

### Release PR not created
- Check commit messages follow conventional commit format
- Ensure there are commits since the last release
- Check GitHub Actions logs for errors

### Release failed
- Check that tests pass: `go test ./...`
- Check that linting passes: `golangci-lint run`
- Verify GITHUB_TOKEN has write permissions

### Wrong version number
- Check commit messages - they determine version bumps
- If needed, manually adjust the version in the release PR

## Best Practices

1. **Write clear commit messages**: They become your CHANGELOG
2. **Batch related changes**: Group related commits before merging to master
3. **Use fix vs feat correctly**: This affects version numbers
4. **Document breaking changes**: Always include `BREAKING CHANGE:` in commit body
5. **Test before merging**: CI runs tests, but test locally first
6. **Review release PRs**: Don't auto-merge release-please PRs without review

## Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Go Modules Version Numbers](https://go.dev/doc/modules/version-numbers)
- [release-please](https://github.com/googleapis/release-please)
