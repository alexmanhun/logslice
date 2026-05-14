# logslice

Stream and filter structured JSON logs from multiple files with a simple query syntax.

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git && cd logslice && go build ./...
```

## Usage

```bash
logslice [options] <file1> [file2 ...]
```

Filter logs by field value:

```bash
logslice --query 'level=error' app.log worker.log
```

Stream logs in real time and match multiple conditions:

```bash
logslice --follow --query 'level=error,service=api' /var/log/app/*.log
```

Output only specific fields:

```bash
logslice --query 'status=500' --fields time,message,trace app.log
```

### Query Syntax

| Expression | Description |
|---|---|
| `key=value` | Exact field match |
| `key!=value` | Field not equal |
| `key~=pattern` | Regex match |
| `key>value` | Numeric comparison |

Multiple conditions are separated by commas and evaluated as logical AND.

## Options

| Flag | Description |
|---|---|
| `--query`, `-q` | Filter expression |
| `--fields`, `-f` | Comma-separated fields to output |
| `--follow`, `-F` | Stream file(s) in real time |
| `--pretty` | Pretty-print JSON output |

## License

MIT © yourname