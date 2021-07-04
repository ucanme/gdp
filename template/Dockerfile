FROM golang:1.13 as builder

WORKDIR /root
ADD ../template .

RUN go env -w GOPROXY=https://goproxy.io,direct

RUN make

FROM golang:1.13
WORKDIR /root
COPY --from=builder /root/bin/{{.app_name}} /root/{{.app_name}}
RUN ln -s /root/{{.app_name}} /bin/{{.app_name}}

CMD ["fastgo", "server"]
