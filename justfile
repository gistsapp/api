build:
  go build -o api -v

test-all:
  cd test && go test

report-all:
  #!/usr/bin/env sh
  for file in gists storage storage server user utils organizations; do
    just report-test $file
  done

report-test PACKAGE:
  cd test && go test -coverprofile=cov-{{PACKAGE}}.out -coverpkg=./../{{PACKAGE}} && go tool cover -html=cov-{{PACKAGE}}.out -o {{PACKAGE}}-coverage.html && rm cov-{{PACKAGE}}.out

test TEST:
  go test ./tests/{{TEST}} -v

migrate: build
  ./api migrate

dev:
  air

