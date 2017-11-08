from setuptools import setup, find_packages

install_requires = [
    'six>=1.10.0',
]

def _load_version():
	import subprocess

	p = subprocess.Popen('../binary/sdk version', stdout=subprocess.PIPE, shell=True)
	(output, error) = p.communicate()
	status = p.wait()

	if error is not None:
		raise Exception('Could not call sdk binary: Error was' + str(err))
	elif status != 0:
		raise Exception('Could not call sdk binary: Status was' + str(status))

	return output.split()[2].decode('utf-8')


setup(
    name = 'flywheel',
    version = _load_version(),
    description = 'Flywheel Python SDK',
    author = 'Nathaniel Kofalt',
    author_email = 'nathanielkofalt@flywheel.io',
    url = 'https://github.com/flywheel-io/sdk',
    license = 'MIT',
    packages = find_packages(),
    package_data = {'': ['flywheelBridge.*']},
    install_requires = install_requires,
)
