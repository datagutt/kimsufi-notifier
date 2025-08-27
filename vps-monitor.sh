#!/bin/bash

# VPS Availability Monitor and Auto-Order Script
# This script checks for VPS availability and automatically places an order when stock is available

set -e  # Exit on any error

# Configuration - Modify these variables as needed
VPS_PLAN_CODE="${VPS_PLAN_CODE:-vps-2025-model2}"
COUNTRY="${COUNTRY:-US}"
ENDPOINT="${ENDPOINT:-ovh-us}"
DATACENTERS="${DATACENTERS:-US-WEST-OR}"  # Comma-separated list of datacenters
OS_OPTION="${OS_OPTION:-os=option-linux}"
CHECK_INTERVAL="${CHECK_INTERVAL:-300}"  # 5 minutes default
MAX_RETRIES="${MAX_RETRIES:-3}"
LOG_FILE="${LOG_FILE:-vps-monitor.log}"
PREFERRED_ORDER="${PREFERRED_ORDER:-true}"  # Try datacenters in order if true, random if false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    local level="$1"
    local message="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "$timestamp [$level] $message" | tee -a "$LOG_FILE"
}

# Check if required environment variables are set
check_env_vars() {
    local missing_vars=()

    if [[ -z "${OVH_APP_KEY:-}" ]]; then
        missing_vars+=("OVH_APP_KEY")
    fi
    if [[ -z "${OVH_APP_SECRET:-}" ]]; then
        missing_vars+=("OVH_APP_SECRET")
    fi
    if [[ -z "${OVH_CONSUMER_KEY:-}" ]]; then
        missing_vars+=("OVH_CONSUMER_KEY")
    fi

    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        log "ERROR" "${RED}Missing required environment variables: ${missing_vars[*]}${NC}"
        log "INFO" "Please set the following environment variables:"
        log "INFO" "  export OVH_APP_KEY='your_app_key'"
        log "INFO" "  export OVH_APP_SECRET='your_app_secret'"
        log "INFO" "  export OVH_CONSUMER_KEY='your_consumer_key'"
        exit 1
    fi
}

