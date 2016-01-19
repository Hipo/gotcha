prepare:
	GOPATH=`pwd` go get github.com/gorilla/mux
	GOPATH=`pwd` go get golang.org/x/crypto/bcrypt
	GOPATH=`pwd` go get gopkg.in/mgo.v2
	GOPATH=`pwd` go get gopkg.in/mgo.v2/bson
	GOPATH=`pwd` go get github.com/hipo/gotcha

build:
	GOPATH=`pwd` go build