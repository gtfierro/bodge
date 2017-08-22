local bw = require('bw')
local math = require('math')
local mod = {}

local new = function(uri)
    local client = {
        uri = uri or "",
        nonce = nil,
        result = {},
    }
    bw.subscribe(uri .. "/s.hod/_/i.hod/signal/result", "2.0.10.2", function(uri, msg)
        if tostring(msg.Nonce) == client.nonce then
            client.nonce = nil
            client.result = msg
            if client.result.Count > 0 then
                for i, r in ipairs(client.result.Rows) do
                    new_r = {}
                    for var, row in pairs(r) do
                        if #row.Namespace > 0 then new_r[var] = row.Namespace .. "#" end
                        new_r[var] = new_r[var] .. row.Value
                    end
                    client.result.Rows[i] = new_r
                end
            end
        end
    end)
    return setmetatable(client, mod)
end

local query = function(self, q)
    while self.nonce ~= nil do bw.sleep(100) end
    self.nonce = tostring(math.random(0,32768))
    bw.publish(self.uri .. "/s.hod/_/i.hod/slot/query", "2.0.10.1", {
        Query= q,
        Nonce= tostring(self.nonce),
    })
    while self.result.Nonce == nil do bw.sleep(100) end
    local r = self.result
    self.result = {}
    return r
end
mod.query = query

mod.__index = mod
local ctor = function(cls, ...)
    return new(...)
end
return setmetatable({}, {__call = ctor})
