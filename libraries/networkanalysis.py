import builtins
import itertools
import time
import threading
from multiprocessing import Process, cpu_count
import pickle
import os
import logging
import json
import networkx as nx
from networkx.readwrite import json_graph
from libraries.mldatabase import *

def complete_graph_from_list(L, create_using=None):
    G2=nx.empty_graph(len(L),create_using)
    edges=itertools.combinations(L,2)
    G2.add_edges_from(edges)
    return G2

finishedSets = []
def worker(group,num,total):
    logger = logging.getLogger('networkanalysis:worker')
    G=nx.empty_graph()
    global finishedSets
    global skipped
    """thread worker function"""
    db = mlDB(group)
    timestamps = db.executeSqlCommand('select timestamp,id from learn indexed by ind_ex3 group by timestamp')
    totalThings = 0
    for timestamp in timestamps:
        totalThings = totalThings + 1.0
    if num == 0:
        logger.debug('Total number of points: %d' % totalThings)
        logger.debug('Collecting edges')
    timestamps = db.executeSqlCommand('select * from (select timestamp,id from learn indexed by ind_ex3 group by timestamp) where id %% %s == %s' %(str(total),str(num)))
    edges = []
    count = 0
    t=time.time()
    displayedPercent = 0
    for timestamp in timestamps:
        count = count + 1.0
        percentDone = int(count/totalThings*total*100.0)

        if num == 0 and int(percentDone) % 10 == 0 and int(percentDone) != displayedPercent:
                logger.debug('%s percent done'%str(round(percentDone)))
                displayedPercent = int(percentDone)
        macs_db = db.executeSqlCommand('select mac_address from learn indexed by ind_ex3 where timestamp=%s'%timestamp[0])
        macs = []
        for mac_db in macs_db:
            macs.append(mac_db[0])
        hashedSet = hash(frozenset(macs))
        if hashedSet not in finishedSets:
            G = nx.compose(G,complete_graph_from_list(macs))
            finishedSets.append(hashedSet)
    pickle.dump(G,open(str(num)+'.p','wb'))
    db.close()
    return



