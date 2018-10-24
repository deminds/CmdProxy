#!/usr/bin/env python3

import argparse
import requests

API_VER = 'v1.0'
SESSIONID_PARAM = 'sessionid'
COMMAND_PARAM = 'command'

def parse_args():
    parser = argparse.ArgumentParser(description="Simple tool for testing and demonstrating CmdProxy work",
        usage='''
        Connect:
        ./testHandler.py --type=console --action=connect 

        Command:
        ./testHandler.py --type=console --action=command --command="ls -lah" --id=219602104153538926

        Disconnect:
        ./testHandler.py --type=console --action=disconnect --id=219602104153538926

        Complex use:
        ./testHandler.py --type=console --action=commandWithConnect --command="ls -lah"
        ''')

    parser.add_argument('--type', type=str, choices=['console', 'telnet'], help='Action type', required=True)
    parser.add_argument('--action', type=str, choices=['connect', 'disconnect', 'command', 'commandWithConnect'], help='Action', required=True)
    
    parser.add_argument('--host', type=str, default="localhost", help='Host of CmdProxy', required=False)
    parser.add_argument('--port', type=int, default=25505, help='Port of CmdProxy', required=False)


    parser.add_argument('--id', type=str, help='Session Id', required=False)

    parser.add_argument('--login', type=str, help='Login', required=False)
    parser.add_argument('--password', type=str, help='Password', required=False)
    parser.add_argument('--telnetHost', type=str, default='localhost', help='Telnet host', required=False)
    parser.add_argument('--telnetPort', type=int, default=23, help='Telnet port', required=False)

    parser.add_argument('--command', type=str, help='Command', required=False)

    args = parser.parse_args()

    return args

def console_command(args, sessId):
    print('[INFO]: console/command')
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
    print('[INFO]: console/disconnect Success. SessId: {}'.format(json[SESSIONID_PARAM]))
    print('[INFO]: raw output: {}'.format(json))

    return
    
def console_disconnect(args, sessId):
    print('[INFO]: console/disconnect')
    url = get_url(args, 'disconnect')

    params = {
        SESSIONID_PARAM: sessId
    }

    resp = requests.get(url, params)
    if resp.status_code != 200:
        print('[ERROR]: GET error. Status: {}'.format(resp.status_code))
        raise Exception()

    print('[INFO]: console/disconnect Success. SessId: {}'.format(sessId))

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


def get_url(args, action):
    url = 'http://{}:{}/api/{}/{}/{}'.format(args.host, args.port, API_VER, args.type, action)
    print("[INFO]: build URL: {}".format(url))

    return url


def console(args):
    print("[INFO]: handle console")
    if args.action == 'connect':
        console_connect(args)
    elif args.action == 'disconnect':
        console_disconnect(args, args.id)
    elif args.action == 'command':
        console_command(args, args.id)
    elif args.action == 'commandWithConnect':
        try:
            sessId = console_connect(args)
            console_command(args, sessId)
            console_disconnect(args, sessId)
        except e:
            return

def telnet(args):
    print("[INFO]: handle telnet")
    print("[ERROR]: not implemented yet")

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