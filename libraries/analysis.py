import builtins
import traceback
import time
import logging
import json
import os
import numpy as np
import random
from math import *

import xml.etree.ElementTree as ET
from multiprocessing import Process, cpu_count
from collections import Counter
import xml.etree.ElementTree as ET
from collections import OrderedDict
from libraries.mldatabase import *
from libraries.posteriors import *

builtins.GETLOCATIONS_CACHE = {}

def erfcc(x):
    """Complementary error function."""
    z = abs(x)
    t = 1. / (1. + 0.5*z)
    r = t * exp(-z*z-1.26551223+t*(1.00002368+t*(.37409196+
    	t*(.09678418+t*(-.18628806+t*(.27886807+
    	t*(-1.13520398+t*(1.48851587+t*(-.82215223+
    	t*.17087277)))))))))
    if (x >= 0.):
    	return r
    else:
    	return 2. - r

def norm_cdf(x):
    return 2.71**(2*x)
    return 1. - 0.5*erfcc(x/(2**0.5))


def getUserLocations(user,group):
    data = {}
    if isinstance(user,list):
        db = mlDB(group)
        users = user
        for user in users:
            data[user] = []
            if group in builtins.fingerprint_cache and user in builtins.fingerprint_cache[group]: 
                for i in range(len(builtins.fingerprint_cache[group][user])):
                    if builtins.fingerprint_cache[group][user][i] is not None:
                        dat = {}
                        dat['location'] = builtins.fingerprint_cache[group][user][i]['location']
                        dat['time'] = builtins.fingerprint_cache[group][user][i]['time']
                        data[user].append(dat)
            else:
                data[user] += db.getLastLocationsFromTracking(user)
                try:
                    data[user] += db.getLastLocationsFromTracking(user)
                except:
                    data[user].append({'location':'unknown','time':int(time.time()*1000)})
        db.close()
    else:
        data[user] = []
        if group in builtins.fingerprint_cache and user in builtins.fingerprint_cache[group]: 
            for i in range(len(builtins.fingerprint_cache[group][user])):
                if builtins.fingerprint_cache[group][user][i] is not None:
                    dat = {}
                    dat['location'] = builtins.fingerprint_cache[group][user][i]['location']
                    dat['time'] = builtins.fingerprint_cache[group][user][i]['time']
                    data[user].append(dat)
        else:
            try:
                db = mlDB(group)
                data[user] += db.getLastLocationsFromTracking(user)
                db.close()
            except:
                data[user].append({'location':'unknown','time':int(time.time()*1000)})
    return data

