import time
import builtins
import numpy
import math
import logging
import copy
import requests
import json
import sys
import hashlib
import os
import subprocess
import numpy as np
import shlex
from uuid import uuid4
from utm import from_latlon
import xml.etree.ElementTree as ET
from libraries.mldatabase import *

builtins.POSTERIOR_VARIABLES = {}
builtins.fingerprint_cache = {}
builtins.fingerprint_count = {}
builtins.BEST_METHOD = ""

def getGPX(group,location):
    if location == None:
        return {1:{'info':'Unknwon','lat':0,'lng':0,'ele':0}}

    tree = ET.parse('data/'+group+'.gpx')
    root = tree.getroot()

    datas = {}
    for child in root:
        data = {}
        data['lat'] = child.attrib['lat']
        data['lng'] =  child.attrib['lon']
        for chil in  child.getchildren():
            try:
                data[chil.tag.split('}')[1]] = int(chil.text)
            except:
                data[chil.tag.split('}')[1]] = chil.text
        datas[data['name']] = {'lat':data['lat'],'lng':data['lng'],'ele':data['ele']}
   
    location_string = location
    if 'floor' in location:
        location_string = location.split('floor')[1].strip()
        location_string = ' '.join(location_string.split(' ')[1:]).title()
    
    try:
        returnData = {1:{'info':location_string,'lat':datas[location]['lat'],'lng':datas[location]['lng'],'ele':datas[location]['ele']}}
    except:
        returnData = {1:{'info':location_string,'lat':0,'lng':0,'ele':0}}

    return returnData


