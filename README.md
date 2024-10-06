# SmartHome Monitoring

This project is a smart home monitoring system that allows you to monitor and track data from various sensors using MQTT. The project is designed to be easy to deploy on devices such as Raspberry Pi, and it provides support for monitoring using Prometheus and Grafana.

## Features
- Monitor sensors via MQTT broker.
- Docker-based deployment for easy installation.
- Prometheus metrics for monitoring the performance.
- Grafana integration for visualizing the data.
- Configurable via YAML and environment variables.

## Requirements
- Docker
- Docker Compose
- Raspberry Pi (or other Linux-based system)
- MQTT Broker (e.g., Mosquitto)

## Installation

### 1. Clone the Repository
```sh
git clone https://github.com/lkobylski/smarthome-monitor.git
cd smarthome-monitoring
```

### 2. Configure the Project

1. **Environment Variables**
    - Copy the example `.env.example` to create a `.env` file and edit it with your configuration.
   ```sh
   cp .env.example .env
   ```
    - Edit `.env` to include your config file, debug, etc.
    - If you are using a different System than Raspberry Pi, change the `GOARCH` variable in .env with your architecutre value (`arm64|amd64|other`).

2. **YAML Configuration**
    - Copy the example `config.yaml.example` to `config.yaml` and modify it as needed.
   ```sh
   cp config.yaml.example config.yaml
   ```
    - Set the broker, client ID, list of topics and devices you wish to monitor.

### 3. Build and Run the Project

To build the Docker images and run the services, use the Makefile:
```sh
make all
```
This command will build the Docker image and start all the services in the background.

### 4. View Logs
To view the logs of the monitoring application:
```sh
make logs
```

### 5. Stop the Services
To stop all running services:
```sh
make stop
```

## Monitoring with Prometheus and Grafana
The project includes configuration for exposing Prometheus metrics. You can use Grafana to visualize these metrics.

- **Prometheus Port**: 2112 (exposed by default)
- **Grafana**: The Grafana service is also available in the `docker-compose.yml` for monitoring purposes.

### Access Grafana Dashboard
1. Open your browser and navigate to the IP address of your device, using the Grafana port `3005`.
2. Use the default credentials (`admin/admin`) to log in, and configure your data sources to include Prometheus.

## Configuration Overview

### Environment Variables (`.env`)
The `.env` file allows you to configure basic parameters like:
- `CONFIG_FILE`: path to your config.yaml file.
- `DEBUG`: Debug mode for logs.

### YAML Configuration (`config.yaml`)
The `config.yaml` file provides a more detailed configuration for the monitoring:
- **broker**: The address of your MQTT broker.
- **client_id**: The MQTT client identifier.
- **topics**: List of topics that should be monitored.
- **disconnect_timeout**: The timeout in milliseconds for disconnecting the client.
- **devices**: A list of devices to be monitored, where each device has:
    - **id**: Unique identifier for the device.
    - **topic**: The MQTT topic for the device.
    - **data_format**: The format of the data being published by the device (e.g., `json`).

### Adding Bluetooth Devices

If you have Bluetooth devices, you can integrate them using Node-RED, which will capture signals from your Bluetooth devices and emit corresponding MQTT messages. You can then add these MQTT topics to your application's monitoring configuration.

**Example Configuration with Node-RED**:
- Set up Node-RED to listen to your Bluetooth devices.
- Use Home Assistant (if available) to receive sensor states and forward them to MQTT through Node-RED automation.
- Add the new topics to the `topics` section of `config.yaml` to include them in the monitoring.

```
topics:
  - "zigbee2mqtt/#"
  - "nodered/#"
devices:
  - id: "xiaomi_temp_mini"
    topic: "zigbee2mqtt/mi_small_temp_temperature_sensor"
    data_format: "json"
  - id: "temp_bt_livingroom"
    topic: "nodered/temp_bt_livingroom"
    data_format: "json"
```

In my case, I use Home Assistant and Node-RED to capture the state of sensors and pass it to MQTT using simple automation. Feel free to reach out if you need help with this setup.
  
## Multi-Platform Build
If you want to build the Docker image for a specific architecture (e.g., Raspberry Pi, which uses `arm64` architecture), change the `GOARCH` variable in .env file:
```sh
GOARCH=arm64
```
This ensures that the image is compatible with the Raspberry Pi architecture.

## Troubleshooting

- **Error: `exec format error`**
    - This error occurs when the application is built for the wrong architecture. Ensure you build the Docker image for the correct platform using `--platform linux/arm64` for Raspberry Pi.

- **Docker Compose Compatibility Issues**
    - Ensure you have Docker Compose v2 installed, as older versions may not support all the features used in this project. To install Docker Compose v2:
  ```sh
  sudo apt-get update
  sudo apt-get install docker-compose-plugin
  ```

## Contributing
Contributions are welcome! Feel free to open issues or submit pull requests to improve the project 💪.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

