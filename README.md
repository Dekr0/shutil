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

## Usage

### `--walker`

- This option will accept a list of directories separated by space. Then, it
will walk each of this directory recursively in parallel. During this process,
it will collect all directories it finds, and pipe it to FzF.
- If no directory is provided, it will use the bookmark directory stored in
`$HOME/.shutil.json`.
- There are two additional parameters you can use,
    - `--walker-depth` specifies the depth of a walk
    - `--walker-worker` specified the maximum go routine

### `--kitty_fzf_tab`

- Fuzzy find a kitty tab in the current active window

### `--kitty_start_session`

- Select a session profile stored in `$HOME/.shutil.json` and start a set of 
new tabs using this profile in the active window.

### `--kitty_store_session`

- Store the tabs information in the active window as a session profile.

### `--wezterm_fzf_tab`

- Fuzzy find a wezterm tab
