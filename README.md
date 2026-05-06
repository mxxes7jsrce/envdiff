# envdiff

A CLI tool to compare `.env` files across environments and surface missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <file1> <file2>
```

**Compare two `.env` files:**

```bash
envdiff .env .env.production
```

**Example output:**

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED values:
  - APP_ENV: "development" vs "production"
  - LOG_LEVEL: "debug" vs "info"

✔ All other keys match.
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--keys-only` | Compare keys only, ignore values |
| `--quiet` | Exit with non-zero status if differences found (useful in CI) |
| `--json` | Output results as JSON |

---

## CI Usage

```bash
# Fail the pipeline if envs are out of sync
envdiff --quiet .env.example .env.production
```

---

## License

MIT © [yourusername](https://github.com/yourusername)