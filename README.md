# `shutil`

- This is a personal utility CLI that aims to simplify / automate the following:
    - `cd` to an arbitrary directory
    - replace all white space with underscore for all filenames under a directory
    recursively
    - kitty integration,
        - launch a new kitty tab under an arbitrary directory
        - fuzzy find a kitty tab (using tab title) in the current active window
- This CLI is intended to use along with other CLI tools such as `fzf`, and with 
a few lines of script (either Bash or Power Shell) without using features or 
syntaxes that are specific to a scripting language.

## Motivation

- I want to reduce amount of mental load from maintaining a set of scripts that 
are written in two different scripting languages (bash and Power Shell), and 
from investing extra time to learn features specific to a scripting language (
mainly Power Shell).

## Usage (WIP)

- Here are some quick examples I currently have 

```bash
fd() {
    # fuzzy find a list of provided directorys (as the third arguments), and
    # `cd` to the selective directory
    local depth="${1:-2}"
    local worker="${2:-0}"
    local dir="${3:-.}"
    cd $(shutil --walker -walker-depth $depth --walker-worker $worker $dir)
}

fdb() {
    # fuzzy find a list of directorys specified in $HOME/.shutil.json, and
    # `cd` to the selective directory
    cd $(shutil --walker --walker-depth 3 --walker-worker 0)
    zle reset-prompt
}

fdb_kitty() {
    # fuzzy find a list of directorys specified in $HOME/.shutil.json, and
    # launch a new kitty tab the selective directory
    select=$(shutil --walker --walker-depth 3 --walker-worker 0)
    kitten @ launch --type=tab --cwd $select
}

kitty_tab_fzf() {
    # fuzzy find a kitty tab in the current active window
    shutil --kitty-fzf-tab
}
```
