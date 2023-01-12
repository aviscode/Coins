# **Coins**

## **Pricing-service**
This service is responsible for fetching USD prices for 3 tokens (BTC, ETH, USDT) every minute from Coinmarketcap API and storing them locally.

this service exposes a grpc endpoint so clients can request coins prices.

## **Client-Service**
This service is responsible for fetching USD prices for 3 tokens (BTC, ETH, USDT) every minute from Pricing-service API and printing it the stdout.


## **Usage**
clone the repo and run ```make``` this will get all dependencies and build the 2 binaries

after running this u will see server and client binary 

just run the binaries and see the output

if needed server adders and api-key can be change via flags just run ```--help```
if needed client adders can be change via flags just run ```--help```

## **Output**
the client out put will look like:
```bash
INFO[2580] BTC: 18854.566030, 18842.675561, -0.06%      
INFO[2580] ETH: 1427.584936, 1426.820354, -0.05%        
INFO[2580] USDT: 1.000098, 1.000102, 0.00%            
```

## **Note**
Make sure the server is running at the same time to get the desired output.