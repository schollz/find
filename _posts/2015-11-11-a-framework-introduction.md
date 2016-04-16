---
layout: post
section-type: post
title: A framework for internal positioning that is open-source
category: tech
tags: [ 'releases' ]
---


We developed a new, open-source framework for internal positioning. 

What is *internal positioning*? Its the ability to monitor your location within a home/office/business to resolution within a single room (or better). Typically that is done using motion sensors, or iBeacons, magnetic fields or Bluetooth dots - but our system allows you to do the same with **only your smartphone**!

What is a *framework* for internal positioning? Our framework, named FIND (the Framework for Internal Navigation and Discovery), is both a program that runs on a server and an app that runs on a smartphone. The smartphone detects the wireless signals and strengths and uses that information to classify your current location. The server program is used to keep track of everything and do the heavy calculations.

We made a majority of the code open-source so you can play around with [the source](https://github.com/schollz/find) (it is written in Python). If you just want to try it out, you can do that too using our own demo servers. Follow [the guide](http://internalpositioning.com/guide/getting-started/) for getting started.

Check it out, let us know what you think! If you have any improvements, fixes, bugs, criticism, anecdotes, we would love to hear them!