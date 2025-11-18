# nvmgr: The Neovim Configuration Manager

`nvmgr` is a command-line tool to help you manage multiple Neovim configurations.

## Installation

You can install `nvmgr` by running the installation script. It will install the binary to `~/.local/bin` and attempt to add this directory to your shell's `PATH` if it isn't already there.

```sh
curl -sSfL https://raw.githubusercontent.com/adamdlear/nvmgr/main/install.sh | sh
```

After installation, you may need to restart your shell or source your shell's configuration file (e.g., `source ~/.zshrc`).

### One-Time Setup

Once `nvmgr` is installed, run the one-time setup command:

```sh
nvmgr setup
```

This command prepares your system by:
1.  Identifying any existing Neovim configurations.
2.  Installing a wrapper script so that running `nvim` will automatically use the configuration managed by `nvmgr`.
3.  Ensuring `~/.local/bin` is in your `PATH`.

## Usage

`nvmgr` provides several commands to manage your configurations.

### `nvmgr list`

List all available Neovim configurations. The currently active configuration will be marked as `(Active)`.

```sh
nvmgr list
```

### `nvmgr new <name>`

Create a new, blank Neovim configuration.

```sh
nvmgr new my-awesome-config
```

You can also clone an existing configuration using the `--from` flag.

```sh
nvmgr new fork-of-awesome --from my-awesome-config
```

### `nvmgr install <distribution | git-url> [name]`

Install a configuration from a recognized distribution or a Git repository. If a name is not provided, it will be inferred from the source.

**Recognized Distributions:**
*   `astronvim`
*   `kickstart`
*   `lazyvim`
*   `nvchad`

**Examples:**

```sh
# Install from a recognized distribution
nvmgr install lazyvim

# Install from a git repository and give it a custom name
nvmgr install https://github.com/your-user/your-repo.git my-custom-config
```

### `nvmgr use <name>`

Switch the active Neovim configuration. The next time you run `nvim`, it will use the specified configuration.

```sh
nvmgr use my-awesome-config
```

### `nvmgr launch <name>`

Launch Neovim with a specific configuration for a single session, without changing the default active configuration.

```sh
nvmgr launch my-awesome-config
```

### `nvmgr restore`

Restore your system to its state before `nvmgr`. This command removes the `nvim` wrapper, allowing you to go back to a standard Neovim setup where it directly uses `~/.config/nvim`.

```sh
nvmgr restore
```

## Issues

If you run into any issues or have a feature request, please [create a GitHub issue](https://github.com/adamdlear/nvmgr/issues) and I'll look into it!
