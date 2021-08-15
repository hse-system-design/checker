import os

import numpy
import pytest
import requests


class Bullet:
    def __init__(self, method, path_url, tag, headers=None, body=None):
        headers = headers or {}
        body = body or ""
        self.method = method
        self.path_url = path_url
        self.tag = tag
        self.headers = headers
        self.body = body

    def to_ammo_format(self) -> str:
        req = "{method} {path_url} HTTP/1.1\r\n{headers}\r\n{body}".format(
            method=self.method,
            path_url=self.path_url,
            headers=''.join('{0}: {1}\r\n'.format(k, v) for k, v in self.headers.items()),
            body=self.body or "",
        )
        return "{req_size} {tag}\n{req}\r\n".format(req_size=len(req), req=req, tag=self.tag)


def test_ping(cluster_url, bullets):
    r = requests.get(str(cluster_url))
    assert r.status_code == 200
    assert len(r.text) > 0

    bullets.append(Bullet(
        method="GET",
        path_url=str(cluster_url),
        tag="root"
    ))
    bullets.append(Bullet(
        method="GET",
        path_url=str(cluster_url),
        tag="root"
    ))


@pytest.mark.order("last")
def test_shoot_from_a_tank(bullets, ammo_writer, tank_load_conf_writer, tank_results_getter, final_results_writer):
    ammo_filename = ammo_writer(bullets)
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
