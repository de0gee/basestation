import subprocess
import re
import json
import os.path 

from flask import Flask, request, send_from_directory,jsonify
app = Flask(__name__, static_url_path='')

wpa_conf = """country=GB
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
network={
    ssid="_ssid_"
    psk="_password_"
}"""

@app.route('/')
def main():
    print("index page")
    return send_from_directory('build', 'index.html')

@app.route('/static/<path:path>')
def send_static(path):
    return send_from_directory('build/static', path)

@app.route('/signin', methods=['POST'])
def signin():
    content = request.get_json(silent=True)
    print(content)
    email = content['currentState']['email']
    ssid = content['currentState']['ssid']
    password = content['currentState']['password']
    print(email,ssid,password)
    with open('wpa.conf','w') as f:
        f.write(wpa_conf.replace('_ssid_',ssid).replace('_password_',password))
    with open('status.json','w') as f:
        f.write(json.dumps({'status':'disconnected'}))
    # subprocess.call("./disable_ap.sh", shell=True)
    return jsonify({'success':True})

if __name__ == "__main__":
    # get status
    s = {'status':'disconnected'}
    if not os.path.isfile('status.json'):
        with open('status.json','w') as f:
            f.write(json.dumps(s))
    else:
        s = json.load(open('status.json'))

    # check connection
    result = subprocess.run(['iwconfig', 'wlan0'], stdout=subprocess.PIPE)
    matches=re.findall(r'\"(.+?)\"',result.stdout.split(b'\n')[0].decode('utf-8'))
    if len(matches) > 0:
        s['status'] = 'connected'
        print("got connected to " + matches[0])

    if s['status'] == 'disconnected':
        s['status'] = 'hostapd'
        with open('status.json','w') as f:
            f.write(json.dumps(s))    
        subprocess.call("./enable_ap.sh", shell=True)
    elif s['status'] == 'hostapd':
        app.run(host="0.0.0.0",port=80)