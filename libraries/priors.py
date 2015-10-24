import builtins
import logging
import numpy
import math
import os
import time
import threading
import copy
import json
from multiprocessing import Process, cpu_count
from collections import OrderedDict
import pickle
from libraries.mldatabase import *
from libraries.analysis import *

def calculatePriors(group):
    logger = logging.getLogger('priors.calculatePriors')
    P={}
    nP={}
    normpdf = [1,0,0,0.0,0.000]
    threads = []

    db = mlDB(group)
    macs = db.getUniqueMacs()
    locs = db.getUniqueLocations()
    connected_components = db.getResource('connected_components')
    connected_components_locs = db.getResource('connected_components_locs')
    connected_components_macs = db.getResource('connected_components_macs')
    num_processors = cpu_count()

    db.close()
    
    
    
    pdf_types = {}
    pdf_types[1] = [1,0,0,0.0,0.000,0.00] # Normal distribution values at 0 SD, 1 SD, 2 SD, 3 SD, and 4 SD
    pdf_types[2] = [.7979,.1080,.0002,0.00,0.00,0.00] # Normal distribution values at 0 SD, 1 SD, 2 SD, 3 SD, and 4 SD
    pdf_types[3] = [.7979,.7979,.1210,0.05,0.0001,0.00] # Normal distribution values at 0 SD, 1 SD, 2 SD, 3 SD, and 4 SD
    pdf_types[4] = [.5319,.2187,.0152,0.001,0.00,0.00] # Normal distribution values at 0 SD, 1 SD, 2 SD, 3 SD, and 4 SD
    pdf_types[5] = [.3989,.2420,.0540,.0044,0.00,0.00] # Normal distribution values at 0 SD, 1 SD, 2 SD, 3 SD, and 4 SD
    pdf_types[6] = [.1995,.1760,.1210,.0648,.027,0.005] # Normal distribution values at 0 SD, 1 SD, 2 SD, 3 SD, and 4 SD
    pdf_types[7] = [.15,.12,.10,.08,.06,0.04] # Normal distribution values at 0 SD, 1 SD, 2 SD, 3 SD, and 4 SD
    



    # Initialize probabilities
    # Allowed RSSI values = -100 to -10 (mapped to 0 to 90)
    MIN_RSSI = -100
    MAX_RSSI = -10 #used to be -25
    RSSI_PARTITIONS=MAX_RSSI-MIN_RSSI+1
    Wdefault = {}
    totalThings = 0
    calculation_parameters = copy.deepcopy(builtins.PARAMETERS[group])
    for graph in range(len(connected_components)):
        # Get calculation parameters from libraries.analysis
        calculation_parameters[graph]['normpdf'] = pdf_types[calculation_parameters[graph]['pdf_type']]
        calculation_parameters[graph]['absentee'] = float(calculation_parameters[graph]['absentee'])
        calculation_parameters[graph]['usefulness'] = float(calculation_parameters[graph]['usefulness'])
        Wdefault[graph]={}
        P[graph]={}
        nP[graph]={}
        for mac in connected_components_macs[graph]:
            Wdefault[graph][mac]=0
            P[graph][mac]={}
            nP[graph][mac]={}
            for loc in connected_components_locs[graph]:
                totalThings = totalThings + 1.0
                P[graph][mac][loc] = numpy.zeros(RSSI_PARTITIONS,dtype=numpy.dtype(numpy.float64) )
                nP[graph][mac][loc] = numpy.zeros(RSSI_PARTITIONS,dtype=numpy.dtype(numpy.float64) )



    bench = time.time()
    numThreads = num_processors
    for i in range(numThreads):
        tThreads = Process(target=calculate_worker, args=(i,numThreads,totalThings,calculation_parameters,group,))
        threads.append(tThreads)
    for each in threads:
        each.start()

    for each in threads:
        each.join()


    for i in range(numThreads):
        logger.debug('loading ' + str(i)+'.p to combine')
        (P0,nP0) = pickle.load(open(str(i)+'.p','rb'))
        for graph in range(len(connected_components)):
            for mac in connected_components_macs[graph]:
                for loc in connected_components_locs[graph]:
                    P[graph][mac][loc] = numpy.add(P[graph][mac][loc],P0[graph][mac][loc])
        os.system('rm ' + str(i)+'.p')


    logger.debug('Organizing priors took %s seconds' % str(time.time()-bench))

    # Normalize the P
    logger.debug('Normalizing P')
    for graph in range(len(connected_components)):
        for mac in connected_components_macs[graph]:
            for loc in connected_components_locs[graph]:
                pTotal = numpy.sum(P[graph][mac][loc])
                if pTotal>0:
                    for i in range(RSSI_PARTITIONS):
                        P[graph][mac][loc][i] = float(P[graph][mac][loc][i])/float(pTotal)

    # Calculate the nP
    logger.debug('Calculating nP')
    for graph in range(len(connected_components)):
        for mac in connected_components_macs[graph]:
            for loc in connected_components_locs[graph]:
                allLocs =  copy.deepcopy(connected_components_locs[graph])
                allLocs.remove(loc)
                for aloc in allLocs:
                    nP[graph][mac][loc] = numpy.add(nP[graph][mac][loc],P[graph][mac][aloc])

    # Normalize the nP
    logger.debug('Normalizing nP')
    for graph in range(len(connected_components)):
        for mac in connected_components_macs[graph]:
            for loc in connected_components_locs[graph]:
                npTotal = numpy.sum(nP[graph][mac][loc])
                if npTotal>0:
                    for i in range(RSSI_PARTITIONS):
                        nP[graph][mac][loc][i] = nP[graph][mac][loc][i]/npTotal


    # Convert back to float16
    for graph in range(len(connected_components)):
        for mac in connected_components_macs[graph]:
            for loc in connected_components_locs[graph]:
                P[graph][mac][loc] = P[graph][mac][loc].astype(numpy.float16)
                nP[graph][mac][loc] = nP[graph][mac][loc].astype(numpy.float16)                    


    # Calculate average signal for each mac in a given location
    average_values = {}
    stds = {}
    allStds = {}
    for graph in range(len(connected_components)):
        rssiRange = numpy.arange(MIN_RSSI,MAX_RSSI+1)
        average_values[graph] = {}
        stds[graph] = {}
        allStds[graph] = []
        for mac in connected_components_macs[graph]:
            average_values[graph][mac] = []
            for loc in connected_components_locs[graph]:
                x=P[graph][mac][loc]
                avg = numpy.sum(numpy.multiply(rssiRange,x))
                if avg == 0:
                    avg = MIN_RSSI
                average_values[graph][mac].append(avg)
            stds[graph][mac] = numpy.std(average_values[graph][mac])
            allStds[graph].append(stds[graph][mac])



    usefulMacs = {}
    for graph in range(len(connected_components)):
        usefulMacs[graph] = {}
        usefulMacs[graph]['bad_macs'] = []
        usefulMacs[graph]['good_macs'] = []
        meanStd = numpy.average(allStds[graph])
        stdStd = numpy.std(allStds[graph])
        for mac in stds[graph].keys():
            if stds[graph][mac]<meanStd+calculation_parameters[graph]['usefulness']*stdStd:
                usefulMacs[graph]['bad_macs'].append(mac)
            else:
                usefulMacs[graph]['good_macs'].append(mac)

    logger.debug('Working on filtering macs')
    good_macs_by_loc = OrderedDict()
    db = mlDB(group)
    for graph in range(len(connected_components)):
        good_macs_by_loc[graph] = {}
        for loc in connected_components_locs[graph]:
            if loc not in good_macs_by_loc[graph]:
                good_macs_by_loc[graph][loc] = {}
                good_macs_by_loc[graph][loc]['min'] = 1.0
                rows = db.executeSqlCommand('select mac_address,count(*) from learn where location_uuid = "%s" group by mac_address order by count(*) desc' % (loc))
                maxCount = -1
                for row in rows:
                    if maxCount == -1:
                        maxCount = float(row[1])
                    weight =  float(row[1])/float(maxCount)
                    good_macs_by_loc[graph][loc][row[0]] = weight
                    if weight< good_macs_by_loc[graph][loc]['min']:
                        good_macs_by_loc[graph][loc]['min'] = weight
    db.close()

    logger.debug('Working on filtering macs')
    bad_macs_by_loc = OrderedDict()
    db = mlDB(group)
    macCount = {}
    for graph in range(len(connected_components)):
        bad_macs_by_loc[graph] = {}
        macCount[graph] = {}
        for loc in connected_components_locs[graph]:
            bad_macs_by_loc[graph][loc] = {}
            macCount[graph][loc] = {}
            # SELECT ONLY MACS FROM LOCATIONS IN THE CURRENT CONNECTED COMPONENT
            alocs = copy.deepcopy(connected_components_locs[graph])
            alocs.remove(loc)
            for aloc in alocs:
                for mac in good_macs_by_loc[graph][aloc].keys():
                    if mac not in bad_macs_by_loc[graph][loc]:
                        bad_macs_by_loc[graph][loc][mac] = 0
                    bad_macs_by_loc[graph][loc][mac] += good_macs_by_loc[graph][aloc][mac]
                    if mac not in macCount[graph][loc]:
                        macCount[graph][loc][mac] = 0
                    macCount[graph][loc][mac] += 1

    for graph in range(len(connected_components)):
        for loc in connected_components_locs[graph]:
            bad_macs_by_loc[graph][loc]['min'] = 1.0
            for mac in macCount[graph][loc].keys():
                bad_macs_by_loc[graph][loc][mac] = bad_macs_by_loc[graph][loc][mac] / macCount[graph][loc][mac]
                if bad_macs_by_loc[graph][loc][mac] < bad_macs_by_loc[graph][loc]['min']:
                    bad_macs_by_loc[graph][loc]['min'] = bad_macs_by_loc[graph][loc][mac]
    db.close()

    # Add in new baseline and renormalize
    for graph in range(len(connected_components)):
        baseline= numpy.ones(RSSI_PARTITIONS,dtype=numpy.dtype(numpy.float64) )*calculation_parameters[graph]['absentee']
        for mac in connected_components_macs[graph]:
            for loc in connected_components_locs[graph]:
                P[graph][mac][loc] = numpy.add(baseline,P[graph][mac][loc])
                nP[graph][mac][loc] = numpy.add(baseline,nP[graph][mac][loc])
                pTotal = numpy.sum(P[graph][mac][loc])
                if pTotal>0:
                    for i in range(RSSI_PARTITIONS):
                        P[graph][mac][loc][i] = float(P[graph][mac][loc][i])/float(pTotal)
                npTotal = numpy.sum(nP[graph][mac][loc])
                if npTotal>0:
                    for i in range(RSSI_PARTITIONS):
                        nP[graph][mac][loc][i] = float(nP[graph][mac][loc][i])/float(npTotal)

    # Pickle new calculations
    db = mlDB(group)
    db.insertResource('P',P)
    db.insertResource('nP',nP)
    db.insertResource('W',Wdefault)
    db.insertResource('usefulMacs',usefulMacs)
    db.insertResource('good_macs_by_loc',good_macs_by_loc)
    db.insertResource('bad_macs_by_loc',bad_macs_by_loc)
    # Load the new data sources into memory
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




