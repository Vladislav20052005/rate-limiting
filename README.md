# rate-limiting
Program rate limits the incoming http requests.
If requests come to often it responds with ```index.html``` file from dummy folder. Otherwise, it redirects request to another resource.
To run program there should be specified 3 arguments:
1. max requests per minute
2. max requests per minute for a single IP
3. url of resource

Repository contains ```Dockerfile```.
Docker-image specified runs redis server that is used for containing recent requests.
Before running program execute these commands
```
docker build -t redis-image .
```
```
docker run --rm -d -p  6379:6379 --name redis-container redis-image
```
