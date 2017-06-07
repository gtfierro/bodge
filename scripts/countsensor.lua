uri = "ciee/*/operative"
count = 0
subscribe(uri, nil, function(uri)
    count = count +1
    print(count)
end)

keeprunning()
