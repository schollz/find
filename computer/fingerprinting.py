"""Fingerprinting for computers

Submits a POST request to the specified server with a scanned fingerprint of the following format:

'''
	{
		"username":"iamauser",
		"time": 1409108787,
		"group": "find",
		"wifi-fingerprint": 
		[
			{
			"mac": "AA:AA:AA:AA:AA:AA",
			"rssi": -45,
			},
			...
			{
			"mac": "BB:BB:BB:BB:BB:BB",
			"rssi": -55,
			}
		]
	}
'''

"""

import subprocess
import json
import time
import requests 
import sys

__author__ = "Zack"
__copyright__ = "Copyright 2014-2015, FIND"
__credits__ = ["Zack Scholl", "Stefan Safranek"]
__version__ = "0.2"
__email__ = "zack@hypercubeplatforms.com"
__status__ = "Development"

def initialize():
	""" Initialization

	Returns a configuration JSON or exits and prints help text
	"""
	help_text = """fingerprint.py

	Submit fingerprints from a wifi-enabled computer.

	Usage: 
		
		sudo python fingerprint.py [options]

	Options:

		--server,-s 'address' (default: localhost)
		--port,-p port (default: 0, use 0 for none)
		--route,-r 'track/learn' (default: track)
		--group,-g 'group' (default: find)
		--user,-u 'user' (default: unknown)
		--location,-l location (default: unknown)
		--continue,-c number of times to run (default: 10)

	"""

	conf = {'server':'localhost','port':8888,'route':'track','group':'find','user':'unknown','location':'unknown','continue':10}

	try:
		for i in range(1,len(sys.argv),2):
			if '-s' in sys.argv[i]:
				conf['server'] = sys.argv[i+1]
			elif '-p' in sys.argv[i]:
				conf['port'] = int(sys.argv[i+1])
			elif '-r' in sys.argv[i]:
				conf['route'] = sys.argv[i+1]
			elif '-g' in sys.argv[i]:
				conf['group'] = sys.argv[i+1]
			elif '-u' in sys.argv[i]:
				conf['user'] = sys.argv[i+1]
			elif '-l' in sys.argv[i]:
				conf['location'] = sys.argv[i+1]
			elif '-c' in sys.argv[i]:
				conf['continue'] = int(sys.argv[i+1])
			else:
					print(help_text)
					sys.exit(-1)
	except:
		print(help_text)
		sys.exit(-1)

	return conf

def fingerprint():
	""" Fingerprinting using the specified scanner
	"""

	conf = initialize()
	url = ''
	if "http" in conf['server'] and conf['port'] == 0:
		url = conf['server']
	elif conf['port'] == 0:
		url = 'http://' + conf['server']
	elif "http" in conf['server']:
		url = conf['server'] + ':' + str(conf['port'])
	elif "http" not in conf['server']:
		url = 'http://' + conf['server'] + ':' + str(conf['port'])
	else:
		print('Error understanding server ' + conf['server'] + ' on port ' + str(conf['port']))
		sys.exit(-1)

	url = url + '/' + conf['route']
	print('url: ' + url)

	headers = {'Content-type': 'application/json', 'Accept': 'text/plain'}
	for i in range(conf['continue']):
		data = {'username': conf['user'],'group':conf['group'],'location':conf['location'],'time': round(time.time()*1000), 'wifi-fingerprint':[]}

		print("Scanning...")
		linux = True
		windows = False
		if linux:
			proc = subprocess.Popen(["iwlist wlan0 scan | grep 'Address\|Signal'"], stdout=subprocess.PIPE, shell=True)
			(out, err) = proc.communicate()
			print("Collecting...")
			macAddress = ""
			signal = ""
			for line in out.splitlines():
				line = str(line)
				if "Address" in line:
					macAddress = line.split(':',1)[1].strip().lower().split("'")[0]
				if "Signal" in line:
					signal = line.split('level=',1)[1].strip().split('dB')[0].split('/')[0]
					data['wifi-fingerprint'].append({'mac':macAddress,'rssi':int(signal)})
		elif windows:
			print("Scanning...")
			proc = subprocess.Popen(["netsh wlan show network mode=bssid"], stdout=subprocess.PIPE, shell=True)
			(out, err) = proc.communicate()
			print("Collecting...")
			macAddress = ""
			signal = ""
			for line in out.split("\n"):
				if "BSSID" in line:
					macAddress = line.split(':',1)[1].strip()
				if "Signal" in line:
					signal = line.split(':',1)[1].split('%')[0].strip()
					data['wifi-fingerprint'].append({'mac':macAddress,'rssi':int(signal)})


		print("Submitting...")
		try:
			print(json.dumps(data))
			r = requests.post(url, data=json.dumps(data))
			print(r.json())
		except:
			print("error submitting fingerprint")
			print(json.dumps(data))

if __name__ == "__main__":
    """Main subroutine
    """
    fingerprint()
