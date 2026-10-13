# fetchswap

Fetch a file over HTTP and atomically swap it into place.

Fetching a file over HTTP and atomically replacing a local copy is a common
job, and it is usually done by writing the same small shell script again
and copying it to every host that needs it. `fetchswap` is that script
written once, as a single static binary.

The file is never left half-written: either the fetch succeeds and the new
file is renamed into place, or the existing file stays as it was.

## Usage

```
fetchswap -url URL -dest DEST -mode MODE -timeout DURATION
```

Each of these flags is required; a missing one stops the run.

- `-url`: the URL to GET, over HTTP or HTTPS. Redirects are followed, and
  `HTTP_PROXY`, `HTTPS_PROXY` and `NO_PROXY` are honored.
- `-dest`: the file to replace. Its directory must already exist and be
  writable.
- `-mode`: permission bits of DEST in octal, such as `0644`.
- `-timeout`: limit for the whole request including the body, such as
  `60s`.

Any non-2xx status or a body cut short is an error. Every successful run
replaces DEST, even when the content did not change.

Exit status:

- 0: DEST was replaced.
- 1: error; DEST is left as it was.
- 2: usage error.

For example, reload a service only after a successful fetch:

```
fetchswap -url https://example.com/list.txt \
  -dest /etc/example/list.txt -mode 0644 -timeout 60s \
  && systemctl reload example.service
```

## Notes

- The replaced file is owned by the user running `fetchswap`.
- If DEST is a symlink, the symlink itself is replaced by a regular file.
  Point `-dest` at the target to keep the link.

## License

This project is licensed under the [MIT License](./LICENSE).
