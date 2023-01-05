SHELL:=/bin/bash

build-bins:
	go build -o ./aux/bins/ ./cmd/*


ttt: build-bins t12

ttt-all: build-bins t1 t2 t3 t4 t5 t6 t7 t8 t9 t10 t11 t12

tt1: build-bins t1
tt2: tt1 t2
tt3: tt2 t3
tt4: tt3 t4
tt5: tt4 t5
tt6: tt5 t6
tt7: tt6 t7
tt8: tt7 t8
tt9: tt8 t9
tt10: tt9 t10
tt11: tt10 t11
tt12: tt11 t12
tt13: tt12 t13




T = ${PWD}/aux/tools/shortenertest-darwin-amd64 -test.v -source-path=. -binary-path=./aux/bins/shortener
R = ${PWD}/aux/tools/random-darwin-amd64
t1:
	$T -test.run=^TestIteration1$$
t2:
	$T -test.run=^TestIteration2$$
t3:
	$T -test.run=^TestIteration3$$
t4:
	$T -test.run=^TestIteration4$$
t5:
	export SERVER_HOST=$(shell $R domain); \
	export SERVER_PORT=$(shell $R unused-port); \
	$T -test.run=^TestIteration5$$ \
	   -server-host=$$SERVER_HOST \
	   -server-port=$$SERVER_PORT \
	   -server-base-url="http://$$SERVER_HOST:$$SERVER_PORT" \
    ;
t6:
    export SERVER_PORT=$(shell $R unused-port); \
  	export TEMP_FILE=$(shell $R tempfile); \
  	$T -test.run=^TestIteration6$$ \
	   -server-port=$$SERVER_PORT \
	   -file-storage-path=$$TEMP_FILE \
    ;
t7:
    export SERVER_PORT=$(shell $R unused-port); \
  	export TEMP_FILE=$(shell $R tempfile); \
  	$T -test.run=^TestIteration7$$ \
	   -server-port=$$SERVER_PORT \
	   -file-storage-path=$$TEMP_FILE \
	;
t8:
	$T -test.run=^TestIteration8$$
t9:
	$T -test.run=^TestIteration9$$

DD = postgres://yandex_practicum:${DB_PASSWORD}@${DB_HOST}:5432/yp_shortener?sslmode=disable
t10:
	$T -test.run=^TestIteration10$$ -database-dsn=${DD}
t11:
	$T -test.run=^TestIteration11$$ -database-dsn=${DD}
t12:
	$T -test.run=^TestIteration12$$ -database-dsn=${DD}
t13:
	$T -test.run=^TestIteration13$$ -database-dsn=${DD}




