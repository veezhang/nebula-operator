FROM alpine:3.18.2

ARG TARGETDIR

RUN addgroup -S -g 65532 ng-user && \
    adduser -S -D -H -u 65532 \
    -s /sbin/nologin -G ng-user -g ng-user ng-user

ADD bin/${TARGETDIR}/controller-manager /usr/local/bin/controller-manager
ADD bin/${BUILDPLATFORM}/autoscaler /usr/local/bin/autoscaler
ADD bin/${TARGETDIR}/scheduler /usr/local/bin/scheduler
USER 65532:65532
