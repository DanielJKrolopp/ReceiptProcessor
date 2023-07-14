# Daniel Krolopp's Receipt Processor

### Usage
Assuming you have Docker installed, you can run my submission with
`docker run --rm -p 8080:8080 -it $(docker build -q .)` from
within the top-level directory. This will bring up the container and 
remove it automatically when you are finished. You can then make any 
requests against `localhost:8080` using cUrl, Postman or another application.

### Design choices
I chose to write this in Go, using the Gin framework for serving up RESTful
endpoints. While my day-to-day work is done in Java and Spring, I've used
Go and Gin in the past and am more familiar with Gin than other Go
backend frameworks.

I containerized the application using Docker, to allow for portability to
systems that don't have Go installed.

### Testing
Should you want to run the test suite and have Go installed, run
`go test` within the top-level source directory.