# Check VPS availability in multiple datacenters
check_availability() {
    local plan_code="$1"
    local country="$2"
    local endpoint="$3"
    local datacenters="$4"  # Comma-separated list

    log "INFO" "${BLUE}Checking availability for $plan_code across datacenters: $datacenters${NC}"

    # Convert comma-separated string to array
    IFS=',' read -ra DC_ARRAY <<< "$datacenters"

    # Run the check command and capture output
    local output
    if ! output=$("./kimsufi-notifier" check --plan-code "$plan_code" --country "$country" --endpoint "$endpoint" 2>&1); then
        log "ERROR" "${RED}Failed to check availability: $output${NC}"
        return 1
    fi

    local available_datacenters=()

    # Check each datacenter
    for datacenter in "${DC_ARRAY[@]}"; do
        # Trim whitespace
        datacenter=$(echo "$datacenter" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

        if echo "$output" | grep -q "$datacenter.*available"; then
            available_datacenters+=("$datacenter")
            log "SUCCESS" "${GREEN}✅ $plan_code is AVAILABLE in $datacenter!${NC}"
        else
            log "DEBUG" "${YELLOW}❌ $plan_code is not available in $datacenter${NC}"
        fi
    done

    # Log all availability status for debugging
    log "DEBUG" "Full availability status:"
    echo "$output" | while read -r line; do
        if echo "$line" | grep -q "datacenter\|status"; then
            log "DEBUG" "  $line"
        fi
    done

    # Return first available datacenter or empty if none available
    if [[ ${#available_datacenters[@]} -gt 0 ]]; then
        if [[ "$PREFERRED_ORDER" == "true" ]]; then
            # Return first datacenter in preferred order
            echo "${available_datacenters[0]}"
        else
            # Return random available datacenter
            echo "${available_datacenters[$RANDOM % ${#available_datacenters[@]}]}"
        fi
        return 0
    else
        log "INFO" "${YELLOW}❌ $plan_code is not available in any of the specified datacenters${NC}"
        return 1
    fi
}

# Place an order
place_order() {
    local plan_code="$1"
    local country="$2"
    local endpoint="$3"
    local datacenter="$4"
    local os_option="$5"

    log "INFO" "${GREEN}🚀 Placing order for $plan_code in $datacenter...${NC}"

    local order_cmd="./kimsufi-notifier order --plan-code \"$plan_code\" --country \"$country\" --endpoint \"$endpoint\" -d \"$datacenter\" --item-option \"$os_option\""

    # Execute the order command
    if output=$($order_cmd 2>&1); then
        log "SUCCESS" "${GREEN}✅ Order placed successfully!${NC}"
        echo "$output"
        # Extract order URL if present
        echo "$output" | grep "order completed" | while read -r line; do
            log "SUCCESS" "${GREEN}$line${NC}"
        done
        return 0
    else
        log "ERROR" "${RED}❌ Order failed: $output${NC}"
        return 1
    fi
}

# Send notification (optional - requires curl)
send_notification() {
    local message="$1"
    local webhook_url="${WEBHOOK_URL:-}"

    if [[ -n "$webhook_url" ]]; then
        curl -s -X POST "$webhook_url" \
            -H 'Content-Type: application/json' \
            -d "{\"text\":\"$message\"}" || true
    fi
}

# Main monitoring loop
main() {
    log "INFO" "${BLUE}🚀 Starting VPS availability monitor${NC}"
    log "INFO" "Configuration:"
    log "INFO" "  Plan Code: $VPS_PLAN_CODE"
    log "INFO" "  Country: $COUNTRY"
    log "INFO" "  Endpoint: $ENDPOINT"
    log "INFO" "  Datacenters: $DATACENTERS"
    log "INFO" "  OS Option: $OS_OPTION"
    log "INFO" "  Check Interval: $CHECK_INTERVAL seconds"
    log "INFO" "  Preferred Order: $PREFERRED_ORDER"
    log "INFO" "  Log File: $LOG_FILE"

    # Check environment variables
    check_env_vars

    # Check if binary exists
    if [[ ! -f "./kimsufi-notifier" ]]; then
        log "ERROR" "${RED}kimsufi-notifier binary not found in current directory${NC}"
        log "INFO" "Please build the binary first: go build -o kimsufi-notifier ."
        exit 1
    fi

    local attempt=0
    local start_time=$(date +%s)

    while true; do
        ((attempt++))
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))

        log "INFO" "${BLUE}=== Attempt #$attempt (Running for ${elapsed}s) ===${NC}"

        # Check availability
        local available_datacenter
        if available_datacenter=$(check_availability "$VPS_PLAN_CODE" "$COUNTRY" "$ENDPOINT" "$DATACENTERS"); then
            # VPS is available in at least one datacenter, try to place order
            local order_success=false

            for ((retry=1; retry<=MAX_RETRIES; retry++)); do
                log "INFO" "Order attempt $retry/$MAX_RETRIES for $available_datacenter"

                if place_order "$VPS_PLAN_CODE" "$COUNTRY" "$ENDPOINT" "$available_datacenter" "$OS_OPTION"; then
                    order_success=true
                    send_notification "✅ VPS Order Successful! $VPS_PLAN_CODE in $available_datacenter"
                    log "SUCCESS" "${GREEN}🎉 Order completed successfully! Exiting monitor.${NC}"
                    exit 0
                else
                    if [[ $retry -lt $MAX_RETRIES ]]; then
                        log "WARN" "${YELLOW}Order attempt $retry failed, retrying in 10 seconds...${NC}"
                        sleep 10
                    fi
                fi
            done

            if [[ "$order_success" == false ]]; then
                send_notification "❌ VPS Order Failed after $MAX_RETRIES attempts: $VPS_PLAN_CODE in $available_datacenter"
                log "ERROR" "${RED}All order attempts failed. Continuing to monitor...${NC}"
            fi
        fi

        # Wait before next check
        log "INFO" "${BLUE}⏰ Waiting $CHECK_INTERVAL seconds before next check...${NC}"
        sleep "$CHECK_INTERVAL"
    done
}

# Handle script interruption
trap 'log "INFO" "${YELLOW}Script interrupted by user${NC}"; exit 130' INT TERM

# Show usage if requested
if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
    echo "VPS Availability Monitor and Auto-Order Script"
    echo ""
    echo "This script monitors VPS availability and automatically places orders when stock is available."
    echo ""
    echo "Configuration (set these environment variables or modify the script):"
    echo "  VPS_PLAN_CODE    - VPS plan code to monitor (default: vps-2025-model2)"
    echo "  COUNTRY          - Country code (default: US)"
    echo "  ENDPOINT         - OVH endpoint (default: ovh-us)"
    echo "  DATACENTERS      - Comma-separated list of datacenters to monitor (default: US-WEST-OR)"
    echo "  PREFERRED_ORDER  - Use datacenter order (true) or random selection (false) (default: true)"
    echo "  OS_OPTION        - OS option for order (default: os=option-linux)"
    echo "  CHECK_INTERVAL   - Seconds between availability checks (default: 300)"
    echo "  MAX_RETRIES      - Max order retry attempts (default: 3)"
    echo "  LOG_FILE         - Log file path (default: vps-monitor.log)"
    echo "  WEBHOOK_URL      - Optional webhook URL for notifications"
    echo ""
    echo "Required OVH API credentials:"
    echo "  OVH_APP_KEY      - Your OVH application key"
    echo "  OVH_APP_SECRET   - Your OVH application secret"
    echo "  OVH_CONSUMER_KEY - Your OVH consumer key"
    echo ""
    echo "Usage:"
    echo "  ./vps-monitor.sh              # Start monitoring"
    echo "  ./vps-monitor.sh --help       # Show this help"
    echo ""
    echo "Example:"
    echo "  export OVH_APP_KEY='your_key'"
    echo "  export OVH_APP_SECRET='your_secret'"
    echo "  export OVH_CONSUMER_KEY='your_consumer_key'"
    echo "  VPS_PLAN_CODE=vps-starter-1-2-20 ./vps-monitor.sh"
    exit 0
fi

# Run main function
main "$@"