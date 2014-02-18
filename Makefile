TARGET = rss

all: $(TARGET)

rss: $(wildcard src/*.go)
	go build $^
