CMD_NAME = qwe
INSTALL_PATH = /usr/local/bin/$(CMD_NAME)


all: bumpversion tidy build install
define PATCH_VERSION
$(shell grep -oE '[0-9]+$$' VERSION)
endef

bumpversion:
	@echo "🔢 Bumping patch version..."
	old=$$(grep -oE '[0-9]+$$' VERSION); \
	new=$$(($$old + 1)); \
	sed -i '' -E "s/(v1\.0\.)([0-9]+)/\\1$$new/" VERSION; \
	echo "🔢 Version updated to $$(cat VERSION)"

tidy:
	@echo "🔍 Running go mod tidy..."
	go mod tidy

build:
	@echo "🔨 Building $(CMD_NAME)..."
	go build -ldflags "-X main.Version=$$(cat VERSION)" -o $(CMD_NAME) main.go

install: build
	@echo "📦 Installing to $(INSTALL_PATH)..."
	sudo mv -f $(CMD_NAME) $(INSTALL_PATH)
	@echo "✅ Command '$(CMD_NAME)' installed to $(INSTALL_PATH)"
	export PATH=$$PATH:/usr/local/go/bin

clean:
	@echo "🧹 Cleaning up..."
	rm -f $(CMD_NAME)