def processTrackingFingerprint(data, toSave = True, testing = False):
    """Takes a fingerprint from tracking, calculates posterior, and saves if nessecary

    Inputs: 
        data (standard fingerprint format)
        toSave (boolean of whether or not to save)

    Uses a global ```builtins.fingerprint_cache``` to keep track of persistence calculations
    """
    if testing:
        data['username'] = 'test'
    dataOrig = copy.deepcopy(data)

    logger = logging.getLogger('posteriors.processTrackingFingerprint')

    group = data['group'].lower()
    
    graph = getGraph(data)
    
    persistence = 1
    trigger_server = 'None'
    if graph > 0:
        persistence = builtins.PARAMETERS[group][graph]['persistence']
        trigger_server = builtins.PARAMETERS[group][graph]['trigger_server']
    

    persistence = persistence + 1
    if persistence < 1:
        persistence = 1

    # Store the fingerprints for bulk upload and persistence calculation
    if group not in builtins.fingerprint_cache:
        builtins.fingerprint_cache[group] = {}
    if data['username'] not in builtins.fingerprint_cache[group] or len(builtins.fingerprint_cache[group][data['username']]) != persistence:
       builtins.fingerprint_cache[group][data['username']] = [None]*persistence

    # Move the fingerprints across
    for i in range(len(builtins.fingerprint_cache[group][data['username']])-1,0,-1):
        if builtins.fingerprint_cache[group][data['username']][i-1] is not None:
            builtins.fingerprint_cache[group][data['username']][i] = builtins.fingerprint_cache[group][data['username']][i-1] 
    builtins.fingerprint_cache[group][data['username']][0] = copy.deepcopy(data)

    # Count incoming fingerprints
    if group not in builtins.fingerprint_count:
        builtins.fingerprint_count[group] = {}
    if data['username'] not in builtins.fingerprint_count[group]:
       builtins.fingerprint_count[group][data['username']] = 0
    builtins.fingerprint_count[group][data['username']] += 1

    # Calculate persistence
    current_macs = []
    for dat in data['wifi-fingerprint']:
        current_macs.append(dat['mac']) 
    num_fingerprints = len(builtins.fingerprint_cache[group][data['username']])
    for i in range(1,persistence):
        if builtins.fingerprint_cache[group][data['username']][i] is not None:
            for dat in builtins.fingerprint_cache[group][data['username']][i]['wifi-fingerprint']:
                if dat['mac'] not in current_macs:
                    data['wifi-fingerprint'].append({'mac':dat['mac'],'rssi':dat['rssi']})
                    current_macs.append(dat['mac']) 
                    if not testing:
                        logger.debug('adding ' + dat['mac'])
                    
    posteriorResult = calculatePosterior(data,dataOrig,graph)
    builtins.fingerprint_cache[group][data['username']][0]['location'] = posteriorResult[0][0]
    if group not in builtins.CURRENT_LOCATIONS:
        builtins.CURRENT_LOCATIONS[group] = {}
    if data['username'] not in builtins.CURRENT_LOCATIONS:
        builtins.CURRENT_LOCATIONS[group][data['username']] = {}
    
    try:
        if not testing:
            builtins.CURRENT_LOCATIONS[group][data['username']] = getGPX(group,posteriorResult[0][0])
            positions = from_latlon(float(builtins.CURRENT_LOCATIONS[group][data['username']][1]['lat']),float(builtins.CURRENT_LOCATIONS[group][data['username']][1]['lng']))
            logger.info(str(time.time()) +
                    ',' + data['username'] +
                    ',' + data['group'] +
                    ',' + posteriorResult[0][0] + 
                    ',' + str(round(posteriorResult[0][1]/posteriorResult[-1][1],3)) +
                    ',' + str(positions[0]) +
                    ',' + str(positions[1]) +
                    ',' + str(builtins.CURRENT_LOCATIONS[group][data['username']][1]['ele']))
    except:
        logger.warn('No GPX for ' + posteriorResult[0][0])
    
            
    # Save the tracking data in bulk 
    if (builtins.fingerprint_count[group][data['username']] % persistence == 0 
            and None not in builtins.fingerprint_cache[group][data['username']] 
            and toSave
            and not testing):
        logger.debug('Saving bulk tracking data for ' + data['username'])
        db = mlDB(group)
        db.insertFingerprints('/track', builtins.fingerprint_cache[group][data['username']])
        db.close()

    # Send the fingerprint to the trigger server
    
    if (trigger_server is not None 
            and "none" not in trigger_server.lower() 
            and 'http' in trigger_server.lower()
            and not testing):
        logger.debug('sending trigger data to ' + trigger_server)
        try:
            payload = {'location':posteriorResult[0][0],'time':time.time(),'1st/2nd':round(posteriorResult[0][1]/posteriorResult[1][1],1),'1st/last':round(posteriorResult[0][1]/posteriorResult[-1][1],1)}
            logger.debug(json.dumps(payload,indent=4))
            r = requests.post(trigger_server, data=json.dumps(payload))
        except Exception as e:
            logger.error(e)

    if not testing:
        logger.debug(posteriorResult)
        ratioDen = 1
        try:
            ratioDen = round(posteriorResult[1][1],1)
        except:
            pass
        logger.warn(data['username'] + 
            ' (' + data['group'] + ') at ' + 
            posteriorResult[0][0] + 
            ' with prob ratio: ' + str(round(posteriorResult[0][1],1)) + 
            '/' + str(ratioDen))
    return posteriorResult
    
class NumpyAwareJSONEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, numpy.ndarray) and obj.ndim == 1:
            return obj.tolist()
        return json.JSONEncoder.default(self, obj)
    
def sha1(name):
    return hashlib.sha1(name.encode('utf-8')).hexdigest()

    
def getGraph(data):
    if data['group'] not in builtins.POSTERIOR_VARIABLES:
        group = data['group']
        db = mlDB(group)
        builtins.POSTERIOR_VARIABLES[group] = {}
        builtins.POSTERIOR_VARIABLES[group]['P'] = db.getResource('P')
        builtins.POSTERIOR_VARIABLES[group]['nP'] = db.getResource('nP')
        builtins.POSTERIOR_VARIABLES[group]['W'] = db.getResource('W')
        builtins.POSTERIOR_VARIABLES[group]['connected_components'] = db.getResource('connected_components')
        builtins.POSTERIOR_VARIABLES[group]['connected_components_locs'] = db.getResource('connected_components_locs')
        builtins.POSTERIOR_VARIABLES[group]['connected_components_macs'] = db.getResource('connected_components_macs')
        builtins.POSTERIOR_VARIABLES[group]['usefulMacs'] = db.getResource('usefulMacs')
        builtins.POSTERIOR_VARIABLES[group]['good_macs_by_loc'] = db.getResource('good_macs_by_loc')
        builtins.POSTERIOR_VARIABLES[group]['bad_macs_by_loc'] = db.getResource('bad_macs_by_loc')
        db.close()
    
    connected_components = builtins.POSTERIOR_VARIABLES[data['group']]['connected_components']
    try:
        graph = -1
        for dat in data['wifi-fingerprint']:
            for i in range(len(connected_components)):
                if dat['mac'] in connected_components[i]:
                    return i

    except:
        pass
    return -1
        
