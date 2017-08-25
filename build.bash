mkdir -p build
go get github.com/gorilla/handlers
for GOOS in linux darwin windows
do
    export GOOS
    for GOARCH in amd64 386
    do
        export GOARCH
        output="build/httpserver-$GOOS-$GOARCH"
        [[ $GOOS == "windows" ]] && output="$output.exe"
        echo $output
        go build -o $output httpserver.go &
    done
done
wait
