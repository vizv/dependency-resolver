SOURCE := cmd/demo-dependency-resolver
TARGET := demo-dependency-resolver

.PHONY: all clean $(TARGET)

all: $(TARGET)
clean:
	rm -f $(TARGET)

$(TARGET): $(SOURCE)
	go build -o $@ ./$<
