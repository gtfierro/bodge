if bw.nargs() == 0 then
    namespace = "ucberkeley"
else
    namespace = bw.arg(1)
end


_suburi = namespace .. "/*"

count = 0
bw.subscribe(_suburi, "", function()
    count = count + 1
end)

bw.invokePeriodically(10000, function()
    per_sec = count / 10
    count = 0
    print(string.format("%0.2f msgs per sec", per_sec))
end)

print("Measuring msgs/sec on " .. namespace)
print("Starting (10sec)...")

bw.loop()
