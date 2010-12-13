DIRS=draw fs \
	cmd/demo cmd/hello cmd/mapscroll cmd/mkfs cmd/palette

clean:
	for i in $(DIRS); do \
		gomake -C $$i clean ;\
	done
