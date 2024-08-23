# Monogo CLI Tool

Monogo is a 0 configuration tool companion for go workspace monorepo. It is like Turborepo but for Go

## Installation

To install Monogo, clone the repository and build the binary:

```sh
go install github.com/nicolasgere/monogo@latest
```

## Usage

### Commands

#### `install`

Install dependencies for every module.

```sh
monogo install [flags]
```

#### `fmt`

Format every module.

```sh
monogo fmt [flags]
```

#### `test`

Run tests for every module.

```sh
monogo test [flags]
```

### Flags

All commands support the following flags:

- `--target, -t`: Specify a targeted module.
- `--dependency, -d`: Run with all dependencies of the target (both descendants and ascendants).
- `--branch, -b`: Compare the current branch with the master branch and find affected modules.
- `--path, -p`: Directory to run the cli in, default .
