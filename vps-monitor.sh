#!/bin/bash

# VPS Availability Monitor and Auto-Order Script
# This script checks for VPS availability and automatically places an order when stock is available

# set -e  # Exit on any error - commented out to allow script to continue on errors

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
STATE_FILE="${STATE_FILE:-vps-monitor.state}"  # State file to prevent duplicate orders

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

# State management functions
create_state_file() {
    local plan_code="$1"
    local datacenter="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    cat > "$STATE_FILE" << EOF
# VPS Monitor State File
# This file prevents duplicate orders
export ORDER_PLACED="true"
export PLAN_CODE="$plan_code"
export DATACENTER="$datacenter"
export TIMESTAMP="$timestamp"
export PID="$$"
EOF
    log "INFO" "${GREEN}✅ State file created: $STATE_FILE${NC}"
}

check_state_file() {
    if [[ -f "$STATE_FILE" ]]; then
        log "INFO" "${YELLOW}📁 Found existing state file: $STATE_FILE${NC}"

        # Read state file
        source "$STATE_FILE" 2>/dev/null || {
            log "WARN" "${YELLOW}Could not read state file, ignoring...${NC}"
            return 1
        }

        if [[ "${ORDER_PLACED:-}" == "true" ]]; then
            log "SUCCESS" "${GREEN}🎉 Order already placed successfully!${NC}"
            log "INFO" "Plan: ${PLAN_CODE:-unknown}"
            log "INFO" "Datacenter: ${DATACENTER:-unknown}"
            log "INFO" "Timestamp: ${TIMESTAMP:-unknown}"
            log "INFO" "${BLUE}To reset and monitor again, delete the state file: rm $STATE_FILE${NC}"
            return 0
        fi
    fi
    return 1
}

cleanup_state_file() {
    if [[ -f "$STATE_FILE" ]]; then
        rm -f "$STATE_FILE"
        log "INFO" "${BLUE}🧹 State file cleaned up${NC}"
    fi
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
        log "WARN" "${YELLOW}⚠️  Missing OVH API credentials: ${missing_vars[*]}${NC}"
        log "INFO" "${BLUE}📊 Running in MONITORING-ONLY mode (no orders will be placed)${NC}"
        log "INFO" "To enable auto-ordering, set the following environment variables:"
        log "INFO" "  export OVH_APP_KEY='your_app_key'"
        log "INFO" "  export OVH_APP_SECRET='your_app_secret'"
        log "INFO" "  export OVH_CONSUMER_KEY='your_consumer_key'"
        return 1  # Return 1 to indicate monitoring-only mode
    fi
    return 0  # Return 0 for full mode
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
    local cmd_exit_code
    if ! output=$("./kimsufi-notifier" check --plan-code "$plan_code" --country "$country" --endpoint "$endpoint" 2>&1); then
        cmd_exit_code=$?
        log "ERROR" "${RED}Failed to check availability (exit code: $cmd_exit_code): $output${NC}"
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

    # Execute the order command directly to avoid shell quoting issues
    if output=$(./kimsufi-notifier order --plan-code "$plan_code" --country "$country" --endpoint "$endpoint" --datacenters "$datacenter" --item-option "$os_option" 2>&1); then
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
    log "INFO" "  State File: $STATE_FILE"

    # Check state file to prevent duplicate orders
    if check_state_file; then
        log "INFO" "${GREEN}Exiting to prevent duplicate orders. Delete $STATE_FILE to reset.${NC}"
        exit 0
    fi

    # Check environment variables
    local full_mode=true
    if ! check_env_vars; then
        full_mode=false
    fi

    # Check if binary exists
    if [[ ! -f "./kimsufi-notifier" ]]; then
        log "ERROR" "${RED}kimsufi-notifier binary not found in current directory${NC}"
        log "INFO" "Please build the binary first: go build -o kimsufi-notifier ."
        exit 1
    fi

    # Display mode
    if [[ "$full_mode" == true ]]; then
        log "INFO" "${GREEN}🔥 Running in FULL mode (monitoring + auto-ordering)${NC}"
    else
        log "INFO" "${BLUE}👀 Running in MONITORING-ONLY mode (no orders will be placed)${NC}"
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
        local check_result
        check_result=$(check_availability "$VPS_PLAN_CODE" "$COUNTRY" "$ENDPOINT" "$DATACENTERS")
        local check_exit_code=$?

        if [[ $check_exit_code -eq 0 ]] && [[ -n "$check_result" ]]; then
            available_datacenter="$check_result"
            if [[ "$full_mode" == true ]]; then
                # VPS is available in at least one datacenter, try to place order
                local order_success=false

                for ((retry=1; retry<=MAX_RETRIES; retry++)); do
                    log "INFO" "Order attempt $retry/$MAX_RETRIES for $available_datacenter"

                    if place_order "$VPS_PLAN_CODE" "$COUNTRY" "$ENDPOINT" "$available_datacenter" "$OS_OPTION"; then
                        order_success=true
                        # Create state file to prevent duplicate orders
                        create_state_file "$VPS_PLAN_CODE" "$available_datacenter"
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
            else
                # Monitoring-only mode - just notify about availability
                log "SUCCESS" "${GREEN}📢 STOCK FOUND: $VPS_PLAN_CODE available in $available_datacenter${NC}"
                send_notification "📢 STOCK ALERT: $VPS_PLAN_CODE available in $available_datacenter (monitoring-only mode)"
                log "INFO" "${BLUE}💡 Set OVH credentials to enable auto-ordering${NC}"
            fi
        fi

        # Wait before next check
        log "INFO" "${BLUE}⏰ Waiting $CHECK_INTERVAL seconds before next check...${NC}"
        sleep "$CHECK_INTERVAL"
    done
}

# Handle script interruption
trap 'log "INFO" "${YELLOW}Script interrupted by user${NC}"; exit 130' INT TERM

# Handle reset command
if [[ "${1:-}" == "--reset" ]]; then
    if [[ -f "$STATE_FILE" ]]; then
        cleanup_state_file
        log "SUCCESS" "${GREEN}✅ State file reset. Monitor will run normally on next start.${NC}"
    else
        log "INFO" "${BLUE}No state file found to reset.${NC}"
    fi
    exit 0
fi

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
    echo "  STATE_FILE       - State file to prevent duplicate orders (default: vps-monitor.state)"
    echo ""
    echo "Required OVH API credentials (for auto-ordering):"
    echo "  OVH_APP_KEY      - Your OVH application key"
    echo "  OVH_APP_SECRET   - Your OVH application secret"
    echo "  OVH_CONSUMER_KEY - Your OVH consumer key"
    echo ""
    echo "Note: If credentials are missing, the script runs in MONITORING-ONLY mode"
    echo ""
    echo "Usage:"
    echo "  ./vps-monitor.sh              # Start monitoring"
    echo "  ./vps-monitor.sh --reset      # Reset state file (allows new orders)"
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