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


import urllib3
import time
import datetime

from misc import *
from valid import *

http = urllib3.PoolManager(timeout=urllib3.Timeout(connect=1))


def print_failure(msg):
    print('\n😱 Grading failed, caused by: ' + msg)

    write2disk(msg, '-1', 0)


def print_success():
    print('\n\n=====================RESULT=====================')

    print('\n⭐ Phase 1 score: {:0.4f}:'.format(board.p1_score.calculate()))
    print(board.p1_score.pilot_snapshot)

    memrate = board.p1_score.mem_used/board.p1_score.mem_total
    print('\n@ M加载内存/M服务内存: {:0.4f}'.format(memrate))
    print('@ 连接标准差: {:0.4f}'.format(board.p1_score.std_conn))
    print('@ 内存标准差: {:0.4f}'.format(board.p1_score.std_mem))
    print('@ 耗时: {:0.4f}ms'.format(board.p1_score.ms_cost))

    print('\n⭐ Phase 2 score: {:0.4f}:'.format(board.p2_score.calculate()))
    print(board.p2_score.pilot_snapshot)
    memrate = board.p2_score.mem_used/board.p2_score.mem_total
    print('\n@ M加载内存/M服务内存: {:0.4f}'.format(memrate))
    print('@ 连接标准差: {:0.4f}'.format(board.p2_score.std_conn))
    print('@ 内存标准差: {:0.4f}'.format(board.p2_score.std_mem))
    print('@ 耗时: {:0.4f}ms'.format(board.p2_score.ms_cost))

    print('\n\n🎉 总得分: {:0.4f}\n'.format(board.total_score()))
    write2disk(msg, '0', board.total_score())


def write2disk(msg, status, score):
    # now convert it to the json object
    obj = {
        'isvalid': 1,
        'message': msg,
        'rank': score,
        'status': status,
        'taskid': '${TASK_ID}',
        'raceId': '${RACE_ID}',
        'timestamp': '0'
    }

    if status == '0':
        obj['scoreJson'] = {'score': score}
    else:
        obj['scoreJson'] = {}

    with open(OUT_DIR + '/result.json', 'w+') as f:
        f.write(json.dumps(obj, separators=(',', ':')))


def check_ready():
    ready = False
    prefix = '⏱  Checking ready(10s timeout)'
    for t0 in range(0, 10):
        try:
            resp = http.request('HEAD', SERVER + APIS.READY, timeout=1)
            if resp.status == 200:
                ready = True
                break
        except Exception as e:
            pass
        finally:
            time.sleep(1)

        print(prefix + '......[{:d}s]\r'.format(t0), end='')

    if not ready:
        print(prefix + '......[TIMEOUT]')
        return False, 'ready check timeout'

    print(prefix + '......[READY]')
    return True, 'ok'


def phase_1():
    # copy inputs to user out dir
    copy_inputs()

    t0 = datetime.datetime.now()
    print('🔥 Phase 1 started: input data generated')

    # generate Pilot json
    pilots = {'pilots': []}
    for p in ctx.pilots.values():
        pilots['pilots'].append(p.id)

    # notify the program
    try:
        resq = http.request('POST', SERVER + APIS.P1,
                            body=json.dumps(pilots),
                            headers={'accept': 'application/json',
                                     'Content-Type': 'application/json;utf-8'},
                            timeout=120)
        if resq.status != 200:
            return False, 'request failed: {:d}'.format(resq.status)

        result = resq.data.decode('utf-8')
    except Exception as e:
        t1 = datetime.datetime.now()
        ms_cost = (t1 - t0).total_seconds() * 1000

        return False, 'phase 1 failed, cost: {:0.2f}ms'.format(ms_cost)

    t1 = datetime.datetime.now()
    ms_cost = (t1 - t0).total_seconds() * 1000
    print('🏁 Phase 1 completed: cost {:0.2f}ms'.format(ms_cost))

    success, result = validate_pilotset(result)
    if not success:
        return False, 'failed to validate data, reason: {}'.format(result)

    board.p1_score.ms_cost = ms_cost
    board.p1_score.std_conn = result[0]
    board.p1_score.std_mem = result[1]
    board.p1_score.mem_used = result[2]
    board.p1_score.mem_total = sum(ctx.srvs.values()) * MB_EACH_NODE

    board.p1_score.pilot_snapshot = ctx.pliot_snapshot()

    return True, 'ok'


def phase_2():
    ix = 1
    for data in ctx.incr_data:
        apps = data['apps']
        deps = data['dependencies']

        update_apps(apps)
        update_deps(deps)

        print('🔥 Phase 2 Round {:d} started: new apps {:d}'.format(
            ix, len(apps)))
        t0 = datetime.datetime.now()

        # send incremental apps and dependencies
        try:
            body_data = {'apps': apps, 'dependencies': deps}
            resq = http.request('POST', SERVER + APIS.P2,
                                body=json.dumps(body_data),
                                headers={'accept': 'application/json',
                                         'Content-Type': 'application/json;utf-8'},
                                timeout=15)
            if resq.status != 200:
                return False, 'request failed: {:d}'.format(resq.status)

            result = resq.data.decode('utf-8')
        except Exception as e:
            t1 = datetime.datetime.now()
            ms_cost = (t1 - t0).total_seconds() * 1000

            return False, 'phase 2 failed, cost: {:0.2f}ms'.format(ms_cost)

        t1 = datetime.datetime.now()
        ms_cost = (t1 - t0).total_seconds() * 1000
        print('🏁 Phase 2 Round {:d} completed: cost {:0.2f}ms'.format(
            ix, ms_cost))

        success, result = validate_pilotset(result)
        if not success:
            return False, 'failed to validate data, reason: {}'.format(result)

        # except for the time cost, use the latest value
        board.p2_score.ms_cost += ms_cost
        board.p2_score.std_conn = result[0]
        board.p2_score.std_mem = result[1]
        board.p2_score.mem_used = result[2]
        board.p2_score.mem_total = sum(ctx.srvs.values()) * MB_EACH_NODE

        board.p2_score.pilot_snapshot = ctx.pliot_snapshot()

        ix += 1

    return True, 'ok'


if __name__ == "__main__":
    load_inputs()

    print('🚀 Grading program started, data loaded')

    try:
        success, msg = check_ready()
        if not success:
            print_failure(msg)
            exit(1)

        success, msg = phase_1()
        if not success:
            print_failure(msg)
            exit(1)

        load_dynamic_inputs()

        success, msg = phase_2()
        if not success:
            print_failure(msg)
            exit(1)
    except Exception as e:
        print_failure('grading failed, try again later or contact xuanyin')
        exit(1)

    # done
    print_success()
