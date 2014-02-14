TARGET = rss

all: $(TARGET)

rss: src/rss.go
	go build $<
