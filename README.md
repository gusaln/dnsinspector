# dnsinspector

A tool that manually checks the DNS records of a domain.

This is a work in progress.

## Description

This tools does health checks of different DNS records by querying different authorities through the DNS protocol.

## Checks perform

(WIP)

## Features

- DNS protocol
  - [RFC1035](https://datatracker.ietf.org/doc/html/rfc1035)
    - [x] Queries and Answers headers
    - [x] Label support
    - [x] TCP length header
    - [x] Record support (A, NS, MX, CNAME, SOA, MB, MG, MR, NULL, WKS, PTR, HINFO, MINFO, TXT)
    - [ ] Truncated message support
    - [ ] Chaos type (what is this even used for?)

- Health checks (soon)

- CLI (soon)
  
- Web application (soon)
