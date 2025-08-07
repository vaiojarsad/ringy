BIN_BASE_NAME=ringy
DAEMON_NAME=ringy
DAEMON_USER=ringy
DAEMON_GROUP=ringy
DAEMON_DESC=Ringy service
INSTALL_DIR=/opt
SERVICE_FILE=/etc/systemd/system/$(DAEMON_NAME).service
SERVICE_ENV=/etc/default/$(DAEMON_NAME)

ENV_AUDIO_FILES_PATH=/var/lib/$(DAEMON_NAME)/
ENV_LOGGER_PREFIX=$(DAEMON_NAME)
ENV_NET_MULTI_CAST_ADDR=239.255.255.12
ENV_NET_PORT=9000
ENV_RING_AUDIO=aguanta.mp3

ifeq ($(OS), Windows_NT)
    BIN_NAME := $(BIN_BASE_NAME).exe
else
    BIN_NAME := $(BIN_BASE_NAME)
endif

build:
	go build -o $(BIN_NAME) .

clean:
	rm -f $(BIN_NAME)

update: 
	sudo cp $(BIN_NAME) $(INSTALL_DIR)/$(BIN_NAME)
	sudo chown $(DAEMON_USER):$(DAEMON_USER) $(INSTALL_DIR)/$(BIN_NAME)

install: build create-group create-user add-to-group add-to-audio
	sudo cp $(BIN_NAME) $(INSTALL_DIR)/$(BIN_NAME)
	sudo chown $(DAEMON_USER):$(DAEMON_USER) $(INSTALL_DIR)/$(BIN_NAME)

	# Service env file
	echo "AUDIO_FILES_PATH=$(ENV_AUDIO_FILES_PATH)" | sudo tee $(SERVICE_ENV) > /dev/null
	echo "LOGGER_PREFIX=$(ENV_LOGGER_PREFIX)" | sudo tee -a $(SERVICE_ENV) > /dev/null
	echo "NET_MULTI_CAST_ADDR=$(ENV_NET_MULTI_CAST_ADDR)" | sudo tee -a $(SERVICE_ENV) > /dev/null
	echo "NET_PORT=$(ENV_NET_PORT)" | sudo tee -a $(SERVICE_ENV) > /dev/null
	echo "RING_AUDIO=$(ENV_RING_AUDIO)" | sudo tee -a $(SERVICE_ENV) > /dev/null

	# Service unit file
	echo "[Unit]" | sudo tee $(SERVICE_FILE) > /dev/null
	echo "Description=$(DAEMON_DESC)" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "After=network.target" | sudo tee -a $(SERVICE_FILE) > /dev/null

	echo "[Service]" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "Type=simple" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "ExecStart=$(INSTALL_DIR)/$(BIN_NAME)" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "EnvironmentFile=$(SERVICE_ENV)" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "Restart=on-failure" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "User=$(DAEMON_USER)" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "Group=$(DAEMON_GROUP)" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "StandardOutput=journal" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "StandardError=journal" | sudo tee -a $(SERVICE_FILE) > /dev/null

	echo "[Install]" | sudo tee -a $(SERVICE_FILE) > /dev/null
	echo "WantedBy=multi-user.target" | sudo tee -a $(SERVICE_FILE) > /dev/null

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
