<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/spec-builddep"><img src="https://kaos.sh/r/spec-builddep.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/y/spec-builddep"><img src="https://kaos.sh/y/83c5070ce37641f19d5bf8174847d430.svg" alt="Codacy badge" /></a>
  <a href="https://kaos.sh/w/spec-builddep/ci"><img src="https://kaos.sh/w/spec-builddep/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/spec-builddep/codeql"><img src="https://kaos.sh/w/spec-builddep/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`spec-builddep` is a simple utility for installing dependencies for building an RPM package (`yum-builddep` _drop-in replacement_).

### Installation

#### From source

To build the `spec-builddep` from scratch, make sure you have a working Go 1.23+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/spec-builddep@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/spec-builddep/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) spec-builddep
```

#### From [ESSENTIAL KAOS Public Repository](https://kaos.sh/kaos-repo)

```bash
sudo dnf install -y https://pkgs.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo dnf install spec-builddep
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo spec-builddep --completion=bash 1> /etc/bash_completion.d/spec-builddep
```

ZSH:
```bash
sudo spec-builddep --completion=zsh 1> /usr/share/zsh/site-functions/spec-builddep
```

Fish:
```bash
sudo spec-builddep --completion=fish 1> /usr/share/fish/vendor_completions.d/spec-builddep.fish
```

### Man documentation

You can generate man page using next command:

```bash
spec-builddep --generate-man | sudo gzip > /usr/share/man/man1/spec-builddep.1.gz
```

### Usage

<img src=".github/images/usage.svg" />

### CI Status

| Branch | Status |
|--------|----------|
| `master` | [![CI](https://kaos.sh/w/spec-builddep/ci.svg?branch=master)](https://kaos.sh/w/spec-builddep/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/spec-builddep/ci.svg?branch=develop)](https://kaos.sh/w/spec-builddep/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/.github/blob/master/CONTRIBUTING.md).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://kaos.dev"><img src="https://raw.githubusercontent.com/essentialkaos/.github/refs/heads/master/images/ekgh.svg"/></a></p>
