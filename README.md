# dnsinspector

A tool that manually checks the DNS records of a domain.

This is a work in progress.

## Description

This tools does health checks of different DNS records by querying different authorities through the [DNS protocol - (RFC1035)](https://datatracker.ietf.org/doc/html/rfc1035).

## Checks perform

(WIP)

## Features

- DNS protocol
  - [x] Queries and Answers headers
  - [x] Label support
  - [x] TCP length header
  - [ ] Truncated message support
  - [ ] Records
    - [x] A
    - [x] NS
    - [x] MX
    - [x] CNAME
  - [ ] Chaos type (what is this even used for?)

- Health checks (soon)

- CLI (soon)
  
- Web application (soon)
