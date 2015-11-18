"""Machine Learning server
Responds to cell phone data and can calculate location classifications
"""

from flask import Flask, request, jsonify, make_response, current_app, render_template, url_for, redirect, send_from_directory, Response
import flask.ext.login as flask_login
from werkzeug import secure_filename
import logging
import builtins
import time
import copy
import json
import os
from glob import glob
from datetime import timedelta
from functools import update_wrapper
from libraries.mldatabase import *
from libraries.networkanalysis import *
from libraries.priors import *
from libraries.posteriors import *
from libraries.analysis import *
from libraries.configuration import *
from libraries.livestats import *

from apscheduler.schedulers.background import BackgroundScheduler

__author__ = "Zack Scholl"
__copyright__ = "Copyright 2015, FIND"
__credits__ = ["Zack Scholl", "Stefan Safranek"]
__license__ = "GNU AFFERO GENERAL PUBLIC LICENSE"
__version__ = "0.2"
__maintainer__ = "Zack Scholl"
__email__ = "zack@hypercubeplatforms.com"
__status__ = "Development"

logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s %(name)-12s %(levelname)-8s %(message)s',
    datefmt='%m-%d %H:%M:%S',
    handlers=[
        logging.StreamHandler(), logging.FileHandler('server.log')])



# Patch because basestring doesn't work for crossdomain()
# https://github.com/oxplot/fysom/issues/1
try:
    unicode = unicode
except NameError:
    # 'unicode' is undefined, must be Python 3
    str = str
    unicode = str
    bytes = bytes
    basestring = (str, bytes)
else:
    # 'unicode' exists, must be Python 2
    str = str
    unicode = unicode
    bytes = str
    basestring = basestring

GENERATE_UNIT_TESTS = False 
UPLOAD_FOLDER = 'data/'
ALLOWED_EXTENSIONS = set(['gpx'])
AVAILABLE_DBS = []



app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER
app.secret_key = 'YOUR_SECRET_KEY'
login_manager = flask_login.LoginManager()
login_manager.init_app(app)

builtins.GROUPDATABASE = {'find':{'last_seen':123},'stefangroup':{'last_seen':123}}

class User(flask_login.UserMixin):
    pass

@login_manager.user_loader
def user_loader(group):
    if group not in builtins.GROUPDATABASE:
        pass

    user = User()
    user.id = group
    return user


def groupExists(group):
    for file in os.listdir("data/"):
        if file.endswith(".db"):
            if group + '.db' == file:
                return True
    return False



@login_manager.request_loader
def request_loader(request):
    group = request.form.get('group')
    if not groupExists(group):
        pass

    user = User()
    user.id = group

    # DO NOT ever store passwords in plaintext and always compare password
    # hashes using constant-time comparison!
    user.is_authenticated = True

    return user

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'GET':
        try:
            group = request.args.get('group').lower()
            user = User()
            user.id = group
            flask_login.login_user(user)
            return redirect('/')
        except:
            return open('libraries/templates/login.html','r').read()
    if groupExists(request.form['group'].lower()):
        user = User()
        user.id = request.form['group'].lower()
        flask_login.login_user(user)
        return redirect('/')

    return open('libraries/templates/login_error.html','r').read()


@app.route('/protected')
@flask_login.login_required
def protected():
    return 'Logged in as: ' + flask_login.current_user.id


@app.route('/logout')
def logout():
    flask_login.logout_user()
    return redirect(url_for('login'))


def allowed_file(filename):
    return '.' in filename and \
        filename.rsplit('.', 1)[1] in ALLOWED_EXTENSIONS


