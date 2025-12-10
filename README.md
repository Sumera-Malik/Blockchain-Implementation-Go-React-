# Blockchain Implementation (Go + React) + MetaMask Sepolia Test Transaction

This repository contains two main parts of assignment:

✔ A simple blockchain implemented in Go with Proof-of-Work and a React frontend

✔ A MetaMask Sepolia test transaction demonstrating wallet usage, testnet ETH transfer, and Etherscan verification

**Part 1 — Blockchain Implementation (Go + React)**

Developed a basic blockchain using Go for the backend and React for the frontend interface.

The blockchain supports:

✔ Adding transactions

✔ Mining new blocks using Proof-of-Work

✔ Viewing all blocks

✔ Searching transactions by keyword

✔ Automatic React UI refresh after actions

**Features**

Genesis block includes: Sumera Malik — Roll No: 21i-1579

REST API endpoints:

/view – View full blockchain

/tx – Add a transaction

/mine – Mine a new block

/search – Search for a keyword

React UI allows:

✔ Adding transactions

✔ Mining blocks

✔ Searching blockchain data

✔ Viewing block hashes, previous hashes, nonce, and data

✔ React dashboard showing blockchain blocks 

✔ Mining and transaction forms

✔ Search interface

✔ Blockchain output 

**Part 2 — MetaMask + Sepolia Test Transaction**

For this part, a real test transaction on the Sepolia Ethereum Testnet.

**Steps Completed**

Installed MetaMask and set up Account 1

Enabled Sepolia Test Network

Claimed test ETH from a Sepolia PoW faucet

Created Account 2

Sent 0.01579 SepoliaETH from Account 1 → Account 2

Verified transaction as Confirmed on Sepolia Etherscan

Verified that Account 2 received the exact amount

Faucet reward page

MetaMask transaction page

Send page showing 0.01579 SepoliaETH

Confirmed transaction in Activity tab

Account 2 balance increase

All of this is visible in the uploaded report. 


**Conclusion**

Both parts of the assignment were successfully completed:

✔ The Go + React blockchain works fully with mining, transactions, viewing, and searching

✔ The Sepolia transaction succeeded and was verified on Etherscan

Screenshots demonstrate the entire process

cd backend-go
go run main.go

Frontend (React)
cd frontend-react
npm install
npm start
