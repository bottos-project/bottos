
bottos:
	go build -o bottos main.go

all: bottos

format:
	gofmt -w main.go

clean:
	rm -rf bottos *.o *.out *exe
