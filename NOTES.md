# BW Lua

Functions to implement:
- [X] query: return a list of messages
- [ ] publish with persist
- [ ] stats:
    - [ ] port https://github.com/montanaflynn/stats
- [ ] kv store API:
    - [ ] port leveldb golang bindings. Get persistent storage for your script!
- [ ] memory-mapped devices:
    - for each standard xbos interface, make a distributed object that is synchronized with the device
      ```lua
      local tstat = Thermostat:new("uri of the thermostat iface")
      tstat.get_temperature()
      tstat.set_heating_setpoint(70)
      ```
    - can we get the equivalent of overriding `__setattr__` in Python?
- [ ] Deploy to, import libraries, run code on BOSSWAVE URIs
- [X] more timer methods:
    - [X] fire a callback periodically (e.g. `invokePeriodically`)
    - [X] fire a callback after a set amount of time

## Concurrency Model

Rather than using Lua coroutines, might be nice to use the `cord` approach

- Do NOT hack the runtime to add locks
- lua coroutines do the heavy lifting. Wrap the function provided to the subscribe
  callback in another function that yields coroutines.
- need something like "cord" to create new threads, and takes care of
  collection the pointers and running through them.

## Clock/Date Scheduling

We want something that is human readable and allows the expression of the following:
- events like
    - every Monday
    - next Monday
    - every day at 10:30am
    - every weekday at 10:30am
    - every weekend at 10:30am
    - every day
    - every hour
    - every minute
- API?
    ```lua
    -- next monday
    on("1/30/2017", cb)

    -- every Monday
    every("monday", cb)

    -- every day at 10:30
    every("10:30am", cb)

    -- every weekday at 10:30am
    every("weekday 10:30am", cb)

    -- every weekend at 10:30am
    every("weekend 10:30am", cb)

    -- every day (00:00)
    every("day", cb)

    -- every hour (00:00)
    every("hour", cb)

    -- every minute (00:00)
    every("minute", cb)
    ```

- also have a "list timers" so we can see what's running and when it will next trigger
