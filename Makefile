setup:
	python -m pip install ctfcli

gen/build:
	cd ctfd-config-generator/cmd/generator && go build -o ../../gen

gen: gen/build
	./ctfd-config-generator/gen

ctfcli/init:
	python -m ctfcli init
	./ctfd-config-generator/cmd/scripts/install.sh

ctfcli/sync:
	./ctfd-config-generator/cmd/scripts/sync.sh

