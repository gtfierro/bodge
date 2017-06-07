uri = "ciee/devices/echola/s.powerup.v0/%d/i.xbos.plug/slot/info"
count = 0
subscribe(uri, nil, function(uri)
    count = count +1
    print(count)
end)

keeprunning()
