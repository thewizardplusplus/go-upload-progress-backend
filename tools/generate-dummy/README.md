# generate-dummy

The utility for generating a dummy file of a specified size, filled with random bytes.

## Usage

```
$ generate_dummy.bash -h | --help
$ generate_dummy.bash [options]
```

Options:

- `-h`, `--help` &mdash; show the help message and exit;
- `-s SIZE`, `--size SIZE` &mdash; a desired size of a generated file (should be in format `^([[:digit:]]+(\.[[:digit:]]+)?)(([GMK]i)?B)$`, e.g.: `5B`, `12.3KiB`, `2.35MiB`, `42.6GiB`);
- `-n TEMPLATE`, `--name TEMPLATE` &mdash; a template for a name of a generated file (may contain placeholder `${SIZE}`, which will be replaced by a specified size; default: `dummy_${SIZE}.bin`).
