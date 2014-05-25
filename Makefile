# Copyright 2014 The Yxorp Authors. All rights reserved.

all: yxorp

clean:
	rm -f yxorp

yxorp: config.go main.go config-sample.tml mond/bootstrap.css mond/bootstrap.js mond/index.html mond/jquery.js
	go build
	zrsc-embed $@ config-sample.tml mond