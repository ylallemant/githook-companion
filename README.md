githook-companion
====

Provides a set of commands to make git-hook configuration and utilisation easy.

> [!CAUTION]
> The project is usable but in its early stages, feedback appreciated.
> Please file issues if you encounter problems

## Features

- lists and installs dependencies needed inside the git-hooks
- checks out centralized configurations repositories if used
- retrieves information from the Git configuration - like origin server hostname
- commit message standardisation through configuration

## Install

### Binary

```sh
curl -fsSL https://github.com/ylallemant/githook-companion/raw/0.5.0/install.sh | bash
```

### Initialize

Initialize githooks in a specific project.
You can do it on a global level in your home directory with the flag `global`.

```Bash
githook-companion init
```

If you want to use a centralized configuration for your project add following flags :

- `minimalistic` : will create a minimal configuration file
- `reference-repository` : specifies a centralized configuation repository.
- `reference-path` : specifies where the centralized repository will be cloned into. The path must be relative to your project folder.


```Bash
githook-companion init --minimalistic --reference-repository <git-repository-url> --reference-path <relative-path>
```

### Configure

The configuration consist of different blocks used for different purposes, like related to commit message standardisation or dependency management.

You can find a list of all blocks and properties in the [CONFIGURATION.md](./CONFIGURATION.md)
