FROM golang:1.11-alpine

WORKDIR /go/src/github.com/applepine1125/homeDeviceController
COPY . . 

RUN apk add --no-cache --update build-base git python py-pip python2-dev \
  && pip install -U setuptools configparser netaddr pycrypto broadlink \
  # && git clone https://github.com/mjg59/python-broadlink.git \
  # && python ./python-broadlink/setup.py install \
  && export GO111MODULE=on \
  && CGO_ENABLED=0 GOOS=linux go build

CMD ["./homeDeviceController"]
