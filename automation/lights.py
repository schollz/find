#!/usr/bin/python3

# Copyright 2015-2017 Zack Scholl. All rights reserved.
# Use of this source code is governed by a AGPL
# license that can be found in the LICENSE file.

import sys
import os
import json
import subprocess
import argparse
import urllib.parse as urlparse
from urllib.parse import urlencode
import logging
import threading
import time

import requests
import schedule
import paho.mqtt.client as mqtt

ON = True
OFF = False

# create logger with 'spam_application'
logger = logging.getLogger('lights.py')
logger.setLevel(logging.DEBUG)


class MQTTThread (threading.Thread):

    def __init__(self,debug):
        threading.Thread.__init__(self)
        self.debug = debug

        # Setup logging
        self.logger = logging.getLogger("mqtt")
        self.logger.setLevel(logging.DEBUG)
        fh = logging.FileHandler('lights.log')
        ch = logging.StreamHandler()
        if debug:
            fh.setLevel(logging.DEBUG)
            ch.setLevel(logging.DEBUG)
        else:
            fh.setLevel(logging.INFO)
            ch.setLevel(logging.INFO)
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(funcName)s:%(lineno)d - %(levelname)s - %(message)s')
        fh.setFormatter(formatter)
        formatterSimple = logging.Formatter(
            '%(asctime)s - %(name)s - %(message)s')
        ch.setFormatter(formatterSimple)
        self.logger.addHandler(fh)
        self.logger.addHandler(ch)

    def run(self):
        self.threads = []
        self.threads.append(DeviceThread("bob",self.debug))
        self.threads.append(DeviceThread("jill",self.debug))

        # Start new Threads
        for thread in self.threads:
            thread.start()
        self.client = mqtt.Client()
        self.client.on_connect = self.on_connect
        self.client.on_message = self.on_message

        r = requests.put(
            "https://ml.internalpositioning.com/mqtt?group=%s" % params['group'])
        self.client.username_pw_set(
            params['group'], password=json.loads(r.text)['password'])
        self.client.connect("ml.internalpositioning.com", 1883, 60)
        self.client.loop_forever()
        for thread in self.threads:
            self.thread.join()

    # The callback for when the client receives a CONNACK response from the
    # server.
    def on_connect(self, client, userdata, flags, rc):
        self.logger.debug("Connected with result code "+str(rc))

        # Subscribing in on_connect() means that if we lose the connection and
        # reconnect then subscriptions will be renewed.
        client.subscribe("rr/test/#")

    # The callback for when a PUBLISH message is received from the server.
    def on_message(self, client, userdata, msg):
        self.logger.debug(msg.topic+" "+str(msg.payload))
        self.threads[0].turn_on_light()


class DeviceThread (threading.Thread):

    def __init__(self, name, debug):
        threading.Thread.__init__(self)
        self.name = name
        self.light = OFF

        # Setup logging
        self.logger = logging.getLogger(self.name)
        self.logger.setLevel(logging.DEBUG)
        fh = logging.FileHandler('lights.log')
        ch = logging.StreamHandler()
        if debug:
            fh.setLevel(logging.DEBUG)
            ch.setLevel(logging.DEBUG)
        else:
            fh.setLevel(logging.INFO)
            ch.setLevel(logging.INFO)
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(funcName)s:%(lineno)d - %(levelname)s - %(message)s')
        fh.setFormatter(formatter)
        formatterSimple = logging.Formatter(
            '%(asctime)s - %(name)s - %(message)s')
        ch.setFormatter(formatterSimple)
        self.logger.addHandler(fh)
        self.logger.addHandler(ch)

    def run(self):
        self.logger.debug("Started")
        while True:
            schedule.run_pending()
            time.sleep(1)

    def turn_off_light(self):
        self.logger.debug("Turning off light")
        schedule.clear(self.name)
        if self.light == ON:
            self.logger.debug("Deactivating light")
            self.light = OFF

    def turn_on_light(self):
        self.logger.debug("Turning on light")

        # Re-schedule to turn it off
        schedule.clear(self.name)
        schedule.every(params['turn_light_off_after_seconds']).seconds.do(
            self.turn_off_light).tag(self.name)

        if self.light == OFF:
            # Turn light on
            self.logger.debug("Activating light")
            self.light = ON


def run_command(c):
    logger.debug("Running command '%s'" % c)
    p = subprocess.Popen(
        c,
        universal_newlines=True,
        shell=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT)
    text = p.stdout.read()
    retcode = p.wait()
    return text, retcode

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-d",
        "--debug",
        action="store_true")
    parser.add_argument(
        "-g",
        "--group",
        type=str,
        default="",
        help="group to use")
    parser.add_argument(
        "-s",
        "--seconds",
        type=int,
        default=2,
        help="seconds to wait before switching off")
    args = parser.parse_args()
    if args.group == "":
        print("Must select a group with -g")
    else:
        params = {'group': args.group,
                  'turn_light_off_after_seconds': args.seconds}
        mqtt_listener = MQTTThread(args.debug)
        mqtt_listener.start()
        mqtt_listener.join()
