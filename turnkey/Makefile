
build:
	npm run build 


start:
	gunicorn -b 0.0.0.0:80 server:app 

deps:
	npm install
	sudo python3 -m pip install gunicorn flask

dev:
	export FLASK_DEBUG=1 && export FLASK_APP=server.py && flask run --debugger --port 8004