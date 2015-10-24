from tornado.wsgi import WSGIContainer
from tornado.httpserver import HTTPServer
from tornado.ioloop import IOLoop
from libraries.configuration import *

conf = getOptions()

from libraries.routing import *

print('running on Tornado with port ' + str(conf['port']))
if len(str(conf['port']))<1:
    print('Need to specify port!')
else:
    http_server = HTTPServer(WSGIContainer(app))
    print(conf['address'])
    http_server.listen(conf['port'])
    IOLoop.instance().start()

# Run with
# uwsgi --http 152.3.53.178 -w server --processes 2
