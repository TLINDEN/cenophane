all:
	make -C upd
	make -C upctl

clean:
	make -C upd clean
	make -C upctl clean
