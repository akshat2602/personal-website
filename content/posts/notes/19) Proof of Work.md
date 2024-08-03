---
title: Proof of Work
date: 13 Nov 2023
tags: [subject/distributed-system]
---
# Notes
## Bitcoin
### Reasoning for making bitcoin
1) Anonymous
2) Decentralized money system
3) No trust in the government money system

> How does this work in a no trust environment?
> 
> Use asymmetric encryption scheme like public/private key pairs.

### Double spend problem
Double spend (make a transaction and then unrecord it and then spend that money again)
#### How do you solve this? 
Similar to the **consensus problem** where everyone must agree on transactions.
We can't use consensus algos like **Paxos** or **RAFT** because its not a closed membership situation(no global view of number of servers.)
#### Solution
Consensus on the basis of compute(ask a proposer to solve a very hard to compute problems).
That computation must have some properties-
1) Hard to compute
2) Easy to verify
3) Should be unique(based on the contents of the transaction)
The math problem is - **HASHING**
#### How do you make this hard to compute?
Make it so that some "n" leading bits should be 0. Time complexity - 2^n
#### How do you make sure that these hashes aren't precomputed? 
Include the hash of the previous block as well so that it needs to be real time.
#### How do you decide the size of n? 
- Store the timestamp in the blocks
- Look at the last 1000 blocks generated and calculate the difference, average it and see if time taken for a new block generation is low or high, increase the value for n or decrease it.
