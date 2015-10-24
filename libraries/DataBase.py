"""DataBase parent class
Carries generic features for connecting to database and adding/removing/searching
"""

import sqlite3
import uuid
import builtins
import logging

__author__ = "Stefan"
__copyright__ = "Copyright 2015, FIND"
__credits__ = ["Zack", "Stefan", "Travis"]
__license__ = "MIT"
__version__ = "0.1"
__status__ = "Development"


class Error(Exception):
    """Base class for exceptions in this module."""
    pass


class DataBaseInvalidInput(Error):
    """Exception raised for errors in the input.

    Attributes:
        expr -- input expression in which the error occurred
        msg  -- explanation of the error
    """

    def __init__(self, expr, msg):
        self.expr = expr
        self.msg = msg


class DataBase(object):
    """ Basic DataBase Object """

    def __init__(self, path):
        self.conn = sqlite3.connect(path)
        # returns row objects instead of plain tuples
        self.conn.row_factory = sqlite3.Row
        self.c = self.conn.cursor()

    def close(self):
        self.conn.close()
        self.conn = None
        self.c = None

    def executeSqlCommand(self, command):
        """ WARNING EXTREMELY DANGEROUS """
        self.c.execute(command)
        return self.c.fetchall()


    def addData(self, data, table):
        """Inserts rows into table 
        
        data given as tuple of values in order of table columns
        returns rowid of new row or -1 if integrity error       
        """
        logger = logging.getLogger('DataBase:addData')
        try:
            # None is for the primary key which auto increments
            if isinstance(data, list) or isinstance(data, tuple):
                self.c.execute(
                    'INSERT INTO %s VALUES (%s)' %
                    (table,
                     ('?,' *
                      len(data))[
                         :-
                         1]),
                    data)
            elif isinstance(data, dict):
                keys = list(data.keys())
                try:
                    values = [data[k] for k in keys]
                except KeyError:
                    raise DataBaseInvalidInput(
                        (data, table), 'sanitized column names don\'t match given column names')
                self.c.execute(
                    'INSERT INTO %s (%s) VALUES (%s)' %
                    (table,
                     ','.join(keys),
                        ','.join(
                         ['?'] *
                         len(data))),
                    values)
            else:
                raise Exception('invalid input type: %s' % type(data))
            id = self.c.lastrowid
            # remember to commit changes so we don't lock the db!
            self.conn.commit()
            return id
        except sqlite3.IntegrityError as e:
            logger.error('dbcommands.addData %s' % e)

    def removeData(self, field, value, table):
        """ Removes data from table given parameter args """
        self.c.execute('DELETE FROM %s WHERE (?)=(?)' % table, (field, value))
        self.conn.commit()

    def getData(self, table, parameters={}):
        """ Retrieves data from table given parameter args """
        query = "SELECT * FROM %s" % table
        keys = parameters.keys()
        query_extra = ' AND '.join(k + ' = ?' for k in keys)
        query += (' WHERE ' + query_extra) if query_extra is not '' else ''
        self.c.execute(query, [parameters[k] for k in keys])
        return self.c.fetchall()

    # def contains(self,id,table):
    #     self.c.execute('SELECT EXISTS(SELECT 1 FROM %s WHERE id=(?) LIMIT 1)'%table,(id,))
    #     return self.c.fetchone()[0] is 1

    def contains(self, field, value, table):
        """ Checks if field value exists in table """
        self.c.execute(
            'SELECT EXISTS(SELECT 1 FROM %s WHERE (?)=(?) LIMIT 1)' %
            table, (field, value))
        return self.c.fetchone()[0] is 1

    def showTables(self):
        """ Returns tables in database """
        return [a[0] for a in self.c.execute(
            "SELECT name FROM sqlite_master WHERE type='table';")]

    def tableExists(self, table_name):
        """ Checks if a table exists. returns boolean """
        self.c.execute(
            'SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type="table" AND name=? LIMIT 1)',
            (table_name,
             ))
        return not self.c.fetchone()[0] is 0

    def createTable(self, table_data):
        """ Creates a database table """
        logger = logging.getLogger('DataBase:createTable')
        if self.tableExists(table_data):
            return False
        else:
            self.c.execute('CREATE TABLE %s;' % table_data)
            self.conn.commit()
            logger.info('Table created')
            return True

    def dropTable(self, table):
        """ Drops a database table """
        if self.tableExists(table):
            self.c.execute('DROP TABLE %s' % table)
            self.conn.commit()
            return True
        else:
            return False

    def showColumns(self, table):
        """ Prints description of columns within a table """
        names = self.columnNames(table)
        types = self.columnTypes(table)
        for i in range(len(names)):
            print("%s - %s" % (names[i], types[i]))

    def columnNames(self, table):
        """ Returns column names within a database table """
        return [r[1] for r in self.c.execute('PRAGMA table_info(%s)' % table)]

    def columnTypes(self, table):
        """ Returns column types within a database table """
        return [r[2] for r in self.c.execute('PRAGMA table_info(%s)' % table)]


