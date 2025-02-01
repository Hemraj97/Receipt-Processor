# Receipt-Processor
Take Home Assignment for Fetch

Steps to run to the dockerfile. Ensure Docker is running on your machine.
1. docker build --no-cache -t receipt-processor .

2. docker run -p 8080:8080 receipt-processor

This will expose your Go application on port 8080 of your local machine, and you can interact with it at http://localhost:8080.