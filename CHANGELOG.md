
<a name="Unreleased"></a>
## [Unreleased](https://github.com/ylallemant/githook-companion/compare/0.8.12...Unreleased) (0001-01-01)


<a name="0.8.12"></a>
## [0.8.12](https://github.com/ylallemant/githook-companion/compare/0.8.11...0.8.12) (2025-04-16)

### Bug Fixes

* truncate message output file

### Documentation

* add language detection comment


<a name="0.8.11"></a>
## [0.8.11](https://github.com/ylallemant/githook-companion/compare/0.8.10...0.8.11) (2025-04-16)

### Bug Fixes

* terminal detection
* long dial timeout on request with credentials

### Features

* add force-default-language flag to validate command
* add debug flag
* add iso-8601 date to default lexemes


<a name="0.8.10"></a>
## [0.8.10](https://github.com/ylallemant/githook-companion/compare/0.8.9...0.8.10) (2025-04-10)

### Bug Fixes

* block "zip slip" arbitrary file access attacks
* remove ignored commit-type from message
* issue-tracker-reference lexeme must begin with a letter

### CI

* add run go tests on pre-push

### Documentation

* add changelog

### Features

* add post-commit and pre-push githooks

### Performance Improvements

* reduced idle connection timeout to 5s


<a name="0.8.9"></a>
## [0.8.9](https://github.com/ylallemant/githook-companion/compare/0.8.8...0.8.9) (2025-04-09)

### Bug Fixes