class geoDB(DataBase):

    """ Initiates DataBase Object for Application """

    def __init__(self, path):
        self.conn = sqlite3.connect(path)
        # returns row objects instead of plain tuples
        self.conn.row_factory = sqlite3.Row
        self.c = self.conn.cursor()

    def close(self):
        self.conn.close()
        self.conn = None
        self.c = None

    def insertResource(self, uuid, type, pickle):
        """ Inserts or updates datasource in table"""
        logger = logging.getLogger('DataBase:insertResource')
        try:
            data = self.getData(
                table='resources',
                parameters={
                    'type': type,
                    'uuid': uuid})
            if len(data) == 0:
                binary = sqlite3.Binary(pickle)
                query = u'''INSERT INTO resources (uuid,type,pickle) VALUES (?,?,?)'''
                payload = (uuid, type, binary)
                self.c.execute(query, payload)
                self.conn.commit()
            else:
                binary = sqlite3.Binary(pickle)
                query = u'''UPDATE resources SET pickle=(?) WHERE uuid=(?) AND type=(?)'''
                payload = (binary, uuid, type)
                self.c.execute(query, payload)
                self.conn.commit()
        except Exception as e:
            logger.error(e)

    def getResources(self, type):
        """ Requests all datasources of given type """
        logger = logging.getLogger('DataBase:getResources')
        data = self.getData(table='resources', parameters={'type': type})
        if len(data) == 0:
            logger.info("no datasources found")
            return False
        else:
            return data

    def getResource(self, uuid):
        """ Requests individual datasource """
        logger = logging.getLogger('DataBase:getResource')
        data = self.getData(table='resources', parameters={'uuid': uuid})
        if len(data) == 0:
            logger.info("no datasources found")
            return False
        else:
            return data

    def setupdb(self):
        """ Sets up database tables"""
        
        print('databaseSetup.setupDatabase')
        print('Setting up database...')
        if not self.tableExists('resources'):
            print('"resources" table not found')
            print('Creating "resources" table...')
            self.createTable('resources ('
                             + 'id INTEGER PRIMARY KEY AUTOINCREMENT, '\
                             #    +'uuid TEXT UNIQUE,'\
                             + 'uuid TEXT,'\
                             + 'type TEXT,'\
                             + 'pickle BLOB '\
                             + ')')
        else:
            print('Table "resources" found')
        print('Database ready')


'''

import sqlite3
from io import StringIO

def launch_sqlitedb_inMemory(database):
    """ Reads sqliteDB to tempfile and builds in memory database """
    # Read database to tempfile
    con = sqlite3.connect(database)
    tempfile = StringIO()
    for line in con.iterdump():
        tempfile.write('%s\n' % line)
    con.close()
    tempfile.seek(0)
    # Create a database in memory and import from tempfile
    sqlite = sqlite3.connect(":memory:")
    sqlite.cursor().executescript(tempfile.read())
    sqlite.commit()
    sqlite.row_factory = sqlite3.Row
    return sqlite

class DataBase(object):
    """ Database object """
    def __init__(self,database,*args, **kwargs):
        self._source = database
        if args:
            print(args)
        if kwargs:
            if kwargs['inMemory'] == True:
                self.conn = launch_sqlitedb_inMemory(database)
        else:
            self.conn = sqlite3.connect(database)
        self.cursor = self.conn.cursor()


'''