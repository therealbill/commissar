FROM busybox
MAINTAINER Bill Anderson <therealbill@me.com>
ADD commissar /commissar
CMD ["/commissar"]
