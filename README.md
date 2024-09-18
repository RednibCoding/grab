
# grepl

`grepl` is a command-line tool written in Go that loosely mimics the functionality of the Linux `grep` tool. It allows you to search for a string within files, with options to perform case-sensitive searches and to exclude subdirectories and hidden files.

## Features

- **Recursive Search**: By default, `grepl` searches through all files and subdirectories starting from the current working directory.
- **Case Insensitivity**: The search is case-insensitive by default, but can be made case-sensitive with the `-c` flag.
- **Excluding Subdirectories and Hidden Files**: Use the `-e` flag to exclude subdirectories and hidden files from the search.
- **Binary File Detection**: `grepl` automatically skips binary files during the search.

## Usage

```bash
grepl [-e] [-c] [-s] <search-string>
```

### Flags

- `-e`: Do not search subdirectories and hidden files.
- `-c`: Perform case-sensitive search.
- `-s`: Show directories where files have been skipped.

### Example

1. **Basic Search (case-insensitive, including subdirectories and hidden files)**:
   ```bash
   grepl "search text"
   ```

2. **Case-sensitive Search**:
   ```bash
   grepl -c "search text"
   ```

3. **Exclude Subdirectories and Hidden Files**:
   ```bash
   grepl -e "search text"
   ```
4. **Show skipped files and directories**:
   ```bash
   grepl -s "search text"
   ```

4. **Case-sensitive Search while Excluding Subdirectories and Hidden Files and showing skipped files and directories**:
   ```bash
   grepl -e -c -s "search text"
   ```

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
grepl hello
main.go (4):
  - main.go:25:2
  - main.go:46:2
  - main.go:58:2
  - main.go:86:2
```

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/grepl.git
   ```

2. Build the binary:
   ```bash
   go build -ldflags="-w -s"
   ```

3. Run `grepl`:
   ```bash
   ./grepl "search text"
   ```

## License

This project is licensed under the MIT License.
