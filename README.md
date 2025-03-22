githook-companion
====

Provides a set of commands to make git-hook configuration and utilisation easy.

> [!CAUTION]
> The project is usable but in its early stages, feedback appreciated.
> Please file issues if you encounter problems.
> Provide the output of the [debug command](#debug).
>
> The main focus has been put on Linux and Darwin.
> Windows has not been tested at all, but support is planed

## Binary

### Installation

```sh
curl -fsSL https://github.com/ylallemant/githook-companion/raw/0.7.0/install.sh | bash
```

### Update

```sh
githook-companion update [--force]
```

### Debug

The debug command outputs information about the environment

```Bash
githook-companion debug
```


## Features

Core consept of the tool is to hardcore as little as possible and to let declaration through configuration do the heavy lifting.
This gives the tool great flexibility and adaptation potential to various use cases.


### Conventional Commits

The core functionality of the tool is the standardization of commit messages.
The standardization allows for subsequent processings like the automation of changelogs.

`githook-companion` allows to freely define any number of `commit-types` you want to use.

Inspiration for the configuration are `Conventional Commits` patterns like [this one](https://gist.github.com/qoomon/5dfcdf8eec66a051ecd85625518cfd13).
Different projects have different people, different needs; so this tool won't force you into any specific pattern, although we propose one in the default configuration.

> [!TIP]
> You can test/tweek your commit configuration effects with following command:
> 
> ```bash
> githook-companion commit validate -m <commit-message> [--debug]
> ```

### Parent Configurations

You can reference a centralised configuration in another Git repository to ease the maintenance of your configuration.

`githook-companion` will automatically checkout the parent configuration if needed when you use the [init command](#initialize).

> [!TIP]
> You may also have different project collections use diffrent parent configurations, for different purposes, customers, ...

### Dependecies

You may want to use specific binaries in your githooks.

`githook-companion` allows you to list them in the configuration and will automatically download when you use the [init command](#initialize).

> [!TIP]
> Using the parent configurations, you can define different tools in the parent and the child project:
> - in the parent for common functionalities (like changelog automation)
> - in the child for project specific behaviours (linting, validation, formatting)

### Convenience Commands

`githook-companion` provides some convenience commands to retrieve "complex" information like from the Git configuration.

This is done to avoid too complex manipulations in the githook scripts themselves.


## Usage

Navigate to your project or parent configuration project folder.

### Initialize

Initialize githooks in a specific project.
You can do it on a global level in your home directory with the flag `global`.

- `init` will checkout a parent configuration if necessary
- `init` will download dependencies if necessary
- `init` will at least install the `prepare-commit-msg` hook
- `init` will set the Git config property `core.hooksPath` to your githook folder path

```Bash
githook-companion init
```

If you want to use a parent configuration for your project add following flags :

- `minimalistic` : will create a minimal configuration file
- `parent-repository` : specifies a centralized configuation repository.
- `parent-path` : specifies where the centralized repository will be cloned into. The path must be relative to your project folder.


```Bash
githook-companion init --minimalistic --parent-repository <git-repository-url> --parent-path <relative-path>
```

> [!WARNING]
> We do not recommend to use the `global` flag. It will result in the Git config property `core.hooksPath` being set globally for all your projects.
> If you have only a few private projects, it may be ok, but not if you contribute to a lot of different projects.
>
> Rather use the parent configuration pattern, which allows to isolate different project governance setups.


### Configure

The configuration consist of different blocks used for different purposes, like related to commit message standardisation or dependency management.

You can find a list of all blocks and properties in the [CONFIGURATION.md](./CONFIGURATION.md)
