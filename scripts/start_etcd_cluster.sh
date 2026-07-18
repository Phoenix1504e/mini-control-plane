#!/bin/bash
# Clean up existing data
rm -rf /tmp/infra*

# Node 1
etcd --name infra1 --data-dir /tmp/infra1.etcd \
  --initial-advertise-peer-urls http://127.0.0.1:2380 \
  --listen-peer-urls http://127.0.0.1:2380 \
  --advertise-client-urls http://127.0.0.1:2379 \
  --listen-client-urls http://127.0.0.1:2379 \
  --initial-cluster infra1=http://127.0.0.1:2380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380 \
  --initial-cluster-token etcd-cluster-1 --initial-cluster-state new &

# Node 2
etcd --name infra2 --data-dir /tmp/infra2.etcd \
  --initial-advertise-peer-urls http://127.0.0.1:22380 \
  --listen-peer-urls http://127.0.0.1:22380 \
  --advertise-client-urls http://127.0.0.1:22379 \
  --listen-client-urls http://127.0.0.1:22379 \
  --initial-cluster infra1=http://127.0.0.1:2380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380 \
  --initial-cluster-token etcd-cluster-1 --initial-cluster-state new &

# Node 3
etcd --name infra3 --data-dir /tmp/infra3.etcd \
  --initial-advertise-peer-urls http://127.0.0.1:32380 \
  --listen-peer-urls http://127.0.0.1:32380 \
  --advertise-client-urls http://127.0.0.1:32379 \
  --listen-client-urls http://127.0.0.1:32379 \
  --initial-cluster infra1=http://127.0.0.1:2380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380 \
  --initial-cluster-token etcd-cluster-1 --initial-cluster-state new &
