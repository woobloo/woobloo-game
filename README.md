# Woobloo Game

The source code of this repository compiles to three executable programs:

1. woobloo-orchestration
2. woobloo-frontman
3. woobloo-game

## Woobloo Orchestration

The source for __woobloo-orchestration__ can be found in [src/c](src/c). __woobloo-orchestration__ is written entirely in C and serves one purpose: it hooks together __woobloo-frontman__ and __woobloo-game__. To be more blunt, __woobloo-orchestration__ simply pipes the stdout of __woobloo-frontman__ to the stdin of __woobloo-game__, and vice versa. This is shown in the diagram below.

[!docs/images/woobloo-orchestration.png](docs/images/woobloo-orchestration.png)

## Woobloo Frontman

The source for __woobloo-frontman__ can be found in [src/go](src/go). __woobloo-frontman__ is written entirely in Go. It does two things:

1. Write all data read from websockets to stdout (which is then piped to __woobloo-game__)
2. Write some data read from stdin (which is the stdout of __woobloo-game__) to some websockets

## Woobloo Game

The source for __woobloo-game__ can be found in [src/haskell](src/haskell). This holds the logic and state of each game instance. It reads from stdin (which is the stdout of __woobloo-frontman__) and writes to stdout (which is the stdin of __woobloo-frontman__).