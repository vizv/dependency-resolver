SOURCE := cmd/demo-task-order
TARGET := $(notdir $(SOURCE))

.PHONY: all clean $(TARGET)

all: $(TARGET)
clean:
	rm -f $(TARGET)

$(TARGET): $(SOURCE)
	go build -o $@ ./$<
