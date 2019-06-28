# S3 CAS File Sync

A tool to sync tagged files to an S3 bucket using content-addressable storage ([CAS](https://en.wikipedia.org/wiki/Content-addressable_storage)).

## Building

Run `make build` to build the binary.

## Running

Run `dist/sync` to sync files to a bucket.

```
→ ./dist/sync
S3 Sync Tool

Usage:
  sync --dir=<directory> --bucket=<bucket>
```

## Input

Files in the specified directory should be named as a space-separated list of tags, which will be used in the final output file (see below). Example:

```
<directory>/house home garden 1.svg
<directory>/car automobile drive 2.svg
<directory>/man woman person 3.svg
```

Each file will be uploaded to the bucket as a content-addressable key, which will be linked against the respective tags in the JSON output file. Numbers in file names are ignored.

## Output

On completion, the tool will write a JSON file that links all keys currently in the S3 bucket to the associated tags. Example:

```json
[
  { "key": "5d41402abc...", "tags": "house home garden" },
  { "key": "ae2b1fca51...", "tags": "car automobile drive" },
  { "key": "6057f13c49...", "tags": "man woman person" }
]
```

## Diffing

The tool also uploads the JSON document to the bucket, using it on subsequent runs to avoid doing unnecessary work.
