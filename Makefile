GOROOT=/Users/bolk/go
GOARCH=amd64

include ${GOROOT}/src/Make.inc

TARG=gd

CGOFILES=\
	gd.go

CGO_CFLAGS=-I/usr/local/Cellar/gd/2.0.36RC1/include
CGO_LDFLAGS+=-L/usr/local/Cellar/gd/2.0.36RC1/lib
CGO_LDFLAGS=-lgd

CLEANFILES+=img

include ${GOROOT}/src/Make.pkg

img: install img.go
	$(GC) img.go
	$(LD) -o $@ img.$O
