Test
**The Framework for Internal Navigation and Discovery** (*FIND*) allows you to use your smartphone or laptop to determine your position within your home or office. You can easily use this system in place of motion sensors as its resoltion will allow your phone to distinguish whether you are in the living room or the kitchen or bedroom etc. Simply put, FIND will allow you to replace tons of motion sensors with a single smartphone. The position information can then be used in a variety of ways including home automation, way-finding, tracking, among a few!

<blockquote>Simply put, FIND will allow you to replace <em>tons</em> of motion sensors with a <em>single</em> smartphone!</blockquote>

The system is built on two main components - a server
and a fingerprinting device. The fingerprinting device (computer or android app) sends the specified data to the machine learning server which stores the fingerprints and analyzes them. It then returns the result to the device and stores the result on the server for accessing via a web browser or triggering via hooks.

More detailed documentation can be found in the [FIND Guide](http://internalpositioning.com/guide/).

# Screenshots

## Sign-in


<center>

![Screenshot of the signin](http://internalpositioning.com/guide/img/signin1.png)

</center>

<br>

<center>

![Landing](http://internalpositioning.com/guide/img/landing2.png)
*Screenshot of the landing page*

</center><br>

## Monitor location in realtime


<center>

![Screenshot of the classifications page](http://internalpositioning.com/guide/img/classifications1.png)

</center><br>

## Visualize accuracy and errors

<center>

![Charts show a clear diagnostics of the accuracy for each room](http://internalpositioning.com/guide/img/stats1.png)

</center><br>

<center>

![Pie charts lets you visualize the classification errors](http://internalpositioning.com/guide/img/pies1.png)

</center><br>

## Visualize raw data


<center>

![In-depth analysis of the raw fingerprint data](http://internalpositioning.com/guide/img/signals1.png)

</center><br>

# Server setup

Installation is very simple. Simply download the latest source code and install:

    git clone https://github.com/schollz/find.git
    cd find/
    sudo python3 setup.py install

After which you will be prompted to enter the `address` and `port` of
your server. If you want to run on a home network run `ifconfig` to
check your address (it will be something like `192.168.X.Y` usually).
If you want to use an public address you can also use that. Note: if you are
using a reverse proxy you can also set the `external address`, but if
not, you can just leave that blank.

To run **FIND** simple use:

    python3 server.py


# App

To use the system, you will need a fingerprinting device. The easiest thing to do is to use [our app from Google Play](https://play.google.com/store/apps/details?id=com.hcp.find&hl=en). 

# Notes

## Backup/restore database

### Backup

```
sqlite3 find.db .sch > schema
sqlite3 find.db .dump > dump
grep -v -f schema dump > data
```

### Restore

```
sqlite3 find.db < data
```

### Copy to new repository

```
rsync -avrP --files-from essential_files ./ ~/find
```

## Style guide


To ensure uniform style in coding and documentation, please take a look
at the following notes on style. Please try to follow these the best you
can when submitting pull requests.

PDF of the latest style guide can be found
[here](http://yperevoznikov.com/wp-content/uploads/2014/09/PEP8-python-styles-guide.pdf). Try to follow it the best you can. Use
[autopep8](https://pypi.python.org/pypi/autopep8/) for fixing anything
you missed. General style takeaways:

- **Use 4 spaces per indentation level**, *not tabs!* 
- **Indent continued lines more often, and appropriately!** 
- **Use leading underscore for non-public methods and instance variables**, 
- Use “**if X is not Y**”, do not use “if not X is Y”.

Module headers should be something like:

```python
"""The name of the module
Short summary that makes sense on its own to describe what this module does.

Longer more detailed summary
"""

import built-in-modules
import third-part-modules
import your-own modules

__author__ = "YOUR NAME"
__copyright__ = "Copyright 2015, FIND"
__credits__ = ["YOUR NAME", "HIS/HER NAME"]
__license__ = "MIT"
__version__ = "1.0.1"
__maintainer__ = "YOUR NAME"
__email__ = "your@email"
__status__ = "Development"


CODE-GOES-HERE
```

Function comments should have doc strings as well:

```python
def complex(real=0.0, imag=0.0):
    """Form a complex number.

    Keyword arguments:
    real -- the real part (default 0.0)
    imag -- the imaginary part (default 0.0)
    """
    if imag == 0.0 and real == 0.0:
        return complex_zero
```

If you decided to include more packages, be sure to add them in a virtualenv and add them to `requirements.txt` using:

```bash
pip freeze > requirements.txt
```
