# My Project

This is an Cell Broadcast Service project.
![PWS_thesis-CBS](https://github.com/anna092/cbs5g/assets/113874435/e841bfc6-baf3-4547-b8a1-7f983c9737b4)


## CBCF-A (Cell Broadcast Center Function - Advanced)
CBCF-A is an application function in the 5G core network and is responsible for all types of emergency broadcast messages, such as PWS (Public Warning System), CMAS (Commercial Mobile Alert System), ETWS (Earthquake and Tsunami Warning. System).

## Features:
* Receive UE registration status
* Handle the warning message in CAP (Common Alerting Protocol)
* Handle disaster prevention commands for equipment
* Do the NonUeN2messageTransfer to transfer Write-Replace-Warning Request to the AMF

## Steps to build CBS
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
