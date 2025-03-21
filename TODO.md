TODOs
===

## Githook Locks

- locks for specific hooks to avoid loops
- lock files timestamp check for validity, if they haven'f been deleted because of an execution error

## Cleanup Commit Message history

- list all commit hashes : `git log --pretty=format:%H`
- collect all commit information : `git show --no-color <commit-hast>`
  - extract message
  - process message
  - ignore untyped

