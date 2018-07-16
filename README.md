# Ember
## About
Ember is an included package in the [Bult Client](https://github.com/BultApp/Client). It helps manage the downloading and installation of addons from 3rd party sites as well as the Bult Addon Marketplace.

## Installation
No installation needed. Bult downloads the latest released version of this software during the install process. You have no need to run this by itself, however if you do...it opens an RPC server which you can connect to.

## Usage
Ember opens up an small web server on port `5650`. It's job is simple. Handle installing/unzipping the files.
```json
{
    "url": "https://download.the/file.zip",
    "filename": "somethingtonameitas",
    "extractTo": "path/to/extract/to/"
}
```
