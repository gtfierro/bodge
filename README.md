# BW Lua

Functions to implement:
- [ ] query: return a list of messages
    - [ ] queryone: just get the first
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
- [ ] more timer methods:
    - [ ] fire a callback periodically (e.g. `invokePeriodically`)
    - [ ] fire a callback after a set amount of time
