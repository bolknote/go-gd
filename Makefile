GOROOT=/usr/local/Cellar/go/r60.3

include ${GOROOT}/src/Make.inc
TARG=gd
CGOFILES=\
	gd.go

CGO_LDFLAGS=-lgd

CLEANFILES+=sample

include ${GOROOT}/src/Make.pkg

sample: install sample.go
	$(GC) sample.go
	$(LD) -o $@ sample.$O
