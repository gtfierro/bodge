local ns = bw.arg(1)

i = 0
while true do
    i = i + 1
    val = string.format("%d",50*math.sin(.20*i)+50)
    print(val)
    bw.sleep(1000)
    bw.publish(ns.."/s.cloud/_/i.cover/signal/percent", "64.0.0.1", val)
end
