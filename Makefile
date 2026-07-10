CMD_NAME = qwe
INSTALL_PATH = /usr/local/bin/$(CMD_NAME)

.PHONY: all bumpversion tidy build prepare-embed-version install dev clean

all: bumpversion tidy install

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

prepare-embed-version:
	cp VERSION internal/web/VERSION

build: prepare-embed-version
	@echo "🔨 Building $(CMD_NAME)..."
	go build -ldflags "-X main.Version=$$(cat VERSION)" -o $(CMD_NAME)

install: build
	@echo "📦 Installing to $(INSTALL_PATH)..."
	sudo mv -f $(CMD_NAME) $(INSTALL_PATH)
	@echo "✅ Command '$(CMD_NAME)' installed to $(INSTALL_PATH)"

dev: prepare-embed-version
	@echo "🚀 Running in dev mode (HTML served from disk)..."
	DEV_MODE=1 go run -ldflags "-X main.Version=$$(cat VERSION)" . web

clean:
	@echo "🧹 Cleaning up..."
	rm -f $(CMD_NAME)
