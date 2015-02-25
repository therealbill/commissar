# Configuration

Quick and dirty during early exploration development.

Environment variables:
```
LATENCY_CLIENTCOUNT=1
LATENCY_ITERATIONS=100
LATENCY_JSONOUT=FALSE
LATENCY_REDISAUTHTOKEN=<auth>
LATENCY_REDISCONNECTIONSTRING=<host:port>
```

If LATENCY_JSONOUT is set to true, only the JSON output will be printed to
stdout. LATENCY_CLIENTCOUNT determines the number of concurrent connections to
run the test with. This will help show how your Redis setup handles various 
levels of concurrency from the clients. 


If storing results in mongo: 
```
LATENCY_MONGOCOLLECTIONNAME=<name>
LATENCY_MONGOCONNSTRING=<connstring>
LATENCY_MONGODBNAME=<dbname>
LATENCY_MONGOUSERNAME=<username>
LATENCY_MONGOPASSWORD=<password>
```


# Results
The latency numbers are in nanoseconds, and represent the point of view of the
client. As such it includes networking.