def calculate_worker(num,total,totalThings,calculation_parameters,group):
    logger = logging.getLogger('priors.calculate_worker')
    db = mlDB(group)
    connected_components = db.getResource('connected_components')
    connected_components_locs = db.getResource('connected_components_locs')
    connected_components_macs = db.getResource('connected_components_macs')
    num_processors = cpu_count()

    P={}
    nP={}

    MIN_RSSI = -100
    MAX_RSSI = -10 #used to be -25
    RSSI_PARTITIONS=MAX_RSSI-MIN_RSSI+1
    for graph in range(len(connected_components)):
        P[graph]={}
        nP[graph]={}
        for mac in connected_components_macs[graph]:
            P[graph][mac]={}
            nP[graph][mac]={}
            for loc in connected_components_locs[graph]:
                P[graph][mac][loc] = numpy.zeros(RSSI_PARTITIONS,dtype=numpy.dtype(numpy.float64) )
                nP[graph][mac][loc] = numpy.zeros(RSSI_PARTITIONS,dtype=numpy.dtype(numpy.float64) )

    """thread worker function"""
    counter = 0;
    lastPrinted = ""
    for graph in range(len(connected_components)):
        normpdf = calculation_parameters[graph]['normpdf']
        for mac in connected_components_macs[graph]:
            for loc in connected_components_locs[graph]:
                counter = counter + 1.0
                if counter%total == num:
                    if round(counter/totalThings*100,1) % 10 == 0 and num == 0:
                        newPrint = "%s percent complete"%str(round(counter/totalThings*100,1))
                        if newPrint != lastPrinted:
                            lastPrinted = newPrint
                            logger.debug(newPrint)

                    rows = db.executeSqlCommand('select signal from learn indexed by ind_ex1 where mac_address = "%s" and location_uuid = "%s"'%(mac,loc))
                    #P[graph][mac][loc][0]=absentee
                    for row in rows:
                        P[graph][mac][loc][row[0]-MIN_RSSI]=P[graph][mac][loc][row[0]-MIN_RSSI]+normpdf[0]
                        P[graph][mac][loc][row[0]-MIN_RSSI-1]=P[graph][mac][loc][row[0]-MIN_RSSI-1]+normpdf[1]
                        P[graph][mac][loc][row[0]-MIN_RSSI+1]=P[graph][mac][loc][row[0]-MIN_RSSI+1]+normpdf[1]
                        P[graph][mac][loc][row[0]-MIN_RSSI-2]=P[graph][mac][loc][row[0]-MIN_RSSI-2]+normpdf[2]
                        P[graph][mac][loc][row[0]-MIN_RSSI+2]=P[graph][mac][loc][row[0]-MIN_RSSI+2]+normpdf[2]
                        P[graph][mac][loc][row[0]-MIN_RSSI-3]=P[graph][mac][loc][row[0]-MIN_RSSI-3]+normpdf[3]
                        P[graph][mac][loc][row[0]-MIN_RSSI+3]=P[graph][mac][loc][row[0]-MIN_RSSI+3]+normpdf[3]
                        P[graph][mac][loc][row[0]-MIN_RSSI-4]=P[graph][mac][loc][row[0]-MIN_RSSI-4]+normpdf[4]
                        P[graph][mac][loc][row[0]-MIN_RSSI+4]=P[graph][mac][loc][row[0]-MIN_RSSI+4]+normpdf[4]
                        P[graph][mac][loc][row[0]-MIN_RSSI+5]=P[graph][mac][loc][row[0]-MIN_RSSI+5]+normpdf[5]
                        P[graph][mac][loc][row[0]-MIN_RSSI+5]=P[graph][mac][loc][row[0]-MIN_RSSI+5]+normpdf[5]



    pickle.dump((P,nP),open(str(num)+'.p','wb'))
    db.close()
    return
