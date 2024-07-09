## CBCF (Cell Broadcast Center Function)
CBCF (Cell Broadcast Centre Function) is a network function in the 5G core network that responsible for all kind of broadcast message, including PWS (Public Warning System). 

## Features:
* Receive UE registration status
* Handle the warning message in CAP (Common Alerting Protocol)
* Handle disaster prevention commands for equipment
* Do the NonUeN2messageTransfer to transfer Write-Replace-Warning Request to the AMF

## Usage
1. Clone github repository
```
git clone https://github.com/anna092/cbs5g.git
```

2. Download all the dependency 
```
go get -t
```

3. Run CBCF Service 
```
./CBCF
```

## CBE Agent
This repository also contains CBE (Cell Broadcast Entity) CAP Alert Agent written in python. 
By running that file, we can send a alert message to the CBCF to be forwarded into the AMF. 
Currently, the alert message is hardcoded for testing.
To run the program
```
python3 CBE.py --messageId="id"
```
