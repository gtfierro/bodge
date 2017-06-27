# Bodge

**Bodge**: *Make or repair (something) badly or clumsily.*

Bodge is a "batteries included" DSL and execution environment for prototyping and scheduling BOSSWAVE interactions. 

Bodge extends the [Lua](https://github.com/yuin/gopher-lua) embedded programming language with:
- `bw`, a module for BOSSWAVE operations
- basic timer functionality
- task scheduling
- ability to execute Bodge code from a file, interactive command line or a BOSSWAVE URI

### Installation

Bodge can be installed simply by running

```
go get github.com/gtfierro/bodge
go install github.com/gtfierro/bodge
```

Then, to run Bodge, just run

```
bodge
```

All library functionality is provided through a `bw` module. This is imported automatically into the Bodge runtime, but you can also explicitly run

```lua
local bw = require('bw')
```

## BOSSWAVE Operations

Bodge implements the following BOSSWAVE operations

### Subscribe

```lua
bw.subscribe(uri, ponum, cb)
```

* `uri`: is a string representing a BOSSWAVE URI. 
    * Examples:
        * `ciee/*`
        * `scratch.ns/devices/s.venstar/+/i.xbos.thermostat/signal/info`
* `ponum`: is a string of the dotted form of which payload objects should be matched. This can optionally use the prefix notation. 
    * Examples:
        * `2.0.0.0/8`
        * `64.0.0.1`
* `cb`: this is a Lua function called upon every received message that takes up to 2 arguments: `uri` (the URI of the published message) and `msg` (the contents of the published message). Bodge makes a decent effort to convert messages into Lua types, regardless of the encoding. It currently translates MsgPack, JSON and YAML, but can be easily extended to support more formats.


## Concurrency

Bodge is single threaded, so concurrent manipulation of data structures is safe.
