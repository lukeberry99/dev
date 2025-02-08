# Dev

An over-the-top dotfiles/mac configuration setup that I created because I just don't care to use Ansible.

This probably isn't useful to you, but the only prerequisite is that you are using a Mac. I will update this to work on Linux when I have an immediate need for it (soon), but it will likely still use Homebrew as a package manager, as then I only need to change how Brew is installed.

## What do?

This will install homebrew, and then install and configure the tools I use day-to-day, including building neovim (btw) from source.

## How do?

There's the `install` script, that runs _other_ scripts that live in the `runs/` directory. It's all bash, it's all straightforward enough.

Then there's the `configure` script, this is _kinda_ destructive as it will first delete your config folders - so check that. You can do `./configure --dry` to see what it will delete. It will then copy my config into the desired folders.

## Usage:

```bash
git clone https://github.com/lukeberry99/dev
cd dev
./install # You can run ./install --help for some additional params
./configure # You can run ./configure --help for some additional params
```
