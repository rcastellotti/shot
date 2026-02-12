# shot

<!--TODO: fake vocabulary entry-->

`shot` is a tool to make `curl | sh` installs (a bit) safer.

replace `sh` with `shot` and:

- shot reads the file for `SHOT_*` variables and prints their content
- execs the shell and passes input to it
- shot verifies install script's integrity

try shot now, run `curl https://todo | sh` [one ~~more~~ last time](https://youtu.be/FGBhQbmPwH8)

now go install something, may i interest you in `mise`?

```sh
curl https://mise.run | shot
```

TODO: make `shot` verify script integrity (needs a server)

# shot variables

```toml
SHOT_NAME="my-tool"
SHOT_VERSION="1.3.0"
SHOT_AUTHOR="alice"
SHOT_SOURCE="https://github.com/alice/my-tool"
SHOT_ROOT="/opt/my-tool"
```

```sh
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
curl -fsSL https://deno.land/install.sh | sh
curl -fsSL https://bun.sh/install | bash
curl -fsSL https://claude.ai/install.sh | bash
curl -LsSf https://astral.sh/uv/install.sh | sh
curl -sS https://starship.rs/install.sh | sh
```
