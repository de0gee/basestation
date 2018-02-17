
build:
	npm run build 

start:
	gunicorn -b 0.0.0.0:80 server:app 

dev:
	export FLASK_DEBUG=1 && export FLASK_APP=server.py && flask run --debugger --port 8004