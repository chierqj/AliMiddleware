rm -fr bin/main
cd ../src
go build -o main src
mv main ../bin
cd ../bin
./main -alsologtostderr -data_path="../input/data.json"