SOURCE := cmd/demo-dependency-resolver
TARGET := demo-dependency-resolver

TEST_IN := $(wildcard test/*.in)
TEST_IN_GV := $(TEST_IN:=.gv)
TEST_GV := $(wildcard test/*.gv)

GV_SCRIPT := scripts/gen-gv.sh

.PHONY: all draw test clean $(TARGET)

all: $(TARGET)
draw: $(TEST_IN_GV) $(TEST_GV)
test: $(TARGET)
	echo $(wildcard test/*.in)
clean:
	rm -f $(TARGET) test/*.png

$(TARGET): $(SOURCE)
	go build -o $@ ./$<

%.in.gv: %.in $(GV_SCRIPT)
	$(GV_SCRIPT) $<
%.png: %.gv
	dot -Tpng -o$@ $<
