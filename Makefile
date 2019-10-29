.PHONY: all build run

all: build run

build:
	@echo "+ $@"
	./scripts/build
build-binary:
	@echo "+ $@"
	./scripts/build-binary
run:
	@echo "+ $@"
	./scripts/run
clean:
	@echo "+ $@"
	@running=$$(docker ps -q -f "label=detection.test=attacker"); [ -z "$$running"  ] || ( docker stop $$running || docker kill $$running || true )
