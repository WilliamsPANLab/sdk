import os

from setuptools import setup, find_packages

setup(
    name = 'flywheel',
    version = '0.0.1',
    description = 'Flywheel Python SDK',
    author = 'Nathaniel Kofalt',
    author_email = 'nathanielkofalt@flywheel.io',
    url = 'https://github.com/flywheel-io/sdk',
    license = 'MIT',
    packages = find_packages(),
    package_data = {'': ['flywheelBridge.*']},
)
