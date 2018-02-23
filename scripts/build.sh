./scripts/clean.sh 2>/dev/null
P="github.com/jialin-li/EasyDB"
OUTPUT="bin"

go build -o ${OUTPUT}/master ${1} ${2} ${P}/master 
go build -o ${OUTPUT}/client ${1} ${2} ${P}/client
go build -o ${OUTPUT}/server ${1} ${2} ${P}/server
