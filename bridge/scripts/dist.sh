#!/usr/bin/env bash
set -euo pipefail
unset CDPATH; cd "$( dirname "${BASH_SOURCE[0]}" )"; cd "$(pwd -P)"

cd ../dist/python

echo "Cleaning..."
rm -rf build/ flywheel.egg-info/

echo "Determining python platform string..."
# This value seems to be identical between python 2 & 3
pyPlatform=$(python -c "import distutils.util; print(distutils.util.get_platform())")
echo " Platform is $pyPlatform"

echo "Building wheel for python 2..."
python2 setup.py bdist_wheel --plat-name "$pyPlatform" --dist-dir wheelhouse

echo "Building wheel for python 3..."
python3 setup.py bdist_wheel --plat-name "$pyPlatform" --dist-dir wheelhouse

echo "Cleaning..."
rm -rf build/ flywheel.egg-info/
