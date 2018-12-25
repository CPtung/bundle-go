GOR=go run
GOB=go build
SRC=main/main.go

run:
	$(GOR) $(SRC)

build: tagservice

tagservice:
	$(GOB) -o $@ $(SRC)

clean:
	rm -rf tagservice

.PHONY= build run clean
