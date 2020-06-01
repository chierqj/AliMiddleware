#
#  Copyright (c) 2020, Alibaba Group. (http://www.alibabagroup.com) All Rights Reserved.
#
#  Alibaba Group. licenses this file to you under the Apache License,
#  Version 2.0 (the "License"); you may not use this file except
#  in compliance with the License.
#  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied. See the License for the
#  specific language governing permissions and limitations
#  under the License.

import json
import shutil
import math as m
import os


class APIS:
    READY = '/ready'
    P1 = '/p1_start'
    P2 = '/p2_start'


SERVER = 'http://127.0.0.1:3355'

root_path = "/Users/chier/mygo/src/AliMiddleware"

IN_DIR = os.getenv("IN_DIR") or '{}/input'.format(root_path)
OUT_DIR = os.getenv("OUT_DIR") or '{}/output'.format(root_path)
ORI_DIR = os.getenv("ORI_DIR") or '{}/data'.format(root_path)
MB_EACH_NODE = 0.01


class Pilot:
    def __init__(self, id):
        self.id = id
        self.apps = []
        self.srvs = {}

        self.cons = 0
        self.mem = 0


class Context:
    def __init__(self):
        self.pilots = {}
        self.apps = {}
        self.appSrvs = {}
        self.srvs = {}

        self.incr_data = []

    def pliot_snapshot(self):
        snapshot = 'Pilot Status:\n'
        for p in self.pilots.values():
            snapshot += '{}-> connections: {:d}, services(apps): {:d}({:d}), memory: {:0.2f}\n'.format(
                p.id, p.cons, len(p.srvs), len(p.apps), p.mem)

        snapshot += '\n[Total connections: {:d}, Apps: {:d}]'.format(
            sum(ctx.apps.values()), len(ctx.apps))

        return snapshot


class Board:
    def __init__(self):
        self.p1_score = Score()
        self.p2_score = Score()

    def total_score(self):
        return 0.5 * self.p1_score.calculate() + 0.5 * self.p2_score.calculate()


class Score:
    def __init__(self):
        self.ms_cost = 0
        self.std_conn = 0
        self.std_mem = 0

        self.mem_used = 0
        self.mem_total = 0

        self.pilot_snapshot = ''

    def calculate(self):
        return (self.mem_used / self.mem_total) * (self.std_mem + self.std_conn)


board = Board()
ctx = Context()


def load_inputs():
    with open(ORI_DIR + '/data.json', 'r') as f:
        data_map = json.loads(f.read())

        update_apps(data_map['apps'])
        update_deps(data_map['dependencies'])

    with open(ORI_DIR + '/pilots.json', 'r') as f:
        data = json.loads(f.read())

        update_pilots(data)


def update_pilots(obj):
    # load pilot cluster, the format should be: {"pilots": ["pilot-indentifier"......]},
    # example: {"pilots": ["p1", "p2", "p3",......]}
    pilots = obj['pilots']
    for p in pilots:
        ctx.pilots[p] = Pilot(p)


def update_apps(obj):
    # update apps, the format should be: {"app-name": node-count, ...},
    # example：{"app-1": 5, "app-2": 6, "app-3": 10}
    ctx.apps.update(obj)


def update_deps(obj):
    # load app deps, the format should be: {"app-name":{"service-identifier", node-count, ...},
    # example：{"app-1" : {"srv1": 10, "srv2": 30, "srv3": 10, "srv10": 15}, "app-2" : {"srv3": 5, "srv4": 8, "srv5": 12}, ......}
    mem = 0
    for app, deps in obj.items():
        if app in ctx.apps:
            ctx.appSrvs[app] = deps
            ctx.srvs.update(deps)


def copy_inputs():
    shutil.copyfile(ORI_DIR + "/data.json", IN_DIR + "/data.json")


def load_dynamic_inputs():
    for filename in os.listdir(ORI_DIR):
        # incr data file is like '001.json'
        if not filename.startswith('0'):
            continue

        with open(ORI_DIR + '/' + filename, 'r') as f:
            data = f.read()

            # load to dynamic
            obj = json.loads(data)
            ctx.incr_data.append(obj)
