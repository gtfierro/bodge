# Bodge

**Bodge**: *Make or repair (something) badly or clumsily.*

Bodge is a "batteries included" DSL and execution environment for prototyping and scheduling BOSSWAVE interactions.

Bodge extends the [Lua](https://github.com/yuin/gopher-lua) embedded programming language with:
- `bw`, a module for BOSSWAVE operations
- basic timer functionality
- task scheduling
- ability to execute Bodge code from a file, interactive command line or a BOSSWAVE URI

### Usage

```
NAME:
   bodge - Simple BOSSWAVE Lua scripts for interaction, exploration and rule building

USAGE:
   bodge [global options] command [command options] [arguments...]

VERSION:
   0.2.0

COMMANDS:
     publish, p, pub  Publish a file to a given URI
     cat              Cat a file on a given URI
     ls               List bodge files
     help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --agent value, -a value   Local BOSSWAVE Agent (default: "127.0.0.1:28589") [$BW2_AGENT]
   --entity value, -e value  The entity to use [$BW2_DEFAULT_ENTITY]
   --help, -h                show help
   --version, -v             print the version
```

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

### Building

Bodge can optionally be built to communicate with a [remote agent](https://github.com/immesys/ragent). Bodge will embed the entity and the ragent client inside the Bodge binary so it can be run on a computer without BOSSWAVE being installed. Use the `build-ragent.sh` script to build this.

By default, Bodge connects to the local BOSSWAVE agent.

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

## Example

Here is a sample weekday and weekend schedule implemented using the Bodge API

```lua
-- 410testbed_schedule.lua
-- run using "bodge 410testbed_schedule.lua"
tstat_class = bw.uriRequire('bodge/lib/xbos_tstat.lua')
tstat1 = tstat_class("410testbed/devices/venstar/s.venstar/420lab/i.xbos.thermostat")

plug_class = bw.uriRequire('bodge/lib/xbos_plug.lua')
plug1 = plug_class("410testbed/devices/tplink2/s.tplink.v0/0/i.xbos.plug")

bw.every("weekday 07:30", function()
    tstat1:heating_setpoint(72)
    tstat1:cooling_setpoint(76)
end)

bw.every("weekday 12:00", function()
    tstat1:heating_setpoint(70)
    tstat1:cooling_setpoint(80)
end)

bw.every("weekday 13:00", function()
    tstat1:heating_setpoint(72)
    tstat1:cooling_setpoint(86)
end)

bw.every("weekday 18:00", function()
    tstat1:heating_setpoint(50)
    tstat1:cooling_setpoint(90)
    plug1:state(0) -- turn off plug at 6pm
end)

bw.every("weekend 00:00", function()
    tstat1:heating_setpoint(50)
    tstat1:cooling_setpoint(90)
end)
```

## API

* [BOSSWAVE operations](#bosswave-operations)
    * [subscribe](#subscribe)
    * [query](#query)
    * [getone](#getone)
    * [publish](#publish)
    * [persist](#persist)
* [Timers](#timers)
    * [sleep](#sleep)
    * [invokePeriodically](#invokeperiodically)
    * [invokeLater](#invokelater)
    * [loop](#loop)
* [Utilities](#utiilies)
    * [dumptable](#dumptable)
    * [nargs](#nargs)
    * [arg](#arg)
    * [uriRequire](#urirequire)
* [Scheduler](#scheduler)
    * [every](#every)
* [XBOS Devices](#xbos-devices)

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
bw.invokeLater(period, cb, ...)
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

## Scheduler

```lua
bw.every(spec, cb, ...)
```

This function executes the callback (with any provided arguments) according to the schedule specification string `spec`.

Bodge supports the following specifications:
* Day/Hour/Minute/Second periodic, specified using Go-style duration strings:
    * Executed from the time that the script is run -- NOT aligned!
    * Examples:
        * every 2 hours: `2h`
        * every day: `1d`
* Day/time periodic:
    * Can use the following to indicate what days the schedule runs:
        * `monday`,`tuesday`, etc... (runs on that day of the week)
        * `weekday`, `weekend` (runs on weekdays, weekends)
    * Within a day, can specify which times to execute:
        * uses 24-hour clock. Does not currently support seconds (rounds up to nearest minute)
        * needs leading zeros for times before 10am
        * `15:30` (run at 3:30pm)
        * `07:00` (run at 7am)
    * Combine these into single schedule specs:
        * Every Weekday at 7am and 8am: `weekday 07:00 08:00`
        * Weekends at 6pm: `weekend 18:00`
        * Mondays and Wednesdays at 3:30: `monday wednesday 15:30`
* **coming soon**: specific dates e.g. "run Wednesday 28 June, 3:30pm"

## XBOS Devices

Bodge makes it easy to write wrapper classes for common interactions.
One example of this are XBOS devices such as plugs and thermostats.

There are currently two such implementations:
* [XBOS plug](https://github.com/gtfierro/bodge/blob/master/lib/xbos_plug.lua) published at `bodge/lib/xbos_plug.lua`
* [XBOS thermostat](https://github.com/gtfierro/bodge/blob/master/lib/xbos_tstat.lua) published at `bodge/lib/xbos_tstat.lua`

Basic usage is as follows:

```lua
-- uses https://github.com/SoftwareDefinedBuildings/XBOS/blob/master/interfaces/xbos_thermostat.yaml
tstat_class = bw.uriRequire('bodge/lib/xbos_tstat.lua')

-- instantiate with the base URI of the xbos tstat interface
-- In the future, this URI will come from a Brick query
tstat1 = tstat_class("410testbed/devices/venstar/s.venstar/420lab/i.xbos.thermostat")

-- wait until the thermostat publishes and we get a state
while tstat1:heating_setpoint() == nil do
    bw.sleep(1000)
end

-- get the heating and cooling setpoints
hsp = tstat1:heating_setpoint()
csp = tstat1:cooling_setpoint()
print("heating setpoint", hsp)
print("cooling setpoint", csp)

-- increase the band by 2 degrees on either side
tstat1:heating_setpoint(hsp-2)
tstat1:cooling_setpoint(csp+2)

-- wait 10 seconds for the report and test that it worked
bw.sleep(10*1000)
print("heating setpoint", tstat1:heating_setpoint())
print("cooling setpoint", tstat1:cooling_setpoint())

-- can also set multiple fields concurrently
tstat1:write({fan=1,mode=3})
```
