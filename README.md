# strongo/cicd

Shared CI/CD for Go repositories across the `strongo`, `dal-go`, `sneat-co`,
`ingitdb`, and `bots-go-framework` orgs: reusable GitHub **workflows** and a
composite **action** that run `get`, `vet`, `build`, `test`, `lint`, optional
coverage, and automatic SemVer tagging.

> **Renamed:** this repo was `strongo/go-ci-action`. GitHub redirects the old
> path, so existing `uses: strongo/go-ci-action/...` references keep working.
> New references should use `strongo/cicd`; the Renovate preset below migrates
> the old name for you.

## What's here

| File | Kind | Purpose |
| --- | --- | --- |
| `.github/workflows/workflow.yml` | Reusable workflow (`workflow_call`) | Full Go CI: lint, test (+coverage), build, and version bump. The primary entry point. |
| `.github/workflows/release.yml` | Reusable workflow (`workflow_call`) | GoReleaser release flow (tag + `goreleaser release`). |
| `action.yml` | Composite action | Single-job CI for callers that want CI steps inline in their own job. |
| `default.json` | Renovate preset | Shareable config consumers `extends` to auto-track this repo's tag (see below). |

## Recommended usage: pin to a tag, not `@main`

Pinning `@main` means one bad commit to the shared workflow breaks **every**
consumer's CI at the same time — there is no blast-radius firebreak. Pin to a
version tag instead:

```yaml
jobs:
  ci:
    uses: strongo/cicd/.github/workflows/workflow.yml@v1   # moving major tag
    secrets:
      GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Two pinning styles are supported:

- **`@v1` — moving major tag (recommended default).** A lightweight tag that the
  maintainer advances deliberately to the latest backward-compatible release.
  You automatically get fixes without opening a PR, but a bad release only
  reaches you when `v1` is advanced (not on every push to `main`).
- **`@v1.x.y` — exact release (maximum control).** Pin an immutable release and
  let **Renovate** open a PR to bump it (see below). Each bump runs through your
  own CI before merging, giving you a full per-repo firebreak and an audit trail.

`@main` still works and stays supported for backward compatibility, but is
discouraged for the reasons above.

> Existing `@main` consumers are **not** being mass-migrated. Adopt `@v1` (or the
> Renovate preset) gradually, per repo, on your own schedule.

## Keep the pin fresh with Renovate

Add the shared preset to a consumer repo's `renovate.json` so Renovate keeps the
`strongo/cicd` reference current — and rewrites any lingering
`strongo/go-ci-action` reference onto `strongo/cicd@v1` — automatically:

```json
{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "extends": [
    "config:recommended",
    "github>strongo/cicd"
  ]
}
```

The preset (`default.json` in this repo):

- Groups and auto-updates the `strongo/cicd` reusable-workflow / action ref,
  advancing `@v1.x.y` pins (and following the moving `@v1`) as releases are cut.
- Auto-merges those bumps **through a PR gated by your CI**, so a broken release
  fails your build and blocks the merge — the firebreak — instead of landing
  silently.
- Replaces legacy `strongo/go-ci-action` references with `strongo/cicd@v1`,
  automating the rename find-and-replace.

Override anything you like (e.g. disable `automerge`) in your own `renovate.json`
after the `extends`.

## Automatic version tagging

On a push/merge to `main`, the `go_bump` job (workflow) / tag step (action) runs
[`mathieudutour/github-tag-action`], which reads **conventional commits since the
last tag** and pushes a new SemVer tag:

| Commit type since last tag | Result |
| --- | --- |
| `fix:` | patch bump (`v1.2.3 → v1.2.4`) |
| `feat:` | minor bump (`v1.2.3 → v1.3.0`) |
| `feat!:` / `BREAKING CHANGE:` | major bump — only if `allow_major_version_bump: true` |
| docs/chore/ci/refactor only | `default_bump` decides (see below) |

### Required repo settings (read this — it's why tags were being missed)

`github-tag-action` derives the bump from the commit range `lastTag..HEAD`. Two
things must be true for it to see your `feat:`/`fix:` commits:

1. **Full checkout.** The job checks out with **`fetch-depth: 0`** (fixed in this
   repo). With the default shallow clone the action only sees `HEAD`'s message —
   so a `Merge pull request #N …` merge commit (which is never
   conventional-commit-shaped) produced **no bump and no tag**, even when the
   merged branch was full of `feat:`/`fix:` commits. This was the cause of
   releases needing hand-cut tags. If you consume the **composite `action.yml`**,
   your *own* workflow must `actions/checkout` with `fetch-depth: 0`.

2. **Conventional commits reach `main`.** Choose one, and set it in your repo's
   **Settings → General → Pull Requests**:
   - **Squash merge + "Default to PR title" for the squash commit message**, and
     write conventional PR titles (`feat: …`, `fix: …`). The single squashed
     commit is then conventional. *(Recommended — simplest and most robust.)*
   - **Merge commits**, with conventional commits on your branches. `fetch-depth:
     0` lets the action read them behind the merge commit.

    Either way, enabling a linear-history / conventional-PR-title convention makes
    tagging deterministic.

### `default_bump` input

Controls what happens when **no** commit since the last tag implies a bump
(docs/chore/ci-only changes):

- Reusable `workflow.yml`: **`default_bump: 'patch'`** by default (a push/merge to
  `main` always cuts at least a patch tag — preserves prior behaviour). Set
  `default_bump: 'false'` to tag *only* on `feat:`/`fix:`/breaking commits.
- Composite `action.yml`: **`default_bump: 'false'`** by default (tags only on
  conventional commits). Set to `'patch'`/`'minor'`/`'major'` to always bump.

## Releasing this repo (maintainers)

Tags are cut automatically by this repo's own CI (`v1.x.y`). To publish or advance
the **moving `v1`** major tag after a release lands on `main`:

```bash
git fetch --tags origin
# point v1 at the newest v1.x.y release (or origin/main)
git tag -f v1 "$(git tag -l 'v1.*.*' | sort -V | tail -1)"
git push -f origin v1
```

Advance `v1` only to releases you've verified are backward-compatible.

<!-- dev-approach:v1 -->
## Our approach to development

We build with our own tooling:

- **[SpecScore](https://specscore.md)** — specify requirements as `SpecScore.md` artifacts
- **[SpecStudio](https://specscore.studio)** — author & manage specs across their lifecycle
- **[inGitDB](https://ingitdb.com)** — store structured data in Git where applicable
- **[DALgo](https://dalgo.io)** — data access layer for Go
- **[cover100.dev](https://cover100.dev)** — drive toward 100% test coverage
- **[DataTug](https://datatug.io)** — query & explore data
<!-- /dev-approach -->

[`mathieudutour/github-tag-action`]: https://github.com/mathieudutour/github-tag-action
