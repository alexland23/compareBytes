# compareBytes

A small CLI tool that compares two files byte-by-byte and reports where they
differ.

It streams both files concurrently, comparing one byte at a time, and logs
the position and hex values of every mismatch. It also reports if one file
is longer than the other.

## Requirements

- Go 1.26+

## Build

```sh
go build -o compareBytes .
```

## Usage

```sh
./compareBytes -f1 <path> -f2 <path> [-p] [-o <path>]
```

### Flags

| Flag  | Description                                              |
| ----- | ---------------------------------------------------------|
| `-f1` | Path of the first file (required)                        |
| `-f2` | Path of the second file (required)                       |
| `-p`  | Print each mismatch to stdout as it's found (default: false) |
| `-o`  | Path to write mismatch results to a file (optional)       |

### Example

```sh
./compareBytes -f1 testData/file1.txt -f2 testData/file2.txt -p
```

Each mismatch is logged in the form:

```
Pos: 12, f1: 41, f2: 42
```

Where `Pos` is the 1-indexed byte position and `f1`/`f2` are the hex values
of the differing bytes.

When the files are done being compared, a summary is logged with the number
of bytes read from each file, how many bytes (if any) remain unread in the
longer file, and the total number of mismatches found.

## Linting

This project uses [golangci-lint](https://golangci-lint.run/) (see
`.golangci.yml`). Run it with:

```sh
golangci-lint run
```
