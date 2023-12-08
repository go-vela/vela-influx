# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-influx:certs .    ##
#########################################################################

FROM alpine:3.18.5@sha256:34871e7290500828b39e22294660bee86d966bc0017544e848dd9a255cdf59e0 as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-influx:local .    ##
##########################################################

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-influx /bin/vela-influx

ENTRYPOINT [ "/bin/vela-influx" ]
