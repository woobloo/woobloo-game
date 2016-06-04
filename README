# Woobloo Game

The source code of this repository compiles to three executable programs:

1. woobloo-orchestration
2. woobloo-frontman
3. woobloo-game

## Woobloo Orchestration

The source for _woobloo-orchestration_ can be found in [src/c](src/c). _woobloo-orchestration_ is written entirely in C and serves one purpose: it hooks together _woobloo-frontman_ and _woobloo-game_. To be more blunt, _woobloo-orchestration_ simply pipes the stdout of _woobloo-frontman_ to the stdin of _woobloo-game_, and vice versa. This is shown in the diagram below.

[!docs/images/woobloo-orchestration.png](docs/images/woobloo-orchestration.png)

## Woobloo Frontman

The source for _woobloo-frontman_ can be found in [src/go](src/go). _woobloo-frontman_ is written entirely in Go. It does two things:

1. Write all data read from websockets to stdout (which is then piped to _woobloo-game_)
2. Write some data read from stdin (which is the stdout of _woobloo-game_) to some websockets

## Woobloo Game

The source for _woobloo-game_ can be found in [src/haskell](src/haskell). This holds the logic and state of each game instance. It reads from stdin (which is the stdout of _woobloo-frontman_) and writes to stdout (which is the stdin of _woobloo-frontman_).