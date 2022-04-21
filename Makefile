# 后期需要放进sgx中，目前仅仅作为一个服务提供测试使用
build:
	go build ./
run:
	nohup ./sgx-sign-service >nohup.out 2>&1 &
stop:
	pkill sgx-sign-service