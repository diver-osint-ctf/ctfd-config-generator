# ctfd-config-generator for DIVER OSINT CTF
[日本語版はこちら](./README.md)

ctfd-config-generator is a tool to support making a CTF with CTFd. This tool provides the following features:

- [x] generates a config file for [ctfcli](https://github.com/CTFd/ctfcli) which manages CTFd as a CLI
- [x] generates files and directories to be used for making CTF challenges

## How to make a new challenge
[Go](https://go.dev/doc/install) and Make are required.

```bash
git init
git submodule add https://github.com/diver-osint-ctf/ctfd-config-generator
echo "include ctfd-config-generator/Makefile" > Makefile
make gen
```

After the generation, the following structure is made.

```bash
genre
└── challenge-name
    ├── README.md       # Write overview of the challenge. (optional)
    ├── build           # Locate files to be run in a challenge server. Their files are not published.(optional)
    ├── challenge.yml   # Update a config for CTFd.
    ├── flag.txt        # Update a flag if you need. If you want to use a regular expression or multiple flags, write them in this file with a newline.
    ├── public          # Locate files to be published. Their files should be set in challge.yml.(optional)
    ├── solver          # Locate solvers. (optional)
    └── writeup
        └── README.md   # Locate a writeup
```

## Contribution
If there are bugs and requests, please make an issue or a pull request. Thanks.
