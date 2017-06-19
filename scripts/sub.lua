count = {count = 0}
uri = "ciee/*/operative"
function trigger(uri)
    count["count"] = count["count"] + 1
end
subscribe(uri, "", trigger)


invokePeriodically(1000, function()
    print("count", count["count"])
end)

loop()
