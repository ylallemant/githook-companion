Centralised Hook
===

## What We Want

Use case :
- you want to define a standard for commit messages accross multiple teams and projects.
- the standardised commit message will enable an automated changelog generation
- you still want the teams to be able to have project specific hook tasks
- you want this standard to be easily maintanable

## How We Get It

This example shows a way to implement those requirements :
- select a tool for the changelog generation
- configure your commit message standard in the `githook-companion`
- define the git-hooks
- create all this setting in a Git repo that teams can checkout and sync
- use following Git command to activate the hooks: `git config --global core.hooksPath "<repo-clone-path>/hooks"`
