build:
	docker build -t forum:1.0 .
run:
	docker run -d --name forum-app -p8023:8023 forum:1.0 && echo "server started at http://localhost:8023/"
