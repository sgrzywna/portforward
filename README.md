# portforward

Simple and small port forwarder with SOCKS5 proxy support.

This is a slighty modified version of [Simple connection fowarder in golang](http://blog.evilissimo.net/simple-port-fowarder-in-golang).

This simple application is an excellent example of the power and simplicity of go - fully functional port forwarder in just a few lines of the source code.

## Installation

```bash
go get github.com/sgrzywna/portforward
```

## Usage examples

### Direct port forwarding

```bash
portforward localhost:2222 192.168.0.1:22
```

### Port forwarding using SOCKS5 proxy

```bash
portforward -proxy localhost:1080 localhost:2222 192.168.0.1:22
```
