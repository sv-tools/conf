# conf

[![Code Analysis](https://github.com/sv-tools/conf/actions/workflows/checks.yaml/badge.svg)](https://github.com/sv-tools/conf/actions/workflows/checks.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sv-tools/conf.svg)](https://pkg.go.dev/github.com/sv-tools/conf)
[![codecov](https://codecov.io/gh/sv-tools/conf/branch/main/graph/badge.svg?token=0XVOTDR1CW)](https://codecov.io/gh/sv-tools/conf)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sv-tools/conf?style=flat)](https://github.com/sv-tools/conf/releases)

The configuration reader with as few dependencies as possible.

The library provides the base code and the interfaces only.
All parsers and readers must be created in the separate repositories to avoid unnecessary dependencies.

```shell
go get github.com/sv-tools/conf
```

## Dependencies

* The [spf13/cast](https://github.com/spf13/cast) has been added as dependency to avoid the code duplication. I will make a hard copy of it if the number of dependencies are increased.
* The [stretchr/testify](https://github.com/stretchr/testify) is used in tests only.

## Addons

* [Go Templates Transformer](https://github.com/sv-tools/conf-transformer-go-template) supports go templates by parsing and applying the templates stored in the configuration manager.
* [JSON Parser](https://github.com/sv-tools/conf-parser-json) reads a data in JSON format.
* [YAML Parser](https://github.com/sv-tools/conf-parser-yaml) reads a data in YAML format.
* [Env reader](https://github.com/sv-tools/conf-reader-env) reads the values from environment variables.
* [Flags reader](https://github.com/sv-tools/conf-reader-flags) reads the command line flags (using [pflag](https://github.com/spf13/pflag))

## Alternatives

* [viper](https://github.com/spf13/viper) is the most know library, it's very heavy and very rich in various features.
* [koanf](https://github.com/knadh/koanf) is an attempt to make a better version of the `viper`, but also contains all parsers in same repo, so the list of dependencies is pretty huge.

## License

MIT licensed. See the bundled [LICENSE](LICENSE) file for more details.
