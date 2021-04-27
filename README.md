# conf

[![Code Analysis](https://github.com/sv-tools/conf/actions/workflows/checks.yaml/badge.svg)](https://github.com/sv-tools/conf/actions/workflows/checks.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sv-tools/conf.svg)](https://pkg.go.dev/github.com/sv-tools/conf)
[![codecov](https://codecov.io/gh/sv-tools/conf/branch/main/graph/badge.svg?token=0XVOTDR1CW)](https://codecov.io/gh/sv-tools/conf)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sv-tools/conf?style=flat)](https://github.com/sv-tools/conf/releases)

The configuration reader with as few dependencies as possible.

The library provides the base code only and the interfaces. All parsers and readers must be places in the separate repositories.


## Dependecies

* The `github.com/spf13/cast` has been added as dependency to avoid the code duplication.
This library has only one dependecy of `github.com/stretchr/testify`.
I will make a hard copy of it if the number of dependencies are increased.
* The `github.com/stretchr/testify` is used for testing only.


## Alternatives

* `viper` (https://github.com/spf13/viper) is the most know library, it's very heavy and very reach in defferent features.
* `koanf` (https://github.com/knadh/koanf) is an attempt to make a better version of the `viper`, but also contains all parsers in same repo, so the list of dependeciies is pretty huge.


## License

MIT licensed. See the bundled [LICENSE](LICENSE) file for more details.
