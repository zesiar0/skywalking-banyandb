# BanyanDB

[![Continuous Integration](https://github.com/apache/skywalking-banyandb/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/apache/skywalking-banyandb/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/skywalking-banyandb)](https://goreportcard.com/report/github.com/apache/skywalking-banyandb)
[![GitHub release](https://img.shields.io/github/tag/apache/skywalking-banyandb.svg?label=release)](https://github.com/apache/skywalking-banyandb/releases)
[![GitHub release date](https://img.shields.io/github/release-date/apache/skywalking-banyandb.svg)](https://github.com/apache/skywalking-banyandb/releases)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/apache/skywalking-banyandb)

![](./assets/banyandb_small.jpg)

BanyanDB, as an observability database, aims to ingest, analyze and store Metrics, Tracing and Logging data.
It's designed to handle observability data generated by observability platform and APM system, like [Apache SkyWalking](https://github.com/apache/skywalking) etc.

## Introduction

BanyanDB, as an observability database, aims to ingest, analyze and store Metrics, Tracing, and Logging data. It's designed to handle observability data generated by Apache SkyWalking. Before BanyanDB emerges, the Databases that SkyWalking adopted are not ideal for the APM data model, especially for saving tracing and logging data. Consequently, There’s room to improve the performance and resource usage based on the nature of SkyWalking data patterns.

The database research community usually uses [RUM Conjecture](http://daslab.seas.harvard.edu/rum-conjecture/) to describe how a database access data. BanyanDB combines several access methods to build a comprehensive APM database to balance read cost, update cost, and memory overhead.

## Contact us

* Submit [an issue](https://github.com/apache/skywalking/issues/new) by selecting the [BanyanDB](https://github.com/apache/skywalking/issues?q=is%3Aopen+is%3Aissue+label%3Adatabase) component.
* Mail list: **dev@skywalking.apache.org**. Mail to dev-subscribe@skywalking.apache.org, follow the reply to subscribe the mail list.
* Send `Request to join SkyWalking slack` mail to the mail list(`dev@skywalking.apache.org`), we will invite you in.
* For Chinese speaker, send `[CN] Request to join SkyWalking slack` mail to the mail list(`dev@skywalking.apache.org`), we will invite you in.
* X (Twitter): [@BanyanDB](https://twitter.com/BanyanDB) and [@ASFSkyWalking](https://twitter.com/ASFSkyWalking)

## Documentation

- [Dev version doc](https://skywalking.apache.org/docs/skywalking-banyandb/next/readme/)
- [Latest release doc](https://skywalking.apache.org/docs/skywalking-banyandb/latest/readme/)
- [Java Client SDK doc](https://skywalking.apache.org/docs/#BanyanDBJavaClient)


## Contributing

For developers who want to contribute to this project, see the [Contribution Guide](CONTRIBUTING.md).

## License

[Apache 2.0 License.](LICENSE)
