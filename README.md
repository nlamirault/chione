# Chione

[![License Apache 2][badge-license]](LICENSE)
[![GitHub version](https://badge.fury.io/gh/nlamirault%2Fchione.svg)](https://badge.fury.io/gh/nlamirault%2Fchione)

* Master : [![Circle CI](https://circleci.com/gh/nlamirault/chione/tree/master.svg?style=svg)](https://circleci.com/gh/nlamirault/chione/tree/master)
* Develop : [![Circle CI](https://circleci.com/gh/nlamirault/chione/tree/develop.svg?style=svg)](https://circleci.com/gh/nlamirault/chione/tree/develop)

This tool is a simple CLI to display informations about skiing resorts.

![Screenshot](chione.png)


## Installation

You can download the binaries :

* Architecture i386 [ [linux](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_linux_386) / [darwin](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_darwin_386) / [freebsd](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_freebsd_386) / [netbsd](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_netbsd_386) / [openbsd](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_openbsd_386) / [windows](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_windows_386.exe) ]
* Architecture amd64 [ [linux](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_linux_amd64) / [darwin](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_darwin_amd64) / [freebsd](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_freebsd_amd64) / [netbsd](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_netbsd_amd64) / [openbsd](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_openbsd_amd64) / [windows](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_windows_amd64.exe) ]
* Architecture arm [ [linux](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_linux_arm) / [freebsd](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_freebsd_arm) / [netbsd](https://bintray.com/artifact/download/nlamirault/oss/chione-0.1.0_netbsd_arm) ]


## Usage

* CLI help:

        $ chione help
        NAME:
        chione - CLI for skiing resorts informations

        USAGE:
            chione [global options] command [command options] [arguments...]

        VERSION:
            0.1.0

        COMMANDS:
            version
            resorts
            resort
            help, h  Shows a list of commands or help for one command

        GLOBAL OPTIONS:
            --debug        Enable debug mode
            --help, -h     show help
            --version, -v  print the version


* Show skiing resorts for a country :

        $ ./chione --debug resorts list --country france
        Resorts:
        - le-markstein [vosges]
        - bussang---larcenaire [vosges]
        - saint-leger-les-melezes [alpes-du-sud]
        - turini-camp-d'argent- [alpes-du-sud]
        [...]


* Display informations about a ski resort:

        $ /chione --debug resort describe --resort val-thorens --region alpes-du-nord
        +----------------------------+--------------------------------+
        | Status                     | Ouverte                        |
        +----------------------------+--------------------------------+
        | Enneigement sur les pistes | - bas: 55cm                    |
        |                            | - milieu:                      |
        |                            | - haut: 130cm                  |
        +----------------------------+--------------------------------+
        | Enneigement hors-pistes    | - bas:                         |
        |                            | - milieu:                      |
        |                            | - haut:                        |
        +----------------------------+--------------------------------+
        | Chute de neige             | - mardi: 0cm                   |
        |                            | - mercredi: 0cm                |
        |                            | - aujourd'hui: 10cm            |
        |                            | - vendredi: 0cm                |
        |                            | - samedi: 0cm                  |
        |                            | - dimanche: 2cm                |
        |                            | - lundi: 0cm                   |
        |                            | - mardi: 0cm                   |
        |                            | - mercredi: 0cm                |
        |                            | - jeudi: 0cm                   |
        +----------------------------+--------------------------------+
        | Domaine skiable            | - Vertes: 10/11                |
        |                            | - Bleues: 29/29                |
        |                            | - Rouges: 23/30                |
        |                            | - Noires: 5/9                  |
        +----------------------------+--------------------------------+


## Development

* Initialize environment

        $ make init

* Build tool :

        $ make build

* Launch unit tests :

        $ make test

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md).


## License

See [LICENSE](LICENSE) for the complete license.


## Changelog

A [changelog](ChangeLog.md) is available


## Contact

Nicolas Lamirault <nicolas.lamirault@gmail.com>

[badge-license]: https://img.shields.io/badge/license-Apache2-green.svg?style=flat
