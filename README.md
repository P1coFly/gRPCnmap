# gRPC сервис обертка над nmap

## Задание
Написать gRPC сервис обертку над nmap с использованием следующего скрипта
https://github.com/vulnersCom/nmap-vulners и предлагаемого API:

```proto
syntax = "proto3"; 

package netvuln_v1;

option go_package = "github.com/P1coFly/gRPCnmap/pkg/netvuln_v1;netvuln_v1";

service NetVulnService {
    rpc CheckVuln(CheckVulnRequest) returns (CheckVulnResponse);
}

message CheckVulnRequest {
    repeated string targets = 1; // IP addresses 
    repeated int32 tcp_port = 2; // only TCP ports
}

message CheckVulnResponse {
    repeated TargetResult results = 1;
}

message TargetResult {
    string target = 1; // target IP 
    repeated Service services = 2;
}

message Service {
    string name = 1;
    string version = 2; 
    int32 tcp_port = 3; 
    repeated Vulnerability vulns = 4;
}

message Vulnerability {
    string identifier = 1; 
    float cvss_score = 2;
}
```
# Запуск проекта

## Установка необходимых зависимостей
Необходимо установить на рабочую станцию, где будет запускаться сервис следующее:
- nmap
- скрипт vulners.nse (инструкция - https://github.com/vulnersCom/nmap-vulners)
- ЯП go (использовалась версия go 1.22.3)
- golangci-lint (для запуска линтера)

## Запуск
1. Необходимо установить зависимости, для это из корня проекта надо выполнить команду - ```go mod download```
2. Скомпелировать выполняемый файл -  ```go build -o <название_выполняемого_файла> ./cmd/grpc-nmap/main.go``` или ```make build```
3. Запуск
    - запуск с флагом config - ```./<название_выполняемого_файла> --config=./path/to/cfg.yml```
    - запуск с переменной окружения. Необходимо задать CONFIG_PATH=./path/to/cfg.yml и запустить - ```./<название_выполняемого_файла>```
Приоритет запуска: флаг > переменной окружения

# Makefile
В makefile реализованы следующие команды:
- generate - кодагенерация по файлу .proto (api/NetVuln_v1/service.proto)
- lint - запуск линтера
- build - сборка проекта golangci-lint
- test - запуск тестов

# Github Actions
При пуше в ветку master запускаются тесты и линтер