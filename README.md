# pico-proxy

A tiny tiny cloud-native HTTP proxy.

## Usage

environment variables:

- _PORT_: set port to listen on
- _PATHS_: key-value pairs (separated by `:`) of endpoints to backend mappings, several can be added divided by `,`, eg: `google:https://google.com,reddit:https://reddit.com`
- _INSECURE_SKIP_VERIFY=yes_: whether to skip verifying TLS certs

## Development

The primary repository is [git.liebler.dev/dominik/pico-proxy](https://git.liebler.dev/dominik/pico-proxy); [github.com/domnikl/pico-proxy](https://github.com/domnikl/pico-proxy) is a push mirror of it. GitHub Actions publish the public multi-arch image `ghcr.io/domnikl/pico-proxy`, Forgejo Actions ([.forgejo/workflows](.forgejo/workflows)) a linux/amd64 image to `git.liebler.dev/dominik/pico-proxy`.
