import json

import requests
import maya
import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
import numpy as np
from scipy.stats import gaussian_kde

r = requests.get('https://ml.internalpositioning.com/location?group=abby3&user=34fcefa92597&n=1000')
with open("out.json","w") as f:
    f.write(r.text)

places = {}
data = json.load(open('out.json','r'))
for dat in data['users']['34fcefa92597']:
    mDate = maya.parse(dat['time'].split('.')[0])
    hour = mDate.hour + mDate.minute/60 + mDate.second/(60*60)
    if len(dat['location'].strip()) == 0:
        continue
    if dat['location'] not in places:
        places[dat['location']] = {'rawdata':[],'kde':[]}
    places[dat['location']]['rawdata'].append(hour)


# http://stackoverflow.com/questions/4150171/how-to-create-a-density-plot-in-matplotlib
xs = np.linspace(0,24,240)
for place in places:
    data = places[place]['rawdata']
    density = gaussian_kde(data)
    density.covariance_factor = lambda : .1
    density._compute_covariance()
    places[place]['kde'] = density(xs)


# http://stackoverflow.com/questions/2225995/how-can-i-create-stacked-line-graph-with-matplotlib
y = []
placeKeys = list(places.keys())
for place in placeKeys:
    data = np.array(places[place]['kde'])
    if len(y) == 0:
        y = data
    else:
        y = np.vstack((y,data))
# this call to 'cumsum' (cumulative sum), passing in your y data, 
# is necessary to avoid having to manually order the datasets
x = xs
y_stack = np.cumsum(y, axis=0)   
print(y_stack)

fig = plt.figure()
ax1 = fig.add_subplot(111)

colors=["#d11141","#00b159","#00aedb","#f37735","#ffc425","#e1f7d5","#ffbdbd","#c9c9ff","#f1cbff"]
patches = []
for i in range(len(y)):
    if i==0:
        ax1.fill_between(x, 0, y_stack[0,:], facecolor=colors[i], alpha=.7)
    else:
        ax1.fill_between(x, y_stack[i-1,:], y_stack[i,:], facecolor=colors[i], alpha=.7)
    patches.append(mpatches.Patch(color=colors[i], label=placeKeys[i]))
plt.legend(handles=patches)
plt.show()
