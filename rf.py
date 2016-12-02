import json
import sys
import os
import pickle
import sklearn
import random
import numpy
from sklearn.ensemble import RandomForestClassifier
from sklearn.feature_extraction import DictVectorizer
from sklearn.pipeline import make_pipeline
from random import shuffle

random.seed(123)
class RF(object):
	#data = []
	def __init__(self):
		self.size = 0
		self.data = []
		self.nameX = []
		self.trainX = numpy.array([])
		self.testX = numpy.array([])
		self.nameY = []
		self.trainY = []
		self.testY = []
		self.macSet = set()
		self.locationSet = set()

	def get_data(self, fname, splitRatio):
		# First go through once and get set of macs/locations
		X = []
		with open("data/" + fname + ".rf.json",'r') as f_in:
			for fingerprint in f_in:
				data = json.loads(fingerprint)
				X.append(data)
				self.locationSet.add(data['location'])
				for signal in data['wifi-fingerprint']:
					self.macSet.add(signal['mac'])

		# Convert them to lists, for indexing
		self.nameX = list(self.macSet)
		self.nameY = list(self.locationSet)

		# Go through the data again, in a random way
		shuffle(X)
		# Split the dataset for training / learning
		trainSize = int(len(X)*splitRatio)
		# Initialize X, Y matricies for training and testing
		self.trainX=numpy.zeros((trainSize, len(self.nameX)))
		self.testX=numpy.zeros((len(X)-trainSize, len(self.nameX)))
		self.trainY =[0]*trainSize
		self.testY =[0]*(len(X)-trainSize)
		curRowTrain = 0
		curRowTest = 0
		for i in range(len(X)):
			newRow = numpy.zeros(len(self.nameX))
			for signal in X[i]['wifi-fingerprint']:
				newRow[self.nameX.index(signal['mac'])] = signal['rssi']
			if i < trainSize: # do training
				self.trainX[curRowTrain,:] = newRow
				self.trainY[curRowTrain] = self.nameY.index(X[i]['location'])
				curRowTrain = curRowTrain + 1
			else:
				self.testX[curRowTest,:] = newRow
				self.testY[curRowTest] = self.nameY.index(X[i]['location'])
				curRowTest = curRowTest + 1


	def learn(self, dataFile,splitRatio):
		self.get_data(dataFile,splitRatio)
		clf = RandomForestClassifier(n_estimators=10, max_depth=None,
								min_samples_split=2, random_state=0)
		clf.fit(self.trainX, self.trainY)
		score = clf.score(self.testX, self.testY)
		with open('data/'+dataFile+'.rf.pkl','wb') as fid:
			pickle.dump([clf,self.nameX,self.nameY],fid)
		return score

	def classify(self,groupName,fingerpintFile):
		with open('data/' + groupName + '.rf.pkl','rb') as pickle_file:
			[clf,self.nameX,self.nameY] = pickle.load(pickle_file)

		# As before, we need a row that defines the macs
		newRow = numpy.zeros(len(self.nameX))
		data = {}
		with open(fingerpintFile,'r') as f_in:
			for line in f_in:
				data = json.loads(line)
		if len(data) == 0:
			return
		for signal in data['wifi-fingerprint']:
			# Only add the mac if it exists in the learning model
			if signal['mac'] in self.nameX:
				newRow[self.nameX.index(signal['mac'])] = signal['rssi']

		prediction = clf.predict_proba(newRow.reshape(1,-1))
		predictionJson = {}
		for i in range(len(prediction[0])):
			predictionJson[self.nameY[i]] = prediction[0][i]
		return predictionJson


# import socket
#
#
# TCP_IP = '127.0.0.1'
# TCP_PORT = 5006
# BUFFER_SIZE = 1024  # Normally 1024, but we want fast response
#
# s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
# s.bind((TCP_IP, TCP_PORT))
# s.listen(1)
#
# while True:
# 	conn, addr = s.accept()
# 	print ('Connection address:', addr)
# 	while 1:
# 		data = conn.recv(BUFFER_SIZE)
# 		if not data: break
# 		data = data.decode('utf-8')
# 		print ("received data:", data)
# 		group = data.split('=')[0].strip()
# 		filename = data.split('=')[1].strip()
# 		randomF = RF()
# 		payload = "error".encode('utf-8')
# 		try:
# 			payload = json.dumps(randomF.classify(group,filename+".rftemp")).encode('utf-8')
# 		except:
# 			pass
#
# 		conn.send(payload)  # echo
# 	conn.close()

import socketserver

class EchoRequestHandler(socketserver.BaseRequestHandler):

	def handle(self):
		# Echo the back to the client
		data = self.request.recv(1024)
		data = data.decode('utf-8').strip()
		print ("received data:'%s'" % data)
		group = data.split('=')[0].strip()
		filename = data.split('=')[1].strip()
		payload = "error".encode('utf-8')
		if len(group) == 0:
			self.request.send(payload)
			return
		randomF = RF()
		if len(filename) == 0:
			payload = json.dumps(randomF.learn(group,0.7)).encode('utf-8')
		else:
			payload = json.dumps(randomF.classify(group,filename+".rftemp")).encode('utf-8')
		self.request.send(payload)
		return

if __name__ == '__main__':
	import socket
	import threading
	socketserver.TCPServer.allow_reuse_address = True
	address = ('localhost', 5009) # let the kernel give us a port
	server = socketserver.TCPServer(address, EchoRequestHandler)
	ip, port = server.server_address # find out what port we were given
	server.serve_forever()

# from flask import Flask, request
# app = Flask(__name__)
#
# @app.route('/learn')
# def learn():
# 	group = request.args.get('group', '')
# 	if len(group) == 0:
# 		return "error"
# 	randomF = RF()
# 	randomF.learn(group,0.7)
# 	return 'done'
#
# @app.route('/track')
# def track():
# 	import time
# 	start_time = time.time()
# 	group = request.args.get('group', '')
# 	if len(group) == 0:
# 		return "error"
# 	filename = request.args.get('filename', '')
# 	if len(filename) == 0:
# 		return "error"
# 	randomF = RF()
# 	print("--- %s seconds ---" % (time.time() - start_time))
# 	return json.dumps(randomF.classify(group,filename+".rftemp"))
# # python3 rf.py groupName
# try:
# 	# randomF = RF()
# 	# randomF.classify(sys.argv[2],sys.argv[3])
# 	# randomF.learn(fname,0.5) # file, and percentage of data to use to learn
# 	if len(sys.argv)==2:
# 		# Learn print("python3 rf.py groupName")
# 		# Requires writing a file to disk, groupName.rf.json
# 		randomF = RF()
# 		randomF.learn(sys.argv[1],0.7)
# 	elif len(sys.argv)==3:
# 		randomF = RF()
# 		randomF.classify(sys.argv[1],sys.argv[2])
# 	else:
# 		print("error")
# except:
# 	print("error")
