FROM scpike/libpostal-docker

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
# RUN go build -o main src/main.go 
COPY ./main /app/main

ENV GIN_MODE release
ENTRYPOINT ["/app/main"]

EXPOSE 8080
