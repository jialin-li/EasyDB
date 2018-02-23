#!/bin/bash

. scripts/util.sh

./scripts/clean.sh 2>/dev/null

# does not work yet with args
#if [ ${#} -ne 0 ]; then
    #go build -o ${OUTPUT}/master ${P}/master 
    #exit 0
#fi

mkdir output -p
cd bin/

#set -x

for t in $(find ../tests -type f -name "*.txt" | grep -v "_out"); do
    dir=$(echo $t | awk -F'/' '{print $(NF-1)}')
    file=$(echo $t | awk -F'/' '{print $NF}')
    name=$(echo $file | awk -F'.' '{print $1}')
    #echo $t | awk -F'/' '{print "Test: "$(NF-1) "\nFile: "$NF}'

    light_cyan "RUNNING $dir/$file"
    cat $t | ./master > ../output/${name}.output
    DIFF=$(diff ../tests/${dir}/${name}_out.txt ../output/${name}.output)

    if [ "$DIFF" != "" ]; then
        red "==========FAIL=========="
        diff -U 100000 ../tests/${dir}/${name}_out.txt ../output/${name}.output
    else
        green "==========PASS=========="
    fi
    ./../scripts/clean.sh 2>/dev/null
    #echo "============================================================"
done

