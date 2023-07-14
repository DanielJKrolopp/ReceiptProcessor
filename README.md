# Daniel Krolopp's Receipt Processor

### Usage

Assuming you have Docker installed, you can run my submission with
```docker run --rm -p 8080:8080 -it $(docker build -q .)```. 
This will bring up the container and remove it automatically when you
are finished. You can then make any requests against localhost:8080
using cUrl, Postman or another application.