def getAllLocations(group):
    logger = logging.getLogger('analysis.getAllLocations')
    if group not in builtins.GETLOCATIONS_CACHE:
        builtins.GETLOCATIONS_CACHE[group] = {}
    db = mlDB(group)
    rows = db.executeSqlCommand('select distinct user_id from track')
    data = {}
    for row in rows:
        data[row[0]] = {}
    for user in data:
        if user not in builtins.GETLOCATIONS_CACHE[group]:
            builtins.GETLOCATIONS_CACHE[group][user] = {}
            builtins.GETLOCATIONS_CACHE[group][user]['time'] = 'none'

        if group in builtins.fingerprint_cache and user in builtins.fingerprint_cache[group]:
            fingerprint = copy.deepcopy(builtins.fingerprint_cache[group][user][0])
            data[user]['time'] = time.strftime("%a, %d %b %Y %H:%M:%S", time.localtime(int(fingerprint['time'])/1000))
            data[user]['location'] = fingerprint['location']
            if data[user]['time'] == builtins.GETLOCATIONS_CACHE[group][user]['time']:
                data[user] = copy.deepcopy(builtins.GETLOCATIONS_CACHE[group][user])
                continue
        else:
            rows = db.executeSqlCommand("select timestamp,location_uuid from track where user_id like '%s' order by id desc limit 1" % user)
            for row in rows:
                timestamp = row[0]
                data[user]['time'] = time.strftime("%a, %d %b %Y %H:%M:%S", time.localtime(int(row[0])/1000))
                data[user]['location'] = row[1]

            if data[user]['time'] == builtins.GETLOCATIONS_CACHE[group][user]['time']:
                data[user] = copy.deepcopy(builtins.GETLOCATIONS_CACHE[group][user])
                continue

            fingerprint = db.retrieveFingerprint('track',timestamp)
            logger.debug('Using database location for ' + user)

        fingerprint['group'] = group
        guesses = processTrackingFingerprint(fingerprint, toSave = False, testing = True)
        data[user]['guesses'] = {}
        total = 0
        for guess in guesses:
            total += norm_cdf(guess[1])
        foo = {}
        for guess in guesses:
            try:
                data[user]['guesses'][guess[0]]=int(norm_cdf(guess[1])/total*100.0)
                foo[guess[0]] = int(norm_cdf(guess[1])/total*100.0)
            except:
                data[user]['guesses'][guess[0]] = 0
                foo[guess[0]] = 0
        
        pcts = []
        labels = []
        for key in sorted(foo.keys()):
            labels.append(key)
            if foo[key] < 5:
                pcts.append(str(0))
            else:
                pcts.append(str(foo[key]))
        data[user]['pcts'] = '[' + ','.join(pcts) + ']'
        data[user]['labels'] = '["' + '" ,"'.join(labels) + '"]'
        builtins.GETLOCATIONS_CACHE[group][user] = copy.deepcopy(data[user])
    db.close() 
    return data

def submitGPX(group):
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

    return datas


def getSimulationResults(group,force=False):
    db = mlDB(group)
    simulation_results_pie_json = db.getResource('simulation_results_pie_json')
    db.close()
    if simulation_results_pie_json == None or force:
        db = mlDB(group)
        dropout_percentage = db.getResource('dropout_percentage')
        if dropout_percentage == None:
            dropout_percentage = 10 # probability to leave out a random fingerprint
            db.insertResource('dropout_percentage',dropout_percentage)
        db.close()
        simulation_results_pie_json = simulationAnalysis(group,force)
        db = mlDB(group)
        db.insertResource('simulation_results_pie_json',simulation_results_pie_json)
        db.close()
    return simulation_results_pie_json

def simulationAnalysis_worker(group,num,total,dropout_percentage,number_simulations):
    logger = logging.getLogger('analysis.simulationAnalysis_worker')
    db = mlDB(group)
    newP = db.getResource('newP')
    db.close()
    MIN_RSSI = -100
    pie_data = {}
    for i in range(number_simulations):
        if num==0 and i%int(100/10)==0:
            logger.debug('%d percent done' % int(100*i/number_simulations))
        if i%total == num:
            for graph in newP:
                for loc in newP[graph]:
                    if loc not in pie_data:
                        pie_data[loc] = {}
                    wifi_fingerprint = []
                    for mac in newP[graph][loc]:
                        rssi = np.array(newP[graph][loc][mac])
                        rssi_ind = np.array(range(len(newP[graph][loc][mac])))
                        rssi_ind = rssi_ind[rssi>0]
                        rssi = rssi[rssi>0]
                        rand_num = random.randint(0,len(rssi)-1)
                        while True:
                            #print('mac: ' + mac + ' with val:' + str(rssi_ind[rand_num]) + '/' + str(rssi_ind[rand_num]+MIN_RSSI) + ' with prob ' + str(rssi[rand_num]))
                            if (rssi[rand_num]*1000>=random.randint(0,1000)):
                                break
                            else:
                                rand_num = random.randint(0,len(rssi)-1)
                        if random.randint(0,100)>dropout_percentage:
                            wifi_fingerprint.append({'mac':mac,'rssi':rssi_ind[rand_num]+MIN_RSSI})
                    fingerprint = {}
                    fingerprint['wifi-fingerprint'] = wifi_fingerprint
                    fingerprint['location'] = loc
                    fingerprint['time'] = int(time.time())
                    fingerprint['group'] = group
                    fingerprint['user'] = 'bootstrapper'
                    results = calculatePosterior(fingerprint)
                    loc_guess = results[0][0]
                    if loc_guess not in pie_data[loc]:
                        pie_data[loc][loc_guess] = 1
                    else:
                        pie_data[loc][loc_guess] += 1

    pickle.dump(pie_data,open(str(num)+'ea.p','wb'))

