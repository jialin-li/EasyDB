for pid in $(ps -ef | grep "\./server" | awk '{print $2}'); do kill -9 $pid; done
for pid in $(ps -ef | grep "\./client" | awk '{print $2}'); do kill -9 $pid; done