def crossdomain(origin=None, methods=None, headers=None,
                max_age=21600, attach_to_all=True,
                automatic_options=True):
    """ Decorator for the HTTP access control

    From http://flask.pocoo.org/snippets/56/

    Cross-site HTTP requests are HTTP requests for resources from a different domain 
    than the domain of the resource making the request. For instance, a resource loaded from 
    Domain A makes a request for a resource on Domain B. The way this is implemented in 
    modern browsers is by using HTTP Access Control headers: Documentation on developer.mozilla.org.
    """
    if methods is not None:
        methods = ', '.join(sorted(x.upper() for x in methods))
    if headers is not None and not isinstance(headers, basestring):
        headers = ', '.join(x.upper() for x in headers)
    if not isinstance(origin, basestring):
        origin = ', '.join(origin)
    if isinstance(max_age, timedelta):
        max_age = max_age.total_seconds()

    def get_methods():
        if methods is not None:
            return methods

        options_resp = current_app.make_default_options_response()
        return options_resp.headers['allow']

    def decorator(f):
        def wrapped_function(*args, **kwargs):
            if automatic_options and request.method == 'OPTIONS':
                resp = current_app.make_default_options_response()
            else:
                resp = make_response(f(*args, **kwargs))
            if not attach_to_all and request.method != 'OPTIONS':
                return resp

            h = resp.headers

            h['Access-Control-Allow-Origin'] = origin
            h['Access-Control-Allow-Methods'] = get_methods()
            h['Access-Control-Max-Age'] = str(max_age)
            if headers is not None:
                h['Access-Control-Allow-Headers'] = headers
            return resp

        f.provide_automatic_options = False
        return update_wrapper(wrapped_function, f)
    return decorator


"""
Routes to handle web pages
"""

@app.route("/")
def landing():
    try:
        group = flask_login.current_user.id
    except:
        return redirect(url_for('login'))
    return render_template('index.html',url=builtins.ADDRESS,group=group)
    
@app.route('/help/<path:path>')
def static_proxy(path):
  if path[-1]=='/':
      path += 'index.html'
  return send_from_directory(os.getcwd() + '/libraries/templates/help/', path)

@app.route('/help/')
def static_proxy2():
  return send_from_directory(os.getcwd() + '/libraries/templates/help/','index.html')

@app.route("/dashboard.html")
def dashboard_html():
    try:
        group = flask_login.current_user.id
    except:
        return redirect(url_for('login'))
    logger = logging.getLogger('routing-dashboard_data')
    if request.method == 'GET':
        message = request.args.get('message')
        if message == None:
            message = 'Logged in as ' + group

        data = dashboard_data(group)
        if len(message)>0:
            data['message'] = message

        return render_template('dashboard.html',data=data)

@app.route("/classification.html")
def classification_html():
    try:
        group = flask_login.current_user.id
    except:
        return redirect(url_for('login'))
    logger = logging.getLogger('routing-classification_html')
    if request.method == 'GET':
        message = request.args.get('message')
        if message == None:
            message = 'Logged in as ' + group

        data = {}
        data['address'] = builtins.ADDRESS
        data['group'] = group
        data['locations'] = getAllLocations(group)
        if len(message)>0:
            data['message'] = message

        return render_template('classification.html',data=data)


@app.route("/dashboard.json")
def dashboard_json():
    if request.method == 'GET':
        group = 'find'
        try:
            group = request.args.get('group').lower()
        except:
            pass
        message = request.args.get('message')
        if message == None:
            message = ''
        data = dashboard_data(group)
        if len(message)>0:
            data['message'] = message
        return jsonify(json.loads(json.dumps(data)))

def dashboard_data(group):
    """GET /dashboard.html?group=GROUP

    Returns a HTML page compiling the statistics from the databases
    as well as iframes containing pie charts.
    """
    data = makeStats(group)
    if data == None:
        generateGraph(group)
        calculatePriors(group)
        evaluateAccuracy(group,['test'])           
        data = makeStats(group)

    data['address'] = builtins.ADDRESS
    data['group'] = group
    #data['locations'] = getAllLocations(group)

    try:
        data['gpx'] = submitGPX(group)
    except:
        pass
    
    try:
        data['current_gpx_txt'] = open('data/'+group+'.gpx').read()
    except:
        pass

    return data

@app.route("/mappingdata.html")
def mapping_data():
    """GET /dashboard.html?group=GROUP

    Returns a HTML page compiling the statistics from the databases
    as well as iframes containing pie charts.
    """
    
    logger = logging.getLogger('routing-mapping_data')
    if request.method == 'GET':
        try:
            group = flask_login.current_user.id
        except:
            return redirect(url_for('login'))
        message = request.args.get('message')
        if message == None:
            message = ''

        data = {}
        data['group'] = group
        data['address'] = builtins.ADDRESS
        if len(message)>0:
            data['message'] = message
        try:
            data['gpx'] = submitGPX(group)
        except:
            pass
        
        try:
            data['current_gpx_txt'] = open('data/'+group+'.gpx').read()
        except:
            pass
        return render_template('mappingdata.html',group=group,data=data)

