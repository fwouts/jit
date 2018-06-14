# jit
A command-line tool to simplify Git workflows for Jira users

## Introduction

`jit` aims to make your life easier if you find yourself constantly switching between your terminal and your browser to
look at Jira tickets.

In particular, use `jit` directly from your shell to:
- see Jira tickets assigned to you
- create a branch named after a Jira ticket and update the ticket's status
- create a new Jira ticket and start work on a corresponding branch

## Installation

- Download `jit` from https://github.com/zenclabs/jit/releases.
- Run `chmod +x jit`.
- Call `jit` from within any Git repository.

## Updating

`jit` checks for updates regularly. You will know when a new version is available.

## Usage

```sh
# Print the list of Jira tickets assigned to you, pick a branch.
jit

# Would you like more commands? Please add feature requests in the Issues section.
```

# Configuration

`jit` will ask for your configuration details the first time you use it. You can adjust your settings anytime by editing
`.jit/config.yaml` at the root of the Git repository.
