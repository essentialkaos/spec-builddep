<p align="center"><a href="#readme"><img src="https://gh.kaos.st/spec-builddep.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/spec-builddep"><img src="https://kaos.sh/r/spec-builddep.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/l/spec-builddep"><img src="https://kaos.sh/l/1008b1e64602a52fa7d7.svg" alt="Code Climate Maintainability" /></a>
  <a href="https://kaos.sh/b/spec-builddep"><img src="https://kaos.sh/b/e1d77494-93c2-4bd7-aee4-c7898dcb2afa.svg" alt="Codebeat badge" /></a>
  <a href="https://kaos.sh/w/spec-builddep/ci"><img src="https://kaos.sh/w/spec-builddep/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/spec-builddep/codeql"><img src="https://kaos.sh/w/spec-builddep/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`spec-builddep` is a simple utility for installing dependencies for building an RPM package (`yum-builddep` _drop-in replacement_).

### Installation

#### From source

To build the `spec-builddep` from scratch, make sure you have a working Go 1.20+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/spec-builddep@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/spec-builddep/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) spec-builddep
```

#### From [ESSENTIAL KAOS Public Repository](https://pkgs.kaos.st)

```bash
sudo yum install -y https://pkgs.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo yum install spec-builddep
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

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
