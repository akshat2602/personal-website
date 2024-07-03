---
title: Percolator
date: 06 Nov 2023
tags: [subject/distributed-system, distributed-computing]
---

# Notes

Bigtable doesn't support multi-row/multi-table transactions.

> Why does Google need multi-table transactions?
> Removing duplicates(Multiple URLs may lead to the same website), calculation of pagerank will get affected.

Built on top of big table, because didn't have that many people working on it and also didn't have source code access to big table.

## Locks

Locks in percolator could have been implemented in two ways -

1. In place(in database)
    1. Problem with this is that you can't maintain complex locks with queues and techniques like wound wait, wait die, etc.
    2. Database overhead
2. Standalone lock service

> Why is Snapshot Isolation used instead of Optimistic Concurrency Control?
> Because <mark style="background: #FFF3A3A6;">OCC requires complex locking mechanisms</mark> like queues to hold pending locks.

Percolator has a primary lock scheme where all the locks point to a single key that had initiated the transaction.
