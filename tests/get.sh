#!/usr/bin/env bash

errcho() {
    (>&2 echo -e "\e[31m$1\e[0m")
}


for url in 'http://api.ipify.org' 'http://google.com' 'google.com' 'google.com:80'; do
  # api.ipify.org can be reeeally slow at times, long tail
    out=$(justget --timeout=10s "$url")
    if [[ $? -ne 0 ]]; then
      FAILED=1
      errcho "Failed url test: ${url}"
    fi

done

if [[ -n "$FAILED" ]]; then
  echo "Failed one or more tests"
  exit 1
fi

echo "PASS!"