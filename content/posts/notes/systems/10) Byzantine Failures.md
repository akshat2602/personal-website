---
title: 10) Byzantine Failures
date: 04 Oct 2023
tags: [subject/distributed-system, research/failures]
---

Readings - [Practical Byzantine Fault Tolerance.pdf](https://pmg.csail.mit.edu/papers/osdi99.pdf) [Byzantine PAXOS.pdf](https://lamport.azurewebsites.net/tla/byzsimple.pdf)

# Notes

## What can be potential byzantine failures?

1. Wrong replies which cannot be detected
2. Halt consensus(same as crash so not really byzantine)
3. Internal state mismanagement(same as crash)
4. Forgery(Use public/private keys to encrypt data)

## Converting paxos to byzantine paxos

### Strategy 1(Outvoting)

1. Increase the number of servers and increase quorum requirement to
    1. N = 3f + 1, Q = 2f + 1
2. Normal paxos wouldn't work where let's say-

> Thought process of increasing servers:
> There are 4 servers. S1, S2, S3 are normal servers and S4 is a malicious one.
> Let's say S1 prepares and accepts a value with itself S1, S3 and S4. v1 is decided for number 1.
> S1 goes down and S2 comes up. S2 tries to prepare a new number and receives 1,v1 from S3 but 2,v2 from S4. How do you decide which is the correct value?
> We outvote it meaning we'll have to increase number of good servers.

Increase N to 4f + 1 and Q to 3f + 1 so that there is at least f+1 good servers.
This isn't ideal though because only f faults can be tolerated in this.

### Strategy 2(Detecting a lie)

1. P2b and P1a malicious replies are kind of redundant.
2. P1b malicious replies can be detected by asking each reply to attach the previous message signature and then it can be checked.
3. P2a can also attach previous message signature and then it can be checked. In case of the very first round of P2a we might also need to let each server talk to each other as there is no previous chain of messages to verify from. So servers talk with each other to get a quorum before replying back to the malicious leader.
