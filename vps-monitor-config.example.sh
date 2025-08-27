# VPS Monitor Configuration Example
# Copy this file to vps-monitor-config.sh and modify the values as needed

# VPS Configuration
export VPS_PLAN_CODE="vps-2025-model2"        # VPS plan to monitor
export COUNTRY="US"                          # Country code
export ENDPOINT="ovh-us"                     # OVH API endpoint
export DATACENTERS="US-WEST-OR,US-EAST-VA"   # Comma-separated list of datacenters to monitor
export PREFERRED_ORDER="true"                # Use datacenter order (true) or random (false)
export OS_OPTION="os=option-linux"           # OS option for ordering

# Monitoring Configuration
export CHECK_INTERVAL="300"                  # Check every 5 minutes (300 seconds)
export MAX_RETRIES="3"                       # Max order retry attempts
export LOG_FILE="vps-monitor.log"            # Log file location

# Optional: Webhook for notifications (Slack, Discord, etc.)
# export WEBHOOK_URL="https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"

# OVH API Credentials (REQUIRED)
# Get these from: https://eu.api.ovh.com/console/
export OVH_APP_KEY="your_ovh_app_key_here"
export OVH_APP_SECRET="your_ovh_app_secret_here"
export OVH_CONSUMER_KEY="your_ovh_consumer_key_here"

# Alternative VPS configurations you can use:
#
# VPS Starter 1-2-20
# export VPS_PLAN_CODE="vps-starter-1-2-20"
# export OS_OPTION="os=option-linux"
#
# VPS Value 1-2-40
# export VPS_PLAN_CODE="vps-value-1-2-40"
# export OS_OPTION="os=option-linux"
#
# VPS Essential 2-4-40
# export VPS_PLAN_CODE="vps-essential-2-4-40"
# export OS_OPTION="os=option-linux"

# Multiple datacenter examples:
#
# US datacenters (use with ENDPOINT="ovh-us")
# export DATACENTERS="US-WEST-OR,US-EAST-VA"  # Both US datacenters
# export DATACENTERS="US-WEST-OR"              # Only Oregon
#
# European datacenters (use with ENDPOINT="ovh-eu")
# export DATACENTERS="GRA,SBG,BHS"             # Top 3 European datacenters
# export DATACENTERS="GRA,SBG,BHS,DE,UK,WAW"   # All European datacenters
# export DATACENTERS="GRA"                     # Only Gravelines, France
# export DATACENTERS="SBG"                     # Only Strasbourg, France
# export DATACENTERS="BHS"                     # Only Beauharnois, Canada
# export DATACENTERS="DE"                      # Only Frankfurt, Germany
# export DATACENTERS="UK"                      # Only London, UK
# export DATACENTERS="WAW"                     # Only Warsaw, Poland