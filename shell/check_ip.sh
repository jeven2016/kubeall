#!/bin/bash

LOG_DIR=/var/log/check_kafka_ip.log
CONFIG_FILE="docker-compose.yaml"


# Function to get network interfaces
get_interfaces() {
    ip link show | grep -E '^[0-9]+:' | awk '{print $2}' | sed 's/://g'
}

# Function to get IP address of an interface
get_ip_address() {
    local interface=$1
    ip addr show "$interface" | grep -oE 'inet [0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' | awk '{print $2}'
}

# Function to get IP from config file
get_config_ip() {
    local config_file=$1
    grep "KAFKA_CFG_ADVERTISED_LISTENERS" "$config_file" | grep -oE 'SSL://[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' | cut -d'/' -f3
}

# Function to update config file
update_config() {
    local ip=$1
    local config_file=$2
    cp "$config_file" "$config_file.bak" || {
        echo "Error: Failed to create backup"
        exit 1
    }
    sed -i "s|SSL://[0-9]\+\.[0-9]\+\.[0-9]\+\.[0-9]\+|SSL://$ip|g" "$config_file"
}

# Function to create script and cronjob
create_cronjob() {
    local interface=$1
    local config_file=$2
    local script_file="/usr/local/bin/check_kafka_ip.sh"

    sudo touch $LOG_DIR

    # Create separate check script file
    cat > "$script_file" << 'EOF'
#!/bin/bash

# Check and update Kafka advertised listener IP
CONFIG_FILE="CONFIG_FILE_PATH"
INTERFACE="SELECTED_INTERFACE"
LOG_DIR="LOG_FILE_DIR"

current_ip=$(ip addr show $INTERFACE | grep -oE 'inet [0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' | awk '{print $2}')
config_ip=$(grep "KAFKA_CFG_ADVERTISED_LISTENERS" "$CONFIG_FILE" | grep -oE 'SSL://[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' | cut -d'/' -f3)

if [ -z "$current_ip" ]; then
    echo "$(date "+%Y-%m-%d %H:%M:%S") warning: No IP found for interface $INTERFACE" | sudo tee $LOG_DIR
    exit 1
fi

if [ "$current_ip" != "$config_ip" ]; then
    cp "$CONFIG_FILE" "$CONFIG_FILE.bak"
    sed -i "s|SSL://[0-9]\+\.[0-9]\+\.[0-9]\+\.[0-9]\+|SSL://$current_ip|g" "$CONFIG_FILE"
    echo "$(date "+%Y-%m-%d %H:%M:%S") IP updated from $config_ip to $current_ip" | sudo tee $LOG_DIR
else
    echo "$(date "+%Y-%m-%d %H:%M:%S") IP matches, no update needed" | sudo tee $LOG_DIR
fi
EOF

    # Replace placeholders in script
    sed -i "s|CONFIG_FILE_PATH|$config_file|g" "$script_file"
    sed -i "s|SELECTED_INTERFACE|$interface|g" "$script_file"
    sed -i "s|LOG_FILE_DIR|$LOG_DIR|g" "$script_file"
    
    chmod +x "$script_file"

    # Create cronjob (runs every 1 minutes, adjust as needed)
    local cron_entry="* * * * * root $script_file"
    
    # Check if cronjob already exists
    if ! crontab -l 2>/dev/null | grep -q "$script_file"; then
        # Add cronjob
        (crontab -l 2>/dev/null; echo "$cron_entry") | crontab -
        echo "Cronjob created: runs every 1 minutes" 
    else
        echo "Cronjob already exists"
    fi
}

main() {
    # Configuration
    echo "Network interfaces list:"
    interfaces=($(get_interfaces))
    local count=${#interfaces[@]}

    echo "Found $count network interfaces:"
    for i in "${!interfaces[@]}"; do
        ip=$(get_ip_address "${interfaces[$i]}")
        if [ -n "$ip" ]; then
            echo "[$i] ${interfaces[$i]}: $ip"
        else
            echo "[$i] ${interfaces[$i]}: No IP address"
        fi
    done

    echo -n "Please select an interface (enter number 0-$(($count-1))): "
    read selection

    if ! [[ "$selection" =~ ^[0-9]+$ ]] || [ "$selection" -ge "$count" ] || [ "$selection" -lt 0 ]; then
        echo "Error: Invalid selection"
        exit 1
    fi

    local selected_interface=${interfaces[$selection]}
    local selected_ip=$(get_ip_address "$selected_interface")

    if [ -z "$selected_ip" ]; then
        echo "Error: Selected interface has no IP address"
        exit 1
    fi

    echo "Selected interface: $selected_interface"
    echo "IP address: $selected_ip"

    if [ ! -f "$CONFIG_FILE" ]; then
        echo "Error: Configuration file $CONFIG_FILE not found"
        exit 1
    fi

    # Create script and cronjob
    echo "Creating check script and cronjob..."
    create_cronjob "$selected_interface" "$CONFIG_FILE"
    
    echo "Setup completed successfully"
    echo "Cronjob will run every 1 minutes"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "Error: Please run as root (sudo)"
    exit 1
fi

main