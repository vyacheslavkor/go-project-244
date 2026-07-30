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
- The root value of each file must be a JSON object or a YAML mapping. Root-level arrays/sequences, scalars, and null are rejected.
- Empty files are rejected.
- Supported input formats: JSON (`.json`) and YAML (`.yml`, `.yaml`).
- Both files must use compatible formats: JSON with JSON, or YAML with YAML. Mixing `.yml` and `.yaml` is allowed; mixing JSON and YAML is not.
- The parser is chosen by file extension; content is not sniffed.
- Usage errors (wrong args/flags/input contract) print help on stdout and a reason on stderr. Operational errors (missing/unreadable files, malformed content) print only the reason on stderr.

### Output Formats

| Format | Description |
| :--- | :--- |
| `stylish` | Nested tree with `+` / `-` markers (default). No changes: pretty-printed empty object `{}`. |
| `plain` | Line-oriented human-readable messages (skips unchanged properties). No changes: empty output (nothing is printed). Complex values are shown as `[complex value]`. |
| `json` | Compact single-line JSON tree of changes. Envelope: `{"key":"","status":"root","children":[...]}` (`children` omitted when empty). Node fields: `key`, `status`, `old_value?`, `value?`, `children?`. Node statuses: `added`, `removed`, `updated`, `nested`, `unchanged`. Status `root` is JSON-wire only. Unchanged nodes are included. |

### Examples

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
Property 'group1.baz' was updated. From 'bas' to 'bars'
Property 'group1.nest' was updated. From [complex value] to 'str'
Property 'group2' was removed
Property 'group3' was added with value: [complex value]
```

**3. JSON format:**
```bash
$ ./bin/gendiff --format=json before.json after.json
{"key":"","status":"root","children":[{"key":"follow","status":"removed","old_value":false},{"key":"host","status":"unchanged","value":"hexlet.io"},{"key":"proxy","status":"removed","old_value":"123.234.53.22"},{"key":"timeout","status":"updated","old_value":50,"value":20},{"key":"verbose","status":"added","value":true}]}
```
