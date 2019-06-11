# traffic-sniffer
Python Traffic Sniffer and Analyzer Framework

# Installation and setting up
## Linux
To install the current LyaWAF Test Web Application and prepare a virtual environment you will need:

Install python version 3.6.x or newer (if it is not already installed - you can check it with “python3 --version” command):
```bash
$ sudo apt-get update 
$ sudo apt-get install python3.6
```
Install pip3 (if is it not already installed - you can check it with “which pip3” or “pip3 --version” commands):
```bash
$ sudo apt-get install python3-pip
```
Install virtualenv (if it is not already installed - you can check it with “which virutalenv” or “virtualenv --version” commands):
```bash
$ sudo apt-get install python3-venv
$ sudo pip3 install virtualenv
```
Create a virtual environment:
```bash
$  virtualenv -p python3.6 venv
```
or
```bash
$ python3.6 -m venv venv
```

Activate a virtual environment:
```bash
$ source venv/bin/activate
```
Install requirements:
```bash
$ pip3 install -r requirements.txt
```
Run:
```bash
$ python3 app.py
```
Open application at http://127.0.0.1:5000/

## MacOS
Slightly different, you will need installed python3 (from brew, for example). 3.7 or any other is ok:
```bash
pip3 install --upgrade virtualenv
```
Then
```
virtualenv -p python3 venv
```
And next steps from linux
