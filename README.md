# Blockchain App Project
update: May 16th, 2019  
  
## • Proposal  
  
***What***: Database for sharing economy  
  
***Why***: Users of the service cannot lie their history of transactions because they are recorded in the blockchain, which is immutable. It is less expensive for users to use the service than a centralized system because there is no middleman. Also, decentralized system has no single point of failure, which means no server down time.  
  
***How***: Users can share their assets, such as bike, room, without middlemen. Lenders put their bike rental advertisement in the blockchain. Borrowers send Start Request to use their bike and deposit a certain amount of money. When the lenders get the Start Request, they record Start Time and send the KEY to use the bike to the borrower’s smartphone. If the borrowers want to stop using the bike, then they send Stop Request. When the lenders get the Stop Request, they record Stop Time, calculate the fee, and send the remaining money back to the borrower’s smartphone.  
  
  
## • Functionalities  
  
- User can submit advertisement.  
- User can submit start request to start using the service and deposit money.  
- Miners can earn TX fees.  
- User ID, the balance and the start time are recorded.  
- User can submit stop request to stop using the service.  
- The end time is recorded, and the fee is calculated.  
  
  
## • Algorithm  
  
1. [Publish] User (Lender) creates a ***Rental Info*** and an initialized ***History*** as a transaction, and sends it to the Peer to Peer (P2P) network.  
2. Miners receive the transaction, verify the transaction, and do Proof of Work (PoW).  
3. If a miner finds a nonce, then the miner creates a block and sends it to the P2P network.  
4. Miners or Users verify the nonce and add the block to the blockchain.  
5. [DisplayAds] User (Borrower) sees the ad in the blockchain.  
6. [SendStartRequest] Borrower creates the updated ***History (Start Request)*** with the rental ID and deposit as a transaction, and sends it to the P2P network.  
7. Miners do 2 ~ 4.  
8. [Permit] Lender checks the new block, and if it is for their bike, the Lender creates the updated ***History*** as a transaction to record the Borrower's ID, the deposit, and the start time. The Lender also creates the updated ***Rental Info*** as a transaction to change the availability flag to false, and sends them to the network.  
9. Miners do 2 ~ 4.  
10. The start time is confirmed. (And the Lender sends the key to use the bike to the Borrower.)  
11. [SendStopRequest] Borrower creates the updated ***History (Stop Request)*** with the rental ID as a transaction, and sends it to the network.  
12. Miners do 2 ~ 4.  
13. [Checkout] Lender checks the new block, and if it is for their bike, the Lender creates the updated ***History*** as a transaction to record the end time, calculate the fee, and also creates the updated ***Rental Info*** as a transaction to change the availability flag back to true.  
14. Miners do 2 ~ 4.  
15. The end time and the fee are confirmed. (And Lender sends remaining money back to the Borrower.)  
  
  
## • Data  
  
- ***Rental Info*** includes rental ID, lender ID, asking price, availability flag, and signature.  
- ***History*** includes history ID, rental ID, borrower ID, start time, end time, deposit, state, and signature.  
The 'state' is used to call the step 8, 10, 13, and 15.   
  
## • Reference  
  
- Slock (https://slock.it/landing.html)  
  
  
  
  
  
