FROM postman/newman:alpine
COPY test.json /etc/newman/
ENTRYPOINT ["newman", "run", "test.json", "--verbose", "--bail"]