from setuptools import setup
from libraries.configuration import *

setup(name='find-ml',
      version='0.2',
      description='Framework for Internal Navigation and Discovery',
      author='Zack',
      author_email='zack@hypercubeplatforms.com',
      url='http://www.python.org/sigs/distutils-sig/',
      install_requires=['APScheduler','Flask','Flask-Login','networkx','numpy','requests','tornado','utm'],
     )
     
import os

os.system('wget https://github.com/schollz/find/releases/download/0.1/calculate.zip')
os.system('unzip calculate.zip')
os.system('rm calculate.zip')