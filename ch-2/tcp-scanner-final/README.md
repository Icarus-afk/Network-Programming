# Port Scanner

This is a simple port scanner written in Go. It scans the specified address for open ports and displays the results.

## Prerequisites

- Go (version 1.16 or later)

## Installation



1. **Install the `color` package**:

```sh
go get github.com/fatih/color
```

## Build

To build the program, run the following command in the project directory:

```sh
go build -o portscanner
```

## Run

To run the program, use the following command:

```sh
./portscanner -address your.target.address
```

Replace your.target.address with the actual address you want to scan.
Example

To scan the default address scanme.nmap.org, you can run:

```sh
./portscanner -address scanme.nmap.org
```
