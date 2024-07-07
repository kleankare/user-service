.PHONY: compose skaffold test

compose:
	cd ../deployments && docker-compose up --build

skaffold:
	skaffold dev -f ../deployments/skaffold.yaml

test:
	go test -v ../test/...