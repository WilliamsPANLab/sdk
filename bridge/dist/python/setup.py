from setuptools import setup, find_packages
import os

script_dir = os.path.dirname(os.path.abspath(__file__))

install_requires = [
    'six>=1.10.0',
]

def _load_version():
	with open(os.path.join(script_dir, 'VERSION'), 'r') as f:
		return f.read()

setup(
    name = 'flywheel',
    version = _load_version(),
    description = 'Flywheel Python SDK',
    author = 'Nathaniel Kofalt',
    author_email = 'nathanielkofalt@flywheel.io',
    url = 'https://github.com/flywheel-io/sdk',
    license = 'MIT',
    packages = find_packages(),
    package_data = {'': ['flywheelBridge.*', 'VERSION']},
    install_requires = install_requires,
)
