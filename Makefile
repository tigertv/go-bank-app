GOC = go
SRC = main.go bank.go

.PHONY: all clean

all: main

testpack: testpack.go
	$(GOC) build -x $<

main: $(SRC) 
	$(GOC) build $^

run: $(SRC)
	$(GOC) run $^

clean: 
	rm -rf main
