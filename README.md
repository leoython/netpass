# netpass
A minimalist intranet penetration, port proxy tool written in go

## How to use

### build
run 
``
go build main.go
``

### run netpass
run
```
./netpass -localPort=9001 -remoteHost=10.0.0.1 -remotePort=8090
```

## TODO
-[] SSH tunnel

-[] TCP heartbeat and reconnect

-[] Local connection pool

-[] Multi-Port Proxy

