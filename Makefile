SOURCE := cmd/demo-task-order
TARGET := a.out

.PHONY: all clean $(TARGET)

all: $(TARGET)
clean:
	rm -f $(TARGET)

$(TARGET): $(SOURCE)
	go build -o $@ ./$<
