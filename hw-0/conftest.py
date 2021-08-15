import csv
import glob
import json
import os
import pathlib
from collections import defaultdict

import furl
import pytest


_tank_yaml_conf_template = """
phantom:
  address: {ip}:{port}
  load_profile:
    load_type: rps
    schedule: line(1, 10, 10m)
console:
  enabled: false
telegraf:
  enabled: false
"""


def pytest_addoption(parser):
    parser.addoption("--cluster_ip", action="store")
    parser.addoption("--workdir", action="store")


@pytest.fixture(scope='session')
def cluster_ip(request):
    cluster_ip_value = request.config.option.cluster_ip
    if cluster_ip_value is None:
        pytest.fail("No cluster ip were provided")
    return cluster_ip_value


@pytest.fixture(scope="session")
def cluster_port():
    return 30030


@pytest.fixture
def cluster_url(cluster_ip, cluster_port):
    return furl.furl(scheme="http", host=cluster_ip, port=cluster_port)


@pytest.fixture(scope="session")
def bullets():
    bullet_list = []
    return bullet_list


@pytest.fixture(scope="session")
def workdir(request):
    workdirname = request.config.option.workdir
    if not os.path.isdir(workdirname):
        os.mkdir(workdirname)

    os.chdir(workdirname)
    return workdirname


@pytest.fixture
def ammo_writer(workdir):
    def _write(bullets):
        ammofilename = str(pathlib.Path(workdir) / 'ammo.txt')
        with open(ammofilename, 'w') as f:
            for b in bullets:
                f.write(b.to_ammo_format())
        return pathlib.Path(ammofilename).name
    return _write


@pytest.fixture
def tank_load_conf_writer(workdir, cluster_ip, cluster_port):
    def _write():
        configfilename = str(pathlib.Path(workdir) / 'load.yaml')
        with open(configfilename, "w") as f:
            f.write(_tank_yaml_conf_template.format(ip=cluster_ip, port=cluster_port))
        return pathlib.Path(configfilename).name
    return _write


@pytest.fixture
def tank_results_getter(workdir):
    def _get():
        phout_filename = glob.glob(str(pathlib.Path(workdir) / "logs/*/phout*.log"))[0]
        results = defaultdict(list)
        with open(phout_filename, 'r') as f:
            r = csv.reader(f, delimiter='\t')
            for row in r:
                tag = row[1]
                latency = float(row[5])
                results[tag].append(latency)
        return results
    return _get


@pytest.fixture
def final_results_writer(workdir):
    def _write(results):
        result_filename = str(pathlib.Path(workdir) / 'tank-results.json')
        with open(result_filename, 'w') as f:
            f.write(json.dumps(results))
        return result_filename
    return _write
