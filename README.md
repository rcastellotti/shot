# shot

`shot` := a safer alternative to `curl | sh` installs.

replace `sh` with `shot`, and `shot` will:

- read the script and extract any `SHOT_*` variables
- store those values in a local SQLite database
- pass the script to a real shell for execution
- verify the script's integrity before running it (TODO)

try shot now, run `curl https://todo | sh` [one ~~more~~ last time](https://youtu.be/FGBhQbmPwH8)

now go install something, may i interest you in `mise`?

```sh
curl https://mise.run | shot
```

# shot variables

```sh
#!/usr/bin/env sh
# /// shot
SHOT_AUTHOR="rcastellotti"
SHOT_NAME="shot"
SHOT_VERSION="1.27.1"
SHOT_WEBSITE="https://rcastellotti.dev"
# todo: add root/prefix to track where files are installed
# ///
```

- TODO: make `shot` verify script integrity (needs a server)
- TODO: release binaries, start with `shot_darwin_arm64`
- TODO: add support for SHOT_SHELL to explicitly set the shell to use
- TODO: interactive mode -> see all variables and inspect script in editor