def generateGraph(group):
    logger = logging.getLogger('networkanalysis:generateGraph')
    connected_components = []
    threads = []
    t=time.time()
    db = mlDB(group)
    try:
        db.executeSqlCommand('drop index ind_ex3')
    except:
        pass
    db.executeSqlCommand('create index ind_ex3 on learn(timestamp)')
    db.close()
    numThreads = cpu_count()
    
    for i in range(numThreads):
        tThread = Process(target=worker, args=(group,i,numThreads,))
        threads.append(tThread)
    for tThread in threads:
        tThread.start()
    for tThread in threads:
        tThread.join()

    G=nx.empty_graph()    
    for i in range(numThreads):
        G0 = pickle.load(open(str(i)+'.p','rb'))
        G = nx.compose(G,G0)
        os.system('rm ' + str(i)+'.p')
        
    '''
    G = nx.compose(G,complete_graph_from_list(["a", "b", "c", "d"]))
    print(G.edges())
    G = nx.compose(G,complete_graph_from_list(["d","e"]))
    print(G.edges())
    '''
    
    '''
    db = mlDB(group)
    timestamps = db.getSqlCommand('select timestamp,id from location_learn group by timestamp')
    totalThings = 0
    for timestamp in timestamps:
        totalThings = totalThings + 1.0
        
    timestamps = db.getSqlCommand('select timestamp,id from location_learn group by timestamp')
    edges = []
    print('collecting edges...')
    count = 0
    
    for timestamp in timestamps:
        count = count + 1.0
        percentDone = count/totalThings*100.0
        if percentDone % 10 == 0:
            print('%s percent done'%str(round(percentDone)))
        macs_db = db.getSqlCommand('select mac_address from location_learn where timestamp=%s'%timestamp[0])
        macs = []
        for mac_db in macs_db:
            macs.append(mac_db[0])
        hashedSet = hash(frozenset(macs))
        if hashedSet not in finishedSets:
            G = nx.compose(G,complete_graph_from_list(macs))
            finishedSets.append(hashedSet)

        for i in range(0,len(macs)-1):
            for j in range(i,len(macs)):
                if i is not j:
                    if (macs[i],macs[j]) not in edges:
                        edges.append((macs[i],macs[j]))
        '''
    logger.debug('Finished collecting edges in %s seconds' % "{0:.2f}".format(round(time.time()-t,2)))
    outdeg = G.degree()
    to_remove = [n for n in outdeg if outdeg[n] == 0]
    G.remove_nodes_from(to_remove)                  
    #G.add_edges_from([('a','b'),('a','c'),('a','d'),('a','e'),('b','c'),('b','d'),('b','e'),('c','d'),('c','e'),('d','e'),('e8:fc:af:81:4f:d4','a'),('22:10:7a:ed:1d:87','b')]) #,('e8:fc:af:81:4f:d4','a')
    cliques = sorted(nx.find_cliques(G), key = len, reverse=True)

    # Need algorithm where you keep cutting until
    # you only cut off one node at a time
    # then stop cutting
    removeEdges = []
    for k in nx.connected_component_subgraphs(G):
        isOne = False
        cuts = nx.minimum_edge_cut(k)
        for edge in cuts:
            k.remove_edge(edge[0],edge[1])
        connected_components = list(nx.connected_components(k))
        for connected_component in connected_components:
            if len(connected_component) == 1:
                isOne = True
        if isOne:
            pass
        else:
            for cut in cuts:
                removeEdges.append(cut)

    if len(removeEdges)>0:
        print("NEED TO REMOVE EDGES!")
        print(removeEdges)
    #for edge in removeEdges:
    #    G.remove_edge(edge[0],edge[1])
    connected_components = list(nx.connected_components(G))
    if len(G.nodes())<50:
        pos=nx.shell_layout(G) # positions for all nodes
    else:
        pos=nx.spring_layout(G,scale=10000) # positions for all nodes
    #pos=nx.spectral_layout(G) # positions for all nodes

    N = len(connected_components)

    labels = {}
    for node in G.nodes():
        labels[node]="%s"%node

    nx.set_node_attributes(G,'pos',pos)
    connected_components_locs={}
    connected_components_macs={}

    for i in range(len(connected_components)):
        connected_components[i] = list(set(connected_components[i]))
        connected_components_locs[i]=[]
        connected_components_macs[i]=[]
        locations = "Cluster " + str(i) + "\n"
        num = 0
        maxX = 0
        minY = 1000000
        for component in connected_components[i]:
            connected_components_macs[i].append(component)
            G.node[component]['name']=component
            if G.node[component]['pos'][0]>maxX:
                maxX = G.node[component]['pos'][0]
            if G.node[component]['pos'][1]<minY:
                minY = G.node[component]['pos'][1]
            num = num + 1
            db = mlDB(group)
            locs_uuid = db.executeSqlCommand('select location_uuid from learn where mac_address like "%s" group by location_uuid'%component)
            db.close()
            for loc in locs_uuid:
                connected_components_locs[i].append(loc[0])

        connected_components_locs[i]=list(set(connected_components_locs[i]))
        connected_components_macs[i]=list(set(connected_components_macs[i]))


    db = mlDB(group)
    db.insertResource('connected_components_macs',connected_components_macs)
    db.insertResource('connected_components_locs',connected_components_locs)
    db.insertResource('connected_components',connected_components)
    db.insertResource('G',G)
    calculation_parameters = db.getResource('calculation_parameters')
    for con in range(len(connected_components)):
        if con not in calculation_parameters.keys():
            calculation_parameters[con] = calculation_parameters[0]
    builtins.PARAMETERS[group] = calculation_parameters
    db.insertResource('calculation_parameters',calculation_parameters)
    
    
    # Create index if it hasn't already been created
    logger.debug('Creating the learning index')
    try:
        db.executeSqlCommand('drop index ind_ex1')
        logger.debug('Dropped previous index')
    except:
        logger.debug('No index to drop')
    db.executeSqlCommand('create index ind_ex1 on learn(mac_address, location_uuid)')
    logger.debug('Generated new index on learn')

    
    logger.debug('Creating the training index')
    try:
         db.executeSqlCommand('drop index ind_ex4')
    except:
        pass
    db.executeSqlCommand('create index ind_ex4 on test(timestamp)')
    
    db.close()
    return len(connected_components)


def makeNetworkJson(group):
    db = mlDB(group)
    G = db.getResource('G')
    db.close()
    d = json_graph.node_link_data(G)
    for i in range(len(d['nodes'])):
        temp =d['nodes'][i]['pos']
        d['nodes'][i]['pos'] = []
        d['nodes'][i]['pos'].append(int(temp[0]))
        d['nodes'][i]['pos'].append(int(temp[1]))
    return d

def determineComponent(testmac):
    foundComponent = False
    componentNum = -1
    for i in range(len(connected_components)):
        if testmac in connected_components[i]:
            foundComponent = True
            componentNum = i
            break
    return componentNum

'''
testmacs = ['90:1a:ca:78:bb:35','c8:b3:73:35:f2:af','90:1a:ca:78:bb:35']
generateGraph()
print(determineComponent(testmacs))
'''
