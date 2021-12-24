SOURCE := cmd/demo-dependency-resolver
TARGET := ./demo-dependency-resolver

TEST_IN := $(wildcard test/*.in)
TEST_IN_GV := $(TEST_IN:=.gv)
TEST_GV := $(filter-out $(TEST_IN_GV), $(wildcard test/*.gv))
TEST_CMP := $(TEST_IN:=.cmp) $(TEST_GV:=.cmp)

GV_SCRIPT := scripts/gen-gv.sh

.PHONY: all draw test clean $(TARGET)

all: $(TARGET)
draw: $(TEST_IN_GV:.gv=.png) $(TEST_GV:.gv=.png)
	@echo 'ALL DRAWN !!!'
test: $(TEST_CMP)
	@echo 'ALL PASSED !!!'
clean:
	rm -f $(TARGET) test/*.in.gv test/*.png test/*.out

%.out: % $(TARGET)
	@touch $@
	-$(TARGET) $< > $@

%.cmp: %.out %.ans
	cmp $^
	@echo

$(TARGET): $(SOURCE)
	go build -o $@ ./$<

%.in.gv: %.in $(GV_SCRIPT)
	$(GV_SCRIPT) $<
%.png: %.gv
	dot -Tpng -o$@ $<
	@echo
