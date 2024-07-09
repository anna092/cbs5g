## CBCF (Cell Broadcast Center Function)
CBCF (Cell Broadcast Centre Function) is a network function in the 5G core network that responsible for all kind of broadcast message, including PWS (Public Warning System). CBCF is the same with CBC in 5G, but already implemented 5G core concept, which is the SBA (Service Based Arhitecture). This repository is the implementation of the CBCF by using Go language. 

## Features:
* Receive the warning message in CAP (Common Alerting Protocol)
* Do the NonUeN2messageTransfer to transfer write replace warning request to the AMF

## Usage
1. Clone github repository
```
git clone https://github.com/anna092/5Gcbs.git
```

2. Download all the dependency 
```
go get -t
```

3. Run the Program
```
./CBCF
```

## CBE Simulator
This repository also contains CBE (Cell Broadcast Entity) simulator written in python. By running that file, we can send a warning message to the CBCF to be forwarded into the AMF. For now, the message is still hardcoded, but in the future, the program will also receive parameter from input to create a dynamic message. 
To run the program
```
python3 cbe.py --messageId="id"
```