@app.route("/map2.html")
def map2_html():
    """GET /map2.html?group=GROUP

    Returns a HTML page wiht a map of the current location
    """
    try:
        group = flask_login.current_user.id
    except:
        return redirect(url_for('login'))
    data = {'group':group,'address':builtins.ADDRESS}
    return render_template('map.html',data=data)
    
    
@app.route("/map.html")
def map_html():
    """GET /map.html?group=GROUP

    Returns a HTML page with a map of the current location
    """
    logger = logging.getLogger('routing-map_html')
    try:
        group = flask_login.current_user.id
    except:
        return redirect(url_for('login'))
    data = {}
    if request.method == 'GET':
        try:
            group = request.args.get('group').lower()
        except:
            pass
    t1 = time.time()
    data['locations'] = getAllLocations(group)
    data['group'] = group
    data['address'] = builtins.ADDRESS
    print(time.time()-t1)
    return render_template('map2.html',data=data)



@app.route("/charts.html")
def charts_html():
    """GET /charts.html?group=GROUP

    Returns a HTML page showing histogram charts showing
    all the macs for each location and their associated signals
    """
    if request.method == 'GET':
        group = request.args.get('group').lower()
        data = makeChartJson(group)
        return render_template('charts.html',json_string=json.dumps(data))



@app.route("/time_charts.html")
def time_charts_html():
    """GET /charts.html?group=GROUP

    Returns a HTML page showing histogram charts showing
    all the macs for each location and their associated signals
    """
    if request.method == 'GET':
        group = request.args.get('group').lower()
        data = makeTimeChartJson(group)
        return render_template('charts_time.html',json_string=json.dumps(data))


@app.route("/pies.html")
def pies_html():
    """GET /pies.html?group=GROUP

    Returns a HTML page showing D3 donuts with
    the true/false positives from real fingerprints in database
    """
    if request.method == 'GET':
        group = request.args.get('group').lower()
        graph = int(request.args.get('graph'))
        data = {}
        data['results']=makePies(group,graph)
        return render_template('pies.html',json_string=json.dumps(data))

@app.route("/pies_simulated.html")
def pies_simulated_html():
    """GET /pies_simulated.html?group=GROUP

    Returns a HTML page showing D3 donuts with
    the true/false positives from simulated fingerprints
    """
    if request.method == 'GET':
        group = request.args.get('group').lower()
        data = {}
        data['results']=getSimulationResults(group)
        return render_template('pies_simulated.html',json_string=json.dumps(data))


"""
Routes to handle interaction with the databases
"""


'''
#@app.route('/upload_gpx_file', methods=['GET', 'POST'])
def upload_file():
    message = 'Error uploading file.'
    if request.method == 'POST':
        file = request.files['file']
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))
            message = 'File saved successfully.'
        else:
            message = 'Wrong format. Needs to be GPX format.'
    return redirect(url_for('.mapping_data',message=message))

#@app.route('/upload_gpx_text', methods=['GET', 'POST'])
def upload_gpx_text():
    message = 'Error uploading text.'
    if request.method == 'POST':
        text = request.form['gpxtext']
        group = request.form['group']
        with open('data/' + group + '.gpx','w') as f:
            f.write(text)
        message = "Successfully updated GPX coordinates."
    return redirect(url_for('.mapping_data',message=message))
'''
    
    
@app.route('/track', methods=['POST', 'GET'])
@app.route('/learn', methods=['POST', 'GET'])
@app.route('/test', methods=['POST', 'GET'])
def fingerprinting():
    """POST /track, /learn, and /test
    Routes for sending fingerprints

    /track sends the information to the track table and automatically performs a posterior calculations

    /learn sends information to the learn table

    /test sends information to the test table

    The POST data should have the standard fingerprint (above) but without the location name

    ```javascript
    {
        "group":"whatevergroup",
        "username":"iamauser",
        "time": 1409108787,
        "location":"office",
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
    ```
    """
    logger = logging.getLogger('routing-fingerprinting')
    start_time = time.time()
    resp = {'success': True}
    if request.method == 'POST':
        request.get_json(force=True)
        for i in range(len(request.json['wifi-fingerprint'])):
            request.json['wifi-fingerprint'][i]['rssi'] = sanitizeSignal(request.json['wifi-fingerprint'][i]['rssi'])

        newdata = copy.deepcopy(request.json)
        newdata['wifi-fingerprint'] = []
        for fingerprint in request.json['wifi-fingerprint']:
            if "00:00:00:00:00:00" not in fingerprint['mac']:
                newdata['wifi-fingerprint'].append({'mac':fingerprint['mac'],'rssi':fingerprint['rssi']})

        if GENERATE_UNIT_TESTS:
            payload = {'route':request.path,'method':request.method,'json':newdata}
            with open('unit_tests','a') as f:
                f.write(json.dumps(payload) + '\n')
            
        group = request.json['group'].lower()
        if group not in builtins.PARAMETERS:
            logger.debug('Getting new builtins parameters')
            db = mlDB(group)
            builtins.PARAMETERS[group] = db.getResource('calculation_parameters')
            db.close()
            
        if request.path == '/learn':
            try:
                builtins.counters += 1
            except:
                builtins.counters = 0
            if counters % 3 == 0:
                request.path = '/test'

        logger.debug('got ' + request.method +
                    ' from ' + request.json['username'] +
                    ' (' + group + ') posting ' + request.json['location'] +
                    ' fingerprint to ' + request.path)


        if 'track' in request.path:
            resp['position'] = processTrackingFingerprint(newdata)
        else:
            db = mlDB(group)
            db.insertFingerprint(request.path, newdata)
            db.close()


            resp['position'] = processTrackingFingerprint(newdata)

        resp['message'] = 'Saved to database ' + request.path

    else:
        resp['success'] = False
        resp['message'] = 'Something went wrong'
    print(resp)
    return jsonify(resp)

