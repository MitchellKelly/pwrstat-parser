FROM golang

RUN wget https://dl4jz3rbrsfum.cloudfront.net/software/powerpanel-132-x86_64.tar.gz && \
	tar xf powerpanel-132-x86_64.tar.gz && \
	cd powerpanel-1.3.2 && \
	./install.sh

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["service", "pwrstatd", "start"]
