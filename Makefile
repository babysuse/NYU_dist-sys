.PHONY: start-cli start-srv init clean
start-cli:
	cd frontend && npm start
stop-cli:
	npx kill-port 3000
# go run web/server.go web/account.go web/operations.go
# ./raft --id 1 --cluster http://127.0.0.1:16048,http://127.0.0.1:16058,http://127.0.0.1:16068 --port 16049
# ./raft --id 2 --cluster http://127.0.0.1:16048,http://127.0.0.1:16058,http://127.0.0.1:16068 --port 16059
# ./raft --id 3 --cluster http://127.0.0.1:16048,http://127.0.0.1:16058,http://127.0.0.1:16068 --port 16069
start-srv:
	go run web/server.go web/account.go web/operations.go
stop-srv:
	npx kill-port 16048
	npx kill-port 16058
	npx kill-port 16068
init:
	cd internal/pkg/raft/src/go.etcd.io/etcd/contrib/raftexample/ && go build -o ../../../../../../../../raft
clean:
	rm -rf raftexample*
