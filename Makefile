GOSRC_SHARED := go.mod 				\
		./internal/config/config.go 	\
		./internal/encoding/decoder.go 	\
		./internal/encoding/encoder.go 	\
		./internal/util/forward.go

GOSRC_CLIENT := ./cmd/client/main.go
GOSRC_SERVER := ./cmd/server/main.go
GOSRC_TEST_REMOTE := ./cmd/test-remote/main.go

CLIENT_BIN := bin/client
SERVER_BIN := bin/server
TEST_REMOTE_BIN := tests/test-remote

all: $(CLIENT_BIN) $(SERVER_BIN)

test: all $(TEST_REMOTE_BIN)
	-cp ./bin/client ./bin/server ./tests/

$(CLIENT_BIN): $(GOSRC_SHARED) $(GOSRC_CLIENT)
	go build -o $@ $(GOSRC_CLIENT)

$(SERVER_BIN): $(GOSRC_SHARED) $(GOSRC_SERVER)
	go build -o $@ $(GOSRC_SERVER)

$(TEST_REMOTE_BIN): $(GOSRC_SHARED) $(GOSRC_TEST_REMOTE)
	go build -o $@ $(GOSRC_TEST_REMOTE)
