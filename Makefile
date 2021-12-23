SOURCE := cmd/demo-dependency-resolver
TARGET := a.out

.PHONY: all clean $(TARGET)

all: $(TARGET)
clean:
	rm -f $(TARGET)

$(TARGET): $(SOURCE)
	go build -o $@ ./$<
