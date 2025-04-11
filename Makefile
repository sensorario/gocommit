CMD_NAME = gocommit
INSTALL_PATH = /usr/local/bin/$(CMD_NAME)

all: build install

build:
	go build -o $(CMD_NAME) main.go

install:
	sudo mv -f $(CMD_NAME) $(INSTALL_PATH)
	@echo "✅ Command '$(CMD_NAME)' installed to $(INSTALL_PATH)"

clean:
	rm -f $(CMD_NAME)
