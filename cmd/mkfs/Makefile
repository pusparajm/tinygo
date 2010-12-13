# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

ifeq ($(GOOS),tiny)
$(error Should only compile and use this tool on the host, not with GOOS=tiny)
endif

TARG=mkfs
GOFILES=\
	mkfs.go\

include $(GOROOT)/src/Make.cmd

