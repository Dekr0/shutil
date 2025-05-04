# `shutil`

- This is a personal utility CLI that aims to simplify / automate the following:
    - `cd` to an arbitrary directory
    - replace all white space with underscore for all filenames under a directory
    recursively
    - integration with terminals (Currently, it supports Kitty and Wezterm) that
    provide native multiplexing:
        - launch a new session in an arbitrary directory, this aribtary directory
        is selected using fuzzy finding
        - fuzzy find an existed session
- This CLI is intended to use along with other CLI tools such as `fzf`, and with 
a few lines of script (either Bash or Power Shell) without using features or 
syntaxes that are specific to a scripting language.

## Motivation

- I want to reduce amount of mental load from maintaining a set of scripts that 
are written in two different scripting languages (bash and Power Shell), and 
from investing extra time to learn features specific to a scripting language (
mainly Power Shell).

## Installation

- Make sure `FzF` is installed.
- Clone this project, and use the build script
```sh
git clone https://github.com/Dekr0/shutil && cd shutil
python build.py
```
- Include `$GOPATH/bin` (Usually, it's under `$HOME/go/bin`) in `$PATH`.

## Usage

### `--walker`

- This option will accept a list of directories separated by space. Then, it
will walk each of this directory recursively in parallel. During this process,
it will collect all directories it finds, and pipe it to FzF.
- If no directory is provided, it will use the bookmark directory stored in
`$HOME/.shutil.json`.
- There are two additional parameters you can use,
    - `--walker_depth` specifies the depth of a walk
    - `--walker_worker` specified the maximum go routine

### `--[kitty|wezterm]_activate_tab`

- Fuzzy find a kitty tab in the current active window

### `--wezterm_new_tab`

- This option do pretty much the same thing as `walker` but it will start a new 
tab in `wezterm` using the selected path.

### `--[kitty|wezterm]_new_sessions`

- Select a session profile stored in `$HOME/.shutil.json` and start a set of 
new tabs using this profile in the active window.

### `--[kitty|wezterm]_create_session_profile`

- Store the tabs information in the active window as a session profile.

## Example 

### Kitty Configuration

```conf
map ctrl+alt+f launch --type overlay shutil --kitty_activate_tab
map ctrl+alt+n launch --type overlay sh -c "kitten @ launch --type=tab --cwd $(shutil --walker --walker_depth 3 --walker_worker 0)"
map ctrl+shift+r launch --type overlay shutil --kitty_new_sessions
map ctrl+shift+s launch --type overlay shutil --kitty_create_session_profile
```

### Wezterm Configuration

```lua
config.keys = {
    {
        key = 'f',
        mods = 'CTRL|ALT',
        action = wezterm.action.SpawnCommandInNewTab {
            args = { 'shutil', '--wezterm_activate_tab' },
        }
    },
    {
        key = 'n',
        mods = 'CTRL|ALT',
        action = wezterm.action.SpawnCommandInNewTab {
            args = { 'shutil', '--wezterm_new_tab' }
        }
    },
    {
        key = 's',
        mods = 'CTRL|ALT',
        action = wezterm.action.SpawnCommandInNewTab {
            args = { 'shutil', '--wezterm_new_sessions' }
        }
    },
    {
        key = 's',
        mods = 'CTRL|ALT|SHIFT',
        action = wezterm.action.SpawnCommandInNewTab {
            args = { 'shutil', '--wezterm_create_session_profile' }
        }
    },
}
```

### ZSH (Alias and Binding)

```sh
# Find directory (general cases)
fd() {
    local depth="${1:-2}"
    local worker="${2:-0}"
    local dir="${3:-.}"
    cd $(shutil --walker -walker_depth $depth --walker_worker $worker $dir)
}

# Find directory (using bookmark)
fdb() {
    cd $(shutil --walker --walker_depth 3 --walker_worker 0)
    zle reset-prompt
}

zle     -N            fdb
bindkey -M emacs '^J' fdb 
bindkey -M vicmd '^J' fdb 
bindkey -M viins '^J' fdb 
```
