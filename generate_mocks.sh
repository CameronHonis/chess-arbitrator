#!/bin/bash

cd helpers || exit
mkdir -p mocks
rm -rf mocks/*
$GOPATH/bin/mockgen -destination mocks/auth_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/auth AuthenticationServiceI &>> mocks/auth_service_mock.go
$GOPATH/bin/mockgen -destination mocks/clients_manager_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/clients_manager ClientsManagerI &>> mocks/clients_manager_mock.go
$GOPATH/bin/mockgen -destination mocks/matcher_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/matcher MatcherServiceI &>> mocks/matcher_service_mock.go
$GOPATH/bin/mockgen -destination mocks/matchmaking_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/matchmaking MatchmakingServiceI &>> mocks/matchmaking_service_mock.go
$GOPATH/bin/mockgen -destination mocks/msg_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/msg_service MessageServiceI &>> mocks/msg_service_mock.go
$GOPATH/bin/mockgen -destination mocks/router_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/router_service RouterServiceI &>> mocks/router_service_mock.go
$GOPATH/bin/mockgen -destination mocks/sub_service_mock.go -package mocks github.com/CameronHonis/chess-arbitrator/sub_service SubscriptionServiceI &>> mocks/sub_service_mock.go
$GOPATH/bin/mockgen -destination mocks/logger_service_mock.go -package mocks github.com/CameronHonis/log LoggerServiceI &>> mocks/logger_service_mock.go