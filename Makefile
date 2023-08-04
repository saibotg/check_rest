# Name des ausführbaren Programms
APP_NAME = check_rest

# Befehle
GOFLAGS = -v

.PHONY: build test

# Standardziel: Build und Test ausführen
all: build test

# Build-Schritt
build:
	go build $(GOFLAGS) -o dist/$(APP_NAME) ./...

# Test-Schritt
test:
	go test $(GOFLAGS) ./...

# "clean" Ziel zum Löschen der generierten Dateien
clean:
	rm -f $(APP_NAME)
