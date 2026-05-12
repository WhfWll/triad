import json
import random
import time
from urllib.parse import urljoin

import requests
from rocketmq.client import PushConsumer
import logging

logging.basicConfig(level=logging.DEBUG)

XIAOZHI_TOKEN = "a699c70d4e080807d48dc1d4b5566f2b"

host = "http://127.0.0.1:8011"
tyjk_host = "http://221.122.84.13:8888"
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36',
    'platform-token': XIAOZHI_TOKEN,
}

random_list = ['z', 'y', 'x', 'w', 'v', 'u', 't', 's', 'r', 'q', 'p', 'o', 'n', 'm', 'l', 'k', 'j', 'i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a']

def add(json_msg):
    logging.info('add a data')
    try:
        roleIdentities = json_msg.get('roleIdentities')[0]
    except Exception as e:
        logging.error("角色不存在")
        return
    password = "0d13cc5a7d97243fcadf8baea30bb053"  # tempuserpass123
    repassword = password
    print(password)
    email = ''.join(random.sample(random_list, 6)) + "@qq.com"
    data = {
        "username": json_msg.get("userName"),
        "email": json_msg.get("email") or email,
        "role": 1,
        "department": '-'.join(json_msg.get("organizations")),
        "remark": json_msg.get('remark'),
        "password": password,
        "repassword": repassword,
        "accountExpireTime": "2035-07-20",
    }
    url = urljoin(host, "/smart/user/save")
    resp = requests.post(url=url, data=data, headers=headers)
    logging.info(resp.text)

    platformUserId = get_user_id(json_msg=json_msg)

    logging.info("向统一接口管理平台同步用户id")
    data2 = {'request': 'UserSyncRequest', 'id': json_msg.get('id'), 'platformUserId': platformUserId,'platformCode': 'WLCSMK'}
    url = urljoin(tyjk_host, "/api/user/sync")
    resp = requests.post(url=url, json=data2)
    logging.info(resp.text)


def delete(json_msg):
    logging.info("delete a data")
    data = {
        "username": json_msg.get("userName"),
    }

    platformUserId = json_msg.get('platformUserId')
    if not platformUserId:
        logging.info("缺少platformUserId字段")
        platformUserId = get_user_id(json_msg=json_msg)

    url = urljoin(host, "/smart/user/del?userIds=" + str(platformUserId))
    resp = requests.get(url=url, data=data, headers=headers)
    logging.info(resp.text)


def update(json_msg):
    logging.info("update a data")
    try:
        roleIdentities = json_msg.get('roleIdentities')[0]
    except Exception as e:
        logging.error(f"角色不存在: {e}")
        return

    platformUserId = json_msg.get('platformUserId')
    if not platformUserId:
        logging.info("缺少platformUserId字段")
        platformUserId = get_user_id(json_msg=json_msg)
    email = ''.join(random.sample(random_list, 6)) + "@qq.com"
    print(platformUserId)
    data = {
        "id": platformUserId,
        "username": json_msg.get("userName"),
        "role": roleIdentities,
        "email": json_msg.get("email") or email,
        "department": '-'.join(json_msg.get("organizations")),
        "remark": json_msg.get("remark"),
    }
    url = urljoin(host, "/smart/user/save")
    resp = requests.post(url=url, data=data, headers=headers)
    logging.info(resp.text)


operator_kwargs = {"ADD": add, "DELETE": delete, "UPDATE": update}


def get_user_id(json_msg):
    platformUserId = None
    resp = requests.get(urljoin(host, "/smart/user/list?page=1&size=1&search=" + json_msg.get("userName","")), headers=headers)
    if resp.status_code != 200:
        logging.info("请求用户列表状态错误")
        return
    results = resp.json().get('data').get("list")
    for one in results:
        if json_msg.get("userName").strip() == one.get('username'):
            platformUserId = one.get('id')
            break
    return platformUserId


def callback(msg):
    logging.info(b"receive a msg from rocketmq: " + msg.body)
    json_msg = json.loads(msg.body)
    operationType = json_msg.get("operationType")
    handle_func = operator_kwargs.get(operationType)
    if not handle_func:
        logging.error("操作类型错误")
        return
    if json_msg.get('platformCode') != "WLCSMK":
        logging.info('非网络测试系统数据, 不进行更新')
        return
    handle_func(json_msg)


consumer = PushConsumer('WLCSMK')
consumer.set_namesrv_addr('118.89.248.24:9876')
consumer.subscribe("userSync", callback)
consumer.start()
while True:
    time.sleep(30)
