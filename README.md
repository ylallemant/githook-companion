githooks-butler
====

Provides a set of tools to make git-hook configuration and utilisation easy.

## Features

- init git-hooks in a repository or globally
- lists and installs tools needed inside the git-hooks
- retrieves information from the Git configuration - like origin server hostname
- interactive commit message standardisation

## Install

```sh
curl -fsSL https://github.com/ylallemant/githooks-butler/raw/0.2.0/install.sh | bash
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
