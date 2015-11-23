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

import shlex
import subprocess
import json
import time
import requests
import sys
import platform

__author__ = "Zack"
__copyright__ = "Copyright 2014-2015, FIND"
__credits__ = ["Zack", "Stefan"]
__version__ = "0.3"
__email__ = "zack@hypercubeplatforms.com"
__status__ = "Development"


def get_network_call(operating_system):
    cmd = ''
    if operating_system == 'Darwin':
        cmd = "/System/Library/PrivateFrameworks/Apple80211.framework/" + \
              "Versions/Current/Resources/airport -I en0"
    elif operating_system == 'Linux':
        cmd = r"iwlist wlan0 scan | grep 'Address\|Signal'"
    elif operating_system == 'Windows':
        cmd = "netsh wlan show network mode=bssid"
    return cmd


def network_call_on_os(operating_system):
    cmd = get_network_call(operating_system)
    if operating_system == 'Windows':
        proc = subprocess.Popen(shlex.split(cmd), stdout=subprocess.PIPE, shell=True)
    else:
        proc = subprocess.Popen(cmd, stdout=subprocess.PIPE, shell=True)
    (out, _) = proc.communicate()
    if len(out) == 0:
        raise ValueError('No ouput from network call')
    return out if operating_system != 'Windows' else out.decode('utf-8')


def initialize():
    """ Initialization
    Returns a configuration JSON or exits and prints help text
    """
    help_text = """fingerprint.py

Submit fingerprints from a wifi-enabled computer.

Usage:

    sudo python3 fingerprint.py [options]

Options:

    --server,-s 'address' (default: localhost)
    --port,-p port (default: 0, use 0 for none)
    --route,-r 'track/learn' (default: track)
    --group,-g 'group' (default: find)
    --user,-u 'user' (default: unknown)
    --location,-l location (default: unknown)
    --continue,-c number of times to run (default: 10)
    """

    conf = {'server': 'localhost', 'port': 8888, 'route': 'track',
            'group': 'find', 'user': 'unknown', 'location': 'unknown', 'continue': 10}
    if len(sys.argv) < 2:
        print(help_text)
        sys.exit(0)
    try:
        for i in range(1, len(sys.argv), 2):
            if '-s' in sys.argv[i]:
                conf['server'] = sys.argv[i + 1]
            elif '-p' in sys.argv[i]:
                conf['port'] = int(sys.argv[i + 1])
            elif '-r' in sys.argv[i]:
                conf['route'] = sys.argv[i + 1]
            elif '-g' in sys.argv[i]:
                conf['group'] = sys.argv[i + 1]
            elif '-u' in sys.argv[i]:
                conf['user'] = sys.argv[i + 1]
            elif '-l' in sys.argv[i]:
                conf['location'] = sys.argv[i + 1]
            elif '-c' in sys.argv[i]:
                conf['continue'] = int(sys.argv[i + 1])
            else:
                print(help_text)
                sys.exit(0)
    except:
        print(help_text)
        sys.exit(-1)

    return conf


def get_network_data(conf):
    data = {'username': conf['user'], 'group': conf['group'], 'location': conf['location'],
            'time': round(time.time() * 1000), 'wifi-fingerprint': []}

    operating_system = platform.system()
    macAddress = ''
    signal = ''

    print("Scanning...")
    out = network_call_on_os(operating_system)

    if operating_system == 'Linux':

        for line in out.splitlines():
            line = line.decode('utf-8')
            if "Address" in line:
                macAddress = line.split(':', 1)[1].strip().lower().split("'")[0]
            if "Signal" in line:
                signal = line.split('level=', 1)[1].strip().split('dB')[0].split('/')[0]
                data['wifi-fingerprint'].append({'mac': macAddress, 'rssi': int(signal)})

    elif operating_system == 'Darwin':

        for line in out.splitlines():
            line = str(line)
            if "BSSID" in line:
                macAddress = line.split(': ', 1)[1].strip("'")
            if "agrCtlRSSI" in line:
                signal = int(line.split(': ', 1)[1].strip("'"))
                data['wifi-fingerprint'].append({'mac': macAddress, 'rssi': int(signal)})

    elif operating_system == 'Windows':

        for line in out.split("\n"):
            if "BSSID" in line:
                macAddress = line.split(':', 1)[1].strip()
            if "Signal" in line:
                signal = line.split(':', 1)[1].split('%')[0].strip()
                data['wifi-fingerprint'].append({'mac': macAddress, 'rssi': int(signal)})

    print("Submitting...")

    newData = {
        'token': '9c3eaa3fff1717',
        'address': 1,
        'wifi': []
    }
    ss = ""
    for dat in data['wifi-fingerprint']:
        newData['wifi'].append({'bssid': dat['mac'], 'signal': dat['rssi']})
        ss += dat['mac'].lower() + ',' + str(dat['rssi']) + ','
    print(json.dumps(newData, indent=2))
    try:
        r = requests.post(conf['url'], data=json.dumps(data))
        print(r.text)
    except requests.exceptions.RequestException as e:
        print(str(e), "error submitting fingerprint")
        print(data)


def fingerprint():
    """ Fingerprinting using the specified scanner """

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

    conf['url'] = url
    for _ in range(conf['continue']):
        get_network_data(conf)

if __name__ == "__main__":
    fingerprint()
