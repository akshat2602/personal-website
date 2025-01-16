---
title: Eventual Consistency
date: 18 Oct 2023
tags: [subject/distributed-system, distributed-storage, consistency]
---

# Notes

Dynamo is similar to a Distributed Hash Table meaning it uses [consistent hashing](https://www.youtube.com/watch?v=zaRkONvyGr8&pp=ygUSY29uc2lzdGVudCBoYXNoaW5n).

Dynamo is called a zero-hop DHT because each node has enough information about the whole **consistent hashing ring** that it can directly transfer any client to the server responsible for serving that key.

## Causes for the system not being linearizable

1. Initializing replication factor but not having the same number of writes/reads
2. Membership changes
    1. Discovery of membership changes is at different times, due to gossip protocol, causing writes to happen at different servers than intended.
    2. When key exchanges are happening and data is being written to the old server.

## Vector Clock
