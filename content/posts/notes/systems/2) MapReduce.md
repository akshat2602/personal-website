---
title: MapReduce
date: 30 Aug 2023
tags: [subject/distributed-system, distributed-computing]
---

# TIL:

**POSIX**

## Reason for MapReduce(Why use a distributed system?)

1. Lots of data(1 PetaByte), machines had 160 GB storage so can't process in one machine
2. I/O speed is very low(performance)
3. Fault tolerance(Tolerate machine and disk failures)
4. Application programmers don't need to work on systems and making sure their job is running in a distributed fashion

## Workflow

| Input      | Map                    | Reduce              | Result |
| ---------- | ---------------------- | ------------------- | ------ |
| key, value | K1 <v1, v1\`, v1\`\` > | K1 R(v1, v1\`, ...) |        |
|            | K2, v2                 |                     |        |
|            | K3, v3                 |                     |        |

### Examples

1. Word Count
2. Sort
3. Reverse Links(Used for page ranking)
    1. Input - key: URL, value: HTML
    2. Map - (<target1, src1>, <target2, src2>,....)
    3. Reduce - target1, <src1, src1\`,....>
    4. Result - rank

## Why MapReduce was developed as a library and not as a service?

Because loading a user defined function into a already existing running process was not possible in case of C++. It is easier to compile the code with the library than dynamically loading the UDF(user defined function).
