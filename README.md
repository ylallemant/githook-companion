githook-companion
====

Provides poweful possibilities to define a commit message standard and a set of commands to make git-hook configuration and utilisation easy.

The main goal is to be still able to use `Git` as your primary tool and to not be forced into using some of the multiple wrapper commands other projects propose.
Sticking with Git means any workflow/tool you build/use on top of it will not break.

> [!CAUTION]
> The project is usable but in its early stages, feedback appreciated.
> Please file issues if you encounter problems.
> Provide the output of the [debug command](#debug).
>
> The main focus has been put on Linux and Darwin.
> Windows has not been tested at all, but support is planed
>
> the threshold configuration for the language detector is still
> being tweaked and may be a nuisance. You can see the confidence value
> by using `--debug` ([see here](#test-commit-standard)) or force the usage
> of the defaut language defined in the `language_detection_options`
> with the `--force-default-language` flag.
>
> If anybody has a good suggestion for a language detection library that is Go compatible, I am interested.
> Main issue is to get reliable detections on short sentences that also have
> good chances to use technical terms in English - which makes it harder.

## Binary

### Installation

```sh
curl -fsSL https://github.com/ylallemant/githook-companion/raw/main/install.sh | bash
```

### Upgrade

```sh
githook-companion upgrade [--force]
```

### Debug

The debug command outputs information about the environment

```Bash
githook-companion debug
```

### Remove

The remove command will remove the local configuration folder `.githook-companion` and unset the Git config property `core.hooksPath` (disable the githooks).

```Bash
githook-companion remove
```


## Features

Core consept of the tool is to hardcore as little as possible and to let declaration through configuration do the heavy lifting.
This gives the tool great flexibility and adaptation potential to various use cases.


### Conventional Commits

The core functionality of the tool is the standardization of commit messages.
The standardization allows for subsequent processings like the automation of changelogs.

`githook-companion` allows to freely define any number of `commit-types` you want to use ([example](https://github.com/ylallemant/githooks/blob/3533e5d6aa7f49a5582a9f133e86728bed3f613a/.githook-companion/config.yaml#L3)).

Inspiration for the configuration are `Conventional Commits` patterns like [this one](https://www.conventionalcommits.org/en/v1.0.0/).
Different projects have different people, different needs; so this tool won't force you into any specific pattern, although we propose one in the [default configuration](https://github.com/ylallemant/githooks/blob/3533e5d6aa7f49a5582a9f133e86728bed3f613a/.githook-companion/config.yaml#L22).

The tool uses tokenization to provide complex checks and formatting capabilities:

- restrict commit message language(s) ([example](https://github.com/ylallemant/githooks/blob/3533e5d6aa7f49a5582a9f133e86728bed3f613a/.githook-companion/config.yaml#L28))
- define regular expression based [lexemes](https://en.wikipedia.org/wiki/Lexical_analysis#Lexical_token_and_lexical_tokenization) for the tokenization ([example](https://github.com/ylallemant/githooks/blob/3533e5d6aa7f49a5582a9f133e86728bed3f613a/.githook-companion/config.yaml#L94))
- use regular expressions and Go templates for normalization ([example](https://github.com/ylallemant/githooks/blob/3533e5d6aa7f49a5582a9f133e86728bed3f613a/.githook-companion/config.yaml#L102))
- define dictionaries for the tokenization. comparation done with cleaned and [lemmatized](https://en.wikipedia.org/wiki/Lemmatization#Description) words. ([example](https://github.com/ylallemant/githooks/blob/3533e5d6aa7f49a5582a9f133e86728bed3f613a/.githook-companion/config.yaml#L31))
- define Go Templates with available [helper functions](https://masterminds.github.io/sprig/) to format the final message. You can use the `camel-case` name of the tokens to reference them in the template ([example](https://github.com/ylallemant/githooks/blob/3533e5d6aa7f49a5582a9f133e86728bed3f613a/.githook-companion/config.yaml#L21))

#### Test Commit Standard

> [!TIP]
> You can test/tweak your commit configuration effects with following command:
> 
> ```bash
> githook-companion git commit validate -m "commit-message" [--debug] [--force-default-language]
> ```

### Parent Configurations

You can reference a centralised configuration in another Git repository to ease the maintenance of your configuration.

In the presence of a child configuration, `githook-companion` will automatically checkout the parent configuration if needed when you use the [init command](#initialize).

The hooks in the parent configuration can also trigger project specific hook-scripts in the child project.

Here you can find an [example of parent configuration](https://github.com/ylallemant/githooks/tree/main/.githook-companion)


> [!TIP]
> You may also have different project collections use diffrent parent configurations, for different purposes, customers, ...
>
> A parent configuration repository can be private. You will have to use a `~/.git-credentials` configuration or an ssh agent.

### Dependecies

You may want to use specific binaries in your githooks.

`githook-companion` allows you to list them in the configuration and it will automatically download them when you use the [init command](#initialize).

> [!TIP]
> Using the parent configurations, you can define different tools in the parent and the child project:
> - in the parent for common functionalities (like changelog automation)
> - in the child for project specific behaviours (linting, validation, formatting)

### Commands

`githook-companion` provides various commands like some convenience commands to retrieve "complex" information like from the Git configuration.

This is done to avoid too complex manipulations in the githook scripts themselves.

#### Environment Commands

```bash
# checks if the command is running within a terminal
# returns a "true" or "false"
githook-companion environment terminal
```


#### Configuration Commands

```bash
# synchronizes the parent configuration if necessary
# very convenient to ensure that the latest configuration is used
githook-companion config sync
```

```bash
# prints the used configuration
githook-companion config show
```

```bash
# returns the absolute path to the parent or child githook directory
# (if no parent configuration is used both parent and child paths are equal)
githook-companion config hook path [--child]
```

```bash
# returns the absolute path to the parent or child dependency directory
# (if no parent configuration is used both parent and child paths are equal)
githook-companion config dependency path [--child]
```

```bash
# returns the absolute path to the parent or child configuration directory
# (if no parent configuration is used both parent and child paths are equal)
githook-companion config directory [--child]
```

```bash
# returns computer's OS
githook-companion config os
```

```bash
# returns computer's CPU architecture
githook-companion config arch
```

#### Git Commit Commands

```bash
# list configured commit types
githook-companion git commit type list
```

```bash
# validate commit
githook-companion git commit validate 'message' [--debug]
```

#### Git Hooks Commands

```bash
# disables the locally the githooks by unsetting "core.hooksPath"
githook-companion git hook disable
```

```bash
# locks a specific githook
githook-companion git hook lock "hook-name"
```

```bash
# unlocks a specific or all githooks
githook-companion git hook unlock  "hook-name" | --all
```

```bash
# checks if a specific githook is locked - returns "true" or "false"
githook-companion git hook locked "hook-name"
```

#### Git Server Commands

```bash
# returns the server's hostname
githook-companion git server hostname
```

```bash
# returns the server's name (shorthand for the provider: "github", ...)
githook-companion git server name
```

```bash
# returns the normalised repository uri
githook-companion git server repository
```



## Usage

Navigate to your project or parent configuration project folder.

### Initialize

### As Full Configuration (And Potential Parent)

Initialize githooks in a specific project.
You can do it on a global level in your home directory with the flag `global`.

- `init` will checkout a parent configuration if necessary
- `init` will download dependencies if necessary
- `init` will create a configuration with our defaults if necessary
- `init` will at least install the `prepare-commit-msg` hook if necessary
- `init` will add `githook-companion` specific entries to your `.gitignore` if necessary
- `init` will set the Git config property `core.hooksPath` to your githook folder path if necessary

```Bash
githook-companion init
```

### As Minimalistic Child Configuration

If you want to use a parent configuration for your project add following flags :

- `minimalistic` : will create a minimal configuration file
- `parent-repository` : specifies a centralized configuation repository.
- `parent-path` : specifies where the centralized repository will be cloned into. The path must be relative to your project folder.

You have an [example in this repository](https://github.com/ylallemant/githook-companion/blob/main/.githook-companion/config.yaml) ("eat your own dog food").

```Bash
githook-companion init --minimalistic --parent-path "../githooks" --parent-repository https://github.com/ylallemant/githooks
```

> [!WARNING]
> We do not recommend to use the `global` flag. It will result in the Git config property `core.hooksPath` being set globally for all your projects.
> If you have only a few private projects, it may be ok, but not if you contribute to a lot of different projects.
>
> Rather use the parent configuration pattern, which allows to isolate different project governance setups.


### Configure

The configuration consist of different blocks used for different purposes, like related to commit message standardisation or dependency management.

You can find a list of all blocks and properties in the [CONFIGURATION.md](./CONFIGURATION.md)
