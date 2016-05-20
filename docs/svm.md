# SVM

Follow these instructions if you are running a FIND server and would like to add SVM to the machine learning routines.

FIND will automatically utilize `libsvm` once it is installed. Here are the instructions to install (you should run with root/sudo):

```
sudo apt-get install g++
wget http://www.csie.ntu.edu.tw/~cjlin/cgi-bin/libsvm.cgi?+http://www.csie.ntu.edu.tw/~cjlin/libsvm+tar.gz
tar -xvf libsvm-*.tar.gz
cd libsvm-*
make
cp svm-scale /usr/local/bin/
cp svm-predict /usr/local/bin/
cp svm-train /usr/local/bin/
```

Then just restart FIND! It will automatically detect whether its installed. When SVM is enabled, you will see SVM data along with the Naive-Bayes information.

_Note_: Currently FIND defaults to use the Naive-Bayes machine learning for the actual guesses. In my experience
SVM is generally inferior, but this may depend on your location.
