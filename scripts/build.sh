P="github.com/jialin-li/EasyDB"
OUTPUT="bin"
go build -o ${OUTPUT}/master ${P}/master 
go build -o ${OUTPUT}/client ${P}/client
go build -o ${OUTPUT}/server ${P}/server
