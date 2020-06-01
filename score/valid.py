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

from misc import *
import math as m


def std(data):
    s = sum(data)
    l = len(data)

    dif = sum(map(lambda x: m.pow(x - s/l, 2), data)) / l

    return m.sqrt(dif)


def validate_pilotset(data):
    apps_left = ctx.apps.copy()

    try:
        ps = json.loads(data)
        for p in ps:
            if p not in ctx.pilots:
                return False, 'pilot {} not found'.format(p)
            if ps[p] is None or len(ps[p]) <= 0:
                continue

            pilot = ctx.pilots[p]

            pilot.apps = list(filter(lambda app: app in ps[p], ctx.apps))
            for app in pilot.apps:
                if app not in apps_left:
                    return False, 'app {} resigned to pilot'.format(app)

                del apps_left[app]

        if len(apps_left) > 0:
            return False, 'some apps are not assigned to any pilot'

        # now calculate the score

        cons_l = []
        mem_l = []

        for pilot in ctx.pilots.values():
            pilot.cons = 0
            # do not reset pilot.srvs here

            for app in pilot.apps:
                pilot.cons += ctx.apps[app]
                pilot.srvs.update(ctx.appSrvs[app])

            # assume every node consumes 1 KiB memory
            pilot.mem = sum(pilot.srvs.values()) * MB_EACH_NODE

            cons_l.append(pilot.cons)
            mem_l.append(pilot.mem)

        return True, (std(cons_l), std(mem_l), sum(mem_l))
    except Exception as e:
        return False, 'returned data does not meet requirements'