@app.route('/delete', methods=['POST', 'GET'])
def deletion():
    """GET /delete?group=GROUP&location=LOCATION
    Deletes location from ``track`` and ``learn`` databases
    """
    logger = logging.getLogger('routing-deletion')
    resp = {'success': False}
    resp = {'message': 'Error processing request. Use GET /delete?group=GROUP&location=LOCATION'}

    if request.method == 'GET':    
            
        group = request.args.get('group').lower()    
        location = str(request.args.get('location'))
        logger.debug('Requesting deletion of "' + location + '"" in group ' + group)
        db = mlDB(group)
        (resp['success'],resp['message']) = db.deleteLocation(location)
        db.close()

    print(resp)
    return jsonify(resp)

@app.route('/editname', methods=['POST', 'GET'])
def editname():
    """GET /editname?group=GROUP&location=LOCATION&newname=NEWNAME
    Updates location in ``track`` and ``learn`` databases to new name
    """
    logger = logging.getLogger('routing-editname')
    resp = {'success': False}
    resp = {'message': 'Error processing request. Use GET /editname?group=GROUP&location=LOCATION&newname=NEWNAME}'}

    if request.method == 'GET':    
        group = request.args.get('group').lower()    
        location = str(request.args.get('location'))
        newname = str(request.args.get('newname'))
        logger.debug('Requesting update of name for "' + location + '" to "' + newname + '" in group ' + group)
        db = mlDB(group)
        (resp['success'],resp['message']) = db.editLocationName(location,newname)
        db.close()


    print(resp)
    return jsonify(resp)
    
@app.route('/getalllocations', methods=['POST', 'GET'])
def get_all_locations():
    """GET /getalllocations?group=GROUP
    """
    logger = logging.getLogger('routing-editname')
    resp = {'success': False}
    resp = {'message': 'Error processing request. Use GET /getalllocations?group=GROUP'}

    if request.method == 'GET':    
        group = request.args.get('group').lower()    
        resp['success'] = True
        resp['message'] = 'Got locations'
        resp['locations'] = getAllLocations(group)

    return jsonify(resp)

@app.route('/calculateBest', methods=['POST', 'GET'])
def calculating():
    """GET /calculate?group=GROUP
    Calculates new posteriors from data and returns JSON of success

    First it analyzes the fingerprint network using libraries/networkanalysis.py
    to determine the number of components.

    Then it calculates priors for each component in the network.

    Then it evaluates the accuracy of the calculation by running known
    fingerprints against the prior calculations.
    """

    logger = logging.getLogger('routing-calculating')
    logger.debug('Requesting new calculations')
    resp = {'success': True}
    if request.method == 'GET':
    
        if GENERATE_UNIT_TESTS:
            payload = {'route':request.path,'method':request.method,'query_string':request.query_string.decode('utf-8')}
            with open('unit_tests','a') as f:
                f.write(json.dumps(payload) + '\n')
            
        group = request.args.get('group').lower()

        logger.debug('Analyzing fingerprint network...')
        generateGraph(group)

        logger.debug('Calculating priors for production...')
        calculatePriors(group)

        logger.debug('Evaluating accuracy...')
        logger.debug(evaluateAccuracy(group,['test']))

    return jsonify(resp)


