setup:
	python -m pip install ctfcli
	python -m ctfcli init

gen/build:
	cd ctfd-config-generator/cmd/generator && go build -o ../../gen

gen: gen/build
	./ctfd-config-generator/gen

ctfcli/build:
	cd ctfd-config-generator/cmd/ctfcli && go build -o ../../ctfcli

sync: ctfcli/build
	./ctfd-config-generator/ctfcli