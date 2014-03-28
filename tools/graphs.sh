#!/bin/sh
write_script() {
cat << EOF > /tmp/grypt_bench.gnuplot
set term png size 1024,768
set output "${1}.png"
set xlabel "Size (K)"
set multiplot layout 3,1 title "$1"
set tmargin 1
set ylabel 'ns/op'
plot '/tmp/grypt_bench.gnuplot.dat' using 2:3 title '' with linespoints
set ylabel 'B/op'
plot '/tmp/grypt_bench.gnuplot.dat' using 2:4 title '' with linespoints
set ylabel 'allocs/op'
plot '/tmp/grypt_bench.gnuplot.dat' using 2:5 title '' with linespoints
unset multiplot
quit
EOF
}

if [ -z "$1" ]; then
	tests="."
else
	tests="$1"
fi

go test -bench "$tests" git.polydawn.net/hank/grypt | sed '1d;$d' |
awk '{
	sz=substr($1,match($1,/..$/))
	sub(/K/,"",sz); sub(/M/,"*1024",sz)
	print sz |& "bc"; "bc"|&getline sz
	sub(/..$/,"",$1)
	print $1" "sz" "$3" "$5" "$7
} ' > /tmp/grypt_bench.dat

for b in $(awk '{print $1}' < /tmp/grypt_bench.dat | uniq | egrep "$tests" | paste -s -d \\t); do
	write_script "$b"
	fgrep "$b" < /tmp/grypt_bench.dat > /tmp/grypt_bench.gnuplot.dat
	gnuplot /tmp/grypt_bench.gnuplot
done
