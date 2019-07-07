# S3 CAS File Sync

A tool to sync tagged files to an S3 bucket using content-addressable storage ([CAS](https://en.wikipedia.org/wiki/Content-addressable_storage)).

## Building

Run `make build` to build the binary.

## Running

Run `dist/sync` to sync files to a bucket.

```
→ dist/sync
usage: sync --dir=<directory> --bucket=<bucket> --region=<region> --ext=<extensions> --silent
```

## Flags

### `--dir`

Specifies the directory of files to synchronise.

Filenames must consist of a space-separated list of tags. Example: `--dir=data`

```
<data>/house home garden 1.svg
<data>/car automobile drive 2.svg
<data>/man woman person 3.svg
```

Each file will be uploaded to the bucket under a content-addressable key, which will be linked against its respective tags in the final JSON output.

Note: numbers in filename are ignored.

### `--bucket`

Specifies the name of the S3 bucket in which to store objects.

Note: before running `sync`, ensure you have an AWS profile that allows read and write access to the bucket .

### `--region`

Sets the AWS region for the bucket. If this flag isn't provided, the default value is `eu-west-1`.

### `--ext`

A comma-separated list of extensions used to filter files in `--dir`. If specified, only files with a matching extension will be processed. Example: `--ext=png,svg`

### `--silent`

Silences all progress output of `sync`.

## Diffing

On completion, `sync` will write an `index.json` object to the bucket, which it will use on subsequent runs to avoid doing unnecessary work.

## Output

Once `sync` completes, it'll output a JSON document to `stdout` that tracks all CAS keys in the bucket against their respective tags.

Example:

```
→ sync --dir=icons --ext=svg --bucket=icons
uploading: f693691218...
uploading: db99ed52a2...
uploading: 27dd8ed44a...
processed 3 objects
uploaded 3 objects

[
  { "key": "f693691218...", "tags": "house home garden" },
  { "key": "db99ed52a2...", "tags": "car automobile drive" },
  { "key": "27dd8ed44a...", "tags": "man woman person" }
]
```

Note: redirect `stdout` to a file to capture only the JSON output. Example:

```
→ sync --dir=icons --ext=svg --bucket=icons > index.json
```
