from libraries.configuration import *
import os
import socket
import builtins
import platform
import multiprocessing
import json
import time
import datetime
import urllib
import re



def get_ip_address():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect(("8.8.8.8", 80))
    return s.getsockname()[0]


def get_external_ip():
    site = urllib.request.urlopen("http://checkip.dyndns.org/").read()
    grab = re.findall('([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)', site.decode("utf-8"))
    address = grab[0]
    return address



if os.name != "nt":
    import fcntl
    import struct

    def get_interface_ip(ifname):
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        return socket.inet_ntoa(fcntl.ioctl(s.fileno(), 0x8915, struct.pack('256s',
                                ifname[:15]))[20:24])

def get_lan_ip():
    try:
        ip = socket.gethostbyname(socket.gethostname())
    except:
        ip = '127.0.0.1'
    if ip.startswith("127.") and os.name != "nt":
        interfaces = [
            "eth0",
            "eth1",
            "eth2",
            "wlan0",
            "wlan1",
            "wifi0",
            "ath0",
            "ath1",
            "ppp0",
            ]
        for ifname in interfaces:
            try:
                ip = get_interface_ip(ifname.encode('utf-8'))
                break
            except IOError:
                pass
    return ip

builtins.SERVER_STATS = {}
builtins.SERVER_STATS['start_time'] = time.time()
builtins.SERVER_STATS['address'] = {}
conf = getOptions()
builtins.SERVER_STATS['address']['internal_ip'] = get_lan_ip()
builtins.SERVER_STATS['address']['specified_address'] = conf['address']
builtins.SERVER_STATS['address']['specified_port'] = conf['port']
builtins.SERVER_STATS['address']['specified_ext_address'] = conf['ext_address']
builtins.SERVER_STATS['server'] = {}
builtins.SERVER_STATS['server']['system'] = platform.system()
builtins.SERVER_STATS['server']['release'] = platform.release()
builtins.SERVER_STATS['server']['version'] = platform.version().split()[0]
builtins.SERVER_STATS['server']['cores'] = multiprocessing.cpu_count()

availableRoutes = []
with open('libraries/routing.py','r') as f:
    for line in f:
        if 'app.route' in line and '#' not in line:
            if '("' in line:
                availableRoutes.append(line.split('("')[1].split('"')[0])
            elif "('" in line:
                availableRoutes.append(line.split("('")[1].split("'")[0])

builtins.SERVER_STATS['routes'] = availableRoutes

print(json.dumps(builtins.SERVER_STATS,indent=2))

def getServerStats():
    errors = 0
    warnings = 0
    with open('server.log','r') as f:
        for line in f:
            if 'ERROR' in line:
                errors += 1
            elif 'WARN' in line:
                warnings += 1

    builtins.SERVER_STATS['up'] = {}
    builtins.SERVER_STATS['up']['seconds'] = int(time.time() - builtins.SERVER_STATS['start_time'])
    builtins.SERVER_STATS['up']['warnings'] = warnings
    builtins.SERVER_STATS['up']['errors'] = errors

    results = {
        "internal_ip": get_ip_address(),
        "external_ip": get_external_ip(),
        "status": "static",
        "port": conf['port'],
        "registered": builtins.START_TIME,
        "num_cores": multiprocessing.cpu_count()
    }

    return results



