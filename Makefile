.PHONY: init clean
init:
	cd internal/pkg/raft/src/go.etcd.io/etcd/contrib/raftexample/ && go build -o raftexample
clean:
	rm -rf raftexample*
