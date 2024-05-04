systemctl start docker
mkdir /sys/fs/cgroup/memory/JOJ-judger

go run main.go