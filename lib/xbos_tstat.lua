local mod = {}

local new = function(uri)
    local obj = {
        uri = uri or "",
        _temperature = nil,
        _relative_humidity = nil,
        _heating_setpoint = nil,
        _cooling_setpoint = nil,
        _override = nil,
        _fan = nil,
        _mode = nil,
        _state = nil
    }
    print("Tstat at", uri)
    subscribe(uri .. "/signal/info", "2.1.1.0", function(uri, msg)
        obj._state = msg["state"]
        obj._temperature = msg["temperature"]
        obj._relative_humidity = msg["relative_humidity"]
        obj._heating_setpoint = msg["heating_setpoint"]
        obj._cooling_setpoint = msg["cooling_setpoint"]
        obj._override = msg["override"]
        obj._fan = msg["fan"]
        obj._mode = msg["mode"]
        obj._state = msg["state"]
    end)
    return setmetatable(obj, mod)
end

-- read state
local temperature = function(self)
    return self._temperature
end
mod.temperature = temperature

local relative_humidity = function(self)
    return self._relative_humidity
end
mod.relative_humidity = relative_humidity

local state = function(self)
    return self._state
end
mod.state = state

-- read/write state
local heating_setpoint = function(self, val)
    if val ~= nil then
        publish(self.uri.."/slot/state","2.1.1.0",{heating_setpoint=val})
    end
    return self._heating_setpoint
end
mod.heating_setpoint = heating_setpoint

local cooling_setpoint = function(self, val)
    if val ~= nil then
        publish(self.uri.."/slot/state","2.1.1.0",{cooling_setpoint=val})
    end
    return self._cooling_setpoint
end
mod.cooling_setpoint = cooling_setpoint

local override = function(self, val)
    if val ~= nil then
        publish(self.uri.."/slot/state","2.1.1.0",{override=val})
    end
    return self._override
end
mod.override = override

local fan = function(self, val)
    if val ~= nil then
        publish(self.uri.."/slot/state","2.1.1.0",{fan=val})
    end
    return self._fan
end
mod.fan = fan

local mode = function(self, val)
    if val ~= nil then
        publish(self.uri.."/slot/state","2.1.1.0",{mode=val})
    end
    return self._mode
end
mod.mode = mode

mod.__index = mod
local ctor = function(cls, ...)
    return new(...)
end

return setmetatable({}, {__call = ctor})