def simulationAnalysis(group,force=False):
    logger = logging.getLogger('analysis.simulationAnalysis')
    db = mlDB(group)
    newP = db.getResource('newP')
    if newP == None or force:
        P = db.getResource('P')
        newP = {}
        for graph in P:
            if graph not in newP:
                newP[graph] = {}
            for mac in P[graph]:
                for loc in P[graph][mac]:
                    if loc not in newP[graph]:
                        newP[graph][loc] = {}
                    if mac not in newP[graph][loc]:
                        newP[graph][loc][mac] = P[graph][mac][loc]
        db.insertResource('newP',newP)
    # Create index if it hasn't already been created
    try:
        db.executeSqlCommand('drop index ind_ex2')
        logger.debug('Dropped previous index')
    except:
        logger.debug('No index to drop')
    db.executeSqlCommand('create index ind_ex2 on test(mac_address, location_uuid)')
    logger.debug('Generated new index on test')
    db.close()

    
    number_simulations = 100
    dropout_percentage = 10

    bench = time.time()
    num_processors = cpu_count()
    numThreads = num_processors
    threads = []
    for i in range(numThreads):
        tThreads = Process(target=simulationAnalysis_worker, args=(group,i,numThreads,dropout_percentage,number_simulations))
        threads.append(tThreads)
    for each in threads:
        each.start()
    for each in threads:
        each.join()

    pie_data = {}
    for i in range(numThreads):
        logger.debug('loading ' + str(i)+'ea.p to combine')
        pie_data2 = pickle.load(open(str(i)+'ea.p','rb'))
        for loc1 in pie_data2.keys():
            for loc2 in pie_data2[loc1].keys():
                if loc1 not in pie_data:
                    pie_data[loc1] = {}
                if loc2 not in pie_data[loc1]:
                    pie_data[loc1][loc2] = pie_data2[loc1][loc2]
                else:
                    pie_data[loc1][loc2] += pie_data2[loc1][loc2]
        os.system('rm ' + str(i)+'ea.p')

    logger.debug('Evaluating accuracy took %s seconds' % str(time.time()-bench))

    pie_json = []
    for loc1 in pie_data.keys():
        for loc2 in pie_data[loc1].keys():
            pie_json.append({'origin':loc1,
                'carrier':loc2,'count':pie_data[loc1][loc2]})

    return pie_json




