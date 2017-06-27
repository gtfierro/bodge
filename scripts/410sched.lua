-- import thermostat class
tstat = bw.uriRequire("bodge/lib/xbos_tstat.lua")
-- instantiate with 410 testbed venstar
venstar = tstat("410testbed/devices/venstar/s.venstar/420lab/i.xbos.thermostat")

-- THE thermostat DRIVER needs to not re-post the old state

bw.every("07:30", function()
    venstar:heating_setpoint(72)
    venstar:cooling_setpoint(76)
end)

bw.every("12:00", function()
    venstar:heating_setpoint(70)
    venstar:cooling_setpoint(80)
end)

bw.every("13:00", function()
    venstar:heating_setpoint(72)
    venstar:cooling_setpoint(76)
end)

bw.every("18:00", function()
    venstar:heating_setpoint(50)
    venstar:cooling_setpoint(90)
end)

bw.loop()
