build: Dockerfile
	docker build -t home_device_controller ./

run:
	docker stop home_device_controller || true && docker rm home_device_controller || true
	docker run -d -p 80:80 --name home_device_controller home_device_controller

stop:
	docker run home_device_controller
