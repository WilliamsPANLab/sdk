#!/usr/bin/env bash
set -euo pipefail
unset CDPATH; cd "$( dirname "${BASH_SOURCE[0]}" )"; cd "$(pwd -P)"

cd ../dist/python/wheelhouse

echo "Testing wheel fpr python 2..."
pip install "$(find -name *py2*.whl)"
python ../test-drive.py

echo "Testing wheel fpr python 3..."
pip3 install "$(find -name *py3*.whl)"
python3 ../test-drive.py

