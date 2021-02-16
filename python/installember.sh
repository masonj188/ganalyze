#!/bin/bash

git clone https://github.com/elastic/ember.git
perl -i -pe's/lief>=0.9.0/lief==0.10.1/' ./ember/requirements.txt

echo "numba==0.52.0" >> ./ember/requirements.txt

cd ember
pip install --upgrade pip
pip install -r requirements.txt
python setup.py install