def makeStats(group):
    db = mlDB(group)
    a = {}
    a['table_data'] = db.databaseStats()
    a['uptime'] = builtins.START_TIME
    connected_components_locs = db.getResource('connected_components_locs')
    connected_components = db.getResource('connected_components')
    if connected_components_locs == None:
        return None
    dropout_percentage = db.getResource('dropout_percentage')
    db.close()
    a['network'] = {}
    for graph in connected_components_locs.keys():
        words = ' '.join(list(set(connected_components_locs[graph]))).replace('floor','')
        words = ''.join([i for i in words if not i.isdigit()])
        word_list = words.split()
        c = Counter(word_list)
        most_common = c.most_common(1)[0][0]
        name = most_common.title()
        if name not in a['network']:
            a['network'][graph] = []
        a['network'][graph] = list(set(connected_components_locs[graph]))

    for loc in a['table_data']['learn']:
        found = False
        for network in a['network']:
            if loc in a['network'][network]:
                found = True
                break
        if not found:
            if 'Uncalculated' not in a['network']:
                a['network']['Uncalculated'] = []
            a['network']['Uncalculated'].append(loc)

    pies = makePies(group)
    pieStats = {}
    for pie in pies:
        loc = pie['origin']
        if loc not in pieStats:
            pieStats[loc] = {}
            pieStats[loc]['right'] = 0
            pieStats[loc]['wrong'] = 0
        if loc == pie['carrier']:
            pieStats[loc]['right'] += pie['count']
        else:
            pieStats[loc]['wrong'] += pie['count']
    allPercentages = []
    for loc in pieStats.keys():
        pieStats[loc]['percentage'] = int(pieStats[loc]['right']*100 / (pieStats[loc]['right']+pieStats[loc]['wrong']))
        allPercentages.append(pieStats[loc]['percentage'])
    a['pie_stats'] = pieStats
    try:
        a['accuracy'] = int(np.mean(allPercentages))
    except:
        a['accuracy'] = 0

    '''
    pies = getSimulationResults(group)
    pieStats = {}
    for pie in pies:
        loc = pie['carrier']
        if loc not in pieStats:
            pieStats[loc] = {}
            pieStats[loc]['right'] = 0
            pieStats[loc]['wrong'] = 0
        if loc == pie['origin']:
            pieStats[loc]['right'] += pie['count']
        else:
            pieStats[loc]['wrong'] += pie['count']
    allPercentages = []
    for loc in pieStats.keys():
        pieStats[loc]['percentage'] = int(pieStats[loc]['right']*100 / (pieStats[loc]['right']+pieStats[loc]['wrong']))
        allPercentages.append(pieStats[loc]['percentage'])
    a['pie_stats_simulation'] = pieStats
    '''
    try:
        a['accuracy_simulation'] = int(np.mean(allPercentages))
    except:
        a['accuracy_simulation'] = 0
    a['dropout_percentage'] = dropout_percentage
    a['calculation_parameters'] = builtins.PARAMETERS[group]
    
    a['accuracies'] = {}
    for graph in range(len(connected_components)):
        totalAccuracy = 0.0
        totalNum = 0.0
        for loc in connected_components_locs[graph]:
            try:
                totalNum += 1.0
                totalAccuracy += a['pie_stats'][loc]['percentage']*1.0
            except:
                pass
        if totalNum == 0:
            totalNum = 1;
        a['accuracies'][graph] = int(round(totalAccuracy/totalNum,0))

    return a


def evaluateAccuracy_worker(group,num,total,timestamps,table):
    logger = logging.getLogger('priors.calculate_worker')
    numTimestamps = len(timestamps)
    db = mlDB(group)
    locs = db.getUniqueLocations()
    db.close()
    pie_data = {}
    for loc in locs:
        pie_data[loc] = {}
        for loc2 in locs:
            pie_data[loc][loc2] = 0

    counter = 0;
    for timestamp in timestamps:
        try:
       	    if num==0 and counter%int(numTimestamps/8)==0:
                logger.debug('%d percent done' % int(100*counter/numTimestamps))
       	except:
       	    pass
        if counter%total == num:
            db = mlDB(group)
            fingerprint = db.retrieveFingerprint(table,timestamp)
            fingerprint['group'] = group
            db.close()
            #guesses = calculatePosterior(fingerprint)
            actualLoc = fingerprint['location']
            guesses = processTrackingFingerprint(fingerprint, toSave = False, testing = True)
            if actualLoc in pie_data:
                pie_data[actualLoc][guesses[0][0]] += 1
                
        counter += 1

    pickle.dump(pie_data,open(str(num)+'ea.p','wb'))
    logger.info('Wrote pie_data to ' + str(num) + 'ea.p')


