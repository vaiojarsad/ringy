BIN_BASE_NAME=ringy
DAEMON_NAME=ringy
DAEMON_USER=ringy
DAEMON_GROUP=ringy
DAEMON_DESC=Ringy service
INSTALL_DIR=/opt

ifeq ($(OS), Windows_NT)
    BIN_NAME := $(BIN_BASE_NAME).exe
else
    BIN_NAME := $(BIN_BASE_NAME)
endif

build:
	go build -o $(BIN_NAME) .

clean:
	rm -f $(BIN_NAME)

install: build create-group create-user add-to-group add-to-audio
	sudo cp $(BIN_NAME) $(INSTALL_DIR)/$(BIN_NAME)
	sudo chown $(DAEMON_USER):$(DAEMON_USER) $(INSTALL_DIR)/$(BIN_NAME)
	sudo sh -c 'tee /etc/systemd/system/$(DAEMON_NAME).service > /dev/null <<EOF \
[Unit] \
Description=$(DAEMON_DESC) \
After=network.target \
 \
[Service] \
ExecStart=$(INSTALL_DIR)/$(BIN_NAME) \
Restart=on-failure \
User=$(DAEMON_USER) \
Group=$(DAEMON_GROUP) \
StandardOutput=journal \
StandardError=journal \
 \
[Install] \
WantedBy=multi-user.target \
EOF'
	sudo systemctl daemon-reload
	sudo systemctl enable $(DAEMON_NAME)

create-group:
	@if ! getent group $(DAEMON_GROUP) > /dev/null 2>&1; then \
		  sudo groupadd $(DAEMON_GROUP); \
		  echo "Group '$(DAEMON_GROUP)' created."; \
	else \
		echo "Group '$(DAEMON_GROUP)' already exists."; \
	fi

create-user:
	@if id "$(DAEMON_USER)" >/dev/null 2>&1; then \
		echo "User '$(DAEMON_USER)' already exists"; \
	else \
		echo "Creating system user '$(DAEMON_USER)'..."; \
		sudo useradd --system --no-create-home --shell /usr/sbin/nologin -g $(DAEMON_GROUP) $(DAEMON_USER); \
	fi

add-to-audio:
	@if id -nG "$(DAEMON_USER)" | tr ' ' '\n' | grep -Fxq "audio"; then \
		echo "$(DAEMON_USER) is already in the audio group"; \
	else \
		echo "Adding $(DAEMON_USER) to audio group..."; \
		sudo usermod -aG audio $(DAEMON_USER); \
	fi

add-to-group:
	@if id -nG "$(DAEMON_USER)" | tr ' ' '\n' | grep -Fxq "$(DAEMON_GROUP)"; then \
		echo "$(DAEMON_USER) is already in the $(DAEMON_GROUP) group"; \
	else \
		echo "Adding $(DAEMON_USER) to $(DAEMON_GROUP) group..."; \
		sudo usermod -aG $(DAEMON_GROUP) $(DAEMON_USER); \
	fi
