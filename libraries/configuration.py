import os
import configparser
import uuid

_file = os.path.join('data', 'settings.cfg')  # config file

def _create():
    # create config file
    config = configparser.ConfigParser()
    config['server'] = {}
    config['server']['address'] = str(input("address: "))
    config['server']['port'] = str(input("port: "))
    config['server']['ext_address'] = str(input("external address (leave blank to use " + config['server']['address'] + ":" + config['server']['port'] + "): "))
    if len(config['server']['ext_address'])<3:
        config['server']['ext_address'] = config['server']['address'] + ":" + config['server']['port']
    config['server']['apikey'] = str(uuid.uuid4())
    config['admin'] = {}
    config['admin']['master'] = str(uuid.uuid4())
    with open(_file, 'w') as configfile:
        config.write(configfile)
    return config

def _open():
    # read config file
    config = configparser.ConfigParser()
    config.read(_file)
    return config

def _getConf():
    # check if file exists
    if not os.path.isfile(_file):
        config = _create()
        return config
    else:
        config = _open()
        return config

def _getSession(conf):
    if len(conf['server']['port'])<1:
        conf['server']['port'] = ''
    session = {
        'master': conf['admin']['master'],
        'address': conf['server']['address'],
        'ext_address': conf['server']['ext_address'],
        'port': conf['server']['port'],
        'apikey': conf['server']['apikey']
    }
    return session

def getOptions():
    # get session
    conf = _getConf()
    session = _getSession(conf)
    return session


