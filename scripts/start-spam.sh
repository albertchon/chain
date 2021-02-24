rm /tmp/stop.spam

for run in {1..5}; do
    ./scripts/spam.sh requester$run &
done

wait