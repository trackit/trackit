FROM scratch
MAINTAINER Victor Schubert <victor@trackit.io>
EXPOSE 80
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY main /main
ENTRYPOINT ["/main"]
