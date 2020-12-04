FROM chromedp/headless-shell

RUN apt update -y && apt upgrade -y && apt install wget -y
RUN wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin && mkdir -p /go/src

EXPOSE 5000

WORKDIR /go/src/truth

COPY . .

RUN /usr/local/go/bin/go build

ENTRYPOINT /go/src/truth/truth
