GOSRC_SHARED := go.mod 				\
		./internal/config/config.go 	\
		./internal/encoding/decoder.go 	\
		./internal/encoding/encoder.go

GOSRC_CLIENT := ./cmd/client/main.go
GOSRC_SERVER := ./cmd/server/main.go

CLIENT_BIN := bin/client
SERVER_BIN := bin/server

all: $(CLIENT_BIN) $(SERVER_BIN)

$(CLIENT_BIN): $(GOSRC_SHARED) $(GOSRC_CLIENT)
	go build -o $@ $(GOSRC_CLIENT)

$(SERVER_BIN): $(GOSRC_SHARED) $(GOSRC_SERVER)
	go build -o $@ $(GOSRC_SERVER)
