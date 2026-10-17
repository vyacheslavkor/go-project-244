# Hexlet Gen Diff

[![Actions Status](https://github.com/vyacheslavkor/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/vyacheslavkor/go-project-244/actions) [![CI Status](https://github.com/vyacheslavkor/go-project-244/actions/workflows/ci.yml/badge.svg)](https://github.com/vyacheslavkor/go-project-244/actions/workflows/ci.yml) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=vyacheslavkor_go-project-244&metric=coverage)](https://sonarcloud.io/summary/new_code?id=vyacheslavkor_go-project-244)

A command-line utility written in Go that compares two configuration files (JSON or YAML) and prints their difference.

## Demo

[![asciicast](https://asciinema.org/a/1262012.svg)](https://asciinema.org/a/1262012)

## Installation

You can clone the repository to any convenient local directory on your machine. There are no strict web-server path requirements.

1. Clone the repository
```bash
git clone https://github.com/vyacheslavkor/go-project-244.git
```

2. Navigate into the project directory
```bash
cd go-project-244
```

3. Build the executable (assuming make is installed)
```bash
make build
```

## Development

| Command | Description |
| :--- | :--- |
| `make build` | Build the `bin/gendiff` executable. |
| `make help` | Show the CLI help screen. |
| `make test` | Run the full test suite. |
| `make lint` | Run `golangci-lint` over the project. |
| `make lint-fix` | Run `golangci-lint` with safe autofixes. |
| `make coverage` | Run tests and write `coverage.out`. |

## Usage

The utility expects exactly two arguments: paths to the files being compared. By default, it prints the difference in `stylish` format.

```bash
./bin/gendiff [global options] <filepath1> <filepath2>
```

### Options / Flags

| Flag | Alias | Default | Description |
| :--- | :--- | :--- | :--- |
| `--format` | `-f` | `stylish` | Output format. Allowed values: `stylish`, `plain`, `json`. |
| `--help` | `-h` | `false` | Show the help screen and usage instructions. |

### Input Rules

- Exactly two positional arguments are required: paths to existing non-empty regular files.
- The root value of each file must be a JSON object or a YAML mapping with string keys. Root-level arrays/sequences, scalars, and null are rejected, and so is a YAML mapping whose keys are not strings (`1: a`).
- Empty files are rejected.
- Supported input formats: JSON (`.json`) and YAML (`.yml`, `.yaml`).
- JSON and YAML may be mixed; each file is parsed by its own extension.
- The parser is chosen by file extension, not by inspecting the file.
- Values are compared in a normalized form, so the same value written differently matches: numbers become `float64` (`1` equals `1.0`), and a timestamp — including a string that parses as one — becomes the shortest string that keeps every part of the value. One moment therefore matches however it is spelled and whichever format it comes from; values that differ as moments (fractional seconds, zone offsets) stay different. Comparison carries `float64` precision, so integers that differ only beyond it (above 2^53) become the same value and can render differently from the text in the file.
- Usage errors (wrong args/flags/input contract) print help on stdout and a reason on stderr. Operational errors (missing/unreadable files, malformed content) print only the reason on stderr.

### Output Formats

| Format    | Description                                                                                                                                                                                                                                                                                                                                   |
|:----------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `stylish` | Nested tree with `+` / `-` markers (default).                                                                                                                                                                                                                                                                                                 |
| `plain`   | Line-oriented human-readable messages (skips unchanged properties). Complex values are shown as `[complex value]`.                                                                                                                                                                                                                            |
| `json`    | Compact single-line JSON tree. Envelope: `{"key":"","status":"root","children":[...]}` (`children` omitted when there are no nodes). Node fields: `key`, `status`, `old_value?`, `value?`, `children?`. Node statuses: `root` (document root, empty key), `added`, `removed`, `updated`, `nested`, `unchanged`. Unchanged nodes are included. |

Identical files are not a special case: `stylish` prints the whole tree with every property unchanged, `json` prints the full `children` list of `unchanged` nodes, and `plain` prints nothing, because it skips unchanged properties.

A document with no properties (both roots are `{}`) gives an empty tree: `stylish` prints a pretty-printed empty object `{}`, `json` prints `{"key":"","status":"root"}` with `children` omitted, and `plain` prints nothing. An empty *file* is a different case and is rejected (see Input Rules).

Numbers are printed in positional notation in every format: a file that says `1000000` shows `1000000`, never `1e+06`. A value looks the same whichever `--format` is selected.

### Examples

All three examples compare the same pair of files, `file1.json` and `file2.json`,
taken from `internal/cli/testdata/fixture/`.

**1. Stylish format (default):**
```bash
$ ./bin/gendiff file1.json file2.json
{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
      - setting3: true
      + setting3: null
      + setting4: blah blah
      + setting5: {
            key5: value5
        }
        setting6: {
            doge: {
              - wow:
              + wow: so much
            }
            key: value
          + ops: vops
        }
      - timeout: 1000000
      + timeout: 123456789
    }
    group1: {
      - baz: bas
      + baz: bars
        foo: bar
      - nest: {
            key: value
        }
      + nest: str
    }
  - group2: {
        abc: 12345
        deep: {
            id: 45
        }
    }
  + group3: {
        deep: {
            id: {
                number: 45
            }
        }
        fee: 100500
        limit: 1000000000000000000000
    }
}
```

**2. Plain format:**
```bash
$ ./bin/gendiff -f plain file1.json file2.json
Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'common.setting4' was added with value: 'blah blah'
Property 'common.setting5' was added with value: [complex value]
Property 'common.setting6.doge.wow' was updated. From '' to 'so much'
Property 'common.setting6.ops' was added with value: 'vops'
Property 'common.timeout' was updated. From 1000000 to 123456789
Property 'group1.baz' was updated. From 'bas' to 'bars'
Property 'group1.nest' was updated. From [complex value] to 'str'
Property 'group2' was removed
Property 'group3' was added with value: [complex value]
```

**3. JSON format:**
```bash
$ ./bin/gendiff --format=json file1.json file2.json
{"key":"","status":"root","children":[{"key":"common","status":"nested","children":[{"key":"follow","status":"added","value":false},{"key":"setting1","status":"unchanged","value":"Value 1"},{"key":"setting2","status":"removed","old_value":200},{"key":"setting3","status":"updated","old_value":true,"value":null},{"key":"setting4","status":"added","value":"blah blah"},{"key":"setting5","status":"added","value":{"key5":"value5"}},{"key":"setting6","status":"nested","children":[{"key":"doge","status":"nested","children":[{"key":"wow","status":"updated","old_value":"","value":"so much"}]},{"key":"key","status":"unchanged","value":"value"},{"key":"ops","status":"added","value":"vops"}]},{"key":"timeout","status":"updated","old_value":1000000,"value":123456789}]},{"key":"group1","status":"nested","children":[{"key":"baz","status":"updated","old_value":"bas","value":"bars"},{"key":"foo","status":"unchanged","value":"bar"},{"key":"nest","status":"updated","old_value":{"key":"value"},"value":"str"}]},{"key":"group2","status":"removed","old_value":{"abc":12345,"deep":{"id":45}}},{"key":"group3","status":"added","value":{"deep":{"id":{"number":45}},"fee":100500,"limit":1000000000000000000000}}]}
```
