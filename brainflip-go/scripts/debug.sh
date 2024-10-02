rm -rf debug/*
go run bf.go -i ../bfcheck/prog-$1.b < ../bfcheck/input.dat > debug/output.dat
cp ../bfcheck/output-$1.dat debug
python transform.py debug/output.dat > debug/output.data
python transform.py debug/output-$1.dat > debug/output-$1.data
