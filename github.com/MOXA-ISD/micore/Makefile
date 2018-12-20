GOR=go run
GOB=go build
SRC=main/main.go

run:
	$(GOR) $(SRC)

build: sample

sample:
	$(GOB) -o $@ $(SRC)

clean:
	rm -rf sample

.PHONY= build run clean
