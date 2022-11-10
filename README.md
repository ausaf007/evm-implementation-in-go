<h1 align="center">Ethereum Virtual Machine</h1>

<h3 align="center">Implementation Of Ethereum Virtual Machine With Limited Instruction Set. </h3>

<!-- TABLE OF CONTENTS -->
<details open>
  <summary>Table of Contents</summary>
  <ul>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#tech-stack">Tech Stack</a></li>
    <li><a href="#prerequisites">Prerequisites</a></li>
    <li><a href="#how-to-use">How to use?</a></li>
  </ul>
</details>

## About The Project

Implementation of Ethereum Virtual Machine in Golang, with the following instructions supported: 
- PUSH1/PUSH2/PUSH3/PUSH32
- MSTORE/MSTORE8 
- ADD/MUL/SDIV/EXP 

## Tech Stack

[![](https://img.shields.io/badge/Built_with-Go-green?style=for-the-badge&logo=Go)](https://go.dev/)

## Prerequisites

Download and install [Golang 1.19](https://go.dev/doc/install) (or higher).  

## How To Use?

1. Navigate to `task-2-ethereum-vm/`:
   ``` 
   cd /path/to/folder/task-2-ethereum-vm/
   ```
2. Get dependencies:
   ``` 
   go mod tidy
   ```
3. Run the app:
   ``` 
   go run main.go 
   # use "--verbose" flag to get additional logs
   go run main.go --verbose 
   ```
4. CD into `evm/` to run tests:
   ``` 
   cd evm/
   go test
   ```

Thank you!