SHELL=/bin/bash
CMD=`ls -l --time-style="long-iso" ./cmd/ | egrep '^d' | awk '{print $$8}'`

build-windows:
	@for d in $(CMD) ; \
	do \
	  echo "Building $$d"; \
	  $$(cd ./cmd/"$$d" && GOOS=windows GOARCH=amd64 go build -o "$$d".exe "$$d".go); \
  	done

.PHONY: build-windows
