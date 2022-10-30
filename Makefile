CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/paksergey94/telegram-bot/cmd/bot

dev:
	go run ${PACKAGE} -dev

prod:
	mkdir -p data
	go run ${PACKAGE} 2>&1 | tee data/log.txt

all: format build test lint

build: bindir
	go build -o ${BINDIR}/bot ${PACKAGE}

test:
	go test ./...

.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

run:
	go run ${PACKAGE}

generate: install-mockgen
	${MOCKGEN} -source=internal/service/messages/incoming_msg.go -destination=internal/service/messages/mocks/incoming_msg.go
	${MOCKGEN} -source=internal/repository/currency_rate/repository.go -destination=internal/repository/currency_rate/mocks/repository.go
	${MOCKGEN} -source=internal/repository/spend/repository.go -destination=internal/repository/spend/mocks/repository.go
	${MOCKGEN} -source=internal/repository/selected_currency/repository.go -destination=internal/repository/selected_currency/mocks/repository.go
	${MOCKGEN} -source=internal/clients/currency_rate/currency_rate.go -destination=internal/clients/currency_rate/mocks/currency_rate.go
	${MOCKGEN} -source=internal/clients/tg/tgclient.go -destination=internal/clients/tg/mocks/tgclient.go
	${MOCKGEN} -source=internal/service/report/report.go -destination=internal/service/report/mocks/report.go
	${MOCKGEN} -source=internal/database/manager.go -destination=internal/database/mocks/manager.go
	${MOCKGEN} -source=internal/repository/budget/repository.go -destination=internal/repository/budget/mocks/repository.go

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

docker-run:
	sudo docker compose up
