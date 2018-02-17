from flask import Flask, request, send_from_directory,jsonify
app = Flask(__name__, static_url_path='')

@app.route('/')
def main():
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
	return jsonify({'success':True})

if __name__ == "__main__":
    app.run()