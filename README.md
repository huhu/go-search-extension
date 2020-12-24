# Go Search Extension

<img align="right" width="280" src="gopher.svg">

### The ultimate search extension for Golang.

![Chrome Web Store](https://img.shields.io/chrome-web-store/v/epanejkfcekejmmfbcpbcbigfpefbnlb.svg)
![Mozilla Add-on](https://img.shields.io/amo/v/go-search-extension?color=%2320123A)
![Microsoft Edge](https://img.shields.io/badge/microsoft--edge-0.2.0-1D4F8C)
[![license-mit](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/huhu/go-search-extension/blob/master/LICENSE)
[![Discord](https://img.shields.io/discord/711895914494558250?label=chat&logo=discord)](https://discord.gg/xucZNVd)

[https://go.extension.sh/](https://go.extension.sh/)

## Installation

- [Chrome Web Store](https://chrome.google.com/webstore/detail/golang-search/epanejkfcekejmmfbcpbcbigfpefbnlb)

- [Firefox](https://addons.mozilla.org/en-US/firefox/addon/go-search-extension/)

- [Microsoft Edge](https://microsoftedge.microsoft.com/addons/detail/ebibclchdmagkhopidkjckjkbhghfehh)


## Features

- Search standard library docs
- Search third party packages on pkg.go.dev
- Search awesome golang resources
- Builtin commands (`book`, `conf`, `meetup`, `social`, and `history`)

## How to use it
   
Input keyword **go** in the address bar, press `Space` to activate the search bar. Then enter any word 
you want to search, the extension will response the related search results instantly.

## Contribution

[jsonnet](https://jsonnet.org/) is required before getting started. To install `jsonnet`, 
please check `jsonnet`'s [README](https://github.com/google/jsonnet#packages). 
For Linux users, the `snap` is a good choice to [install jsonnet](https://snapcraft.io/install/jsonnet/ubuntu).

```bash
$ git clone --recursive https://github.com/huhu/go-search-extension
Cloning into 'go-search-extension'...
$ cd go-search-extension

$ make chrome # For Chrome version

$ make firefox # For Firefox version

$ make edge # For Edge version
```

## Get involved

- You can contact us on Discord Channel: https://discord.gg/xucZNVd
- Or by adding the Wechat ID: `huhu_io`, we'll invite you to our Wechat group.