@app.route('/calculate', methods=['POST', 'GET'])
def calculating_best():
    """GET /calculate?group=GROUP
    Calculates new posteriors from data and returns JSON of success

    First it analyzes the fingerprint network using libraries/networkanalysis.py
    to determine the number of components.

    Then it calculates priors for each component in the network.

    Then it evaluates the accuracy of the calculation by running known
    fingerprints against the prior calculations.
    """

    logger = logging.getLogger('routing-calculatingbest')
    logger.debug('Requesting new calculations')
    resp = {'success': True}
    if request.method == 'GET':
    
        if GENERATE_UNIT_TESTS:
            payload = {'route':request.path,'method':request.method,'query_string':request.query_string.decode('utf-8')}
            with open('unit_tests','a') as f:
                f.write(json.dumps(payload) + '\n')
            
        group = request.args.get('group').lower()
        
        logger.debug('Analyzing fingerprint network...')
        num_graphs = generateGraph(group)
        calculation_parameters_best = copy.deepcopy(builtins.PARAMETERS[group])
        
        absentee = [1e-6]
        usefulness = [-1,0,1]
        mixinVals = [0.05,0.25,0.5,0.75,0.95]
        
        bestAccuracy = {}
        for graph in range(num_graphs):
            bestAccuracy[graph] = 0
            
        for k in range(len(absentee)):
            for j in range(len(usefulness)):
                for i in range(len(mixinVals)):
                    for graph in range(num_graphs):
                        builtins.PARAMETERS[group][graph]['mix_in'] = mixinVals[i]
                        builtins.PARAMETERS[group][graph]['usefulness'] = usefulness[j]
                        builtins.PARAMETERS[group][graph]['absentee'] = absentee[k]


                    logger.debug('Calculating priors for production...')
                    calculatePriors(group)

                    logger.debug('Evaluating accuracy...')
                    logger.debug(evaluateAccuracy(group,['test']))

                    data = makeStats(group)
                    logger.debug(data['accuracies'])
                    logger.debug(bestAccuracy)
                    for graph in range(num_graphs):
                        if data['accuracies'][graph] > bestAccuracy[graph]:
                            logger.debug("Graph %2.0d %.2f %.2f %1.2f" % (graph,usefulness[j],mixinVals[i],data['accuracies'][graph]))
                            calculation_parameters_best[graph] = copy.deepcopy(builtins.PARAMETERS[group][graph])
                            bestAccuracy[graph] = data['accuracies'][graph]

        db = mlDB(group)
        db.insertResource('calculation_parameters',calculation_parameters_best)
        db.close()
        logger.debug(json.dumps(calculation_parameters_best,indent=4))
        builtins.PARAMETERS[group] = copy.deepcopy(calculation_parameters_best)
        calculatePriors(group)
        logger.debug(evaluateAccuracy(group,['test']))        



    return jsonify(resp)


@app.route("/revaluate")
def reevaluate():
    """GET /revalute?group=GROUP&data=TYPE

    Revaluates accuracy
    """


    results = {"success":True}
    if request.method == 'GET':
    
        if GENERATE_UNIT_TESTS:
            payload = {'route':request.path,'method':request.method,'query_string':request.query_string.decode('utf-8')}
            with open('unit_tests','a') as f:
                f.write(json.dumps(payload) + '\n')
        
        group = request.args.get('group').lower()
        datatype = request.args.get('data').lower()
        if datatype == 'real':
            evaluateAccuracy(group,['test'])           
        elif datatype == 'simulated':
            foo = getSimulationResults(group,force=True)
    return jsonify(results)


