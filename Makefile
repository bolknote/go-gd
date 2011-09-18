GOROOT=/usr/local/Cellar/go/r59
GOARCH=amd64

include ${GOROOT}/src/Make.inc

TARG=gd

CGOFILES=\
	gd.go

CGO_CFLAGS=-I/usr/local/Cellar/gd/2.0.36RC1/include
CGO_LDFLAGS+=-L/usr/local/Cellar/gd/2.0.36RC1/lib
CGO_LDFLAGS=-lgd

CLEANFILES+=sample

include ${GOROOT}/src/Make.pkg

sample: install sample.go
	$(GC) sample.go
	$(LD) -o $@ sample.$O
