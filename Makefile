SHELL := /bin/bash 
BASEDIR = $(shell pwd)
NEWCMD=${GOPATH}/src/${name}

versionDir="cmdctl/pkg/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
buildDate=$(shell TZ=Asia/Shanghai date +%FT%T%z)
	 
all:
	@echo compiling ...
	@go build -v -ldflags "-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}" -o cmdctl

gotool:
	@echo formating ...
	@-gofmt -w  .
	@-go tool vet . |& grep -v vendor

cmd: clean
ifneq (${name},)
	$(call cmd)
else
	@echo "Please specify the command name, like: make cmd name=newctl"
endif

install: all
	$(call install)

clean:
	rm -f cmdctl
	find . -name "[._]*.s[a-w][a-z]" | xargs -i rm -f {}

help:
	@echo "make                 - compile"
	@echo "make gotool          - run go tool"
	@echo "make cmd name=newctl - create new command named: newctl"
	@echo "make clean           - do some clean job"
	@echo "make install         - install command"

# define functions
define install
	@echo install cmdctl.yaml
	@mkdir -p ${HOME}/.cmdctl && cp cmdctl.yaml ${HOME}/.cmdctl

	@echo install cmdctl
    @mkdir -p ${HOME}/bin && mv cmdctl ${HOME}/.cmdctl/
endef

define cmd
	@rm -f cmdctl
	@mkdir -p ${GOPATH}/src
	@cd .. && cp -a cmdctl ${NEWCMD} && cd ${NEWCMD} && sed -i 's/cmdctl/${name}/g' `grep -Rl cmdctl *`
	@cd ${NEWCMD} && mv cmdctl.yaml ${name}.yaml
	@echo "new cmd locate in: ${NEWCMD}"
endef

.PHONY: all gotool clean cmd install
