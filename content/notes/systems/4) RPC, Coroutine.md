---
title: RPC, Couroutines
date: 11 Sep 2023
tags: [subject/distributed-system]
---

# TIL

-   **htop**
-   C100k problem
-   Timer interrupt(A periodic interrupt that allows the OS to stop the current thread's execution and schedule a different thread)

# Coroutine Notes

## 1-thread and n-thread:

### Advantages

| 1 Thread           | n Threads                               |
| ------------------ | --------------------------------------- |
| Memory/Performance | Fairness in scheduling                  |
| No locking         | Easier to understand and write the code |

### Disadvantages

| 1 Thread       | n Threads           |
| -------------- | ------------------- |
| Stack Ripping  | Locking             |
| Debugging hard | Hard on performance |

## Fibre vs Coroutine vs Async:

-   Fibre/Coroutines have their own stacks(benefits like stack trace/shared variables).
-   Async is stackless. Also called as stackless coroutines.

Stackless coroutines are implemented in - JS, C#, Rust, C++ 20
Stackfull are implemented in - Golang, C++

# RPC Notes

Concepts like -

1. At most once
2. At least once
3. Exactly once(technically impossible)
