# Get the necessary packages if not installed
go get golang.org/x/net/html
go get gopkg.in/mgo.v2

# Make sure env variables are set
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin

# Run the program
go run main.go -url=http://localhost:8080/ -search=dfs