def evaluateAccuracy(group,tables):
    """ Evaluates accuracy of the calcuations

    Input:
    List of tables to use for the calculation
    """
    logger = logging.getLogger('analysis.evaluateAccuracy')
    db = mlDB(group)
    locs = db.getUniqueLocations()
    db.close()
    pie_data = {}
    for loc in locs:
        pie_data[loc] = {}
        for loc2 in locs:
            pie_data[loc][loc2] = 0

    tables = ['test']

    for table in tables:
        db = mlDB(group)
        timestamps = db.getUniqueTimestamps(table)
        db.close()
        bench = time.time()
        num_processors = cpu_count()
        numThreads = num_processors
        threads = []
        for i in range(numThreads):
            tThreads = Process(target=evaluateAccuracy_worker, args=(group,i,numThreads,timestamps,table,))
            threads.append(tThreads)
        for each in threads:
            each.start()
        for each in threads:
            each.join()

        for i in range(numThreads):
            logger.debug('loading ' + str(i)+'ea.p to combine')
            try:
                pie_data2 = pickle.load(open(str(i)+'ea.p','rb'))
                for loc1 in pie_data2.keys():
                    for loc2 in pie_data2[loc1].keys():
                        pie_data[loc1][loc2] += pie_data2[loc1][loc2]
                os.system('rm ' + str(i)+'ea.p')
            except Exception:
                logger.error(traceback.format_exc())

    db = mlDB(group)
    db.insertResource('pie_data',pie_data)
    db.close()

    stats = makeStats(group)
    logger.debug('Accuracy of %s took %s seconds' % (str(stats['accuracy']),str(time.time()-bench)))

    return True
 
def makePies(group,graph=-1):
    db = mlDB(group)
    if not db.haveResource('pie_data'):
        return False
    pie_data = db.getResource('pie_data')
    connected_components_locs = db.getResource('connected_components_locs')
    db.close()

    pie_json = []
    for loc1 in pie_data.keys():
        for loc2 in pie_data.keys():
            if pie_data[loc1][loc2] > 0:
                count = pie_data[loc1][loc2]
                loc1name = loc1
                loc2name = loc2
                if graph < 0 or loc1name in connected_components_locs[graph]:
                    pie_json.append({'origin':loc1name,
                        'carrier':loc2name,
                        'count':count})

    return pie_json


def makeTimeChartJson(group):
    logger = logging.getLogger('analysis.makeTimeChartJson')    
    charts = {}
    db = mlDB(group)
    connected_components = db.getResource('connected_components')
    connected_components_locs = db.getResource('connected_components_locs')
    connected_components_macs = db.getResource('connected_components_macs')
    usefulMacs = db.getResource('usefulMacs')
    for graph in range(len(connected_components)):
        for loc in connected_components_locs[graph]:
            charts[loc]={}
            charts[loc][0] = {}
            charts[loc][0][0] = {}
            charts[loc][0][0]['type']='line'
            charts[loc][0][0]['data']={}
            for mac in connected_components_macs[graph]:
                if mac in usefulMacs[graph]['good_macs']:
                    x=[]
                    y=[] 
                    rows = db.executeSqlCommand('select timestamp,signal from learn indexed by ind_ex1 where mac_address = "%s" and location_uuid = "%s" order by timestamp'%(mac,loc))
                    num = 0
                    for row in rows:
                        x.append(num)
                        num += 1
                        y.append(row[1])
                    if len(x) > 0:
                        charts[loc][0][0]['data'][mac] = {}
                        charts[loc][0][0]['data'][mac] = [x,y]
    db.close()
    return charts
   
                                   

