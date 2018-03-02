#!/bin/bash

. scripts/util.sh

./scripts/clean.sh 2>/dev/null || true

# does not work yet with args
#if [ ${#} -ne 0 ]; then
    #go build -o ${OUTPUT}/master ${P}/master 
    #exit 0
#fi

mkdir output/make -p
mkdir output/release -p
cd bin/

#set -x

for t in $(find ../tests/make -type f -name "*.txt" | grep -v "_out"); do
    dir=$(echo $t | awk -F'/' '{print $(NF-1)}')
    file=$(echo $t | awk -F'/' '{print $NF}')
    name=$(echo $file | awk -F'.' '{print $1}')
    #echo $t | awk -F'/' '{print "Test: "$(NF-1) "\nFile: "$NF}'
    if [ "$dir" == "performance" ]; then
        light_cyan "RUNNING TIMING TEST $file"
        time cat $t | ./master >/dev/null
    else

        light_cyan "RUNNING $dir/$file"
        cat $t | ./master > ../output/make/${name}.output
        DIFF=$(diff ../tests/make/${dir}/${name}_out.txt ../output/make/${name}.output)

        if [ "$DIFF" != "" ]; then
            red "==========FAIL=========="
            diff -U 100000 ../tests/make/${dir}/${name}_out.txt ../output/make/${name}.output
        else
            green "==========PASS=========="
        fi
    fi
        ./../scripts/clean.sh 2>/dev/null || true
        #echo "============================================================"
done

