GOCMD=go

# 目标程序名称
Target=wnlc

# 默认目标
all: $(Target)

# 构建wnlc
$(Target): main.go
	$(GOCMD) build -o $(Target) main.go

install: $(Target)
	sudo cp $(Target) /usr/bin

.PHONY: install

# 清理目标
clean:
	sudo rm -rf /usr/bin/$(Target)
	rm -rf $(Target)

# 声明clean为伪目标
.PHONY: clean