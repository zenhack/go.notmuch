.PHONY: test ci fmtcheck

ci: test fmtcheck

test: fixtures/database-v1
	@go test -v -race -cover ./...
	@rm -rf fixtures/database-v1

fmtcheck:
	# Verify that everything is properly gofmt'd.
	@[ -z "$$(gofmt -d .)" ] || ( \
		gofmt -d . >&2; \
		echo "Formatting descrepency; did you forget to run gofmt?" >&2; \
		exit 1 \
	)


fixtures/database-v1: fixtures/database-v1.tar.xz
	@tar -C fixtures -xf fixtures/database-v1.tar.xz

fixtures/database-v1.tar.xz:
	@mkdir -p fixtures
	@curl -Lo fixtures/database-v1.tar.xz http://notmuchmail.org/releases/test-databases/database-v1.tar.xz
