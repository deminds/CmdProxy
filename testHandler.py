#!/usr/bin/env python3

import argparse
import requests

API_VER = 'v1.0'
SESSIONID_PARAM = 'sessionid'
COMMAND_PARAM = 'command'

def parse_args():
    parser = argparse.ArgumentParser(description="Simple tool for testing and demonstrating CmdProxy work",
        usage='''
        Console:
            Connect:
            ./testHandler.py --type=console --action=connect

            Command:
            ./testHandler.py --type=console --action=command --command="ls -lah" --id=219602104153538926

            Disconnect:
            ./testHandler.py --type=console --action=disconnect --id=219602104153538926

            Complex use:
            ./testHandler.py --type=console --action=commandWithConnect --command="ls -lah"

        Telnet:
            Conplex use:
            ./testHandler.py --type=telnet --action=commandWithConnect --targetHost=172.16.5.10 --targetPort=23 --targetLogin=userName --targetPassword="PasSWoRd" --loginExpectedString="Username" --passwordExpectedString="Password" --hostnameExpectedString="host-name" --continueCommandExpected=" --More--" --command="show version"
        ''')

    parser.add_argument('--type', type=str, choices=['console', 'telnet'], help='Action type', required=True)
    parser.add_argument('--action', type=str, choices=['connect', 'disconnect', 'command', 'commandWithConnect'], help='Action', required=True)
    
    parser.add_argument('--host', type=str, default="localhost", help='Host of CmdProxy', required=False)
    parser.add_argument('--port', type=int, default=25505, help='Port of CmdProxy', required=False)


    parser.add_argument('--id', type=str, help='Session Id', required=False)

    parser.add_argument('--targetLogin', type=str, help='User Login on you target device', required=False)
    parser.add_argument('--targetPassword', type=str, help='User Login on you target device', required=False)
    parser.add_argument('--targetHost', type=str, default='localhost', help='Telnet host of target host', required=False)
    parser.add_argument('--targetPort', type=int, default=23, help='Telnet port of target host', required=False)

    parser.add_argument('--loginExpectedString', type=str, help='Marker start of string for enter login. Depends on you device. Usually: "Login:"', required=False)
    parser.add_argument('--passwordExpectedString', type=str, help='Marker start of string for enter password. Depends on you device. Usually: "Password:"', required=False)
    parser.add_argument('--hostnameExpectedString', type=str, help='Marker start of string for enter command. Depends on you device. Usually: is hostname', required=False)
    parser.add_argument('--continueCommandExpectedString', type=str, help='Marker start of string for enter continue command. Then listing to long in that case device separate it on several parts and you must press any key for see next part. Depend on your device.', required=False)

    parser.add_argument('--command', type=str, help='Command for execute on you target device', required=False)

    args = parser.parse_args()

    return args

def command(args, sessId):
    print('[INFO]: {}/command'.format(args.type))
    print('[INFO]: exec command: "{}"'.format(args.command))
    url = get_url(args, 'command')

    params = {
        SESSIONID_PARAM: sessId,
        COMMAND_PARAM: args.command
    }

    resp = requests.post(url, json=params)
    if resp.status_code != 200:
        print('[ERROR]: GET error. Status: {}'.format(resp.status_code))
        raise Exception()

    json = resp.json()
    print('[INFO]: {}/command Success. SessId: {}'.format(args.type, json[SESSIONID_PARAM]))
    print('[INFO]: raw output: {}'.format(json))

    return
    
def disconnect(args, sessId):
    print('[INFO]: {}/disconnect'.format(args.type))
    url = get_url(args, 'disconnect')

    params = {
        SESSIONID_PARAM: sessId
    }

    resp = requests.get(url, params)
    if resp.status_code != 200:
        print('[ERROR]: GET error. Status: {}'.format(resp.status_code))
        raise Exception()

    print('[INFO]: {}/disconnect Success. SessId: {}'.format(args.type, sessId))

    return


def console_connect(args):
    print('[INFO]: console/connect')
    url = get_url(args, 'connect')

    params = {}

    resp = requests.get(url, params)
    if resp.status_code != 200:
        print('[ERROR]: GET error. Status: {}'.format(resp.status_code))
        raise Exception()

    json = resp.json()
    print('[INFO]: console/connect Success. SessId: {}'.format(json[SESSIONID_PARAM]))

    return json[SESSIONID_PARAM]


def console(args):
    print("[INFO]: handle console")
    if args.action == 'connect':
        console_connect(args)
    elif args.action == 'disconnect':
        disconnect(args, args.id)
    elif args.action == 'command':
        command(args, args.id)
    elif args.action == 'commandWithConnect':
        try:
            sessId = console_connect(args)
            command(args, sessId)
            disconnect(args, sessId)
        except e:
            return


def telnet_connect(args):
    print('[INFO]: telnet/connect')
    url = get_url(args, 'connect')

    params = {
        "host": args.targetHost,
        "port": args.targetPort,
        "login": args.targetLogin,
        "password": args.targetPassword,
        "loginExpectedString": args.loginExpectedString,
        "passwordExpectedString": args.passwordExpectedString,
        "hostnameExpectedString": args.hostnameExpectedString,
        "continueCommandExpectedString": args.continueCommandExpectedString
    }

    resp = requests.post(url, json=params)
    if resp.status_code != 200:
        print('[ERROR]: GET error. Status: {}'.format(resp.status_code))
        raise Exception()

    json = resp.json()
    print('[INFO]: telnet/connect Success. SessId: {}'.format(json[SESSIONID_PARAM]))

    return json[SESSIONID_PARAM]


def telnet(args):
    print("[INFO]: handle telnet")
    if args.action == 'connect':
        telnet_connect(args)
    elif args.action == 'disconnect':
        disconnect(args, args.id)
    elif args.action == 'command':
        command(args, args.id)
    elif args.action == 'commandWithConnect':
        try:
            sessId = telnet_connect(args)
            command(args, sessId)
            disconnect(args, sessId)
        except e:
            return

def get_url(args, action):
    url = 'http://{}:{}/api/{}/{}/{}'.format(args.host, args.port, API_VER, args.type, action)
    print("[INFO]: build URL: {}".format(url))

    return url


def main(args):
    if args.type == 'console':
        console(args)
    elif args.type == 'telnet':
        telnet(args)
    else:
        print("[ERROR]: unhandled action type: {}".format(args.type))

if __name__ == "__main__":
    args = parse_args()
    main(args)