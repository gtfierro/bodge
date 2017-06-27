state = 0
bw.invokePeriodically(1000, function()
    state = state + 1
end)
bw.invokePeriodically(1500, function()
    state = state + 1
end)

bw.invokePeriodically(2000, function()
    print(state)
end)

bw.loop()
