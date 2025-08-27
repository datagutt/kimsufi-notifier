# VPS Availability Monitor & Auto-Order Script

A shell script that continuously monitors VPS availability and automatically places orders when stock becomes available.

## Features

- 🔍 **Continuous Monitoring** - Checks VPS availability at configurable intervals
- 🚀 **Auto-Ordering** - Automatically places orders when VPS becomes available
- 📊 **Detailed Logging** - Comprehensive logging with timestamps and colored output
- 🔄 **Retry Logic** - Handles order failures with configurable retry attempts
- 🔔 **Notifications** - Optional webhook notifications for order success/failure
- ⚙️ **Flexible Configuration** - Easy configuration via environment variables
- 🛡️ **Error Handling** - Robust error handling and recovery

## Quick Start

### 1. Setup OVH API Credentials

First, you need to obtain OVH API credentials:

1. Go to [OVH API Console](https://eu.api.ovh.com/console/)
2. Create a new application or use existing credentials
3. Note down your `Application Key`, `Application Secret`, and generate a `Consumer Key`

### 2. Configure the Script

```bash
# Copy the example configuration
cp vps-monitor-config.example.sh vps-monitor-config.sh

# Edit the configuration file with your details
nano vps-monitor-config.sh
```

### 3. Run the Monitor

```bash
# Source your configuration and run the monitor
source vps-monitor-config.sh
./vps-monitor.sh
```

## Configuration Options

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `OVH_APP_KEY` | Your OVH application key | `abc123def456` |
| `OVH_APP_SECRET` | Your OVH application secret | `def789ghi012` |
| `OVH_CONSUMER_KEY` | Your OVH consumer key | `jkl345mno678` |

### Optional Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `VPS_PLAN_CODE` | `vps-2025-model2` | VPS plan to monitor |
| `COUNTRY` | `US` | Country code (US, FR, etc.) |
| `ENDPOINT` | `ovh-us` | OVH API endpoint |
| `DATACENTERS` | `US-WEST-OR` | Comma-separated list of datacenters to monitor |
| `PREFERRED_ORDER` | `true` | Use datacenter order (true) or random selection (false) |
| `OS_OPTION` | `os=option-linux` | OS option for ordering |
| `CHECK_INTERVAL` | `300` | Seconds between checks |
| `MAX_RETRIES` | `3` | Max order retry attempts |
| `LOG_FILE` | `vps-monitor.log` | Log file location |
| `WEBHOOK_URL` | *(empty)* | Webhook URL for notifications |

## Usage Examples

### Basic Usage

```bash
# Monitor default VPS with environment variables
export OVH_APP_KEY="your_key"
export OVH_APP_SECRET="your_secret"
export OVH_CONSUMER_KEY="your_consumer_key"
./vps-monitor.sh
```

### Custom VPS Configuration

```bash
# Monitor VPS Starter in European datacenters (multiple)
export VPS_PLAN_CODE="vps-starter-1-2-20"
export COUNTRY="FR"
export ENDPOINT="ovh-eu"
export DATACENTERS="GRA,SBG,BHS"
export OS_OPTION="os=option-linux"
./vps-monitor.sh
```

### Multiple Datacenter Monitoring

```bash
# Monitor across all US datacenters
export DATACENTERS="US-WEST-OR,US-EAST-VA"
export PREFERRED_ORDER="true"
./vps-monitor.sh

# Monitor European datacenters with random selection
export ENDPOINT="ovh-eu"
export DATACENTERS="GRA,SBG,BHS,DE,UK,WAW"
export PREFERRED_ORDER="false"
./vps-monitor.sh
```

### With Slack Notifications

```bash
# Add webhook URL for Slack notifications
export WEBHOOK_URL="https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
./vps-monitor.sh
```

### Faster Checking

```bash
# Check every 60 seconds instead of 5 minutes
export CHECK_INTERVAL="60"
./vps-monitor.sh
```

## Popular VPS Configurations

### VPS Starter Series
```bash
export VPS_PLAN_CODE="vps-starter-1-2-20"    # 1 vCPU, 2GB RAM, 20GB SSD
export VPS_PLAN_CODE="vps-starter-1-2-40"    # 1 vCPU, 2GB RAM, 40GB SSD
export VPS_PLAN_CODE="vps-starter-1-2-80"    # 1 vCPU, 2GB RAM, 80GB SSD
```

### VPS Value Series
```bash
export VPS_PLAN_CODE="vps-value-1-2-40"      # 1 vCPU, 2GB RAM, 40GB SSD
export VPS_PLAN_CODE="vps-value-1-2-80"      # 1 vCPU, 2GB RAM, 80GB SSD
export VPS_PLAN_CODE="vps-value-1-4-80"      # 1 vCPU, 4GB RAM, 80GB SSD
```

### VPS Essential Series
```bash
export VPS_PLAN_CODE="vps-essential-2-4-40"   # 2 vCPU, 4GB RAM, 40GB SSD
export VPS_PLAN_CODE="vps-essential-2-4-80"   # 2 vCPU, 4GB RAM, 80GB SSD
```

## Multiple Datacenter Support

The script supports monitoring multiple datacenters simultaneously, which provides several benefits:

- **Higher Success Rate**: Increases chances of finding available stock
- **Flexible Selection**: Choose between ordered preference or random selection
- **Geographic Diversity**: Spread across different regions for better availability

### Selection Strategies

- **Preferred Order (default)**: Tries datacenters in the order you specify
- **Random Selection**: Randomly selects from available datacenters (useful for load balancing)

### Example Scenarios

```bash
# High-priority monitoring (try best datacenter first)
export DATACENTERS="GRA,SBG,BHS,DE"
export PREFERRED_ORDER="true"

# Load-balanced monitoring (random selection)
export DATACENTERS="US-WEST-OR,US-EAST-VA"
export PREFERRED_ORDER="false"

# Worldwide monitoring
export DATACENTERS="GRA,SBG,BHS,DE,UK,WAW,US-WEST-OR,US-EAST-VA"
```

## Available Datacenters

### US (ovh-us)
- `US-EAST-VA` - Virginia, USA
- `US-WEST-OR` - Oregon, USA

### Europe (ovh-eu)
- `GRA` - Gravelines, France
- `SBG` - Strasbourg, France
- `BHS` - Beauharnois, Canada
- `DE` - Frankfurt, Germany
- `UK` - London, UK
- `WAW` - Warsaw, Poland
- `SYD` - Sydney, Australia
- `SGP` - Singapore

## Log Output

The script provides detailed logging with colored output:

```
2024-01-15 10:30:15 [INFO] 🚀 Starting VPS availability monitor
2024-01-15 10:30:15 [INFO] Configuration:
2024-01-15 10:30:15 [INFO]   Plan Code: vps-2025-model2
2024-01-15 10:30:15 [INFO]   Country: US
2024-01-15 10:30:15 [INFO]   Datacenter: US-WEST-OR
2024-01-15 10:30:20 [INFO] 🔍 Checking availability for vps-2025-model2 in US-WEST-OR...
2024-01-15 10:30:25 [INFO] ❌ vps-2025-model2 is not available in US-WEST-OR
2024-01-15 10:35:25 [SUCCESS] ✅ vps-2025-model2 is AVAILABLE in US-WEST-OR!
2024-01-15 10:35:30 [INFO] 🚀 Placing order for vps-2025-model2 in US-WEST-OR...
2024-01-15 10:35:35 [SUCCESS] ✅ Order placed successfully!
2024-01-15 10:35:35 [SUCCESS] 🎉 Order completed successfully! Exiting monitor.
```

## Safety Features

- **Dry-run support** - Test configuration without placing real orders
- **Order confirmation** - Script will show order details before completion
- **Retry logic** - Handles temporary API failures gracefully
- **Rate limiting** - Respects OVH API rate limits with configurable intervals

## Troubleshooting

### Common Issues

1. **"Missing required environment variables"**
   - Ensure all OVH API credentials are properly set
   - Check that variables are exported in your shell

2. **"kimsufi-notifier binary not found"**
   - Build the binary first: `go build -o kimsufi-notifier .`
   - Ensure you're in the correct directory

3. **"Order failed" errors**
   - Check your OVH account balance
   - Verify datacenter availability
   - Ensure your OVH account is in good standing

4. **Permission denied**
   - Make the script executable: `chmod +x vps-monitor.sh`

### Debug Mode

For additional debugging, you can modify the script to enable verbose logging or run individual commands manually:

```bash
# Test availability check manually
./kimsufi-notifier check --plan-code vps-2025-model2 --country US --endpoint ovh-us

# Test order placement manually (dry-run)
./kimsufi-notifier order --plan-code vps-2025-model2 --country US --endpoint ovh-us -d US-WEST-OR --item-option os=option-linux --dry-run
```

### Testing the Monitor Script

Before running the monitor in production, test it with a dry-run approach:

```bash
# Test configuration without placing orders (set MAX_RETRIES=0)
export MAX_RETRIES=0
export CHECK_INTERVAL=60
./vps-monitor.sh

# This will check availability every minute and log results without attempting orders
```

## Security Notes

- Never commit API credentials to version control
- Use environment variables or secure credential storage
- Consider using OVH's IP restrictions for additional security
- Regularly rotate your API credentials

## License

This script is provided as-is for educational and personal use. Please ensure compliance with OVH's Terms of Service when using automated ordering.