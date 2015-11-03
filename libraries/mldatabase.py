import builtins
import logging
import pickle
import os.path
import json
import time
from collections import OrderedDict
from libraries.DataBase import *

LoadedMemory = {}

def sanitizeSignal(signal):
    """Sanitizes RSSI values

    Converts RSSI signals to -100:0 range (standard on Android)
    """
    logger = logging.getLogger('fingerprinting:sanitizeSignal')
    if signal >= 0:
        signal = float(signal*0.5)-100
    return signal

    
        
class mlDB(DataBase):

    """ Initiates DataBase Object for Application """

    def __init__(self, group):
        logger = logging.getLogger('mlDB:init')
        #logger.debug('Connecting to database "' + path + '"')
        self.conn = sqlite3.connect(builtins.DATABASE_PATH_PREFIX + group + '.db')
        self.conn.row_factory = sqlite3.Row
        self.c = self.conn.cursor()
        self.group = group
        if group not in LoadedMemory:
            self.setupdb()
        LoadedMemory[group] = True
        

    def close(self):
        logger = logging.getLogger('mlDB:close')
        self.conn.close()
        self.conn = None
        self.c = None
        
        
    def insertFingerprint(self,route,data):
        """Handles database interactions for inserting a fingerprint

        Determines the table based on the route. The data is screened
        and the fingerprints are sanitized (in the case of computer-based
        RSSI signals) and 00:...:00 macs are ignored.
        """
        logger = logging.getLogger('mldatabase:insertFingerprint')
        route = route.replace('/','')
        gotFingerprint = False
        self.conn.isolation_level = None
        self.c = self.conn.cursor()
        self.c.execute('BEGIN')
        added = []
        timeC = int(time.time()*1000)
        for fingerprint in data['wifi-fingerprint']:
            if "00:00:00:00:00:00" not in fingerprint['mac']:
                if not gotFingerprint:
                    gotFingerprint = True
                username = data['username']
                rssi = sanitizeSignal(float(fingerprint['rssi']))
                mac = fingerprint['mac']
                if mac not in added:
                    if 'track' in route:
                        self.c.execute('INSERT into track (user_id,timestamp,mac_address,signal,location_uuid) values (?,?,?,?,?)',(username,timeC,mac,rssi,data['location']))             
                    elif 'learn' in route:
                        self.c.execute('INSERT into learn (user_id,timestamp,mac_address,signal,location_uuid) values (?,?,?,?,?)',(username,timeC,mac,rssi,data['location']))             
                    elif 'test' in route:
                        self.c.execute('INSERT into test (user_id,timestamp,mac_address,signal,location_uuid) values (?,?,?,?,?)',(username,timeC,mac,rssi,data['location']))
                added.append(mac)                

        if gotFingerprint:
            logger.debug('Successfully inserted fingerprints to "' + route + '"')
        a = self.c.execute('COMMIT')
        return gotFingerprint
        
    def insertFingerprints(self,route,datas):
        """Handles database interactions for inserting a fingerprint

        Determines the table based on the route. The data is screened
        and the fingerprints are sanitized (in the case of computer-based
        RSSI signals) and 00:...:00 macs are ignored.
        """
        logger = logging.getLogger('mldatabase:insertFingerprint')
        route = route.replace('/','')
        gotFingerprint = False
        self.conn.isolation_level = None
        self.c = self.conn.cursor()
        self.c.execute('BEGIN')
        for data in datas:
            for fingerprint in data['wifi-fingerprint']:
                if "00:00:00:00:00:00" not in fingerprint['mac']:
                    if not gotFingerprint:
                        gotFingerprint = True
                    username = data['username']
                    time = float(data['time'])
                    rssi = sanitizeSignal(float(fingerprint['rssi']))
                    mac = fingerprint['mac']
                    if 'track' in route:
                        self.c.execute('INSERT into track (user_id,timestamp,mac_address,signal,location_uuid) values (?,?,?,?,?)',(username,time,mac,rssi,data['location']))             
                    elif 'learn' in route:
                        self.c.execute('INSERT into learn (user_id,timestamp,mac_address,signal,location_uuid) values (?,?,?,?,?)',(username,time,mac,rssi,data['location']))             
                    elif 'test' in route:
                        self.c.execute('INSERT into test (user_id,timestamp,mac_address,signal,location_uuid) values (?,?,?,?,?)',(username,time,mac,rssi,data['location']))
                
        if gotFingerprint:
            logger.debug('Successfully inserted ' + str(len(datas)) + ' fingerprints to "' + route + '"')
        self.c.execute('COMMIT')
        return gotFingerprint
        
    def getUniqueMacs(self):
        self.c.execute('SELECT DISTINCT mac_address FROM learn')
        rows =  self.c.fetchall()
        foos = []
        for row in rows:
            foos.append(str(row[0]))
        return foos
        
    def getNthLastFingerprint(self,user,n):
        query = """select mac_address,signal from track where timestamp=(select distinct timestamp from (select timestamp from track where user_id like '%s' order by id desc limit 500) limit 1 offset %d)""" % (user,n)
        data = {}
        self.c.execute(query)
        rows =  self.c.fetchall()
        for row in rows:
            data[row[0]] = row[1]
        if len(data)>0:
            return data
        else:
            return None
        
    def getLastLocationFromTracking(self,user):
        query = """select location_uuid,timestamp from track where user_id like '%s' order by id desc limit 1""" % (user)
        data = {}
        self.c.execute(query)
        rows =  self.c.fetchall()
        for row in rows:
            data['location'] = row[0]
            data['timestamp'] = row[1]
        if len(data)>0:
            return data
        else:
            return None   

    def getLastLocationsFromTracking(self,user,num=10):
        query = """select name from sqlite_master where type='index' and name='track_idx'"""
        self.c.execute(query)
        rows =  self.c.fetchall()
        if len(rows)==0:
            self.executeSqlCommand('create index track_idx on track(user_id,timestamp)')

        query = """select max(id) from track indexed by track_idx where user_id='%s'""" % user
        self.c.execute(query)
        rows =  self.c.fetchall()
        maxId = -1
        for row in rows:
            maxId = int(row[0])
        if maxId == -1:
            return [{'location':'unknown','time':int(time.time()*1000)}]

        query = """select location_uuid,timestamp from track indexed by track_idx 
            where id > %s and
            user_id like '%s' 
            group by timestamp order by timestamp desc limit %s""" % (str(maxId-50000),user,str(num))
        self.c.execute(query)
        rows =  self.c.fetchall()
        locs = []
        for row in rows:
            data = {}
            data['location'] = row[0]
            data['time'] = row[1]
            locs.append(data)
        if len(locs)>0:
            return locs
        else:
            return [{'location':'unknown','time':int(time.time()*1000)}]
            
    def databaseStats(self):
        tables = ['learn','test']
        data = {}
        for table in tables:
            query = """select location_uuid,count(*) from (select location_uuid from %s group by timestamp) group by location_uuid""" % table
            data[table] = {}
            self.c.execute(query)
            rows =  self.c.fetchall()
            for row in rows:
                data[table][row[0]] = row[1]
        print(data)
        return data

    def getUniqueLocations(self):
        self.c.execute('SELECT DISTINCT location_uuid FROM learn')
        rows =  self.c.fetchall()
        foos = []
        for row in rows:
            foos.append(str(row[0]))
        return foos
    
    def getUsers(self):
        self.c.execute('SELECT DISTINCT user_id FROM track')
        rows =  self.c.fetchall()
        foos = []
        for row in rows:
            foos.append(str(row[0]))
        return foos
    
    def getUniqueTimestamps(self,table):
        self.c.execute('SELECT DISTINCT timestamp FROM %s' % table)
        rows =  self.c.fetchall()
        foos = []
        for row in rows:
            foos.append(int(row[0]))
        return foos

    def deleteLocation(self,location):
        logger = logging.getLogger('DataBase:deleteLocation')
        try:
            self.c.execute("DELETE FROM learn WHERE location_uuid like '%s'" % str(location))
            self.conn.commit()
            self.c.execute("DELETE FROM test WHERE location_uuid like '%s'" % str(location))
            self.conn.commit()
            return True,"Successfully deleted. Re-calculate to update parameters."
        except Exception as e:
            logger.error(e)
            return False,str(e)

    def editLocationName(self,location,newname):
        logger = logging.getLogger('DataBase:deleteLocation')
        try:
            self.c.execute("UPDATE learn set location_uuid = '%s' where location_uuid = '%s'" % (str(newname),str(location)))
            self.conn.commit()
            self.c.execute("UPDATE track set location_uuid = '%s' where location_uuid = '%s'" % (str(newname),str(location)))
            self.conn.commit()
            self.c.execute("UPDATE test set location_uuid = '%s' where location_uuid = '%s'" % (str(newname),str(location)))
            self.conn.commit()
            return True,"Successfully renamed."
        except Exception as e:
            logger.error(e)
            return False,str(e)
        
    def retrieveFingerprint(self,table,timestamp):
        query = """SELECT user_id,timestamp,mac_address,signal,location_uuid FROM %s where timestamp=%s""" % (table,str(timestamp))
        if table=='test':
            query = """SELECT user_id,timestamp,mac_address,signal,location_uuid FROM %s indexed by ind_ex4 where timestamp=%s""" % (table,str(timestamp))
        self.c.execute(query)
        data = self.c.fetchall()
        fingerprint = {}
        fingerprint['time'] = timestamp
        fingerprint['wifi-fingerprint'] = []
        for dat in data:
            fingerprint['location'] = dat[4]
            fingerprint['wifi-fingerprint'].append({'mac':dat[2],'rssi':dat[3]})
        return fingerprint
            
            

    def insertFingerprintData(self,\
            username,\
            location_info_id,\
            mac,\
            rssi,\
            time,\
            table_type):
        if location_info_id == None:
            return self.addDataBulk(data={\
                'user_id':username,\
                'timestamp':time,\
                'mac_address':mac,\
                'signal':rssi,\
                },table=table_type)
        else:
            return self.addDataBulk(data={\
                'user_id':username,\
                'timestamp':time,\
                'mac_address':mac,\
                'signal':rssi,\
                'location_uuid':location_info_id\
                },table=table_type)
   
    def insertResource(self, uuid, obj):
        """ Inserts or updates datasource in table"""
        logger = logging.getLogger('DataBase:insertResource')
        pObj = pickle.dumps(obj)
        try:
            data = self.getData(
                table='resources',
                parameters={
                    'uuid': uuid})
            if len(data) == 0:
                binary = sqlite3.Binary(pObj)
                query = u'''INSERT INTO resources (uuid,pickle) VALUES (?,?)'''
                payload = (uuid, binary)
                self.c.execute(query, payload)
                self.conn.commit()
            else:
                binary = sqlite3.Binary(pObj)
                query = u'''UPDATE resources SET pickle=(?) WHERE uuid=(?)'''
                payload = (binary, uuid)
                self.c.execute(query, payload)
                self.conn.commit()
        except Exception as e:
            logger.error(e)

    def getResource(self, uuid):
        """ Requests individual datasource """
        logger = logging.getLogger('DataBase:getResource')
        #logger.debug('Getting resource: %s' % uuid)
        data = self.getData(table='resources', parameters={'uuid': uuid})
        realData = None
        for row in data:
            realData = pickle.loads(row[2])
        if realData == None:
            logger.warn("No data source")
            return None
        else:
            return realData

    def archiveTrack(self):
        """ Archives tracking database """
        logger = logging.getLogger('DataBase:archiveTrack')
        cmd = """SELECT DISTINCT user_id FROM track"""
        #logger.debug(cmd)
        self.c.execute(cmd)
        rows = self.c.fetchall()
        users = []
        for row in rows:
            users.append(row[0])
        print(users)
        if len(users)==0:
            return True

        for user in users:
            cmd = """SELECT id FROM track WHERE user_id='%s' GROUP BY timestamp ORDER BY timestamp DESC LIMIT 1 OFFSET 100""" % user
            #logger.debug(cmd)
            self.c.execute(cmd)
            rows = self.c.fetchall()
            maxid = "-1"
            for row in rows:
                maxid = str(row[0])
            cmd = """INSERT INTO track_archive SELECT * from track WHERE user_id='%s' AND id<%s""" % (user,maxid)
            #logger.debug(cmd)
            self.c.execute(cmd)
            self.conn.commit()
            cmd = """DELETE from track WHERE user_id='%s' AND id<%s""" % (user,maxid)
            #logger.debug(cmd)
            self.c.execute(cmd)
            self.conn.commit()
            #logger.debug('Archived rows with id < ' + maxid + ' for ' + user)

        return True
    


    def haveResource(self, uuid):
        """ Requests individual datasource """
        logger = logging.getLogger('DataBase:haveResource')
        data = self.getData(table='resources', parameters={'uuid': uuid})
        return len(data)>0

    def setupdb(self):
        """ Sets up database tables"""
        logger = logging.getLogger('DataBase:setupdb')
        logger.debug('Setting up database...')
        if not self.tableExists('resources'):
            logger.debug('Creating "resources" table...')
            self.createTable('resources ('
                             + 'id INTEGER PRIMARY KEY AUTOINCREMENT, '\
                             #    +'uuid TEXT UNIQUE,'\
                             + 'uuid TEXT,'\
                             + 'pickle BLOB '\
                             + ')')
        else:
            logger.debug('Table "resources" found')
        
        cons = [0]
        calculation_parameters = {}
        try: 
            cons = self.getResource('connected_components')
            calculation_parameters = self.getResource('calculation_parameters')
            foo = {}
            for graph in range(len(cons)):
                foo[graph] = calculation_parameters[graph]
            logger.debug('Calculation parameters exist')
            logger.debug(json.dumps(foo,indent=4))
        except: 
            logger.debug('Adding calculation default parameters')
            parameters = {}
            parameters['pdf_type'] = 6
            parameters['absentee'] = 0.00001
            parameters['usefulness'] = 0
            parameters['dropout_percentage'] = 10
            parameters['persistence'] = 1
            parameters['number_of_simulations'] = 3
            parameters['mix_in'] = 0.5
            parameters['trigger_server'] = 'None'
            calculation_parameters = {}
            calculation_parameters[0] = parameters
            logger.debug(calculation_parameters)
            self.insertResource('calculation_parameters',calculation_parameters)
        
        try:
            builtins.PARAMETERS[self.group] = calculation_parameters
        except:
            builtins.PARAMETERS = {}
            builtins.PARAMETERS[self.group] = calculation_parameters
            
        if not self.tableExists('learn'):
            logger.debug('Creating "learn" table...')
            self.createTable('learn ('\
                +'id INTEGER PRIMARY KEY AUTOINCREMENT, '\
                +'user_id TEXT, '\
                +'timestamp INTEGER, '\
                +'mac_address INTEGER, '\
                +'signal INTEGER, '\
                +'location_uuid TEXT'\
                +')')
        else:
            logger.debug('Table "learn" found')
            
        if not self.tableExists('track'):
            logger.debug('Creating "track" table...')
            self.createTable('track ('\
                +'id INTEGER PRIMARY KEY AUTOINCREMENT, '\
                +'user_id TEXT, '\
                +'timestamp INTEGER, '\
                +'mac_address INTEGER, '\
                +'signal INTEGER, '\
                +'location_uuid TEXT'\
                +')')
        else:
            logger.debug('Table "track" found')
            
        if not self.tableExists('track_archive'):
            logger.debug('Creating "track_archive" table...')
            self.createTable('track_archive ('\
                +'id INTEGER PRIMARY KEY, '\
                +'user_id TEXT, '\
                +'timestamp INTEGER, '\
                +'mac_address INTEGER, '\
                +'signal INTEGER, '\
                +'location_uuid TEXT'\
                +')')
        else:
            logger.debug('Table "track" found')
            
        if not self.tableExists('test'):
            logger.debug('Creating "test" table...')
            self.createTable('test ('\
                +'id INTEGER PRIMARY KEY AUTOINCREMENT, '\
                +'user_id TEXT, '\
                +'timestamp INTEGER, '\
                +'mac_address INTEGER, '\
                +'signal INTEGER, '\
                +'location_uuid TEXT'\
                +')')
        else:
            logger.debug('Table "test" found')

        logger.debug('Database ready')
