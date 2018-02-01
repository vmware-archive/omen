# OMEN 

## Installation

Dependencies aren't committed to this repo, but are expected to be fetched with `dep`.

- [Install `dep`](https://github.com/golang/dep/#setup)
- `dep ensure`
- `go install`

## Usage

The tool uses `OPSMAN_HOSTNAME`, `OPSMAN_USER`, and `OPSMAN_PASSWORD` environment variables.

### Grab Ops Manager diagnostic report with:

```sh
omen diagnostics
```

### Write out staged tiles report with:

```sh
omen staged-tiles -o outputDir
```
This command will write out the config of all tiles (including BOSH).

### Grab the manifest report with:

```sh
omen manifests
```

### Apply changes with:

```sh
omen apply-changes
```

## Running tests

- The `.envrc` sources the enemytest opsmanager credentials from the `secrets-cf-cloudops-sandbox` repo. Make sure you have that checked out in your workspace.
- Run `direnv allow`
- Run `ginkgo -r` in the root of this repo.
