TODOs
===

## init command
- ~~install githooks~~
- ~~append to .gitignore~~

## debug command
- ~~add tool version~~

## configuration
- adapt commit-types to this norm : https://www.conventionalcommits.org/en/v1.0.0/
- test marshalling (check for null)/unmarschalling (formatter error)

## path commands (dependency, hook)
- add `child` flag to force the path for the child

## Githooks
- lock hook
- unlock hook
- unlock all
- disable hooks
- enable hooks

## Cleanup Commit Message history

- list all commit hashes : `git log --pretty=format:%H`
- collect all commit information : `git show --no-color <commit-hast>`
  - extract message
  - process message
  - ignore untyped

