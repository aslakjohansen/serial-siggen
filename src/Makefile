TARGETS = \
	serial-siggen \
	serial-siggen-win64.exe \
	serial-siggen-win32.exe \


all: ${TARGETS}

clean:
	touch ${TARGETS}
	rm    ${TARGETS}

mrproper: clean
	touch dummy~
	rm        *~


serial-siggen: serial-siggen.go
	go get github.com/tarm/serial
	go build serial-siggen.go

serial-siggen-win64.exe: serial-siggen.go
	go get github.com/tarm/serial
	GOOS=windows GOARCH=amd64 go build -o serial-siggen-win64.exe serial-siggen.go
	
serial-siggen-win32.exe: serial-siggen.go
	go get github.com/tarm/serial
	GOOS=windows GOARCH=386 go build -o serial-siggen-win32.exe serial-siggen.go
