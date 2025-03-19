LOCAL_BIN := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/.bin

EMBED4ASSETS_VERSION=v0.1.0
BUF_VERSION=1.28.1

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# Assume Linux by default.
SED_INPLACE := -i''

# If macOS, change the flags.
ifeq ($(UNAME_S),Darwin)    
	SED_INPLACE := -i ''
endif

export PATH := ${LOCAL_BIN}:$(PATH)

tools: gen-bin-dir buf-install embed4assets-install

.PHONY: gen-bin-dir
gen-bin-dir:
	@mkdir -p $(LOCAL_BIN)

.PHONY: embed4assets-install
embed4assets-install:
	@if [ ! -x "$(LOCAL_BIN)/generate-go-embed4assets" ]; then \
		echo "Installing generate-go-embed4assets..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/f9n/generate-go-embed4assets/cmd/generate-go-embed4assets@$(EMBED4ASSETS_VERSION); \
    fi

.PHONY: buf-install
buf-install:
	@if [ ! -x "$(LOCAL_BIN)/buf" ]; then \
		echo "Installing buf..."; \
		curl -SL "https://github.com/bufbuild/buf/releases/download/v$(BUF_VERSION)/buf-$(UNAME_S)-$(UNAME_M)" -o "$(LOCAL_BIN)/buf" && chmod +x "$(LOCAL_BIN)/buf"; \
	fi

format:
	buf format -w

lint:
	buf lint

lock:
	buf mod update

clean:
	rm -rf gen

generate: tools format lint clean
	$(LOCAL_BIN)/buf generate && $(LOCAL_BIN)/generate-go-embed4assets -directory gen
	find ./gen/js -type f -name '*.js' -exec sed $(SED_INPLACE) -e "/require('..\/..\/..\/..\/google\/api\/annotations_pb.js');/d" {} \;
	find ./gen/js -type f -name '*.js' -exec sed $(SED_INPLACE) -e "/require('..\/..\/..\/..\/protoc-gen-openapiv2\/options\/annotations_pb.js');/d" {} \;
	find ./gen/js -type f -name '*.js' -exec sed $(SED_INPLACE) -e "/goog.object.extend(proto, google_api_annotations_pb);/d" {} \;
	find ./gen/js -type f -name '*.js' -exec sed $(SED_INPLACE) -e "/goog.object.extend(proto, protoc\$gen\$openapiv2_options_annotations_pb);/d" {} \;
	find ./gen/js -name 'swagger*' -type f -exec rm {} \;
	rm -rf ./gen/js
	