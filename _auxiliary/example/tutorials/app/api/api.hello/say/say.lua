wrk.method = "POST"

wrk.body = '{"ping":"ping"}'

wrk.headers["Content-Type"] = "application/json"

function request()

 return wrk.format('POST', nil, nil, body)

end