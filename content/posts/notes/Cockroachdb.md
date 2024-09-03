---
title: CockroachDB, The Resilient Geo-Distributed SQL Database
tags: [subject/distributed-system, distributed-storage, database, paper-summary]
date: 03 Aug 2024
description: Personal notes on the CockroachDB paper that was presented in SIGMOD 2020.
---

These are my personal notes on [CockroachDB paper](https://dl.acm.org/doi/10.1145/3318464.3386134)

## CockroachDB: What is it? 
- CockroachDB is a scalable geo-distributed SQL DBMS. 
- Built using global OLTP workloads in mind with high availability and strong consistency in mind.
- Just like its namesake its supposed to be resilient to disasters(through replication and automatic recovery).


## Requirements 
- Data must be horizontally partitioned(sharding) for ex - GDPR compliance.
- "Always on" experience(fault tolerance).
- Low latency support(data must reside near the users).
- Avoid data anomalies by implementing serializable transactions at the SQL level.


## Features
- Uses multi-version concurrency control
- Fault tolerance by maintaining at minimum three replicas of every partition. High availability maintained using automatic recovery mechanisms when a node failure is detected.
- Geo-distributed partitioning and data placement with automatic and manual data placement strategies.
- Serializable transactions using commodity hardware using hybrid logical clocks for versioning.

> {{< figure src="/data_placement.jpeg" align="center" >}}

## Architecture Overview
- Since, data is partitioned and partitions are replicated, it makes sense to have a RAFT group at the level of partitions(or ranges as the paper calls it). 
- Replication happens at the range level instead of a node level. 
- Layered architecture as illustrated below. Each layer acts as an abstraction to the one above it. 

> {{< figure src="/crdb_layers.jpeg" align="center" >}}

Note: Replication is consensus based, specifically RAFT based.

## Distributed Transaction
- Support serializable transactions as a default.
- Transactions can span across shards/ranges. 
- An SQL gateway(basically a node) is used as the Transaction coordinator which is assigned according to proximity from the client.
- Each transaction has a **transaction record** which specifies the transaction status - **committed**, **aborted**, **pending**, **staging**.
- This transaction record also goes through RAFT, serves to atomically change the visibility of all the intents at once, and is durably stored in the same range as the first write of the transaction.
- All operations in a transaction are provisional until committed(they are accompanied by a pointer to the **transaction record** as metadata).

> {{< figure src="/crdb_transaction_step_1.jpeg" align="center" >}}

> {{< figure src="/crdb_transaction_step_2.jpeg" align="center" >}}

- Above images show an example transaction. 
- As a final step the Gateway node commits the transaction by updating this transaction record from pending to committed via another round of consensus.
- This differs from spanner as spanner acquires locks in its first phase and then applies the changes. 

## Optimized transactions
> {{< figure src="/write_pipeline_parallel_commit.jpeg" align="center" >}}

- This is where the **staged** transaction status is used to indicate that the status of transaction is still unknown but the operations have been pipelined. 

## Concurrency Control
A lot of the text below is just copy pasted from the paper as I don't wanna get something wrong about concurrency(its a heavy and hard to reason about topic as it is).
- CRDB uses an MVCC system and each transaction performs its operations at its commit timestamp which results in a serializable history.
- **Write-read conflicts:** A read running into an uncommitted intent with a lower timestamp will wait for the earlier transaction to finalize. A read running into an uncommitted intent with a higher timestamp ignores the intent and does not need to wait.
- **Read-write conflicts:** A write to a key at timestamp *ta* cannot be performed if there’s already been a read on the same key at a higher timestamp *tb >= ta*. CRDB forces the writing transaction to advance its commit timestamp past *tb*.
- **Write-write conflicts:** A write running into an uncommitted intent with a lower timestamp will wait for the earlier transaction to finalize (similar to write-read conflicts). If it runs into a committed value at a higher timestamp, it advances its timestamp past it (similar to read-write conflicts). Write-write conflicts may also lead to deadlocks in cases where different transactions have written intents in different orders. 
- There's a subtle problem introduced when you advance the commit timestamp. To maintain serializability, read timestamp must also be advanced. But, advancing the read timestamp can cause issues as keys could've gotten updated between time *tb* and *ta*. 
- The transaction coordinator must maintain readset and then check or perform a "read refresh" to validate that timestamp can be advanced. 
- If validation fails then transaction must be retried. 
- Linearizability is also guaranteed per-key by tracking an uncertainty interval(you can sort of model this to true time's [earliest,latest] where consider current time to be earliest and uncertainty interval to be latest - earliest). In case of CRDB, the uncertainty interval is *500ms*. 
- When a transaction encounters a value on a key at a timestamp above its provisional commit timestamp but within its uncertainty interval, it performs an uncertainty restart, moving its provisional commit timestamp above the uncertain value but keeping the upper bound of its uncertainty interval fixed. 
- This corresponds to treating all values in a transaction’s uncertainty window as past writes. As a result, the operations on each key performed by transactions take place in an order consistent with the real time ordering of those transactions hence linearizabilty.

There's some more stuff about how they handle clock skew and what guarantees there hybrid logical clocks provide in the paper that's an interesting read. 


One thing the paper mentions is they had trouble implementing snapshot isolation with strong consistency without the use of pessimistic locking which seems somewhat weird to me because from what I understand about snapshot isolation is that you can implement it using Optimistic schemes. There's a [paper](https://dl.acm.org/doi/10.1145/1071610.1071615) I was linked to regarding this statement that talks about something similar which I plan on checking out.

Some interesting learnings from the devs -
1) RAFT is not as easy as people claim it is. 
	1) Membership changes
	2) Too many RPCs because of shard level replication
2) Initially the database used Postgres wire protocol to make existing clients easily compatible with Crdb but turns out modifications must be made for retry mechanisms and Crdb is exploring introducing Crdb-specific client drivers.

FYI the related work section of the paper is a gold-mine for database architecture research. 


# References
1) CockroachDB: The Resilient Geo-Distributed SQL Database https://dl.acm.org/doi/10.1145/3318464.3386134
2) The resilient, geo-distributed database: A SIGMOD Conference Talk https://www.youtube.com/watch?v=ivVFAd9erfo
3) Cockroach Labs Live: The Architecture of a Geo Distributed Database https://www.youtube.com/watch?v=b3XotjcEmJ0
4) https://muratbuffalo.blogspot.com/2022/03/cockroachdb-resilient-geo-distributed.html