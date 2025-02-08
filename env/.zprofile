eval "$(/opt/homebrew/bin/brew shellenv)"
eval "$(starship init zsh)"
eval "$(fnm env --use-on-cd --shell zsh)"

addToPath() {
    if [[ "$PATH" != *"$1"* ]]; then
        export PATH=$PATH:$1
    fi
}

addToPath $HOME/.local/scripts
addToPath $HOME/.local/bin