def calculatePosterior(data,dataOrig,graph):
    """Takes a fingerprint and returns a tuple of the posterior calculation
    
    data:
    {
        "time": 1409108787
        "location":"office/unknown"
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
    
    returns:
     [('office', -12.509194917163677), ('hall', -63.08724478882834), ('bedroom', -73.37576617048872), ('living', -97.50435920469138)]
    """
    logger = logging.getLogger('posteriors.calculatePosterior')
    
    if graph < 0:
        P_bayes = {}
        P_bayes['unknown'] = 0
        return sorted(P_bayes.items(), key=lambda x: x[1],reverse=True)

        
    data['group'] = data['group'].lower()
    if data['group'] not in builtins.POSTERIOR_VARIABLES or builtins.POSTERIOR_VARIABLES[data['group']]['connected_components']==False:
        logger.debug('loading data for ' + data['group'] + ' for posterior calculations')
        builtins.POSTERIOR_VARIABLES[data['group']] = {}
        db = mlDB(builtins.DATABASE_PATH_PREFIX + data['group'] + '.db')
        builtins.POSTERIOR_VARIABLES[data['group']]['P'] = db.getResource('P')
        builtins.POSTERIOR_VARIABLES[data['group']]['nP'] = db.getResource('nP')
        builtins.POSTERIOR_VARIABLES[data['group']]['W'] = db.getResource('W')
        builtins.POSTERIOR_VARIABLES[data['group']]['connected_components'] = db.getResource('connected_components')
        builtins.POSTERIOR_VARIABLES[data['group']]['connected_components_locs'] = db.getResource('connected_components_locs')
        builtins.POSTERIOR_VARIABLES[data['group']]['connected_components_macs'] = db.getResource('connected_components_macs')
        builtins.POSTERIOR_VARIABLES[data['group']]['usefulMacs'] = db.getResource('usefulMacs')
        builtins.POSTERIOR_VARIABLES[data['group']]['good_macs_by_loc'] = db.getResource('good_macs_by_loc')
        builtins.POSTERIOR_VARIABLES[data['group']]['bad_macs_by_loc'] = db.getResource('bad_macs_by_loc')
        db.close()
    P = builtins.POSTERIOR_VARIABLES[data['group']]['P']
    nP = builtins.POSTERIOR_VARIABLES[data['group']]['nP']
    W = builtins.POSTERIOR_VARIABLES[data['group']]['W']
    connected_components = builtins.POSTERIOR_VARIABLES[data['group']]['connected_components']
    connected_components_macs = builtins.POSTERIOR_VARIABLES[data['group']]['connected_components_macs']
    connected_components_locs = builtins.POSTERIOR_VARIABLES[data['group']]['connected_components_locs']
    usefulMacs = builtins.POSTERIOR_VARIABLES[data['group']]['usefulMacs']
    good_macs_by_loc = builtins.POSTERIOR_VARIABLES[data['group']]['good_macs_by_loc']
    bad_macs_by_loc = builtins.POSTERIOR_VARIABLES[data['group']]['bad_macs_by_loc']




    usefulMacs = usefulMacs[graph]['good_macs']

    macs = list(P[graph].keys())
    W={}
    W2={}
    for mac in macs:
        W[mac]=-100
        W2[mac]=-100
    
    for dat in data['wifi-fingerprint']:
        W[dat['mac']]=dat['rssi']

    for dat in dataOrig['wifi-fingerprint']:
        W2[dat['mac']]=dat['rssi']

    numberLocations =  len(connected_components_locs[graph])
    P_A = 1.0/numberLocations;
    P_notA = (numberLocations-1.0)/numberLocations;
    skipped = 0
    used = 0
   

    P_bayes = {}
    P_bayes2 = {}
    t1 = time.time()
    fileName = 'dat' + str(uuid4()) + '.temp'
    with open(fileName,'w') as f:
        f.write(str(P_A) + '\n')
        f.write(str(P_notA) + '\n')
        for loc in connected_components_locs[graph]:
            P_bayes[loc] = 0;
            P_bayes2[loc] = 0;
            for mac in macs:
                try:
                    weight = (good_macs_by_loc[graph][loc][mac])
                except:
                    weight = (good_macs_by_loc[graph][loc]['min'])
                try:
                    badweight = (bad_macs_by_loc[graph][loc][mac])   
                except:
                    badweight = (bad_macs_by_loc[graph][loc]['min'])   
                pInda = int(W[mac]+100)
                P_B_A = P[graph][mac][loc][pInda]
                P_B_notA = nP[graph][mac][loc][pInda]
                in_useful = 0
                if mac in usefulMacs:
                    in_useful = 1

                f.write(loc + '\n')
                f.write(str(W2[mac]) + '\n')
                f.write(str(weight) + '\n')
                f.write(str(badweight) + '\n')
                f.write(str(in_useful) + '\n')
                f.write(str(P_B_A) + '\n')
                f.write(str(P_B_notA) + '\n')
    if len(builtins.BEST_METHOD) < 1:
        for (dir, _, files) in os.walk("calculate"):
            for f in files:
                path = os.path.join(dir, f)
                if 'src' not in path and 'build.py' not in path:            
                    try:
                        logger.debug('trying ' + path)
                        subprocess.call(shlex.split('./' + path + ' ' + fileName))
                        with open(fileName + '.dumped','r') as f:
                            lineI = 0
                            for line in f:
                                lineI += 1
                                if lineI == 1:
                                    P_bayes = json.loads(line)
                                else:
                                    P_bayes2 = json.loads(line)
                        os.system('rm ' + fileName + ' & rm ' + fileName + '.dumped')
                        builtins.BEST_METHOD = './' + path
                    except OSError:
                        logger.debug(path + ' failed')
                        builtins.BEST_METHOD = ''
                if len(builtins.BEST_METHOD)>0:
                    break
            if len(builtins.BEST_METHOD)>0:
                break
    else:
        #logger.debug('using ' + builtins.BEST_METHOD)
        subprocess.call(shlex.split(builtins.BEST_METHOD + ' ' + fileName))
        with open(fileName + '.dumped','r') as f:
            lineI = 0
            for line in f:
                lineI += 1
                if lineI == 1:
                    P_bayes = json.loads(line)
                else:
                    P_bayes2 = json.loads(line)
        os.system('rm ' + fileName + ' & rm ' + fileName + '.dumped')
    #print("MODIFIED ROUTINE TOOK " + str(time.time()-t1) + " ")


    #logger.debug('Skipped ' + str(skipped) + ' and used ' + str(used))
    #print(sorted(P_bayes.items(), key=lambda x: x[1],reverse=True))
    #print(sorted(P_bayes2.items(), key=lambda x: x[1],reverse=True)   )


    P_bayes = normalizeDict(P_bayes)
    P_bayes2 = normalizeDict(P_bayes2)
    P_bayes3 = {}
    biasWeight = builtins.PARAMETERS[data['group']][graph]['mix_in']
    for key in P_bayes.keys():
        P_bayes3[key] = biasWeight*P_bayes[key] + (1-biasWeight)*P_bayes2[key]
    P_bayes = sorted(P_bayes.items(), key=lambda x: x[1],reverse=True)
    P_bayes2 = sorted(P_bayes2.items(), key=lambda x: x[1],reverse=True)   
    P_bayes3 = sorted(P_bayes3.items(), key=lambda x: x[1],reverse=True)   
    #print(P_bayes)
    #print(P_bayes2)
    #print(P_bayes3)
    return P_bayes3
    
def num_there(s):
    return any(i.isdigit() for i in s)

def normalizeDict(d):
    values = []
    for key in d.keys():
        values.append(d[key])
    values = np.array(values)
    mean = np.mean(values)
    sd = np.std(values)
    d2 = copy.deepcopy(d)
    for key in d.keys():
        d2[key] = (d[key]-mean)/sd
        if np.isnan(d2[key]):
            d2[key] = 0
    return d2