@app.route("/whereami2")
@crossdomain(origin='*')
def get_where_am_i_old():
    """GET /whereami

    Returns JSON encoding of GPX location of user in group
    """
    logger = logging.getLogger('routing-whereami')
    if request.method == 'GET':  
        if GENERATE_UNIT_TESTS:
            payload = {'route':request.path,'method':request.method,'query_string':request.query_string.decode('utf-8')}
            with open('unit_tests','a') as f:
                f.write(json.dumps(payload) + '\n')
        group = request.args.get('group').lower()
        user = request.args.get('user').lower()
        try:
            data = getGPX(group,builtins.fingerprint_cache[group][user][0]['location'])
        except:
            logger.debug('Getting /whereami from database')
            db = mlDB(group)
            locinfo = db.getLastLocationFromTracking(user)
            db.close()
            if locinfo == None:
                data = getGPX(group,None)
            else:
                data = getGPX(group,locinfo['location'])
    return jsonify(data)

@app.route("/whereami")
@crossdomain(origin='*')
def get_where_am_i():
    """GET /whereami?group=GROUP&user=USER

    Returns location user in group
    """
    logger = logging.getLogger('routing-whereami')
    if request.method == 'GET':  
        if GENERATE_UNIT_TESTS:
            payload = {'route':request.path,'method':request.method,'query_string':request.query_string.decode('utf-8')}
            with open('unit_tests','a') as f:
                f.write(json.dumps(payload) + '\n')
        group = request.args.get('group')
        if group == None:
            return jsonify({'message':'No group specified','success':False})
        else:
            group = group.lower()
        user = request.args.get('user')
        data = {}
        if user is not None:
            user = user.lower()
            data = getUserLocations(user,group)
        else:
            logger.debug('Getting /whereami from database')
            db = mlDB(group)
            users = db.getUsers()
            db.close()
            data = getUserLocations(users,group)
    return jsonify(data)
    

@app.route("/find")
@crossdomain(origin='*')
def find():
    """GET /find

    Returns JSON encoding location of user in group
    """
    logger = logging.getLogger('routing-find')
    data = {'success':False,'message':'GET request with group and (user or users). users must be seperated by comma.\nExample: /find?group=find&users=zack,stefan OR /find?group=find&user=zack'}
    if request.method == 'GET':  
        if GENERATE_UNIT_TESTS:
            payload = {'route':request.path,'method':request.method,'query_string':request.query_string.decode('utf-8')}
            with open('unit_tests','a') as f:
                f.write(json.dumps(payload) + '\n')
        try:
            group = request.args.get('group').lower()
        except:
            group = 'find'
        try:
            user = request.args.get('user').lower()
        except:
            user = ''
        try:
            users = request.args.get('users').split(',')
            for i in range(len(users)):
                users[i] = users[i].lower()
        except:
            users = []
        if len(user)>0:
            users.append(user)
        data['results'] = []
        if len(users)>0:
            db = mlDB(group)
            data['message'] = 'Searched database for ' + ' and '.join(users)
            for user in users:
                locinfo = db.getLastLocationFromTracking(user)
                if locinfo == None:
                    data['results'].append({'user':user,'location':'Not found','timestamp':0})
                else:
                    data['results'].append({'user':user,'location':locinfo['location'],'timestamp':locinfo['timestamp']})
                    data['success'] = True
            db.close()
    return jsonify(data)

@app.route("/parameters", methods=['POST', 'GET'])
def set_parameters():
    """POST /parameters

    Sets the corresponding parameters. All the parameters are optional except for the group.
    """
    logger = logging.getLogger('routing-parameters')
    results = {"success":True}
    if request.method == 'POST':
        request.get_json(force=True)
        
        if GENERATE_UNIT_TESTS:
            payload = {'route':request.path,'method':request.method,'json':request.json}
            with open('unit_tests','a') as f:
                f.write(json.dumps(payload) + '\n')
                
        group = request.json['group'].lower()
        graph = request.json['graph']

        calculation_parameters = copy.deepcopy(builtins.PARAMETERS[group][graph])
        logger.debug('request: ' + json.dumps(request.json))
        for key in request.json.keys():
            if key in calculation_parameters.keys():
                logger.debug('inserting new ' + key + ' with value ' + str(request.json[key]))
                if ('absentee' in key.lower()
                    or 'dropout' in key.lower()
                    or 'useful' in key.lower()
                    or 'mix' in key.lower()):
                    val = float(request.json[key])
                elif ('trigger' in key.lower()):
                    val = request.json[key]
                else:
                    val = int(request.json[key])
                calculation_parameters[key] = val       
        builtins.PARAMETERS[group][graph] = copy.deepcopy(calculation_parameters)
        logger.debug(builtins.PARAMETERS[group][graph])
        
        db = mlDB(group)
        db.insertResource('calculation_parameters',builtins.PARAMETERS[group])
        db.close()
    return jsonify(results)

