githook-companion
====

Provides a set of commands to make git-hook configuration and utilisation easy.

## Features

- lists and installs dependencies needed inside the git-hooks
- retrieves information from the Git configuration - like origin server hostname
- commit message standardisation through configuration

## Install

```sh
curl -fsSL https://github.com/ylallemant/githook-companion/raw/0.5.0/install.sh | bash
```

## Tests

### Run

```bash
go test -cover ./...
```

### Coverage

Enable visual coverage feedback in `vscode` :

```json
{
    "go.coverOnSave": true
}
```

## Cleanup Commit Message history

- list all commit hashes : `git log --pretty=format:%H`
- collect all commit information : `git show --no-color <commit-hast>`
  - extract message
  - process message
  - ignore untyped

