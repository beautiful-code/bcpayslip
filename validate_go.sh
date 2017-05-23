#!/usr/bin/env bash
go_files=`find . -name "*.go" -not -path "./vendor/*"`
exit_status=0

for go_file in $go_files
do
  output=$(golint $go_file)
  if [ "$output" != "" ]; then
  echo $go_file
  echo "$output"
  echo
  fi
done

exit 1
