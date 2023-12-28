FROM scratch

LABEL org.opencontainers.image.source=https://github.com/gomicro/concord
LABEL org.opencontainers.image.authors="dev@gomicro.io"

ADD concord concord

CMD ["/concord"]
