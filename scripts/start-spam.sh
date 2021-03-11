rm /tmp/stop.spam

for run in {1..6}; do
    ./scripts/spam.sh $run &
done

wait