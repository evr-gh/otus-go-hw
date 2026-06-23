#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-envdir

export HELLO="SHOULD_REPLACE"
export FOO="SHOULD_REPLACE"
export UNSET="SHOULD_REMOVE"
export ADDED="from original env"
export EMPTY="SHOULD_BE_EMPTY"

result=$(./go-envdir "$(pwd)/testdata/env" "/bin/bash" "$(pwd)/testdata/echo.sh" arg1=1 arg2=2)
expected='HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2'

[ "${result}" = "${expected}" ] || (echo -e "invalid output: ${result}" && exit 1)

echo "-- ADDITIONAL TEST 1 ----------------------------------------------------------------------------------"
result=$(./go-envdir "$(pwd)/testdata/env2" "/bin/bash" "$(pwd)/testdata/echo.sh" arg1=1 arg2=2 arg3 arg4=4; echo $?)
expected='HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2 arg3 arg4=4
0'

[ "${result}" = "${expected}" ] || (echo -e "invalid output: ${result}" && exit 1)

echo "-- ADDITIONAL TEST 2 ----------------------------------------------------------------------------------"
result=$(./go-envdir "$(pwd)/testdata/env_err" "/bin/bash" "$(pwd)/testdata/echo.sh" arg1=1 arg2=2; echo $?)
expected='Error executing command: invalid env variable name (ENV=PARAM)
255'

[ "${result}" = "${expected}" ] || (echo -e "invalid output: ${result}" && exit 1)


echo "-- ADDITIONAL TEST 3 ----------------------------------------------------------------------------------"
result=$(./go-envdir "/no_exist" "/bin/bash" "$(pwd)/testdata/echo.sh" arg1=1 arg2=2; echo $?)
expected='Error executing command: open /no_exist: no such file or directory
255'

[ "${result}" = "${expected}" ] || (echo -e "invalid output: ${result}" && exit 1)


echo "-- ADDITIONAL TEST 4 ----------------------------------------------------------------------------------"
result=$(./go-envdir "$(pwd)/testdata/env"; echo $?)
expected='Not all arguments are set.
Usage: go-envdir <path to env dir> <command> (<argbment>)*
255'

[ "${result}" = "${expected}" ] || (echo -e "invalid output: ${result}" && exit 1)


echo "-- ADDITIONAL TEST 5 ----------------------------------------------------------------------------------"
result=$(./go-envdir "$(pwd)/testdata/env" "not_exist"  arg1=1 arg2=2; echo $?)
expected='Error executing command: exec: "not_exist": executable file not found in $PATH
255'

[ "${result}" = "${expected}" ] || (echo -e "invalid output: ${result}" && exit 1)


echo "-- ADDITIONAL TEST 6 ----------------------------------------------------------------------------------"
result=$(./go-envdir "$(pwd)/testdata/env" "/bin/bash"  "$(pwd)/testdata/err.sh"  arg1=1 arg2=2 2>err;  echo $?; cat err; rm -f err)
expected='start executing
127
error executing'

[ "${result}" = "${expected}" ] || (echo -e "invalid output: ${result}" && exit 1)


rm -f go-envdir
echo "PASS"
