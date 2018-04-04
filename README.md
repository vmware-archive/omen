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

### Toggle product errands

The `toggle-errands` command requires the `--errand-type` option, which can currently 
only accepts `post-deploy` as its value, the target errand state `--action` which can
be one of `enable`, `disable` or `default`.

Optionally, a comma-delimited list of product guids can be supplied as a value for the 
`--products` option.

Some example calls:
```sh
# Enable post-deploy errands for p-bosh-e43806ad5d741db12345
$ omen toggle-errands --errand-type post-deploy --action enable --products p-bosh-e43806ad5d741db12345

# Reset post-deploy errands for all products to their default settings
$ omen toggle-errands --errand-type post-deploy --action default
```
## Running tests

- The `.envrc` sources the enemytest opsmanager credentials from the `secrets-cf-cloudops-sandbox` repo. Make sure you have that checked out in your workspace.
- Run `direnv allow`
- Run `ginkgo -r` in the root of this repo.