* typo
* (#7) config directory command always returns an absolute path

### Code Refactoring

* simplified PARENT_CONFIG_REPOSITORY computation

### Documentation

* add information about commit message hook log
* add binary command descriptions

### Features

* add sync configuration
* add timeout to interactive commit type selection
* add pre-commit githook to full configuration init
* add githook-companion configuration


<a name="0.8.8"></a>
## [0.8.8](https://github.com/ylallemant/githook-companion/compare/0.8.7...0.8.8) (2025-04-08)


<a name="0.8.7"></a>
## [0.8.7](https://github.com/ylallemant/githook-companion/compare/0.8.6...0.8.7) (2025-04-02)

### Bug Fixes

* check in file if .local/bin is specified
* shorter timeouts for binary version checks
* make release list fetching non-blocking
* git repository string manipulations
* build ldflags path

### Features

* add global non-blocking flag option
* add commands to lock, unlock, check githooks
* add temporary network locks if connectivity is not available
* add installation directory to .profile if needed
* add locks to mitigate git server api limits


<a name="0.8.6"></a>
## [0.8.6](https://github.com/ylallemant/githook-companion/compare/0.8.5...0.8.6) (2025-03-28)

### Bug Fixes

* made sync workflows non-blocking
* version command information

### Features

* add git hook disable command


<a name="0.8.5"></a>
## [0.8.5](https://github.com/ylallemant/githook-companion/compare/0.8.4...0.8.5) (2025-03-28)

### Features

* enabled binary version check


<a name="0.8.4"></a>
## [0.8.4](https://github.com/ylallemant/githook-companion/compare/0.8.3...0.8.4) (2025-03-28)

### Features

* added authentication for private configuration repositories


<a name="0.8.3"></a>
## [0.8.3](https://github.com/ylallemant/githook-companion/compare/0.8.2...0.8.3) (2025-03-28)

### Features

* add version check and auto-update for the binary
* add auto update of the parent configuration if existing


<a name="0.8.2"></a>
## [0.8.2](https://github.com/ylallemant/githook-companion/compare/0.8.1...0.8.2) (2025-03-26)

### Bug Fixes

* ensure rules are added to both parent and child .gitignore

### Test Coverage

* add pre-commit hook with tests


<a name="0.8.1"></a>
## [0.8.1](https://github.com/ylallemant/githook-companion/compare/0.8.0...0.8.1) (2025-03-26)

### Bug Fixes

* missing commit type on user input


<a name="0.8.0"></a>
## [0.8.0](https://github.com/ylallemant/githook-companion/compare/0.7.8...0.8.0) (2025-03-26)

### Code Refactoring

* restructure command tree for git and dependency

### Features

* add config directory command
* add child flag to config hook and dependency commands
* add new githook templates
* add dictionary weigth to make commit-type assertion configurable
* add automatic latest version detection (#4)
* add githook-companion configuration


<a name="0.7.8"></a>
## [0.7.8](https://github.com/ylallemant/githook-companion/compare/0.7.7...0.7.8) (2025-03-23)

### Bug Fixes

* fixed minimal word-length tokenization error
* fixed formatter unmarshling error

### Features

* add parent configuration repository example

### Test Coverage

* added test todo


<a name="0.7.7"></a>
## [0.7.7](https://github.com/ylallemant/githook-companion/compare/0.7.6...0.7.7) (2025-03-22)

### Features

* add githooks and .gitignore to init command
* add tool version to output


<a name="0.7.6"></a>
## [0.7.6](https://github.com/ylallemant/githook-companion/compare/0.7.5...0.7.6) (2025-03-21)


<a name="0.7.5"></a>
## [0.7.5](https://github.com/ylallemant/githook-companion/compare/0.7.4...0.7.5) (2025-03-21)


<a name="0.7.4"></a>
## [0.7.4](https://github.com/ylallemant/githook-companion/compare/0.7.3...0.7.4) (2025-03-21)


<a name="0.7.3"></a>
## [0.7.3](https://github.com/ylallemant/githook-companion/compare/0.7.2...0.7.3) (2025-03-21)


<a name="0.7.2"></a>
## [0.7.2](https://github.com/ylallemant/githook-companion/compare/0.7.1...0.7.2) (2025-03-21)


<a name="0.7.1"></a>
## [0.7.1](https://github.com/ylallemant/githook-companion/compare/0.7.0...0.7.1) (2025-03-21)


<a name="0.7.0"></a>
## [0.7.0](https://github.com/ylallemant/githook-companion/compare/0.6.3...0.7.0) (2025-03-20)


<a name="0.6.3"></a>
## [0.6.3](https://github.com/ylallemant/githook-companion/compare/0.6.2...0.6.3) (2025-02-28)


<a name="0.6.2"></a>
## [0.6.2](https://github.com/ylallemant/githook-companion/compare/0.6.1...0.6.2) (2025-02-28)


<a name="0.6.1"></a>
## [0.6.1](https://github.com/ylallemant/githook-companion/compare/0.6.0...0.6.1) (2025-02-28)


<a name="0.6.0"></a>
## [0.6.0](https://github.com/ylallemant/githook-companion/compare/0.5.1...0.6.0) (2025-02-28)


<a name="0.5.1"></a>
## [0.5.1](https://github.com/ylallemant/githook-companion/compare/0.5.0...0.5.1) (2025-02-26)


<a name="0.5.0"></a>
## [0.5.0](https://github.com/ylallemant/githook-companion/compare/0.4.0...0.5.0) (2025-02-20)


<a name="0.4.0"></a>
## [0.4.0](https://github.com/ylallemant/githook-companion/compare/0.3.2...0.4.0) (2025-02-19)


<a name="0.3.2"></a>
## [0.3.2](https://github.com/ylallemant/githook-companion/compare/0.3.1...0.3.2) (2025-02-18)


<a name="0.3.1"></a>
## [0.3.1](https://github.com/ylallemant/githook-companion/compare/0.3.0...0.3.1) (2025-02-17)


<a name="0.3.0"></a>
## [0.3.0](https://github.com/ylallemant/githook-companion/compare/0.2.1...0.3.0) (2025-02-17)


<a name="0.2.1"></a>
## [0.2.1](https://github.com/ylallemant/githook-companion/compare/0.2.0...0.2.1) (2025-02-17)


<a name="0.2.0"></a>
## [0.2.0](https://github.com/ylallemant/githook-companion/compare/0.1.0...0.2.0) (2025-02-17)


<a name="0.1.0"></a>
## [0.1.0](https://github.com/ylallemant/githook-companion/compare/0.0.1...0.1.0) (2025-02-15)


<a name="0.0.1"></a>
## 0.0.1 (2025-02-14)

