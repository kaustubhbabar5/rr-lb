# rr-lb
Round robin load balancer implementation

# Run
### start redis container
`docker run -it -p 6379:6379 --rm --name myredis redis:7.0.2`

### start load balancer
`make run`

### start client to test out load balancer
`make client`



# APIs
## 1. Register replica
0.0.0.0:8080/url/register

method:POST

Content-Type: application/json

Body

    
    {
        "endpoint": "https://www.amazon.com/"
    }
    
    
## 2. proxy
url: 0.0.0.0:8080/proxy/<path>

method: any

Content-Type: application/json

Body

    {
        "endpoint": "https://www.amazon.com/"
    }


Note: have to wait 