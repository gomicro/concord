FROM scratch
MAINTAINER dev@gomicro.io

ADD concord concord

CMD ["/concord"]
