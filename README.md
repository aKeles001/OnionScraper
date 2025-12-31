# OnionScraper


## Overview
A basic web scrapper designed to scrap Tor based sites.


## Installation
```bash
git clone https://github.com/aKeles001/OnionScraper.git
```


## Usage

Targets must be provided with a targets.yaml file on data directory.

```golang
go run main.go
```


## Features
- Scrap all the websites provided on targets.yaml  file.
- Get the HTML contents of the targets.
- Get the screenshot of the targets.


## Requirements
- Tor daemon must be running.