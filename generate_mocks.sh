#!/bin/bash

cd helpers || exit

mkdir -p mocks
rm -rf mocks/*
#$GOPATH/bin/mockgen -destination mocks/auth_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/auth AuthenticationServiceI &>> mocks/auth_service_mock.go
#$GOPATH/bin/mockgen -destination mocks/clients_manager_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/clients_manager ClientsManagerI &>> mocks/clients_manager_mock.go
#$GOPATH/bin/mockgen -destination mocks/matcher_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/matcher MatcherServiceI &>> mocks/matcher_service_mock.go
#$GOPATH/bin/mockgen -destination mocks/matchmaking_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/matchmaking MatchmakingServiceI &>> mocks/matchmaking_service_mock.go
#$GOPATH/bin/mockgen -destination mocks/msg_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/msg_service MessageServiceI &>> mocks/msg_service_mock.go
#$GOPATH/bin/mockgen -destination mocks/router_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/router_service RouterServiceI &>> mocks/router_service_mock.go
#$GOPATH/bin/mockgen -destination mocks/sub_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/sub_service SubscriptionServiceI &>> mocks/sub_service_mock.go
#$GOPATH/bin/mockgen -destination mocks/logger_service_mock.go -package mocks github.com/CameronHonis/log LoggerServiceI &>> mocks/logger_service_mock.go


$GOPATH/bin/mockgen -source=../auth/auth_service.go -destination mocks/auth_service_mock.go -package mocks &>> mocks/auth_service_mock.go
$GOPATH/bin/mockgen -source=../clients_manager/clients_manager.go -destination mocks/clients_manager_mock.go -package mocks &>> mocks/clients_manager_mock.go
$GOPATH/bin/mockgen -source=../matcher/matcher_service.go -destination mocks/matcher_service_mock.go -package mocks &>> mocks/matcher_service_mock.go
$GOPATH/bin/mockgen -source=../matchmaking/matchmaking_service.go -destination mocks/matchmaking_service_mock.go -package mocks &>> mocks/matchmaking_service_mock.go
$GOPATH/bin/mockgen -source=../router_service/router_service.go -destination mocks/router_service_mock.go -package mocks &>> mocks/router_service_mock.go
$GOPATH/bin/mockgen -source=../sub_service/sub_service.go -destination mocks/sub_service_mock.go -package mocks &>> mocks/sub_service_mock.go
$GOPATH/bin/mockgen -source=../../log/logger_service.go -destination mocks/logger_service_mock.go -package mocks &>> mocks/logger_service_mock.go
