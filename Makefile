CMD_NAME = qwe
INSTALL_PATH = /usr/local/bin/$(CMD_NAME)

all: tidy build install

tidy:
	@echo "🔍 Running go mod tidy..."
	go mod tidy

build:
	@echo "🔨 Building $(CMD_NAME)..."
	go build -o $(CMD_NAME) main.go

install:
	@echo "📦 Installing to $(INSTALL_PATH)..."
	sudo mv -f $(CMD_NAME) $(INSTALL_PATH)
	@echo "✅ Command '$(CMD_NAME)' installed to $(INSTALL_PATH)"
	export PATH=$$PATH:/usr/local/go/bin

clean:
	@echo "🧹 Cleaning up..."
	rm -f $(CMD_NAME)
