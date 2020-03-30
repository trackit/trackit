FROM scratch
MAINTAINER Victor Schubert <victor@trackit.io>
EXPOSE 80
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY cost_type_patterns.json /etc/config/cost_type_patterns.json
COPY main /main
ENTRYPOINT ["/main"]