"""
Routes to handle data files
"""


@app.route('/charts.json', methods=['POST', 'GET'])
@crossdomain(origin='*')
def get_charts():
    """GET /charts.json?group=GROUP

    Makes and returns a JSON with the histograms obtained from
    the prior calculations in calculatePriors()
    """
    data = {}
    if request.method == 'GET':
        group = request.args.get('group').lower()
        data = makeChartJson(group)
    return jsonify(data)

@app.route('/time_charts.json', methods=['POST', 'GET'])
@crossdomain(origin='*')
def get_time_charts():
    """GET /time_charts.json?group=GROUP

    Makes and returns a JSON with the histograms obtained from
    the prior calculations in calculatePriors()
    """
    data = {}
    if request.method == 'GET':
        group = request.args.get('group').lower()
        data = makeTimeChartJson(group)
    return jsonify(data)


@app.route('/pies.json', methods=['POST', 'GET'])
@crossdomain(origin='*')
def get_pies():
    """GET /pies.json?group=GROUP

    Returns a JSON with the accuracy information
    from the evaluateAccuracy() function in a format
    that can be rendered with D3 pie charts
    """
    data = {}
    if request.method == 'GET':
        group = request.args.get('group').lower()
        graph = int(request.args.get('graph'))
        data =makePies(group,graph)
    return jsonify(results=data)


@app.route('/server.json', methods=['POST', 'GET'])
@crossdomain(origin='*')
def get_server_info():
    """GET /server.json?

    Returns a JSON with the status of the server
    """
    return jsonify(getServerStats())


@app.route('/network.json', methods=['POST', 'GET'])
@crossdomain(origin='*')
def get_network():
    """GET /network.json?group=GROUP

    Returns a JSON with network information for the
    D3 rendering of the network
    """
    data = {}
    if request.method == 'GET':
        group = request.args.get('group').lower()
        data = makeNetworkJson(group)
    return jsonify(data)


@app.route('/stats.json', methods=['POST', 'GET'])
@crossdomain(origin='*')
def get_stats():
    """GET /stats.json?group=GROUP

    Returns a JSON containing relevant stats
    about the status of the server and number of fingerprints
    in the database.
    """
    if request.method == 'GET':
        group = request.args.get('group').lower()
        data = makeStats(group)
    return jsonify(data)

"""
Inner workings
"""

def cleanDBs():
    logger = logging.getLogger('routing-cleanDBs')
    for name in glob('data/*.db'):
        print(name)
        group = name.split('data')[1][1:].split('.db')[0]
        logger.info('Archiving ' + group)
        try:    
            db = mlDB(group)
            db.archiveTrack()
            db.close()
            AVAILABLE_DBS.append(group)
        except:
            logger.error('Could not archive ' + group)





def launch():
    """Initiates start-up routines

    Checks to see whether its on OPENSHIFT or Local
    and chooses respective databases to startup.
    """
    conf = getOptions()
    print('[SERVER][INFO]:',conf)
    builtins.MASTER = conf['master']
    builtins.APIKEY = conf['apikey']
    if len(conf['port'])>0:
        builtins.ADDRESS = conf['address'] + ':' + str(conf['port'])
    else:
        builtins.ADDRESS = conf['address']
    builtins.ADDRESS = conf['ext_address']
    OPENSHIFT = False
    if OPENSHIFT:
        builtins.DATABASE_PATH_PREFIX = os.path.join(os.environ.get('HOME'),'app-root/data/')
    else:
        builtins.DATABASE_PATH_PREFIX = 'data/'
    builtins.START_TIME = time.time()
    builtins.PARAMETERS = {}
    builtins.CURRENT_LOCATIONS = {}
    cleanDBs()
    print('Available groups are: ' + ','.join(AVAILABLE_DBS))

# Schedules job_function to be run once each hour
scheduler = BackgroundScheduler()
scheduler.add_job(cleanDBs, 'interval', hours=1)
scheduler.start()
# Launch
launch()


if __name__ == "__main__":
    """Main app

    Uses the Flask internals
    """
    conf = getOptions()
    logger = logging.getLogger('routing-main')
    app.run(host=conf['address'], port=conf['port'],debug=False)
    # app.run()
