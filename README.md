# pico-proxy

A tiny tiny cloud-native HTTP proxy.

## Usage

environment variables:

- _PORT_: set port to listen on
- _PATHS_: key-value pairs (separated by `:`) of endpoints to backend mappings, several can be added divided by `,`, eg: `google:https://google.com,reddit:https://reddit.com`
