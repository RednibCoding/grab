
# grab

`grab` is a command-line tool written in Go that loosely mimics the functionality of the Linux `grep` tool. It allows you to search for a string within files, with options to perform case-sensitive searches and to exclude subdirectories and hidden files.

## Features

- **Recursive Search**: By default, `grab` searches through all files and subdirectories starting from the current working directory.
- **Case Insensitivity**: The search is case-insensitive by default, but can be made case-sensitive with the `-c` flag.
- **Excluding Subdirectories and Hidden Files**: Use the `-d` flag to exclude subdirectories and `-h` to exclude hidden files from the search.
- **Binary File Detection**: `grab` automatically skips binary files during the search.

## Usage

```bash
grab [-d] [-h] [-c] [-s] <search-string>
```

### Flags

- `-d`: Do not search subdirectories.
- `-h`: Do not search hidden files.
- `-c`: Perform case-sensitive search.
- `-s`: Show directories where files have been skipped.

## Output Format

The results are displayed in the following format:

```
filename (occurrences):
  - filename:line:column
  - filename:line:column
  ...
```

Example:

```
grab hello
main.go (4):
  - main.go:25:2
  - main.go:46:2
  - main.go:58:2
  - main.go:86:2
```

## Installation

1. Clone this repository:
   ```bash
   git clone git@github.com:RednibCoding/grab.git
   ```

2. Build the binary:
   ```bash
   go build -ldflags="-w -s"
   ```

3. Run `grab`:
   ```bash
   ./grab "search text"
   ```

## Prebuild Binary
 - a prebuild binary for windows can be found in the `bin` directory

## License

This project is licensed under the MIT License.
