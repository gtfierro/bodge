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

## Concurrency

Bodge is single threaded, so concurrent manipulation of data structures is safe.

It is possible for Bodge to "return from main()", that is, finish executing all of the code in the input file. Normally, this would cause the runtime to exit and your program to stop running, which is *not* what you want if you are trying to have a persistent process!

To this end, Bodge implements the `bw.loop()` function, which keeps the Bodge scheduler and runtime active, but stops executing further lines of code in the file. Make sure this is the last line in whatever file you are executing.

```lua
bw.invokePeriodically(1000, function()
    print("say hi every 1 second")
end)

-- keep the timer running and don't exit!
bw.loop()
```

## API

* [BOSSWAVE operations](#bosswave-operations)


## BOSSWAVE Operations

Bodge implements the following BOSSWAVE operations

### Subscribe

```lua
bw.subscribe(uri, ponum, cb)
```

Subscribes to a BOSSWAVE resource and invokes a callback on every received message.

* `uri`: is a string representing a BOSSWAVE URI.
    * Examples:
        * `ciee/*`
        * `scratch.ns/devices/s.venstar/+/i.xbos.thermostat/signal/info`
* `ponum`: is a string of the dotted form of which payload objects should be matched. This can optionally use the prefix notation. `nil` matches all payload objects.
    * Examples:
        * `2.0.0.0/8`
        * `64.0.0.1`
* `cb`: this is a Lua function called upon every received message that takes up to 2 arguments: `uri` (the URI of the published message) and `msg` (the contents of the published message). Bodge makes a decent effort to convert messages into Lua types, regardless of the encoding. It currently translates MsgPack, JSON and YAML, but can be easily extended to support more formats.

```lua
count = 0
bw.subscribe("scratch.ns", nil, function(uri)
    print("Received message from", uri)
    count = count + 1
    print("Received",count,"messages so far")
end)

bw.loop()
```

### Query

Queries a BOSSWAVE resource for persisted messages and invokes a callback on every received message.

`bw.query(uri, ponum, cb)`

Like `subscribe`, but performs a query operation. The callback is invoked for each persisted message.

### GetOne

`bw.getone(uri, ponum)`

Like `subscribe`, but returns after receiving a single published message. Returns only the message and not the URI.

```lua
msg = bw.getone(`scratch.ns/*`)
print("first message is:", msg)
```

### Publish

Publishes an object on a BOSSWAVE resource.

```lua
bw.publish(uri, ponum, obj)
```

* `uri`: is a string representing a BOSSWAVE URI. It must be fully qualified because we are publishing
    * Examples:
        * `scratch.ns/devices/s.venstar/0/i.xbos.thermostat/slot/info`
* `ponum`: is a string of the dotted for0 of which payload objects should be matched. Bodge will use this to do the encoding of the payload object for you
    * Examples:
        * `2.1.1.2`
        * `64.0.0.1`
* `obj`: is the Lua object to be published. Bodge encodes this using the payload object number. `2.0.0.0/8` is MsgPack, `64.0.0.0/8` is string, `65.0.0.0/8` is JSON (broken currently) and `67.0.0.0/8` is YAML

```lua
bw.publish("scratch.ns/testme","2.0.0.0", {x=3,y=4})
```

### Persist

Persists an object on a BOSSWAVE resource.

```lua
bw.persist(uri, ponum, obj)
```

Like `publish`, but persists the message at the specified URI.

## Timers

Bodge also implements a suite of timer operations

### Sleep

Pauses the current "thread" for the given period. You can have multiple overlapping sleeps in different callbacks.

```lua
bw.sleep(period)
```

* `period` is an integer representing milliseconds

### InvokePeriodically

Invokes a callback upon every elapsed period.

```lua
bw.invokePeriodically(period, cb, arg1, arg2)
```

* `period` is an integer representing milliseconds
* `cb` is a function that takes an arbitrary number of arguments, which are specified after the callback

```lua
count = 0
function x()
    count = count + 1
    print(count,"seconds have elapsed")
end
bw.invokePeriodically(1000, x)

bw.loop()
```

### InvokeLater

```lua
bw.invokeLater(period, cb, arg1, arg2)
```

Like `invokePeriodically`, but only executes once.

### Loop

```lua
bw.loop()
```

Keeps the runtime active, but does not execute any further lines of code in the file. Used to keep a Bodge program persistently running.


## Utilities

### DumpTable

```lua
bw.dumptable(table)
```

Outputs the content of a Lua table

### NArgs

```lua
bw.nargs()
```

Returns the number of command line arguments provided to the invocation of Bodge

```lua
local n = bw.nargs()
print("Got",n,"arguments")
```

### Arg

```lua
bw.arg(n)
```

This is 1-indexed (like Lua). Returns the nth argument provided to the invocation of Bodge

```lua
local first = bw.arg(1)
print("First argument was", first)
```

### URIRequire

```lua
bw.uriRequire(uri)
```

Imports and runs the Bodge code persisted at the specified URI under PO numbers `64.0.2.0/24`.

```lua
tstat_class = bw.uriRequire('bodge/lib/xbos_tstat.lua')

my_tstat = tstat_class("uri for my tstat")
```
