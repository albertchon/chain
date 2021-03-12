rm /tmp/stop.spam

for run in {1..10}; do
    sleep 20;
    ./scripts/spam.sh $run &
done

wait