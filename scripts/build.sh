./scripts/clean.sh 2>/dev/null
P="github.com/jialin-li/EasyDB"
OUTPUT="bin"
if [ ${#} -ne 1 ]; then
    go build -o ${OUTPUT}/master ${P}/master 
else
    go build -o ${OUTPUT}/master -tags ${1} ${P}/master 
fi
go build -o ${OUTPUT}/client ${P}/client
go build -o ${OUTPUT}/server ${P}/server
