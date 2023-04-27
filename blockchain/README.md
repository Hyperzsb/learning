# Blockchain

> "A blockchain is a distributed ledger with growing lists of records (blocks) that are securely linked together via cryptographic hashes.
> Each block contains a cryptographic hash of the previous block, a timestamp, and transaction data (generally represented as a Merkle tree, where data nodes are represented by leaves). 
> The timestamp proves that the transaction data existed when the block was created. 
> Since each block contains information about the previous block, they effectively form a chain (compare linked list data structure), with each additional block linking to the ones before it. 
> Consequently, blockchain transactions are irreversible in that, once they are recorded, the data in any given block cannot be altered retroactively without altering all subsequent blocks."
> 
> *From [Wikipedia](https://en.wikipedia.org/wiki/Blockchain)*

## Table of Contents

*This table of contents is for this directory of this repo, not README itself. Anyway, this README is not too long.*

- [**erc3525**](erc3525): This is a Web 3 community credit (certification) system built upon the Semi-Fungible Token, which is proposed in the ERC-3525 standard.
- [**hardhat**](hardhat): This directory contains [Hardhat](https://hardhat.org/) related tutorials and resources.
- [**truffle**](truffle): This directory contains [Truffle Suite](https://trufflesuite.com/) related tutorials and resources.

Blockchain, the key infrastructure of Web 3 ecology, has been attracting more and more attention in recent years.
As a passionate developer, learning Web 3 development is a really intriguing thing for me, but it still took some time to get started.
Fortunately, thanks to the project of my Master's degree, I finally have to learn it and start my journey in the world of Web 3.

## Frameworks

Due to the finanical nature of Blockchain and Web 3, developing and testing your code on the Mainnet of Ethereum and other production environment is extremely expensive.
Additionally, we also need basic SDKs and tools to facilitate our development process.
That is why we need frameworks.

There are some famous and widely-used frameworks you can try, including:

- [Hardhat](https://hardhat.org/): Hardhat is a development environment to compile, deploy, test, and debug your Ethereum software. Get Solidity stack traces & console.log.
- [Truffle Suite](https://trufflesuite.com/): Truffle is the most comprehensive suite of tools for smart contract development.

## Learning Resources

From my point of view, the technology stack of Web 3 is relatively simple compared with the traditional Web 2 (at least for the newbie).
So I think the official docs of [Ethereum](https://ethereum.org/en/developers/docs/) and [Solidity](https://docs.soliditylang.org/en/latest/) are really helpful.

Additionally, [Solidity by Example](https://solidity-by-example.org/) is also a great place where you can learn fundamental concepts about Solidity and more.

To make the boring learning more fun, there are many game-based tutorials available online.
For example, [Cryptozombies](https://cryptozombies.io/) teaches you various basic concepts about smart contracts and Web 3 applications by building a collecting game step by step.
