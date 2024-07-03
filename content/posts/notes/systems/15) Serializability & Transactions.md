---
title: Serializability
date: 30 Oct 2023
tags: [subject/distributed-system]
---

# Important Reading to help understand:

https://jepsen.io/consistency

# Notes

A - Atomicity (All or nothing of a transaction)
C - Consistency (Consistent with the constraints on database tables)
I - Isolation
D - Durability (Handling reboots/failures)

## Isolation

3 types of dependencies that can help determine serializability -

1. Write - Write Dependency
2. Write - Read Dependency
3. Read - Write Dependency (Anti dependency)

> Anti Dependency:
> An anti-dependency, also known as write-after-read (WAR), occurs when an instruction requires a value that is later updated. In the following example, instruction 2 anti-depends on instruction 3 — the ordering of these instructions cannot be changed, nor can they be executed in parallel (possibly changing the instruction ordering), as this would affect the final value of A.
>
> 1. B = 3
> 2. A = B + 1
> 3. B = 7

A single operation transaction isn't linearizable because single operation transaction still doesn't guarantee real time ordering which is something that's required for linearizability.

## Concurrency

How to detect deadlock?

1. Lock everything (Doesn't work because unnecessary locks acquired)
2. Timeouts (Timeout if a transaction is blocked beyond a certain amount of time, it is assumed that a deadlock has occurred)
3. Graph detection (Create a graph of transactions, if it is cyclic then deadlock. Problem is performance)
4. Give priorities to transactions (Could be based on TID or some kind of timestamp)
    1. Wound wait
    2. Wait die
