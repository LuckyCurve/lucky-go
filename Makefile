# Makefile

# 查找所有 go.mod 文件
GOMODS := $(shell find . -name "go.mod")

# 升级所有依赖
upgrade:
	@for mod in $(GOMODS); do \
		echo "===> Updating dependencies in $$(dirname $$mod)"; \
		(cd $$(dirname $$mod) && go get -u ./... && go mod tidy); \
	done

# 安装当前项目
install:
	@echo "===> Installing current module"
	go install ./...

# 一键执行：升级依赖 + 安装
all: upgrade install
