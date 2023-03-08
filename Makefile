VERSION := 2.4.0
VERSION_LDFLAGS="-X=go.1password.io/eventsapibeat/version.Version=$(VERSION)"

.PHONY: eventsapibeat
eventsapibeat: ## Builds the elastic beats binary for the current os and arch
	@go build -ldflags $(VERSION_LDFLAGS) -mod vendor -o ./bin/eventsapibeat main.go

.PHONY: clean ## Cleans the binaries directory
clean:
	@rm -rf bin/*

.PHONY: build_all_apps
build_all_apps: clean ## Clean, builds then packages all elastic beats binaries
	@go mod download
	@go mod vendor
	@gox -arch="amd64 arm arm64" -os="linux windows freebsd openbsd" -osarch="darwin/amd64" -output="bin/{{.OS}}_{{.Arch}}/eventsapibeat" -ldflags "-s $(VERSION_LDFLAGS)" .
	@for d in bin/*/; do cp -a eventsapibeat-sample.yml $${d}; cp -a logstash-sample.conf $${d}; done
	@cd bin && for d in */; do \
  		COPYFILE_DISABLE=1 tar --exclude='.DS_Store' --exclude='.gitignore' --exclude='.travis.yml' -cvzf "eventsapibeat_$(VERSION)_$${d%/}.tar.gz" $${d}; \
  		rm -rf $${d}; \
  	done
