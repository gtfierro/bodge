buttonuri = "amplab/sensors/s.hamilton/01b2/i.temperature/signal/operative"
pluguri = "410testbed/devices/tplink1/s.tplink.v0/0/i.xbos.plug/slot/state"

oldstate = 0
events = 0
bw.subscribe(buttonuri, nil, function(uri, msg)
    if msg["button_events"] > events then
        events = msg["button_events"]
        newstate=1-oldstate
        print("turning to", newstate)
        bw.publish(pluguri, "2.1.1.2", {state=newstate})
        oldstate=newstate
    end
end)

bw.loop()
