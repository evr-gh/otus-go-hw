#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-cp

./go-cp -from testdata/input.txt -to out.txt
cmp out.txt testdata/out_offset0_limit0.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10
cmp out.txt testdata/out_offset0_limit10.txt

./go-cp -from testdata/input.txt -to out.txt -limit 1000
cmp out.txt testdata/out_offset0_limit1000.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10000
cmp out.txt testdata/out_offset0_limit10000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 100 -limit 1000
cmp out.txt testdata/out_offset100_limit1000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 6000 -limit 1000
cmp out.txt testdata/out_offset6000_limit1000.txt


./go-cp -from testdata/no_exist_input.txt -to out.txt > err.txt
cmp err.txt testdata/err2.txt

./go-cp -from testdata/same_input.txt -to testdata/same_input.txt > err.txt
cmp err.txt testdata/err3.txt

./go-cp -from /dev/urandom -to out.txt > err.txt
cmp err.txt testdata/err4.txt

./go-cp -from testdata -to out.txt > err.txt
cmp err.txt testdata/err5.txt

./go-cp -from testdata/input.txt -to /inv/dir/out.txt > err.txt
cmp err.txt testdata/err6.txt

./go-cp -from testdata/input.txt -to out.txt -offset 7000 -limit 1000 > err.txt
cmp err.txt testdata/err7.txt


./go-cp -from testdata/tolstoy_voyna-i-mir.txt -to out.txt > err.txt
cmp out.txt testdata/tolstoy_voyna-i-mir.txt
cmp err.txt testdata/empty_err.txt

./go-cp -from testdata/tolstoy_voyna-i-mir.txt -to out.txt -offset 5350525 -limit 1000 > err.txt
cmp out.txt testdata/empty_out.txt
cmp err.txt testdata/empty_err.txt

./go-cp -from testdata/tolstoy_voyna-i-mir.txt -to out.txt -limit -1 > err.txt
cmp out.txt testdata/tolstoy_voyna-i-mir.txt
cmp err.txt testdata/empty_err.txt

rm -f go-cp out.txt
rm -f go-cp err.txt
echo "PASS"
