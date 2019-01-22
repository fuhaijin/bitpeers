# bitpeers

	Usage of bitpeers:
		--addressonly       outputs only addresses if specified
		--filepath string   the path to peers.dat
		--format string     the output format {json|text} (default "json")

download from github，compile and install package

    go get -u github.com/fuhaijin/bitpeers/cmd/bitpeers

windows build:
    
    go build .\cmd\bitpeers\bitpeers.go

    ./bitpeers --filepath ./peers.dat --format text
    
    ./bitpeers --filepath ./peers.dat --addressonly --format text