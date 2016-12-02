import json
import os
import pickle
import sklearn
import random
import numpy
import copy
from sklearn.ensemble import RandomForestClassifier
from sklearn.feature_extraction import DictVectorizer
from sklearn.pipeline import make_pipeline
from random import shuffle

__author__ = "Huy Tu"
__email__  = "hqtu@ncsu.edu"
__githubID__ = "HuyTu7"

#random.seed(123)
class RF(object):
	def __init__(self):
		self.size = 0
		self.data = []
		self.trainX = numpy.array([])
		self.trainY = []
		self.testX = numpy.array([])
		self.testY = []
		self.macs = set()
		self.locations = set()

	def get_data(self, fname):
		item = {}
		with open(fname, 'r') as f_in:
			for line in f_in:
				signal_data = json.loads(line)
				item["wifi-fingerprint"] = signal_data["wifi-fingerprint"]
				item["location"] = signal_data["location"]
				self.locations.add(signal_data["location"])
				self.data.append(copy.deepcopy(item))
				for signal in signal_data["wifi-fingerprint"]:
					self.macs.add(signal["mac"])
		self.size = len(self.data)
		self.macs = list(self.macs)
		self.locations = list(self.locations)
		return self.data


	def splitDataset(self, dataset, splitRatio):
		trainSize = int(len(dataset)*splitRatio)
		self.trainX.shape=(0, len(self.macs))
		self.testX.shape=(0, len(self.macs))
		index = 0
		xs = [i for i in range(len(dataset))]
		shuffle(xs)
		while index < len(xs):
			item = numpy.zeros(len(self.macs))
			for signal in dataset[xs[index]]['wifi-fingerprint']:
				item[self.macs.index(signal['mac'])] = signal['rssi']
			if index < trainSize:				
				self.trainX = numpy.concatenate((self.trainX, [item]),axis=0)
				self.trainY.append(self.locations.index(dataset[xs[index]]["location"]))
			else:
				self.testX = numpy.concatenate((self.testX, [item]),axis=0)				
				self.testY.append(self.locations.index(dataset[xs[index]]["location"]))
			index += 1
	
	'''def makeMatrix(self, dataset, index):
		item = [] 
		dataT = dataset[index]
		dataTest = dataT["wifi-fingerprint"]
		value = -1
		for i in range(len(self.macs)):
			for j in range(len(dataTest)):
				if self.macs[i] == dataTest[j]["mac"]:
					value = dataTest[j]["rssi"] 
					break
				else:
					value = 0
			item.append(value)
		return item '''

	def randomFC(self):
		clf = RandomForestClassifier(n_estimators=500, n_jobs = -1)
		clf.fit(self.trainX, self.trainY)
		print(self.locations)
		print(clf.score(self.testX, self.testY))		
		
	
randomF = RF()
data = randomF.get_data("data/hackduke.rf.data")
randomF.splitDataset(data, 0.6)
randomF.randomFC()