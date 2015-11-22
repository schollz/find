from setuptools import setup
from libraries.configuration import *

setup(name='find-ml',
      version='0.21',
      description='Framework for Internal Navigation and Discovery',
      author='Zack',
      author_email='zack@hypercubeplatforms.com',
      url='http://www.python.org/sigs/distutils-sig/',
      install_requires=['apscheduler','Flask','Flask-Login','networkx','numpy','requests','tornado','utm'],
     )
     
import os
import shutil
import stat
import urllib.request
import zipfile,os.path

if os.path.isdir("calculate"):
	shutil.rmtree('calculate')

if not os.path.isdir("calculate"):

    # Download the file from `url` and save it locally under `file_name`:
    print('Downloading the calculation binaries...')
    with urllib.request.urlopen('https://github.com/schollz/find/releases/download/0.1/calculate.zip') as response, open('calculate.zip', 'wb') as out_file:
        shutil.copyfileobj(response, out_file)


    def unzip(source_filename, dest_dir):
        with zipfile.ZipFile(source_filename) as zf:
            zf.extractall(dest_dir)


    print('Unzipping the calculation binaries...')
    unzip('calculate.zip','./')
    for (dir, _, files) in os.walk("calculate"):
        for f in files:
            path = os.path.join(dir, f)
            st = os.stat(path)
            os.chmod(path, st.st_mode | 0o111)
    print('Removing the zipfile')
    os.remove('calculate.zip')
