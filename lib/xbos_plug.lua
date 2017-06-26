local mod = {}

local new = function(uri)
    local obj = {
        uri = uri or "",
        _state = nil,
        _voltage = nil,
        _current = nil,
        _power = nil,
        _cumulative = nil
    }
    print("Tstat at", uri)
    subscribe(uri .. "signal/info", "2.1.1.2", function(uri, msg)
        obj._state = msg["state"]
        obj._voltage = msg["voltage"]
        obj._current = msg["current"]
        obj._power = msg["power"]
        obj._cumulative = msg["cumulative"]
    end)
    return setmetatable(obj, mod)
end

local state = function(self, val)
    if val == nil then
        return self._state
    end
    publish("410testbed/devices/tplink2/s.tplink.v0/0/i.xbos.plug/slot/state", "2.1.1.2", {state=val})
    return self._state
end
mod.state = state

local voltage = function(self)
    return self._voltage
end
mod.voltage = voltage

local current = function(self)
    return self._current
end
mod.current = current

local power = function(self)
    return self._power
end
mod.power = power

local cumulative = function(self)
    return self._cumulative
end
mod.cumulative = cumulative

mod.__index = mod
local ctor = function(cls, ...)
    return new(...)
end

return setmetatable({}, {__call = ctor})