def makeChartJson(group):
    logger = logging.getLogger('analysis.makeChartJson')
    charts = {}
    MIN_RSSI = -100
    MAX_RSSI = -10
    RSSI_PARTITIONS=MAX_RSSI-MIN_RSSI+1

    db = mlDB(group)
    P = db.getResource('P')
    nP = db.getResource('nP')
    W = db.getResource('W')
    connected_components = db.getResource('connected_components')
    connected_components_locs = db.getResource('connected_components_locs')
    connected_components_macs = db.getResource('connected_components_macs')
    usefulMacs = db.getResource('usefulMacs')
    good_macs_by_loc = db.getResource('good_macs_by_loc')
    db.close()


    for graph in range(len(connected_components)):
        charts['cluster'+str(graph)]={}

        charts['cluster'+str(graph)]['All Macs'] = {}
        for loc in connected_components_locs[graph]:
            location = loc
            charts['cluster'+str(graph)]['All Macs'][location] = {}
            charts['cluster'+str(graph)]['All Macs'][location]['type']='line'
            charts['cluster'+str(graph)]['All Macs'][location]['data']={}
            for mac in connected_components_macs[graph]:
                x=[]
                y=[]
                for i in range(RSSI_PARTITIONS):
                    x.append(i)
                    g = float("{0:.2f}".format(P[graph][mac][loc][i]))
                    y.append(g)
                if mac not in charts['cluster'+str(graph)]['All Macs'][location]['data']:
                    charts['cluster'+str(graph)]['All Macs'][location]['data'][mac] = {}
                charts['cluster'+str(graph)]['All Macs'][location]['data'][mac]=[x,y]

        charts['cluster'+str(graph)]['Useful Macs'] = {}
        charts['cluster'+str(graph)]['Useful Macs-'] = {}
        charts['cluster'+str(graph)]['Useless Macs'] = {}
        for loc in connected_components_locs[graph]:
            location = loc
            charts['cluster'+str(graph)]['Useful Macs'][location] = {}
            charts['cluster'+str(graph)]['Useful Macs'][location]['type']='line'
            charts['cluster'+str(graph)]['Useful Macs'][location]['data']={}
            charts['cluster'+str(graph)]['Useful Macs-'][location] = {}
            charts['cluster'+str(graph)]['Useful Macs-'][location]['type']='line'
            charts['cluster'+str(graph)]['Useful Macs-'][location]['data']={}
            charts['cluster'+str(graph)]['Useless Macs'][location] = {}
            charts['cluster'+str(graph)]['Useless Macs'][location]['type']='line'
            charts['cluster'+str(graph)]['Useless Macs'][location]['data']={}
            for mac in connected_components_macs[graph]:
                if mac in usefulMacs[graph]['good_macs']:
                    x=[]
                    y=[]
                    for i in range(RSSI_PARTITIONS):
                        x.append(i)
                        g = float("{0:.2f}".format(P[graph][mac][loc][i]))
                        try:
                            y.append(g*good_macs_by_loc[graph][loc][mac])
                        except:
                            y.append(0)
                    if mac not in charts['cluster'+str(graph)]['Useful Macs'][location]['data']:
                        charts['cluster'+str(graph)]['Useful Macs'][location]['data'][mac] = {}
                    charts['cluster'+str(graph)]['Useful Macs'][location]['data'][mac]=[x,y]

                    x=[]
                    y=[]
                    for i in range(RSSI_PARTITIONS):
                        x.append(i)
                        g = float("{0:.2f}".format(nP[graph][mac][loc][i]))
                        y.append(g)

                    if mac not in charts['cluster'+str(graph)]['Useful Macs-'][location]['data']:
                        charts['cluster'+str(graph)]['Useful Macs-'][location]['data'][mac] = {}
                    charts['cluster'+str(graph)]['Useful Macs-'][location]['data'][mac]=[x,y]
                else:
                    x=[]
                    y=[]
                    for i in range(RSSI_PARTITIONS):
                        x.append(i)
                        g = float("{0:.2f}".format(P[graph][mac][loc][i]))
                        try:
                            y.append(g*good_macs_by_loc[graph][loc][mac])
                        except:
                            y.append(0)

                    if mac not in charts['cluster'+str(graph)]['Useless Macs'][location]['data']:
                        charts['cluster'+str(graph)]['Useless Macs'][location]['data'][mac] = {}
                    charts['cluster'+str(graph)]['Useless Macs'][location]['data'][mac]=[x,y]


    return charts
