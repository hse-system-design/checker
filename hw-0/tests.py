import json
import os
import random
import string

import numpy
import pytest
import requests


class Bullet:
    def __init__(self, method, path_url, tag, headers=None, body=None, json_=None):
        headers = headers or {}
        body = body or ""
        self.method = method
        self.path_url = path_url
        self.tag = tag
        self.headers = headers
        self.body = body
        if json_:
            self.body = json.dumps(json_)

    def to_ammo_format(self) -> str:
        req = "{method} {path_url} HTTP/1.1\r\n{headers}\r\n{body}".format(
            method=self.method,
            path_url=self.path_url,
            headers=''.join('{0}: {1}\r\n'.format(k, v) for k, v in self.headers.items()),
            body=self.body or "",
        )
        return "{req_size} {tag}\n{req}\r\n".format(req_size=len(req), req=req, tag=self.tag)


def test_ping(cluster_url):
    r = requests.get(str(cluster_url))
    assert r.status_code == 200
    assert len(r.text) > 0


def test_short_link(cluster_url):
    add_url = str(cluster_url / "api/add")
    r = requests.post(add_url, json={
        "url": "http://yandex.ru"
    })
    assert r.status_code == 200

    data = r.json()
    assert 'ShortUrl' in data
    assert data['ShortUrl'] != ""

    r = requests.get(str(cluster_url / data["ShortUrl"]), allow_redirects=False)
    assert r.status_code == 307
    assert r.headers['Location'] == 'http://yandex.ru'


def generate_requests(cluster_url, number):
    for _ in range(number):
        yield Bullet(method="GET", path_url=str(cluster_url), tag="root")

        rnd = ''.join(random.choice(string.ascii_uppercase) for _ in range(10))
        yield Bullet(method="POST", path_url=str(cluster_url / "api/add"),
                     json_={"url": "http://yandex.ru/{}".format(rnd)}, tag="add")

        rnd = ''.join(random.choice(string.ascii_uppercase) for _ in range(8))
        yield Bullet(method="GET", path_url=str(cluster_url / rnd), tag="get")


@pytest.mark.order("last")
def test_shoot_from_a_tank(ammo_writer, tank_load_conf_writer, tank_results_getter, final_results_writer):
    ammo_filename = ammo_writer(generate_requests, 10000)
    load_filename = tank_load_conf_writer()
    exit_code = os.system("yandex-tank -q -c {load} {ammo}".format(load=load_filename, ammo=ammo_filename))
    assert exit_code == 0

    results = tank_results_getter()
    q50, q90 = numpy.quantile(results['root'], [0.5, 0.9])

    final_results_writer({
        'root': {
            'q50': q50,
            'q90': q90,
        }
    })
