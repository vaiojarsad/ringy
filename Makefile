ifeq ($(OS), Windows_NT)
    BIN := ringy.exe
else
    BIN := ringy
endif

# Build
build:
	go build -o $(BIN) .

# Clean
clean:
	rm -f ringy ringy.exe
