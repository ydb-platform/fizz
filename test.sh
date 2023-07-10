#!/bin/bash

set -e

# NOTE: See also docker-compose.yml and database.yml to configure database
# properties.
export MYSQL_PORT=3307
export COCKROACH_PORT=26258

COMPOSE=docker-compose
which docker-compose || COMPOSE="docker compose"

args=$@

function cleanup {
    echo "Cleanup resources..."
    $COMPOSE down
    docker volume prune -f
    find ./tmp -name *.sqlite* -delete || true
}
# defer cleanup, so it will be executed even after premature exit
trap cleanup EXIT

function test {
  export SODA_DIALECT=$1

  echo ""
  echo "######################################################################"
  echo "### Running unit tests for $SODA_DIALECT"
  ./tsoda drop -e $SODA_DIALECT
  ./tsoda create -e $SODA_DIALECT
  ./tsoda migrate -e $SODA_DIALECT -p ./testdata/migrations
  go test -tags sqlite -count=1 $args ./...

  echo ""
  echo "######################################################################"
  echo "### Running e2e tests for $1"
  ./tsoda drop -e $SODA_DIALECT
  ./tsoda create -e $SODA_DIALECT
  pushd testdata/e2e; go test -tags sqlite,e2e -count=1 $args ./...; popd
}


$COMPOSE up --wait

go build -v -tags sqlite -o tsoda ../pop/soda

test "sqlite"
test "postgres"
test "cockroach"
test "mysql"
test "ydb"

# Does not appear to be implemented in pop:
# test "sqlserver"
