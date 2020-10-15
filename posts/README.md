# Posts Service

This is the Posts service

Generated with

```
micro new posts
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- Alias: posts

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./posts
```

Build a docker image
```
make docker
```